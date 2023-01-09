package goc

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

func GetRequest(ctx context.Context, isPOST bool, url string, body []byte) (*http.Request, error) {
	method := "GET"
	if isPOST {
		method = "POST"
	}

	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if len(body) > 0 {
		req.Body = io.NopCloser(bytes.NewBuffer(body))
	}

	req.Header.Set("Content-Type", "application/json")

	return req, err
}

func SendRequest(req *http.Request, dataPointer any) error {
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, dataPointer)
	if err != nil {
		return err
	}

	return nil
}
