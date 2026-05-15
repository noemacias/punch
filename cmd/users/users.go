package users

import (
	"github.com/spf13/cobra"
)

func NewUsersGroup() *cobra.Command {

	cmd := cobra.Command{
		Use:   `user`,
		Short: `Manage users`,
	}

	cmd.AddCommand(NewUserListCommand())
	return &cmd
}
