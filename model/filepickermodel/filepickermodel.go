package filepickermodel

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
)

type FilepickerModelConfig struct {
	Title            string
	CurrentDirectory string
	AllowedTypes     []string
	EnableFastSelect bool
}

type FilepickerModel struct {
	filepicker       filepicker.Model
	quitting         bool
	err              error
	enableFastSelect bool
	title            string
	SelectedFile     string
}

type clearErrorMsg struct{}

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

func (m FilepickerModel) Init() tea.Cmd {
	return m.filepicker.Init()
}

func (m FilepickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		m.SelectedFile = path
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
		m.SelectedFile = ""
		return m, tea.Batch(cmd, clearErrorAfter(2*time.Second))
	}

	return m, cmd
}

func (m FilepickerModel) View() string {
	if m.quitting {
		return ""
	}
	var s strings.Builder
	s.WriteString("\n  ")
	if m.err != nil {
		s.WriteString(m.filepicker.Styles.DisabledFile.Render(m.err.Error()))
	} else if m.SelectedFile == "" {
		title := "Pick a file:"
		if len(m.title) > 0 {
			title = m.title
		}
		s.WriteString(title)
	} else {
		s.WriteString("Selected file: " + m.filepicker.Styles.Selected.Render(m.SelectedFile))
	}
	s.WriteString("\n\n" + m.filepicker.View() + "\n")
	return s.String()
}

func NewFilepickerModel(config FilepickerModelConfig) FilepickerModel {
	var dir string
	if len(config.CurrentDirectory) > 0 {
		dir = config.CurrentDirectory
	} else {
		dir, _ = os.Getwd()
	}

	fp := filepicker.New()
	fp.AllowedTypes = config.AllowedTypes
	fp.CurrentDirectory = dir

	m := FilepickerModel{
		filepicker:       fp,
		enableFastSelect: config.EnableFastSelect,
		title:            config.Title,
	}

	return m
}

func NewFilepickerModelProgram(config FilepickerModelConfig) *tea.Program {
	m := NewFilepickerModel(config)

	return tea.NewProgram(&m)
}
