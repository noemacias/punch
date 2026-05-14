package time

import (
	"fmt"
	"log/slog"
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

	cmd.Flags().StringP("user", "u", "", "User")
	cmd.Flags().StringP("size", "s", "100", "The amount of entries for each page")
	cmd.Flags().StringP("activity", "a", "", "Activity ID to filter timesheets")
	cmd.Flags().String("begin", "", "Only records after this date will be included (format: HTML5 datetime-local, e.g. YYYY-MM-DDThh:mm:ss)")
	cmd.Flags().String("end", "", "Only records before this date will be included (format: HTML5 datetime-local, e.g. YYYY-MM-DDThh:mm:ss)")
	return &cmd
}

func (o *TimeListCommand) Run(cmd *cobra.Command, args []string) {

	user, _ := cmd.Flags().GetString("user")
	limit, _ := cmd.Flags().GetString("activity")
	configFile, _ := cmd.Flags().GetString("config")
	pageSize, _ := cmd.Flags().GetString("size")
	begin, _ := cmd.Flags().GetString("begin")
	end, _ := cmd.Flags().GetString("end")

	settings := config.NewSettings(configFile)

	timesheet := track.NewTimeSheet(settings)
	timesheets, err := timesheet.List(begin, end, pageSize, user, limit)

	if err != nil {
		slog.Error("Failed to list timesheets", "error", err.Error())
		return
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

	// Total
	grandTotal := map[int]map[int]int{}
	projectsTotal := map[int]int{}

	fmt.Printf("%-10v %-8v %-8v %-8v %-8v %v\n", "Date", "Begin", "End", "Duration", "Project", "Activity")
	for _, t := range timesheets {

		date := t.ParseTimeStamp(t.Begin).Format(time.DateOnly)
		begin := t.ParseTimeStamp(t.Begin)
		end := t.ParseTimeStamp(t.End)

		activity, _ := activitiesMap[t.Activity]

		fmt.Printf("%-10v %-8v %-8v %-8v %-8v %v\n", date, begin.Format(time.TimeOnly), end.Format(time.TimeOnly), end.Sub(begin), t.Project, fmt.Sprintf("%-4v %v", t.Activity, activity))

		if grandTotal[t.Project] == nil {
			grandTotal[t.Project] = map[int]int{}
		}
		grandTotal[t.Project][t.Activity] += t.Duration

		projectsTotal[t.Project] += t.Duration
	}

	fmt.Println()
	fmt.Println("Time spent on each project and activity")
	fmt.Printf("%-8v %-10v %v\n", "Project", "Duration", "Activity")
	for p, a := range grandTotal {
		for act, dur := range a {
			activity, _ := activitiesMap[act]
			fmt.Printf("%-8v %-10v %v\n", p, time.Duration(dur)*time.Second, fmt.Sprintf("%-4v %v", act, activity))
		}

	}

	fmt.Println()
	fmt.Println("Total time spent on each project")
	fmt.Printf("%-8v %-10v\n", "Project", "Duration")
	for p, dur := range projectsTotal {
		fmt.Printf("%-8v %-10v\n", p, time.Duration(dur)*time.Second)
	}
}
