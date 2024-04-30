package infoapp

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/demingongo/my-ecs-helper/aws"
	formmmodel "github.com/demingongo/my-ecs-helper/model/formmodel"
)

func generateFormLoadBalancer(list []aws.ContainerPortMapping) *huh.Form {

	options := []huh.Option[aws.ContainerPortMapping]{
		huh.NewOption("(None)", aws.ContainerPortMapping{}),
	}

	for _, cpm := range list {
		var text string

		if cpm.Name != "" && cpm.PortMapping.Name != "" {
			if len(cpm.PortMapping.Name) > 0 {
				text += " (" + cpm.PortMapping.Name + ")"
			}
			options = append(options, huh.NewOption(cpm.Name+text, cpm))
		}
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[aws.ContainerPortMapping]().
				Title("Select a \"container (port)\" to load balance?:").
				Key("loadbalancer").
				Options(
					options...,
				).Height(6),
		),
	).
		WithTheme(theme).
		WithWidth(formWidth)

	return form
}

func runFormLoadBalancer(list []aws.ContainerPortMapping) *huh.Form {

	form := generateFormLoadBalancer(list)
	fModel := formmmodel.NewModel(formmmodel.ModelConfig{
		Form:       form,
		InfoBubble: info,
	}).Width(width)

	tea.NewProgram(&fModel).Run()

	return form
}
