package waffleio

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

type ProjectSource struct {
	ID       string `json:"_id"`
	Hooked   bool   `json:"hooked"`
	Private  bool   `json:"private"`
	Provider struct {
		ID         string `json:"_id"`
		Type       string `json:"type"`
		BaseAPIURL string `json:"baseApiUrl"`
		BaseURL    string `json:"baseUrl"`
		ClientID   string `json:"clientId"`
		Default    bool   `json:"default"`
		Public     bool   `json:"public"`
	} `json:"provider"`
	RepoPath          string    `json:"repoPath"`
	Type              string    `json:"type"`
	V                 int       `json:"__v"`
	LowerCaseRepoPath string    `json:"lowerCaseRepoPath"`
	UniqueID          string    `json:"uniqueId,omitempty"`
	LastSyncedAt      time.Time `json:"lastSyncedAt,omitempty"`
}

type Project struct {
	ID      string `json:"_id"`
	Account struct {
		ID              string    `json:"_id"`
		CreatedAt       time.Time `json:"createdAt"`
		V               int       `json:"__v"`
		Owner           string    `json:"owner"`
		LowerCaseOwner  string    `json:"lowerCaseOwner"`
		Type            string    `json:"type"`
		IsGrandfathered bool      `json:"isGrandfathered"`
	} `json:"account"`
	Name                    string          `json:"name"`
	LowerCaseName           string          `json:"lowerCaseName"`
	V                       int             `json:"__v"`
	Viewers                 []interface{}   `json:"viewers"`
	Team                    []string        `json:"team"`
	Sources                 []ProjectSource `json:"sources"`
	Integrations            []interface{}   `json:"integrations,omitempty"`
	RallyIntegrationEnabled bool            `json:"rallyIntegrationEnabled,omitempty"`
}

type Projects []Project

func (c *Client) GetProjects() (Projects, error) {

	req, err := http.NewRequest("GET", "https://api.waffle.io/user/projects", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var ps Projects
	err = json.Unmarshal(body, &ps)

	return ps, err
}

func (c *Client) DeleteSource(projectID string, sourceID string) error {
	req, err := http.NewRequest("DELETE", "https://api.waffle.io/projects/"+projectID+"/sources/"+sourceID, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New("Expected 200 from DELETE, got " + resp.Status)
	}
	return nil
}

func (c *Client) AddSource(projectID string, sourceName string) error {
	type Payload struct {
		Provider string `json:"provider"`
		Private  bool   `json:"private"`
		RepoPath string `json:"repoPath"`
		Type     string `json:"type"`
	}

	data := Payload{
		Provider: "5399ba229c4ef5963e000508",
		Private:  false,
		RepoPath: sourceName,
		Type:     "github",
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "https://api.waffle.io/projects/"+projectID+"/sources", body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("Expected 200 from POST, got " + resp.Status)
	}
	return nil
}
