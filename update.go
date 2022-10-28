package goc

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type ReleaseRes struct {
	Version string  `json:"tag_name"`
	Link    string  `json:"html_url"`
	Assets  []Asset `json:"assets"`
}

type Asset struct {
	DownloadLink string `json:"browser_download_url"`
}

func getLatestRelease() (*ReleaseRes, error) {
	url := "https://api.github.com/repos/leandergangso/goc/releases/latest"

	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch latest release: %v", err)
	}
	defer res.Body.Close()

	release := &ReleaseRes{}
	err = json.NewDecoder(res.Body).Decode(release)
	if err != nil {
		return nil, fmt.Errorf("unable to read response: %v", err)
	}

	return release, nil
}

func isBreakingVersionUpdate(current, latest string) bool {
	currentList := strings.Split(current, ".")
	latestList := strings.Split(latest, ".")

	return currentList[0] != latestList[0]
}

func downloadExe(url string) error {
	path, err := getExePath()
	if err != nil {
		return err
	}

	fmt.Println(path)

	data, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("unable to download update: %v", err)
	}
	defer data.Body.Close()

	f, err := os.OpenFile(path, os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("unable to open file: %v", err)
	}
	defer f.Close()

	_, err = io.Copy(f, data.Body)
	if err != nil {
		return fmt.Errorf("unable to copy: %v", err)
	}

	return nil
}

func getExePath() (string, error) {
	if runtime.GOOS == "windows" {
		return "", fmt.Errorf("can't execute commands on windows machine")
	}

	out, err := exec.Command("which", "goc").Output()
	if err != nil {
		return "", fmt.Errorf("unable to execute command: %v", err)
	}

	return string(out), nil
}
