package projects

import "github.com/spf13/cobra"

func NewProjectGroup() *cobra.Command {
	cmd := cobra.Command{
		Use:   `project`,
		Short: `Manage projects`,
	}

	cmd.AddCommand(NewProjectListCommand())

	return &cmd
}
