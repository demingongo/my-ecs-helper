package infoapp

import (
	"errors"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/demingongo/my-ecs-helper/aws"
	formmmodel "github.com/demingongo/my-ecs-helper/model/formmodel"
)

func selectTargetGroupJSON(info string) string {
	value := selectJSONFile("Pick a target group (.json):", "", info)
	return value
}

func generateFormTargetgroup() *huh.Form {
	options := []huh.Option[aws.TargetGroup]{
		huh.NewOption("(None)", aws.TargetGroup{}),
	}

	arnTextMaxSize := 12

	for _, tg := range aws.DescribeTargetGroups() {
		var arnText string

		if tg.Arn != "" {
			if len(tg.Arn) > arnTextMaxSize {
				arnText += " (..." + tg.Arn[len(tg.Arn)-arnTextMaxSize:] + ")"
			} else {
				arnText += " (" + tg.Arn + ")"
			}
		}
		options = append(options, huh.NewOption(tg.Name+arnText, tg))
	}

	confirm := true

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[aws.TargetGroup]().
				Title("Select a target group:").
				Key("targetgroup").
				Options(
					options...,
				).Height(6),

			huh.NewConfirm().
				Key("confirm").
				Title("Are you sure?").
				Value(&confirm).
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

func runFormTargetgroup() *huh.Form {

	firstSelectForm := generateFormTargetgroup()
	firstSelect := formmmodel.NewModel(formmmodel.ModelConfig{
		Form:       firstSelectForm,
		InfoBubble: info,
	}).Width(width)

	tea.NewProgram(&firstSelect).Run()

	return firstSelectForm
}
