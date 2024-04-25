package main

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

type model struct {
	filepicker       filepicker.Model
	selectedFile     string
	quitting         bool
	err              error
	enableFastSelect bool
	title            string
}

type clearErrorMsg struct{}

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

func (m model) Init() tea.Cmd {
	return m.filepicker.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		}
	case clearErrorMsg:
		m.err = nil
	}

	var cmd tea.Cmd
	m.filepicker, cmd = m.filepicker.Update(msg)

	// Did the user select a file?
	if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
		// Get the path of the selected file.
		m.selectedFile = path
		if m.enableFastSelect {
			m.quitting = true
			return m, tea.Quit
		}
	}

	// Did the user select a disabled file?
	// This is only necessary to display an error to the user.
	if didSelect, path := m.filepicker.DidSelectDisabledFile(msg); didSelect {
		// Let's clear the selectedFile and display an error.
		m.err = errors.New(path + " is not valid.")
		m.selectedFile = ""
		return m, tea.Batch(cmd, clearErrorAfter(2*time.Second))
	}

	return m, cmd
}

func (m model) View() string {
	if m.quitting {
		return ""
	}
	var s strings.Builder
	s.WriteString("\n  ")
	if m.err != nil {
		s.WriteString(m.filepicker.Styles.DisabledFile.Render(m.err.Error()))
	} else if m.selectedFile == "" {
		title := "Pick a file:"
		if len(m.title) > 0 {
			title = m.title
		}
		s.WriteString(title)
	} else {
		s.WriteString("Selected file: " + m.filepicker.Styles.Selected.Render(m.selectedFile))
	}
	s.WriteString("\n\n" + m.filepicker.View() + "\n")
	return s.String()
}

func selectRuleFile(title string) string {
	fp := filepicker.New()
	fp.AllowedTypes = []string{".json"}
	fp.CurrentDirectory, _ = os.Getwd()

	m := model{
		filepicker: fp,
	}
	m.enableFastSelect = true
	m.title = title
	tm, _ := tea.NewProgram(&m).Run()
	mm := tm.(model)

	return mm.selectedFile
}

func targetGroupSelectForm(value *string, confirm *bool) *huh.Form {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select a target group:").
				Options(
					huh.NewOption("dev-candidates", "arn:targetgroup/dev-candidates"),
					huh.NewOption("dev-websites", "arn:targetgroup/dev-websites"),
					huh.NewOption("prod-candidates", "arn:targetgroup/prod-candidates"),
					huh.NewOption("prod-websites", "arn:targetgroup/prod-websites"),
				).
				Value(value),

			huh.NewConfirm().
				Title("Are you sure?").
				Value(confirm),
		),
	)

	return form
}

func targetGroupActionForm(value *int) *huh.Form {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title("Create/select a target group?").
				Options(
					huh.NewOption("Create", 1),
					huh.NewOption("Select", 2),
					huh.NewOption("None", 0),
				).
				Value(value),
		),
	)

	return form
}

func confirmForm(title string, value *bool) *huh.Form {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(title).
				Value(value),
		),
	)

	return form
}

func createTargetGroup(logger *log.Logger) {
	logger.Info("createTargetGroup")
}

func selectTargetGroup(logger *log.Logger) string {
	var targetGroupArn string
	var ok bool = true
	err := targetGroupSelectForm(&targetGroupArn, &ok).Run()
	if err != nil {
		logger.Fatal(err)
	}

	if !ok || len(targetGroupArn) == 0 {
		logger.Fatal("Bye")
	}

	return targetGroupArn
}

func askTargetGroup(logger *log.Logger) string {
	var tgAction int = 1
	var targetGroupArn string

	form := targetGroupActionForm(&tgAction)

	err := form.Run()
	if err != nil {
		logger.Fatal(err)
	}

	if tgAction == 1 {
		createTargetGroup(logger)
	}

	if tgAction != 0 {
		targetGroupArn = selectTargetGroup(logger)
	}

	return targetGroupArn
}

func askRules(logger *log.Logger, targetGroup string) {

	var ok bool
	var files []string
	err := confirmForm("Create rules?", &ok).Run()
	if err != nil {
		logger.Fatal(err)
	}

	if ok {
		for i := 0; i < 10; i++ {
			file := selectRuleFile(fmt.Sprintf("Pick a file (%d):", len(files)))
			if len(file) > 0 {
				if slices.Contains(files, file) {
					break
				} else {
					files = append(files, file)
				}
			} else {
				break
			}
		}
	}

	fmt.Println(files)

}

func main() {
	// Override the default info level style.
	styles := log.DefaultStyles()
	styles.Levels[log.InfoLevel] = lipgloss.NewStyle().
		SetString("INFO").
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("50")).
		Foreground(lipgloss.Color("0"))
	// Add a custom style for key `err`
	styles.Keys["err"] = lipgloss.NewStyle().Foreground(lipgloss.Color("50"))
	styles.Values["err"] = lipgloss.NewStyle().Bold(true)
	logger := log.New(os.Stderr)
	logger.SetStyles(styles)

	targetGroup := askTargetGroup(logger)

	logger.Info("Target group: " + targetGroup)

	if len(targetGroup) > 0 {
		askRules(logger, targetGroup)
	}

	fmt.Println("Done!")
}
