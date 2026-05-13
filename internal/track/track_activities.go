package track

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/noemacias/punch/internal/config"
)

type Activity struct {
	ParentTitle string `json:"parentTitle,omitempty"`
	Project     int    `json:"project,omitempty"`
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Comment     string `json:"comment,omitempty"`
	Visible     bool   `json:"visible,omitempty"`
	Billable    bool   `json:"billable,omitempty"`
	MetaFields  []any  `json:"metaFields,omitempty"`
	Teams       []any  `json:"teams,omitempty"`
	Number      any    `json:"number,omitempty"`
	Color       string `json:"color,omitempty"`
}

type Activities struct {
	Settings *config.Settings
}

func NewActitivies(settings *config.Settings) *Activities {
	return &Activities{
		Settings: settings,
	}
}

func buildUrl(baseUrl, path string, params map[string]string) (*url.URL, error) {
	url, err := url.Parse(baseUrl)

	if err != nil {
		return nil, err
	}

	url.Path = path
	query := url.Query()

	for k, v := range params {
		query.Add(k, v)
	}

	url.RawQuery = query.Encode()

	return url, nil
}

func (a *Activities) List(limit string) ([]Activity, error) {

	params := map[string]string{
		"project": a.Settings.CustomerId,
		"term":    limit,
	}

	url, err := buildUrl(a.Settings.TrackingUrl, "/api/activities", params)

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
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", a.Settings.APIToken))

	resp, err := client.Do(req)

	if err != nil {
		return nil, err

	}

	defer resp.Body.Close()

	activities := []Activity{}

	err = json.NewDecoder(resp.Body).Decode(&activities)

	if err != nil {
		return nil, err

	}

	return activities, nil
}

func (a *Activities) Add(activity *Activity) error {

	url, err := buildUrl(a.Settings.TrackingUrl, "/api/activities", map[string]string{})

	if err != nil {
		return err
	}

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	data, err := json.Marshal(activity)

	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url.String(), bytes.NewBuffer(data))

	if err != nil {
		return err

	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", a.Settings.APIToken))

	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New("failed to crerate activitiy")
	}

	return nil
}
