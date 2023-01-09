package goc

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// option to search self w/jql?
// e.g goc jira -search 'assignee=currentuser() and status in ("In Progress")'

// option to get single issue by cc-code
// e.g goc jira cc-1234 (sets it in progress)
// e.g goc jira cc-1234 -info (shows issue info)

// get available transitions
// GET https://houseofcontrol.atlassian.net/rest/api/3/issue/CC-14871/transitions

// update issue state
// POST https://houseofcontrol.atlassian.net/rest/api/3/issue/CC-14871/transitions
// body:
// {
//     "transition": {
//         "id": 621
//     }
// }

const baseURL string = "https://houseofcontrol.atlassian.net/rest/api/3"

type issueRes struct {
	Issues []struct {
		Id     string `json:"key"`
		Fields struct {
			Summary string `json:"summary"`
			Status  struct {
				Name string `json:"name"`
			} `json:"status"`
		} `json:"fields"`
	} `json:"issues"`
}

func addAuthToRequest(data *FileData, req *http.Request) {
	if data.Jira == nil {
		setJiraAuth(data)
	}

	req.SetBasicAuth(data.Jira.Username, data.Jira.Token)
}

func setJiraAuth(data *FileData) {
	// get jira username
	fmt.Print("What is your Jira username (email): ")
	reader := bufio.NewReader(os.Stdin)
	username, _ := reader.ReadString('\n')
	username = strings.Replace(username, "\n", "", -1)

	// get jira toeken
	fmt.Print("Also need an API token (https://id.atlassian.com/manage-profile/security/api-tokens): ")
	reader = bufio.NewReader(os.Stdin)
	token, _ := reader.ReadString('\n')
	token = strings.Replace(token, "\n", "", -1)
	fmt.Println()

	// set data
	data.Jira = &JiraAuth{
		Username: username,
		Token:    token,
	}

	writeToFile(data)
}

func JiraGetOwnIssues(ctx context.Context, data *FileData) (*issueRes, error) {
	jql := url.PathEscape("assignee=currentuser() AND status IN ('In Progress','To Do','In Review')")
	url := baseURL + "/search?maxResults=10&fields=summary,status&jql=" + jql

	req, err := GetRequest(ctx, false, url, nil)
	if err != nil {
		return nil, err
	}
	addAuthToRequest(data, req)

	body := &issueRes{}
	SendRequest(req, body)

	return body, nil
}
