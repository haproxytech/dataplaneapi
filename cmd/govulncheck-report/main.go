package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Note struct {
	Body string `json:"body"`
	ID   int    `json:"id"`
}

type Issue struct {
	Title string `json:"title"`
	State string `json:"state"`
	ID    int    `json:"id"`
	IID   int    `json:"iid"`
}

type Thread struct {
	Body string `json:"body"`
	ID   int    `json:"id"`
	IID  int    `json:"iid"`
}

const baseURL = "https://gitlab.int.haproxy.com/api/v4"

func main() {
	fmt.Print(hello)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == "source" {
				x := a.Value
				src := x.Any().(*slog.Source)
				path := strings.Split(src.File, "/")
				src.File = path[len(path)-1]
				return slog.Attr{
					Key:   "source",
					Value: slog.AnyValue(src),
				}
			}
			return a
		},
	}))
	slog.SetDefault(logger)

	slog.Info("Starting GoVulnCheck")
	mergeRequestMode := false
	currentBranch := os.Getenv("CI_COMMIT_BRANCH")
	if currentBranch == "" {
		currentBranch = os.Getenv("CI_MERGE_REQUEST_SOURCE_BRANCH_NAME")
		mergeRequestMode = true
	}
	if currentBranch == "" {
		cmd := exec.Command("git", "branch", "--show-current")
		out, err := cmd.Output()
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		currentBranch = string(out)
	}
	slog.Info("Current branch: " + currentBranch)

	cmd := exec.Command("govulncheck", "./...")
	out, _ := cmd.Output()

	vulnMessage := string(out)
	fmt.Println(vulnMessage)
	noVuln := false
	if !strings.Contains(vulnMessage, "Vulnerability #") {
		noVuln = true
	}

	if currentBranch == "" {
		if strings.Contains(vulnMessage, "Vulnerability #") {
			slog.Error(vulnMessage)
			os.Exit(1)
		}
		slog.Info("Current branch is empty, exiting...")
		os.Exit(0)
	}

	if mergeRequestMode {
		if strings.Contains(vulnMessage, "Vulnerability #") {
			os.Exit(1)
		}
		slog.Info("no vulnerabilities found")
		os.Exit(0)
	}

	token := os.Getenv("GITLAB_GOPHERS_TOKEN")
	projectID := "534"
	title := "Data Plane API: GoVulnCheck: Branch: " + strings.Trim(currentBranch, "\n")

	userID, err := fetchUserID(token)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	issues, err := fetchOpenIssues(projectID, userID, token)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	found := false
	var issueFound Issue
	for _, issue := range issues {
		if issue.Title == title && issue.State == "opened" {
			found = true
			issueFound = issue
			break
		}
	}
	vulnMessage = "```\n" + vulnMessage + "\n```"
	if found {
		if noVuln {
			closeTheIssue(baseURL, token, projectID, issueFound.IID, "No vulnerability found.")
		} else {
			addCommentToIssue(baseURL, token, projectID, issueFound.IID, vulnMessage)
		}
	} else if !noVuln {
		createIssue(baseURL, token, projectID, title, vulnMessage)
	}
	slog.Info("GoVulnCheck done.")
}

func createIssue(baseURL, token, projectID string, title, commentBody string) {
	slog.Info("Active issue with title '" + title + "' not found in project " + projectID)
	// Create the issue here
	issueData := map[string]interface{}{
		"title":       title,
		"description": commentBody,
		"labels":      "bot,critical",
		// Add other fields as needed
	}
	issueDataBytes, err := json.Marshal(issueData)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/projects/%s/issues", baseURL, url.PathEscape(projectID)), bytes.NewBuffer(issueDataBytes))
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	req.Header.Add("Private-Token", token)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	var issue Issue
	err = json.Unmarshal(body, &issue)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	slog.Info("Issue created with ID " + strconv.Itoa(issue.ID))
}

func closeTheIssue(baseURL, token, projectID string, issueIID int, commentBody string) {
	addCommentToIssue(baseURL, token, projectID, issueIID, commentBody)

	client := &http.Client{}
	issueData := map[string]interface{}{
		"state_event": "close",
	}
	issueDataBytes, err := json.Marshal(issueData)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/projects/%s/issues/%d", baseURL, url.PathEscape(projectID), issueIID), bytes.NewBuffer(issueDataBytes))
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	req.Header.Add("Private-Token", token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	var issue Issue
	err = json.Unmarshal(body, &issue)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	slog.Info("Issue " + strconv.Itoa(issue.IID) + " closed")
}

func addCommentToIssue(baseURL, token, projectID string, issueIID int, commentBody string) {
	client := &http.Client{}
	noteData := map[string]interface{}{
		"body": commentBody,
	}
	noteDataBytes, err := json.Marshal(noteData)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/projects/%s/issues/%d/notes", baseURL, url.PathEscape(projectID), issueIID), bytes.NewBuffer(noteDataBytes))
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	req.Header.Add("Private-Token", token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	var note Note
	err = json.Unmarshal(body, &note)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	slog.Info("Comment added with ID " + strconv.Itoa(note.ID))
}

func fetchOpenIssues(projectID string, userID int, accessToken string) ([]Issue, error) {
	perPage := 20 // Number of issues to fetch per page

	var allIssues []Issue
	page := 1

	for {
		url := fmt.Sprintf("%s/projects/%s/issues?state=opened&author_id=%s&page=%d&per_page=%d", baseURL, projectID, strconv.Itoa(userID), page, perPage)

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Authorization", "Bearer "+accessToken)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var issues []Issue
		err = json.Unmarshal(body, &issues)
		if err != nil {
			return nil, err
		}

		allIssues = append(allIssues, issues...)

		// Check if there are more pages
		linkHeader := resp.Header.Get("Link")
		if !strings.Contains(linkHeader, `rel="next"`) {
			break
		}

		page++
	}

	return allIssues, nil
}

func fetchUserID(accessToken string) (int, error) {
	url := baseURL + "/user"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var user struct {
		ID int `json:"id"`
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}

const hello = `
                        _            _               _
  __ _  _____   ___   _| |_ __   ___| |__   ___  ___| | __
 / _` + "`" + ` |/ _ \ \ / / | | | | '_ \ / __| '_ \ / _ \/ __| |/ /
| (_| | (_) \ V /| |_| | | | | | (__| | | |  __/ (__|   <
 \__, |\___/ \_/  \__,_|_|_| |_|\___|_| |_|\___|\___|_|\_\
 |___/
`
