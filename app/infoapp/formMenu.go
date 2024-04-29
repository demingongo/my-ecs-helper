package infoapp

import (
	"errors"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	formmmodel "github.com/demingongo/my-ecs-helper/model/formmodel"
)

func generateFormMenu() *huh.Form {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("What do we do:").
				Key("operation").
				Options(
					huh.NewOption("Create a target group", "create-targetgroup"),
					huh.NewOption("Select a target group", "select-targetgroup"),
					huh.NewOption("Create a service", "create-service"),
					huh.NewOption("Nothing", "none"),
				),

			huh.NewConfirm().
				Key("confirm").
				Title("Are you sure?").
				Validate(func(b bool) error {
					if !b {
						return errors.New("waiting till you confirm")
					}
					return nil
				}),
		),
	).
		WithTheme(theme).
		WithWidth(formWidth)

	return form
}

func runFormMenu() *huh.Form {

	firstSelectForm := generateFormMenu()
	firstSelect := formmmodel.NewModel(formmmodel.ModelConfig{
		Form:       firstSelectForm,
		InfoBubble: info,
	}).Width(width)

	tea.NewProgram(&firstSelect).Run()

	return firstSelectForm
}
