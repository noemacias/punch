package track

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/noemacias/punch/internal/config"
)

type User struct {
	APIToken      bool   `json:"apiToken,omitempty"`
	Initials      string `json:"initials,omitempty"`
	ID            int    `json:"id,omitempty"`
	Alias         string `json:"alias,omitempty"`
	Title         string `json:"title,omitempty"`
	Username      string `json:"username,omitempty"`
	AccountNumber string `json:"accountNumber,omitempty"`
	Enabled       bool   `json:"enabled,omitempty"`
	Color         any    `json:"color,omitempty"`
}

type UserList []User

func (u UserList) Get(id int) User {

	for _, user := range u {
		if user.ID == id {
			return user
		}
	}

	return User{}
}

type Users struct {
	Settings *config.Settings
}

func NewUsers(settings *config.Settings) *Users {

	return &Users{
		Settings: settings,
	}
}

func (u *Users) List(term string) (UserList, error) {

	params := map[string]string{
		"term": term,
	}

	url, err := buildUrl(u.Settings.TrackingUrl, "/api/users", params)

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
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", u.Settings.APIToken))

	resp, err := client.Do(req)

	if err != nil {
		return nil, err

	}

	defer resp.Body.Close()

	users := UserList{}

	err = json.NewDecoder(resp.Body).Decode(&users)

	if err != nil {
		return nil, err
	}

	return users, nil
}
