package activities

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/noemacias/punch/internal/config"
	"github.com/noemacias/punch/internal/track"
	"github.com/noemacias/punch/internal/utils"
	"github.com/spf13/cobra"
)

type ActivitiesAddCmd struct {
}

func NewActitiviesAdd() *cobra.Command {

	o := ActivitiesAddCmd{}

	cmd := cobra.Command{
		Use: `add`,
		Run: func(cmd *cobra.Command, args []string) {
			o.Run(cmd, args)
		},
	}

	return &cmd
}

func (o *ActivitiesAddCmd) Run(cmd *cobra.Command, args []string) {
	configFile, _ := cmd.Flags().GetString("config")

	settings := config.NewSettings(configFile)

	activities := track.NewActitivies(settings)

	name, err := utils.ReadInput("Name: ")

	if err != nil {
		return
	}

	if name == "" {
		return
	}

	comment, err := utils.ReadInput("Comment: ")

	if err != nil {
		return
	}

	if comment == "" {
		return
	}

	project, _ := strconv.Atoi(settings.CustomerId)

	activity := track.Activity{
		Name:     name,
		Project:  project,
		Comment:  comment,
		Billable: true,
		Visible:  true,
	}

	err = activities.Add(&activity)

	if err != nil {
		slog.Error(err.Error())
		return
	}

	fmt.Println("Activity successfully added.")
}
