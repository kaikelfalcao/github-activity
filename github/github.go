package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const API_URL = "https://api.github.com/"

type GitHubService struct {
	username string
	client http.Client
}

func NewGitHubService(username string) *GitHubService{
	return &GitHubService{
		username: username,
		client: http.Client{},
	}
}

func (ghs *GitHubService) UserExists() (bool, error) {
	response, err := ghs.client.Get(fmt.Sprintf(API_URL + "users/%s", ghs.username))

	if err != nil {
		return false, err
	}

	if response.StatusCode == http.StatusNotFound {
		return false, fmt.Errorf("User Not Found")
	}

	return true , nil
}

// Generated with https://transform.tools/json-to-go
type EventResponse struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Actor struct {
		ID           int    `json:"id"`
		Login        string `json:"login"`
		DisplayLogin string `json:"display_login"`
		GravatarID   string `json:"gravatar_id"`
		URL          string `json:"url"`
		AvatarURL    string `json:"avatar_url"`
	} `json:"actor"`
	Repo struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"repo"`
	Payload json.RawMessage `json:"payload"`
	Public    bool      `json:"public"`
	CreatedAt time.Time `json:"created_at"`
	Org       struct {
		ID         int    `json:"id"`
		Login      string `json:"login"`
		GravatarID string `json:"gravatar_id"`
		URL        string `json:"url"`
		AvatarURL  string `json:"avatar_url"`
	} `json:"org"`
}

func (ghs *GitHubService) GetActivities() (string, error) {
	var eventsUrl = fmt.Sprintf(API_URL + "users/%s/events", ghs.username)

	response, err := ghs.client.Get(eventsUrl)

	if err != nil {
		return "", err 
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
        return "", fmt.Errorf("HTTP request failed with status %d", response.StatusCode)
    }

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return "", err 
	}

	var events []EventResponse
	if err := json.Unmarshal(body, &events); err != nil {
		return "", fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	responseString, err := classifyEvents(ghs.username, events)

	if err != nil {
		return "", err 
	}

	return responseString, nil
}

func classifyEvents(username string, events []EventResponse) (string, error) {
    response := fmt.Sprintf("The last %d activities from %s\n", len(events), username)

	pushEvents := make(map[string]int)
	
    for _, event := range events {
        var payload map[string]interface{}
        if err := json.Unmarshal([]byte(event.Payload), &payload); err != nil {
            return "", fmt.Errorf("error unmarshaling JSON: %w", err)
        }

        switch event.Type {
        case "WatchEvent":
            response += fmt.Sprintf("- Starred %s\n", event.Repo.Name)
        case "PushEvent":
			pushEvents[event.Repo.Name] += int(payload["size"].(float64))
        case "CreateEvent":
            refType := payload["ref_type"].(string)
            if refType == "repository" {
                response += fmt.Sprintf("- Create %s\n", event.Repo.Name)
            } else {
                response += fmt.Sprintf("- Create %s on %s\n", refType, event.Repo.Name)
            }
		case "IssuesEvent":
			response += fmt.Sprintf("- %s a issue in %s\n", payload["action"], event.Repo.Name)
		case "PullRequestEvent":
			response += fmt.Sprintf("- %s a PR in %s\n", payload["action"], event.Repo.Name)
        }
    }

	for repo, commits := range pushEvents {
		response += fmt.Sprintf("- Pushed %d commits to %s\n", commits, repo)
	}

    return response, nil
}
