package cmd

import (
	"github.com/noemacias/punch/cmd/activities"
	"github.com/noemacias/punch/cmd/time"
	"github.com/noemacias/punch/cmd/users"
	"github.com/spf13/cobra"
)

func Execute() {

	cmd := cobra.Command{
		Use: `punch`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
		},
	}

	cmd.PersistentFlags().String("config", "~/.config/punch.yaml", "Configration file")
	cmd.AddCommand(activities.NewActivityGroup())
	cmd.AddCommand(time.NewTimeGroup())
	cmd.AddCommand(users.NewUsersGroup())
	cmd.Execute()
}
