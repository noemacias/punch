package cmd

import (
	"github.com/noemacias/punch/cmd/activities"
	"github.com/noemacias/punch/cmd/time"
	"github.com/spf13/cobra"
)

func Execute() {

	cmd := cobra.Command{
		Use: `punch`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
		},
	}

	cmd.PersistentFlags().String("config", "config.yaml", "Configration file")
	cmd.AddCommand(activities.NewActiviGroup())
	cmd.AddCommand(time.NewTimeGroup())
	cmd.Execute()
}
