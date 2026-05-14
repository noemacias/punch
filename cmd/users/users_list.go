package users

import (
	"fmt"
	"log/slog"

	"github.com/noemacias/punch/internal/config"
	"github.com/noemacias/punch/internal/track"
	"github.com/spf13/cobra"
)

type UsersListCommand struct {
}

func NewUserListCommand() *cobra.Command {

	o := UsersListCommand{}
	cmd := cobra.Command{
		Use: `list`,
		Run: func(cmd *cobra.Command, args []string) {

			o.Run(cmd, args)
		},
	}

	cmd.Flags().StringP("limit", "l", "", "Limit search term")

	return &cmd
}

func (u *UsersListCommand) Run(cmd *cobra.Command, args []string) {
	configFile, _ := cmd.Flags().GetString("config")
	limit, _ := cmd.Flags().GetString("limit")

	settings := config.NewSettings(configFile)

	users := track.NewUsers(settings)

	usersList, err := users.List(limit)

	if err != nil {
		slog.Error("Failed to get userlist", "error", err.Error())
		return
	}

	fmt.Printf("%-5v %-7v %-8v %-20v %v\n", "ID", "Enabled", "ApiToken", "Alias", "Username")

	for _, u := range usersList {
		fmt.Printf("%-5v %-7v %-8v %-20v %v\n", u.ID, u.Enabled, u.APIToken, u.Alias, u.Username)
	}
}
