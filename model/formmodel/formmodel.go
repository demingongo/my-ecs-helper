package formmmodel

import (
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

type ModelConfig struct {
	Title      string
	Key        string
	InfoBubble string
	Form       *huh.Form
}

var (
	// Status Bar.

	statusNugget = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Padding(0, 1)

	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
			Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"})

	statusStyle = lipgloss.NewStyle().
			Inherit(statusBarStyle).
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#FF5F87")).
			Padding(0, 1).
			MarginRight(1)

	encodingStyle = statusNugget.Copy().
			Background(lipgloss.Color("#A550DF")).
			Align(lipgloss.Right)

	statusText = lipgloss.NewStyle().Inherit(statusBarStyle)

	fishCakeStyle = statusNugget.Copy().Background(lipgloss.Color("#6124DF"))

	// Page.

	docStyle = lipgloss.NewStyle().Padding(1, 2, 1, 2)
)

type Model struct {
	form       *huh.Form // huh.Form is just a tea.Model
	quitting   bool
	infoBubble string
	key        string
	width      int
	State      huh.FormState
}

func NewModel(config ModelConfig) Model {
	return Model{
		form: config.Form,

		State: huh.StateNormal,

		infoBubble: config.InfoBubble,

		// @TODO: will be used in a function
		// to know what was updated
		// for updating the infoBubble for example (when it will be a class and not just a string)
		key: config.Key,
	}
}

func (m Model) Width(width int) Model {
	m.width = width
	return m
}

func (m Model) Init() tea.Cmd {
	return m.form.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// user did what?
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.quitting = true
			m.State = huh.StateAborted
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
		m.State = huh.StateCompleted
		return m, tea.Quit
	}
	// is it aborted?
	if form.(*huh.Form).State == huh.StateAborted {
		m.quitting = true
		m.State = huh.StateAborted
		return m, tea.Quit
	}

	return m, cmd
}

func (m Model) View() string {
	if m.quitting {
		return ""
	}

	physicalWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))

	doc := strings.Builder{}

	// Form
	{
		var formView string
		var infoView string

		// Form / Form
		{
			formView = m.form.View()
		}

		// Form / Info
		{
			if m.infoBubble != "" && (m.width == 0 || (m.width > 0 && physicalWidth >= m.width*4/5)) {
				infoView = m.infoBubble
			}
		}

		doc.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, formView, infoView))
		doc.WriteString("\n\n")
	}

	// Status bar
	{
		w := lipgloss.Width

		statusKey := statusStyle.Render("STATUS")
		encoding := encodingStyle.Render("INFO")
		fishCake := fishCakeStyle.Render("💾 My ECS helper")
		statusVal := statusText.Copy().
			Width(m.width - w(statusKey) - w(encoding) - w(fishCake)).
			Render("Normal")

		bar := lipgloss.JoinHorizontal(lipgloss.Top,
			statusKey,
			statusVal,
			encoding,
			fishCake,
		)

		doc.WriteString(statusBarStyle.Width(m.width).Render(bar))
	}

	if physicalWidth > 0 {
		docStyle = docStyle.MaxWidth(physicalWidth)
	}

	return docStyle.Render(doc.String())
}
