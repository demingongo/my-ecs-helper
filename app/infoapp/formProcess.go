package infoapp

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	formmmodel "github.com/demingongo/my-ecs-helper/model/formmodel"
)

func generateFormProcess() *huh.Form {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Key("confirm").
				Title("").
				Negative("Cancel").
				Affirmative("Proceed").
				Inline(true),
		),
	).
		WithTheme(theme).
		WithWidth(formWidth)

	return form
}

func runFormProcess() *huh.Form {

	form := generateFormProcess()
	fModel := formmmodel.NewModel(formmmodel.ModelConfig{
		Form:         form,
		InfoBubble:   info,
		VerticalMode: true,
	}).Width(width)

	tea.NewProgram(&fModel).Run()

	return form
}
