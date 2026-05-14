package time

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/noemacias/punch/internal/config"
	"github.com/noemacias/punch/internal/track"
	"github.com/noemacias/punch/internal/utils"
	"github.com/spf13/cobra"
)

type TimeAddCommand struct {
}

func NewTimeAddCommand() *cobra.Command {

	o := TimeAddCommand{}
	cmd := cobra.Command{
		Use: `add`,
		Run: func(cmd *cobra.Command, args []string) {
			o.Run(cmd, args)
		},
	}

	cmd.Flags().String("begin", "", "Only records after this date will be included (format: HTML5 datetime-local, e.g. YYYY-MM-DDThh:mm:ss)")
	cmd.Flags().String("end", "", "Only records before this date will be included (format: HTML5 datetime-local, e.g. YYYY-MM-DDThh:mm:ss)")
	return &cmd
}

func (o *TimeAddCommand) Run(cmd *cobra.Command, args []string) {

	configFile, _ := cmd.Flags().GetString("config")
	begin, _ := cmd.Flags().GetString("begin")
	end, _ := cmd.Flags().GetString("end")
	settings := config.NewSettings(configFile)

	if begin == "" || end == "" {
		slog.Error("Missing required flags --begin --end")
		os.Exit(1)
	}
	days, err := track.WeekdaysBetween(begin, end)

	if err != nil {
		slog.Error("Failed to parse time", "error", err.Error())
		return
	}

	if len(days) == 0 {
		return
	}

	timeSheetList := track.TimeSheetlist{}

	for _, day := range days {

		day = day.Add(time.Hour * 9)

		for _, a := range settings.Activities {

			endtime := day.Add(time.Hour * time.Duration(a.Duration))

			timeSheetList = append(timeSheetList, track.TimesheetEntry{
				Activity: a.Id,
				Project:  1,
				Begin:    day.Format(track.Timelayout2),
				End:      endtime.Format(track.Timelayout2),
			})

			day = day.Add(time.Hour * time.Duration(a.Duration))

		}

	}

	activity := track.NewActitivies(settings)

	activities, err := activity.List("", "")
	activitiesMap := map[int]string{}

	if err != nil {
		slog.Error("Failed to list acitivies", "error", err)
	}

	for _, a := range activities {
		activitiesMap[a.ID] = a.Name
	}

	fmt.Printf("%-8v %-10v %-16v %-16v %v\n", "Project", "Weekday", "Begin", "End", "Activity")
	for _, t := range timeSheetList {
		day, _ := time.Parse(track.Timelayout2, t.Begin)
		activity, _ := activitiesMap[t.Activity]
		fmt.Printf("%-8v %-10v %-16v %-16v %v\n", t.Project, day.Weekday(), t.Begin, t.End, fmt.Sprintf("%v - %v", t.Activity, activity))

	}

	fmt.Println()
	confirm, _ := utils.ReadInput("Conitnue: [y/n] ")
	confirm = strings.ToLower(confirm)

	if confirm != "y" && confirm != "yes" {
		return
	}

	timesheet := track.NewTimeSheet(settings)

	for _, t := range timeSheetList {
		activity, _ := activitiesMap[t.Activity]
		fmt.Println("Creating timesheet: ", t.Begin, activity)

		err := timesheet.Add(&t)

		if err != nil {
			slog.Error("Failed to create timesheet", "error", err.Error())
			continue
		}
	}

}
