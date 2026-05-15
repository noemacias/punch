package cmd

import (
	"log/slog"

	"github.com/noemacias/punch/cmd/activities"
	"github.com/noemacias/punch/cmd/projects"
	"github.com/noemacias/punch/cmd/time"
	"github.com/noemacias/punch/cmd/users"
	"github.com/spf13/cobra"
)

func Execute() {

	cmd := cobra.Command{
		Use: `punch`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {

			debug, _ := cmd.Flags().GetBool("debug")

			if debug {
				slog.SetLogLoggerLevel(slog.LevelDebug)
			}
		},
	}

	cmd.PersistentFlags().String("config", "~/.config/punch.yaml", "Configration file")
	cmd.PersistentFlags().Bool("debug", false, "Debug")
	cmd.AddCommand(activities.NewActivityGroup())
	cmd.AddCommand(time.NewTimeGroup())
	cmd.AddCommand(users.NewUsersGroup())
	cmd.AddCommand(projects.NewProjectGroup())
	cmd.Execute()
}
