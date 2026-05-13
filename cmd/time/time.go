package time

import "github.com/spf13/cobra"

func NewTimeGroup() *cobra.Command {

	cmd := cobra.Command{
		Use: "time",
	}

	cmd.AddCommand(NewTimeListCommand())
	cmd.AddCommand(NewTimeAddCommand())
	return &cmd
}
