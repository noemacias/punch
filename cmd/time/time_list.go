package time

import (
	"fmt"
	"log/slog"
	"sort"
	"time"

	"github.com/noemacias/punch/internal/config"
	"github.com/noemacias/punch/internal/track"
	"github.com/spf13/cobra"
)

type TimeListCommand struct {
}

func NewTimeListCommand() *cobra.Command {

	o := TimeListCommand{}
	cmd := cobra.Command{
		Use: `list`,
		Run: func(cmd *cobra.Command, args []string) {
			o.Run(cmd, args)
		},
	}

	cmd.Flags().Int("last", 0, "Last days")
	cmd.Flags().StringP("user", "u", "", "User")
	cmd.Flags().StringSlice("users", []string{}, "Users --users 37,90")
	cmd.Flags().StringP("size", "s", "100", "The amount of entries for each page")
	cmd.Flags().StringP("activity", "a", "", "Activity ID to filter timesheets")
	cmd.Flags().String("begin", "", "Only records after this date will be included (format: HTML5 datetime-local, e.g. YYYY-MM-DDThh:mm:ss)")
	cmd.Flags().String("end", "", "Only records before this date will be included (format: HTML5 datetime-local, e.g. YYYY-MM-DDThh:mm:ss)")
	return &cmd
}

func (o *TimeListCommand) Run(cmd *cobra.Command, args []string) {

	user, _ := cmd.Flags().GetString("user")
	users, _ := cmd.Flags().GetStringSlice("users")
	limit, _ := cmd.Flags().GetString("activity")
	configFile, _ := cmd.Flags().GetString("config")
	pageSize, _ := cmd.Flags().GetString("size")
	begin, _ := cmd.Flags().GetString("begin")
	end, _ := cmd.Flags().GetString("end")
	inLastDays, _ := cmd.Flags().GetInt("last")

	settings := config.NewSettings(configFile)

	// Check if --last is set, it takes precedence over --begin and --end
	if inLastDays > 0 {
		begin, end = track.LastNDaysRangeStr(inLastDays, track.TimeLayoutSecond)
	}

	// Activities map
	activity := track.NewActitivies(settings)
	activitiesMap := map[int]string{}
	activities, err := activity.List("", "")

	if err != nil {
		slog.Error("Failed to list acitivies", "error", err)
		return
	}

	for _, a := range activities {
		activitiesMap[a.ID] = a.Name
	}

	// User map
	usersList, err := track.NewUsers(settings).List("")

	if err != nil {
		slog.Debug("Failed to list users", "error", err)
	}

	// Total
	grandTotal := map[int]map[int]int{}
	projectsTotal := map[int]int{}

	timesheet := track.NewTimeSheet(settings)
	timesheets := track.TimeSheetlist{}

	if len(users) > 0 {

		for _, u := range users {
			ts, err := timesheet.List(begin, end, pageSize, u, limit, []string{})
			if err != nil {
				slog.Debug("Failed to list timesheets for user", "user", u, "error", err)
				continue
			}
			timesheets = append(timesheets, ts...)
		}
	} else {
		timesheets, err = timesheet.List(begin, end, pageSize, user, limit, []string{})
		if err != nil {
			slog.Error("Failed to list timesheets", "error", err.Error())
			return
		}
	}

	sort.Slice(timesheets, func(i, j int) bool {
		return timesheets[i].Begin > timesheets[j].Begin

	})

	fmt.Printf("%-20v %-10v %-8v %-8v %-8v %-8v %v\n", "User", "Date", "Begin", "End", "Duration", "Project", "Activity")
	for _, t := range timesheets {

		date := t.ParseTimeStamp(t.Begin).Format(time.DateOnly)
		begin := t.ParseTimeStamp(t.Begin)
		end := t.ParseTimeStamp(t.End)

		activity, _ := activitiesMap[t.Activity]

		username := ""
		userInfo := usersList.Get(t.User)

		if userInfo.ID != 0 {
			username = userInfo.Alias
		}

		fmt.Printf("%-20v %-10v %-8v %-8v %-8v %-8v %v\n", username, date, begin.Format(time.TimeOnly), end.Format(time.TimeOnly), end.Sub(begin), t.Project, fmt.Sprintf("%-4v %v", t.Activity, activity))

		if grandTotal[t.Project] == nil {
			grandTotal[t.Project] = map[int]int{}
		}
		grandTotal[t.Project][t.Activity] += t.Duration

		projectsTotal[t.Project] += t.Duration
	}

	fmt.Println()
	fmt.Printf("Time range %v - %v\n", begin, end)
	fmt.Printf("Total time entries: %v\n", len(timesheets))
	fmt.Println()
	fmt.Println("Time spent on each project and activity")
	fmt.Printf("%-8v %-10v %v\n", "Project", "Duration", "Activity")
	for p, a := range grandTotal {

		entries := make([]ActiviDuration, 0, len(a))

		for k, v := range a {
			entries = append(entries, ActiviDuration{
				Id:       k,
				Duration: v,
			})
		}

		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Duration > entries[j].Duration
		})

		for _, e := range entries {
			activity, _ := activitiesMap[e.Id]
			fmt.Printf("%-8v %-10v %v\n", p, time.Duration(e.Duration)*time.Second, fmt.Sprintf("%-4v %v", e.Id, activity))
		}

	}

	fmt.Println()
	fmt.Println("Total time spent on each project")
	fmt.Printf("%-8v %-10v\n", "Project", "Duration")
	for p, dur := range projectsTotal {
		fmt.Printf("%-8v %-10v\n", p, time.Duration(dur)*time.Second)
	}

}

type ActiviDuration struct {
	Id       int
	Duration int
}
