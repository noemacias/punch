package projects

import (
	"fmt"
	"log/slog"

	"github.com/noemacias/punch/internal/config"
	"github.com/noemacias/punch/internal/track"
	"github.com/spf13/cobra"
)

type ProjectListCommand struct{}

func NewProjectListCommand() *cobra.Command {

	o := ProjectListCommand{}

	cmd := cobra.Command{
		Use: `list`,
		Run: func(cmd *cobra.Command, args []string) {
			o.Run(cmd, args)
		},
	}
	return &cmd
}

func (o *ProjectListCommand) Run(cmd *cobra.Command, args []string) {
	configFile, _ := cmd.Flags().GetString("config")

	settings := config.NewSettings(configFile)

	project := track.NewProject(settings)

	projects, err := project.List()

	if err != nil {
		slog.Error("Failed to get project list", "error", err.Error())
		return
	}

	fmt.Printf("%-4v %-20v %-25v %v\n", "ID", "Name", "Parent Name", "Team(s)")
	for _, p := range projects {
		fmt.Printf("%-4v %-20v %-25v\n", p.ID, p.Name, p.ParentTitle)

		for _, t := range p.Teams {
			fmt.Printf("%-4v %-20v %-25v %v\n", "", "", "", fmt.Sprintf("%-4v %v", t.ID, t.Name))

		}
	}
}
