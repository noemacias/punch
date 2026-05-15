package time

import "github.com/spf13/cobra"

func NewTimeGroup() *cobra.Command {

	cmd := cobra.Command{
		Use:   "time",
		Short: `Manage time entries and reports`,
	}

	cmd.AddCommand(NewTimeListCommand())
	cmd.AddCommand(NewTimeAddCommand())
	cmd.AddCommand(NewTimeReportCommand())
	return &cmd
}
