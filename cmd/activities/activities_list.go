package activities

import (
	"fmt"
	"log/slog"

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
	return &cmd
}

func (o *ActivitiesListCmd) Run(cmd *cobra.Command, args []string) {

	configFile, _ := cmd.Flags().GetString("config")
	limit, _ := cmd.Flags().GetString("limit")

	settings := config.NewSettings(configFile)

	activity := track.NewActitivies(settings)
	activities, err := activity.List(limit)

	if err != nil {
		slog.Error("Failed to list acitivies", "error", err)
	}

	if len(activities) == 0 {
		return
	}

	fmt.Printf("%-6v %-8v %-8v %-8v %v\n", "ID", "Project", "Billable", "Visible", "Name")
	for _, a := range activities {
		fmt.Printf("%-6v %-8v %-8v %-8v %v\n", a.ID, a.Project, a.Billable, a.Visible, a.Name)

	}
}
