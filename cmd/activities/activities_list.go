package activities

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/noemacias/punch/internal/config"
	"github.com/noemacias/punch/internal/track"
	"github.com/spf13/cobra"
)

type ActivitiesListCmd struct {
}

func NewActitiviesList() *cobra.Command {

	o := ActivitiesListCmd{}

	cmd := cobra.Command{
		Use: `list`,
		Run: func(cmd *cobra.Command, args []string) {
			o.Run(cmd, args)
		},
	}

	cmd.Flags().StringP("limit", "l", "", "Limit search term")
	cmd.Flags().StringP("project", "p", "", "Limit project")
	cmd.Flags().String("sort-by", "name", "Sortby (id,name)")
	return &cmd
}

func (o *ActivitiesListCmd) Run(cmd *cobra.Command, args []string) {

	sort, _ := cmd.Flags().GetString("sort-by")
	configFile, _ := cmd.Flags().GetString("config")
	limit, _ := cmd.Flags().GetString("limit")
	project, _ := cmd.Flags().GetString("project")

	settings := config.NewSettings(configFile)

	activity := track.NewActitivies(settings)
	activities, err := activity.List(limit, project)

	if err != nil {
		slog.Error("Failed to list acitivies", "error", err)
	}

	if len(activities) == 0 {
		return
	}

	switch sort {
	case "name":
		slices.SortFunc(activities, func(a, b track.Activity) int {
			return strings.Compare(a.Name, b.Name)
		})
	case "id":
		slices.SortFunc(activities, func(a, b track.Activity) int {
			return a.ID - b.ID
		})
	case "billable":
		slices.SortFunc(activities, func(a, b track.Activity) int {
			switch {
			case a.Billable == b.Billable:
				return 0
			case a.Billable:
				return -1
			default:
				return 1
			}
		})
	}

	fmt.Printf("%-6v %-8v %-8v %-8v %v\n", "ID", "Project", "Billable", "Visible", "Name")
	for _, a := range activities {
		fmt.Printf("%-6v %-8v %-8v %-8v %v\n", a.ID, a.Project, a.Billable, a.Visible, a.Name)

	}
}
