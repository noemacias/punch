package time

import (
	"fmt"
	"log/slog"
	"sort"
	"strconv"
	"time"

	"github.com/noemacias/punch/internal/config"
	"github.com/noemacias/punch/internal/track"
	"github.com/spf13/cobra"
)

type TimeReportCommand struct{}

func NewTimeReportCommand() *cobra.Command {

	o := TimeReportCommand{}
	cmd := cobra.Command{
		Use: `gaps`,
		Run: func(cmd *cobra.Command, args []string) {
			o.Run(cmd, args)
		},
	}

	cmd.Flags().Int("last", 0, "Last days")
	cmd.Flags().StringSlice("users", []string{}, "Users --users 37,90")
	cmd.Flags().String("begin", "", "Only records after this date will be included (format: HTML5 datetime-local, e.g. YYYY-MM-DDThh:mm:ss)")
	cmd.Flags().String("end", "", "Only records before this date will be included (format: HTML5 datetime-local, e.g. YYYY-MM-DDThh:mm:ss)")

	return &cmd
}

func (o *TimeReportCommand) Run(cmd *cobra.Command, args []string) {

	debug, _ := cmd.Flags().GetBool("debug")
	users, _ := cmd.Flags().GetStringSlice("users")
	configFile, _ := cmd.Flags().GetString("config")
	begin, _ := cmd.Flags().GetString("begin")
	end, _ := cmd.Flags().GetString("end")
	inLastDays, _ := cmd.Flags().GetInt("last")

	settings := config.NewSettings(configFile)

	var (
		endTime   time.Time
		beginTime time.Time
	)

	// Check if --last is set, it takes precedence over --begin and --end
	if inLastDays > 0 {
		begin, end = track.LastNDaysRangeStr(inLastDays, track.TimeLayoutSecond)
	}

	if begin == "" || end == "" {
		slog.Error("Missing required flags --begin, --end or use --last ")
		return
	}

	endTime, err := time.Parse(track.TimeLayoutSecond, end)

	if err != nil {
		slog.Error("Failed to parse --end", "error", err.Error())
		return
	}

	beginTime, err = time.Parse(track.TimeLayoutSecond, begin)

	if err != nil {
		slog.Error("Failed to parse --begin", "error", err.Error())
		return
	}

	if debug {
		fmt.Println("Time Range", begin, "-", end)
	}

	days, err := track.WeekdaysBetween(beginTime.Format(time.DateOnly), endTime.Format(time.DateOnly))

	if err != nil {
		slog.Error("Failed to get weekdays for time range --begin - --end", "error", err.Error())
		return
	}

	daysMapSlice := []string{}
	userDayMap := map[int]map[string]track.TimesheetEntry{}
	userHourMap := map[int]map[string]int{}

	for _, d := range days {
		daysMapSlice = append(daysMapSlice, d.Format(time.DateOnly))
	}

	sort.Sort(sort.Reverse(sort.StringSlice(daysMapSlice)))

	// Users
	usersvc := track.NewUsers(settings)

	usersList, err := usersvc.List("")

	if err != nil {
		slog.Error("Failed to get users list", "error", err)
		return
	}

	timesheet := track.NewTimeSheet(settings)

	for _, user := range users {
		userTimeSheet, err := timesheet.List(begin, end, "1000", user, "", []string{})

		if err != nil {
			return
		}

		for _, u := range userTimeSheet {

			if userDayMap[u.User] == nil {
				userDayMap[u.User] = map[string]track.TimesheetEntry{}
			}

			uBegin, _ := time.Parse(track.TimeLayoutRFC3339TZ, u.Begin)
			userDayMap[u.User][uBegin.Format(time.DateOnly)] = u

			if userHourMap[u.User] == nil {
				userHourMap[u.User] = map[string]int{}
			}

			userHourMap[u.User][uBegin.Format(time.DateOnly)] += u.Duration

			if debug {
				userinfo := usersList.Get(u.User)
				fmt.Println("Timesheet", userinfo.ID, userinfo.Alias, u.Begin)
			}
		}

	}

	// Check missing dates
	for _, d := range daysMapSlice {

		for _, user := range users {
			userId, _ := strconv.Atoi(user)
			userDayMap, ok := userDayMap[userId]

			if !ok {
				continue

			}

			userInfo := usersList.Get(userId)

			_, ok = userDayMap[d]

			if !ok {
				fmt.Println("Missing Entry:", d, userInfo.ID, userInfo.Alias, userInfo.Username)
				continue
			}

		}
	}

	fmt.Println()
	// Check missing hours
	for _, d := range daysMapSlice {

		for _, user := range users {
			userId, _ := strconv.Atoi(user)
			userHourMap, ok := userHourMap[userId]

			if !ok {
				continue
			}

			userInfo := usersList.Get(userId)

			duration, ok := userHourMap[d]

			if !ok {
				continue
			}

			if duration < 28800 {
				fmt.Println("Insufficient Hours:", d, "expected=8h", "actual=", time.Duration(time.Second*time.Duration(duration)), userInfo.ID, userInfo.Alias, userInfo.Username)
				continue
			}

		}
	}

}
