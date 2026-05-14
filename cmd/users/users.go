package users

import (
	"github.com/spf13/cobra"
)

func NewUsersGroup() *cobra.Command {

	cmd := cobra.Command{
		Use: `users`,
	}

	cmd.AddCommand(NewUserListCommand())
	return &cmd
}
