package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Activity struct {
	Id        int
	Duration  int64
	ProjectId int `yaml:"project_id"`
}

type Settings struct {
	APIToken    string     `yaml:"api_token"`
	TrackingUrl string     `yaml:"tracking_url"`
	CustomerId  string     `yaml:"project_id"`
	Activities  []Activity `yaml:""activities`
}

func NewSettings(configPath string) *Settings {

	s := Settings{}

	s.ReadConfigFile(configPath)

	s.OverrideFromEnv()
	return &s
}

func (s *Settings) ReadConfigFile(configPath string) {

	if configPath == "~/.config/punch.yaml" {
		home, err := os.UserHomeDir()

		if err != nil {
			return
		}

		configPath = filepath.Join(home, ".config", "punch.yaml")
	}

	data, err := os.ReadFile(configPath)

	if err != nil {
		return
	}

	err = yaml.Unmarshal(data, s)

	if err != nil {
		return
	}

}

func (s *Settings) OverrideFromEnv() {

	API_TOKEN := os.Getenv("TRACK_API_TOKEN")

	if API_TOKEN != "" {
		s.APIToken = API_TOKEN
	}

	TRACK_URL := os.Getenv("TRACK_URL")

	if TRACK_URL != "" {
		s.TrackingUrl = TRACK_URL
	}

}
