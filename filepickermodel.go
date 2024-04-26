package main

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
)

type filepickerModelConfig struct {
	title            string
	currentDirectory string
	allowedTypes     []string
	enableFastSelect bool
}

type filepickerModel struct {
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

func (m filepickerModel) Init() tea.Cmd {
	return m.filepicker.Init()
}

func (m filepickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m filepickerModel) View() string {
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

func NewFilepickerModel(config filepickerModelConfig) filepickerModel {
	var dir string
	if len(config.currentDirectory) > 0 {
		dir = config.currentDirectory
	} else {
		dir, _ = os.Getwd()
	}

	fp := filepicker.New()
	fp.AllowedTypes = config.allowedTypes
	fp.CurrentDirectory = dir

	m := filepickerModel{
		filepicker:       fp,
		enableFastSelect: config.enableFastSelect,
		title:            config.title,
	}

	return m
}

func NewFilepickerModelProgram(config filepickerModelConfig) *tea.Program {
	m := NewFilepickerModel(config)

	return tea.NewProgram(&m)
}
