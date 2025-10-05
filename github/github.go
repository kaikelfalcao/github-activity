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
	Payload struct {
		Action string `json:"action"`
	} `json:"payload"`
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

type Events struct {
	events []EventResponse
}

func (ghs *GitHubService) GetActivities() (string, error) {
	var eventsUrl = fmt.Sprintf(API_URL + "%s/events", ghs.username)

	response, err := ghs.client.Get(eventsUrl)

	if err != nil {
		return "", err 
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		body, err := io.ReadAll(response.Body)

		if err != nil {
			return "", err 
		}

		var events Events

		json.Unmarshal(body, &events)
	}

	return "", nil
}

func classifyEvents(events EventResponse) {
}