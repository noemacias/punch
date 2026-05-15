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

const (
	TimeLayoutRFC3339TZ = "2006-01-02T15:04:05-0700"
	TimeLayoutMinute    = "2006-01-02T15:04"
	TimeLayoutSecond    = "2006-01-02T15:04:05"
)

func WeekdaysBetween(beginStr, endStr string) ([]time.Time, error) {
	begin, err := time.Parse("2006-01-02", beginStr)
	if err != nil {
		return nil, err
	}

	end, err := time.Parse("2006-01-02", endStr)
	if err != nil {
		return nil, err
	}

	var days []time.Time

	for d := begin; !d.After(end); d = d.AddDate(0, 0, 1) {
		switch d.Weekday() {
		case time.Saturday, time.Sunday:
			continue
		default:
			days = append(days, d)
		}
	}

	return days, nil
}

type TimeSheetlist []TimesheetEntry

type TimesheetEntry struct {
	Activity     int     `json:"activity,omitempty"`
	Project      int     `json:"project,omitempty"`
	User         int     `json:"user,omitempty"`
	Tags         []any   `json:"tags,omitempty"`
	ID           int     `json:"id,omitempty"`
	Begin        string  `json:"begin,omitempty"`
	End          string  `json:"end,omitempty"`
	Duration     int     `json:"duration,omitempty"`
	Description  string  `json:"description,omitempty"`
	Rate         float64 `json:"rate,omitempty"`
	InternalRate float64 `json:"internalRate,omitempty"`
	Exported     bool    `json:"exported,omitempty"`
	Billable     bool    `json:"billable,omitempty"`
	MetaFields   []any   `json:"metaFields,omitempty"`
}

func (t *TimesheetEntry) ParseTimeStamp(timestamp string) time.Time {

	time, _ := time.Parse(TimeLayoutRFC3339TZ, timestamp)
	return time

}

type TimeSheet struct {
	Settings *config.Settings
}

func NewTimeSheet(settings *config.Settings) *TimeSheet {

	return &TimeSheet{
		Settings: settings,
	}
}

func (t *TimeSheet) List(begin, end string, pageSize string, user string, activity string, users []string) (TimeSheetlist, error) {

	params := url.Values{}
	params.Set("begin", begin)
	params.Set("end", end)
	params.Set("size", pageSize)
	params.Set("user", user)
	params.Set("activity", activity)

	for _, u := range users {
		params.Add("users[]", u)
	}

	url, err := buildUrl2(t.Settings.TrackingUrl, "/api/timesheets", params)

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
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", t.Settings.APIToken))

	resp, err := client.Do(req)

	if err != nil {
		return nil, err

	}

	defer resp.Body.Close()

	activities := TimeSheetlist{}

	err = json.NewDecoder(resp.Body).Decode(&activities)

	if err != nil {
		return nil, err

	}

	return activities, nil
}

func (t *TimeSheet) Add(timesheet *TimesheetEntry) error {

	url, err := buildUrl(t.Settings.TrackingUrl, "/api/timesheets", map[string]string{})

	if err != nil {
		return err
	}

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	data, err := json.Marshal(timesheet)

	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url.String(), bytes.NewBuffer(data))

	if err != nil {
		return err

	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", t.Settings.APIToken))

	resp, err := client.Do(req)

	if err != nil {
		return err

	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("failed to create timesheet")
	}

	return nil
}
