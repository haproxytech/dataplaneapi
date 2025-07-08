// Copyright 2019 HAProxy Technologies LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/semver"
	"github.com/joho/godotenv"
)

type Issue struct {
	Title string `json:"title"`
	State string `json:"state"`
	ID    int    `json:"id"`
	IID   int    `json:"iid"`
}

type Note struct {
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Body         string    `json:"body"`
	Attachment   string    `json:"attachment"`
	Author       Author    `json:"author"`
	ID           int       `json:"id"`
	ProjectID    int       `json:"project_id"`
	System       bool      `json:"system"`
	Resolvable   bool      `json:"resolvable"`
	Confidential bool      `json:"confidential"`
	Internal     bool      `json:"internal"`
}

type Author struct {
	CreatedAt time.Time `json:"created_at"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	State     string    `json:"state"`
	ID        int       `json:"id"`
}

type Thread struct {
	Body string `json:"body"`
	ID   int    `json:"id"`
	IID  int    `json:"iid"`
}

// GitlabLabel defines the structure for a GitLab label.
// It includes common fields; you are primarily using Name.
type GitlabLabel struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Color       string `json:"color"`
	Description string `json:"description,omitempty"` // omitempty handles cases where description might be null or absent
}

type Support struct {
	Versions []string `yaml:"versions"`
}

type MergeRequest struct {
	Description string `json:"description"`
}

var baseURL string

const LABEL_COLOR = "#8fbc8f" //nolint:stylecheck

func main() {
	_ = godotenv.Overload()
	fmt.Print(hello) //nolint:forbidigo

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

	baseURL = os.Getenv("CI_API_V4_URL")
	if baseURL == "" {
		slog.Error("CI_API_V4_URL not set")
		os.Exit(1)
	}

	docs, err := GetBranches()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	var versions []*semver.Version
	for _, r := range docs {
		v, err := semver.NewVersion(r)
		if err != nil {
			slog.Debug("could not parse branch name as semver, skipping", "branch", r, "error", err)
			continue
		}
		versions = append(versions, v)
	}

	sort.Sort(semver.Collection(versions))
	// leave only last three since only those are maintained
	if len(versions) > 3 {
		versions = versions[len(versions)-3:]
	}

	gitlabToken := os.Getenv("GITLAB_TOKEN")

	CI_MERGE_REQUEST_IID_STR := os.Getenv("CI_MERGE_REQUEST_IID") //nolint:stylecheck
	if CI_MERGE_REQUEST_IID_STR == "" {
		slog.Error("CI_MERGE_REQUEST_IID not set")
		os.Exit(1)
	}
	CI_MERGE_REQUEST_IID, err := strconv.Atoi(CI_MERGE_REQUEST_IID_STR) //nolint:stylecheck
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	CI_PROJECT_ID := os.Getenv("CI_PROJECT_ID") //nolint:stylecheck
	if CI_PROJECT_ID == "" {
		slog.Error("CI_PROJECT_ID not set")
		os.Exit(1)
	}
	question := `<!-- MR BACKPORT QUESTION -->` + "\n" + "Does this needs backport ? \n| Version | label |\n|:--:|:--:|"
	backportLabels := map[string]struct{}{
		"backport-ee": {},
	}
	for _, version := range versions {
		ver := strconv.Itoa(int(version.Major())) + "." + strconv.Itoa(int(version.Minor()))
		question += "\n" + "| " + ver + " | " + "backport-" + ver + " |"
		backportLabels["backport-"+ver] = struct{}{}
	}
	// ee
	question += "\n" + "| EE | " + "backport-ee |"
	question += "\n\n" + "please add labels for backporting."

	mr, err := getMergeRequest(baseURL, gitlabToken, CI_PROJECT_ID, CI_MERGE_REQUEST_IID)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	if strings.Contains(mr.Description, "<!-- BOT DEPENDABOT -->") {
		slog.Info("Dependabot MR detected, skipping backport check.")
		os.Exit(0)
	}

	notes, err := getMergeRequestComments(baseURL, gitlabToken, CI_PROJECT_ID, CI_MERGE_REQUEST_IID)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	index := slices.IndexFunc(notes, func(note Note) bool {
		return strings.Contains(note.Body, "<!-- MR BACKPORT QUESTION -->")
	})
	if index == -1 {
		// add missing labels
		err = getProjectlabels(backportLabels, CI_PROJECT_ID)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		slog.Info("No backport question found, creating one as thread")
		startThreadOnMergeRequest(baseURL, gitlabToken, CI_PROJECT_ID, CI_MERGE_REQUEST_IID, question)
	}
}

func startThreadOnMergeRequest(baseURL, token, projectID string, mergeRequestIID int, threadBody string) {
	client := &http.Client{}
	threadData := map[string]interface{}{
		"body": threadBody,
	}
	threadDataBytes, err := json.Marshal(threadData)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost,
		fmt.Sprintf("%s/projects/%s/merge_requests/%d/discussions", baseURL, url.PathEscape(projectID), mergeRequestIID), bytes.NewBuffer(threadDataBytes))
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	req.Header.Add("PRIVATE-TOKEN", token) //nolint:canonicalheader
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	defer resp.Body.Close()
}

func getMergeRequest(baseURL, token, projectID string, mergeRequestIID int) (*MergeRequest, error) {
	client := &http.Client{}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet,
		fmt.Sprintf("%s/projects/%s/merge_requests/%d", baseURL, url.PathEscape(projectID), mergeRequestIID), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("PRIVATE-TOKEN", token) //nolint:canonicalheader

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get merge request: status %s, body: %s", resp.Status, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var mr MergeRequest
	err = json.Unmarshal(body, &mr)
	if err != nil {
		return nil, err
	}

	return &mr, nil
}

func getMergeRequestComments(baseURL, token, projectID string, mergeRequestIID int) ([]Note, error) {
	client := &http.Client{}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet,
		fmt.Sprintf("%s/projects/%s/merge_requests/%d/notes", baseURL, url.PathEscape(projectID), mergeRequestIID), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("PRIVATE-TOKEN", token) //nolint:canonicalheader

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var notes []Note
	err = json.Unmarshal(body, &notes)
	if err != nil {
		return nil, err
	}

	return notes, nil
}

func getProjectlabels(backportLabels map[string]struct{}, projectID string) error {
	client := &http.Client{}
	token := os.Getenv("GITLAB_TOKEN")
	if token == "" {
		return errors.New("GITLAB_TOKEN not set")
	}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet,
		fmt.Sprintf("%s/projects/%s/labels", baseURL, url.PathEscape(projectID)), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Add("PRIVATE-TOKEN", token) //nolint:canonicalheader
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to get project labels: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get project labels: status %s, body: %s", resp.Status, string(body))
	}

	var projectLabels []GitlabLabel
	err = json.Unmarshal(body, &projectLabels)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response body (status %s): %w. Body: %s", resp.Status, err, string(body))
	}

	for _, label := range projectLabels {
		_, ok := backportLabels[label.Name]
		if ok {
			delete(backportLabels, label.Name)
		}
	}
	for label := range backportLabels {
		// Create the label if it doesn't exist
		labelData := map[string]string{
			"name":        label,
			"color":       LABEL_COLOR,
			"description": "Label for backporting to " + label + " branch",
		}
		labelDataBytes, err := json.Marshal(labelData)
		if err != nil {
			return fmt.Errorf("failed to marshal label data: %w", err)
		}
		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost,
			fmt.Sprintf("%s/projects/%s/labels", baseURL, url.PathEscape(projectID)), bytes.NewBuffer(labelDataBytes))
		if err != nil {
			return fmt.Errorf("failed to create request to create label: %w", err)
		}
		req.Header.Add("PRIVATE-TOKEN", token) //nolint:canonicalheader
		req.Header.Add("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to create label %s: %w", label, err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusCreated {
			return fmt.Errorf("failed to create label %s, status code: %d", label, resp.StatusCode)
		}
	}

	return nil
}

func GetBranches() ([]string, error) {
	projectID := os.Getenv("CI_PROJECT_ID")
	token := os.Getenv("GITLAB_TOKEN")

	if baseURL == "" || projectID == "" || token == "" {
		return nil, errors.New("one or more required environment variables are not set: CI_API_V4_URL, CI_PROJECT_ID, GITLAB_TOKEN")
	}

	var branches []string
	client := &http.Client{}

	nextPageURL := fmt.Sprintf("%s/projects/%s/repository/branches", baseURL, url.PathEscape(projectID))

	for nextPageURL != "" {
		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, nextPageURL, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}
		req.Header.Add("PRIVATE-TOKEN", token) //nolint:canonicalheader

		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to get branches: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			return nil, fmt.Errorf("failed to get branches: status %s, body: %s", resp.Status, string(body))
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			resp.Body.Close()
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}
		resp.Body.Close()

		type Branch struct {
			Name string `json:"name"`
		}
		var gitlabBranches []Branch
		err = json.Unmarshal(body, &gitlabBranches)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
		}

		for _, b := range gitlabBranches {
			branches = append(branches, b.Name)
		}

		// Check for the next page using Link header
		linkHeader := resp.Header.Get("Link")
		if linkHeader == "" {
			nextPageURL = ""
			continue
		}

		links := strings.Split(linkHeader, ",")
		nextPageURL = ""
		for _, link := range links {
			parts := strings.Split(strings.TrimSpace(link), ";")
			if len(parts) == 2 && strings.TrimSpace(parts[1]) == `rel="next"` {
				nextPageURL = strings.Trim(parts[0], "<>")
				break
			}
		}
	}

	return branches, nil
}

const hello = `
 __  __ ____         _               _
|  \/  |  _ \    ___| |__   ___  ___| | _____ _ __
| |\/| | |_) |  / __| '_ \ / _ \/ __| |/ / _ \ '__|
| |  | |  _ <  | (__| | | |  __/ (__|   <  __/ |
|_|  |_|_| \_\  \___|_| |_|\___|\___|_|\_\___|_|

`
