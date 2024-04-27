package infoapp

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

const (
	formWidth = 60
)

type Model struct {
	form     *huh.Form // huh.Form is just a tea.Model
	quitting bool
}

func NewModel() Model {
	return Model{
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Key("class").
					Options(huh.NewOptions("Warrior", "Mage", "Rogue")...).
					Title("Choose your class"),

				huh.NewSelect[int]().
					Key("level").
					Options(huh.NewOptions(1, 20, 9999)...).
					Title("Choose your level"),
			),
		).WithTheme(huh.ThemeDracula()),
	}
}

func (m Model) Init() tea.Cmd {
	m.form.WithWidth(formWidth)
	return m.form.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// user did what?
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			m.form.State = huh.StateAborted
			return m, tea.Quit
		}
	}

	// update the form
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}

	// is it completed?
	if form.(*huh.Form).State == huh.StateCompleted {
		m.quitting = true
		return m, tea.Quit
	}

	return m, cmd
}

// This returns a string !!!!!!!!!!!!!!!!!! EUREKA
func (m Model) View() string {
	if m.quitting {
		return ""
	}
	formView := m.form.View()

	// Render a block of text.
	var style = lipgloss.NewStyle()
	//Width(formWidth).
	//Padding(2).
	//Background(lipgloss.Color("0"))
	var block string = style.Render(formView)

	//width, height := lipgloss.Size(block)

	t := table.New().
		Border(lipgloss.HiddenBorder()).
		Row(block, "Details here")

	return t.Render()
}

func Run() {
	m := NewModel()
	tm, _ := tea.NewProgram(&m).Run()
	mm := tm.(Model)

	if mm.form.State == huh.StateCompleted {
		fmt.Printf("So yea basically, you selected: %s, Lvl. %d\n", mm.form.GetString("class"), mm.form.GetInt("level"))
	}

	fmt.Println("Done")
}
