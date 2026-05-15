package activities

import "github.com/spf13/cobra"

func NewActivityGroup() *cobra.Command {

	cmd := cobra.Command{
		Use:   `activity`,
		Short: `Manage activities`,
	}

	cmd.AddCommand(NewActitiviesList())
	cmd.AddCommand(NewActitiviesAdd())
	return &cmd
}
