package track

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/noemacias/punch/internal/config"
)

type ProjectList []Project

type Project struct {
	ParentTitle string `json:"parentTitle,omitempty"`
	Customer    int    `json:"customer,omitempty"`
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Start       any    `json:"start,omitempty"`
	End         any    `json:"end,omitempty"`
	Comment     any    `json:"comment,omitempty"`
	Visible     bool   `json:"visible,omitempty"`
	Billable    bool   `json:"billable,omitempty"`
	MetaFields  []any  `json:"metaFields,omitempty"`
	Teams       []struct {
		ID    int    `json:"id,omitempty"`
		Name  string `json:"name,omitempty"`
		Color string `json:"color,omitempty"`
	} `json:"teams,omitempty"`
	GlobalActivities bool   `json:"globalActivities,omitempty"`
	Number           any    `json:"number,omitempty"`
	Color            string `json:"color,omitempty"`
}

type Projects struct {
	Settings *config.Settings
}

func NewProject(settings *config.Settings) *Projects {
	return &Projects{
		Settings: settings,
	}
}

func (p *Projects) List() (ProjectList, error) {
	params := map[string]string{
		// "project": project,
		// "term":    limit,
	}

	url, err := buildUrl(p.Settings.TrackingUrl, "/api/projects", params)

	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("GET", url.String(), nil)

	if err != nil {
		return nil, err

	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", p.Settings.APIToken))

	resp, err := client.Do(req)

	if err != nil {
		return nil, err

	}

	defer resp.Body.Close()

	projects := ProjectList{}

	err = json.NewDecoder(resp.Body).Decode(&projects)

	if err != nil {
		return nil, err

	}

	return projects, nil
}
