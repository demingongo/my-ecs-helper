package infoapp

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/demingongo/my-ecs-helper/aws"
	"github.com/demingongo/my-ecs-helper/model/filepickermodel"
	"github.com/spf13/viper"
)

type TargetGroupConfig struct {
	New      bool   // must create a target group if it's true
	Filepath string // must be filled if New=true
	Arn      string // target group's arn
	Name     string // target group's name
}

type Config struct {
	targetGroup TargetGroupConfig
	rules       []string
	service     string

	targetGroupDescription string
	rulesDescription       string
	serviceDescription     string
}

type filepickerStyleStruct struct {
	cursor    lipgloss.Style
	directory lipgloss.Style
	file      lipgloss.Style
	selected  lipgloss.Style
	symlink   lipgloss.Style
}

const (
	formWidth = 60
	infoWidth = 38

	width = 100
)

var (

	// General.

	config Config
	info   string

	theme = huh.ThemeBase()

	subtle  = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	special = lipgloss.AdaptiveColor{Light: "230", Dark: "#010102"}

	subtleText = lipgloss.NewStyle().Foreground(subtle).Render

	// Titles.

	titleStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Background(lipgloss.Color("7")).
			Foreground(special)

	subtitleStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderTop(true).
			BorderForeground(subtle).
			Foreground(lipgloss.Color("6"))

	// Info block.

	infoStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("7")).
			BorderTop(true).
			BorderLeft(true).
			BorderRight(true).
			BorderBottom(true).
			Width(infoWidth)

	// filepicker
	filepickerStyle = filepickerStyleStruct{
		cursor:    lipgloss.NewStyle().Foreground(lipgloss.Color("6")),
		directory: lipgloss.NewStyle().Foreground(lipgloss.Color("7")).Underline(true),
		file:      lipgloss.NewStyle().Foreground(lipgloss.Color("15")),
		selected:  lipgloss.NewStyle().Foreground(lipgloss.Color("6")),
		symlink:   lipgloss.NewStyle().Foreground(lipgloss.Color("4")).Italic(true),
	}

	/*
		// Summary block.
		summaryStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#bd93f9")).
				BorderTop(true).
				BorderLeft(true).
				BorderRight(true).
				BorderBottom(true).
				Width(summaryWidth)

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
	*/
)

func selectJSONFile(title string, currentDirectory string, info string) string {
	m := filepickermodel.NewFilepickerModel(filepickermodel.FilepickerModelConfig{
		AllowedTypes:     []string{".json"},
		CurrentDirectory: currentDirectory,
		EnableFastSelect: true,
		Title:            title,
		InfoBubble:       info,
	}).
		ShowPermissions(false).
		ShowSize(false).
		Height(8).
		Width(width).
		FilepickerWidth(formWidth).
		StyleDirectory(filepickerStyle.directory).
		StyleFile(filepickerStyle.file).
		StyleSymlink(filepickerStyle.symlink).
		StyleCursor(filepickerStyle.cursor).
		StyleSelected(filepickerStyle.selected)

	tm, _ := tea.NewProgram(&m).Run()

	mm := tm.(filepickermodel.FilepickerModel)

	return mm.SelectedFile
}

func generateInfo() string {

	var (
		tgInfo      string
		rulesInfo   string
		serviceInfo string
	)

	if config.targetGroupDescription != "" {
		tgInfo = config.targetGroupDescription
	} else {
		if config.targetGroup.New {
			tgInfo = config.targetGroup.Filepath
		} else {
			tgInfo = config.targetGroup.Name
		}
	}

	if config.rulesDescription != "" {
		rulesInfo = config.rulesDescription
	} else {
		rulesInfo = strings.Join(config.rules, ", ")
	}

	if config.serviceDescription != "" {
		serviceInfo = config.serviceDescription
	} else {
		serviceInfo = config.service
	}

	if len(tgInfo) == 0 {
		tgInfo = subtleText("-")
	}
	if len(rulesInfo) == 0 {
		rulesInfo = subtleText("-")
	}
	if len(serviceInfo) == 0 {
		serviceInfo = subtleText("-")
	}

	content := lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render("SUMMARY"),
		subtitleStyle.Render("Target group"),
		tgInfo,
		subtitleStyle.Render("Rules"),
		rulesInfo,
		subtitleStyle.Render("Service"),
		serviceInfo,
	)

	return infoStyle.Render(content)
}

func createLogger() *log.Logger {
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
	logger.SetLevel(log.GetLevel())

	return logger
}

func generateTargetGroupDescription(name string, filepath string) string {

	r := name

	filepathMaxSize := 57

	if filepath != "" {
		if len(filepath) > filepathMaxSize {
			r += " (..." + filepath[len(filepath)-filepathMaxSize:] + ")"
		} else {
			r += " (" + filepath + ")"
		}
	}

	return r
}

func process(logger *log.Logger) {
	if config.targetGroup.New && config.targetGroup.Filepath != "" {
		logger.Debug(fmt.Sprintf("create target group \"%s\"", config.targetGroup.Name))
		_, err := aws.CreateTargetGroup(config.targetGroup.Filepath)
		if err != nil {
			logger.Fatal("CreateTargetGroup", err)
		}
	}

	if len(config.rules) > 0 {
		logger.Debug(fmt.Sprintf("create rules for target group \"%s\"", config.targetGroup.Name))
		for _, v := range config.rules {
			_, err := aws.CreateRule(v, config.targetGroup.Arn)
			if err != nil {
				logger.Fatal("CreateRule", err)
			}
		}
	}
}

func Run() {

	logger := createLogger()

	info = generateInfo()

	menuForm := runFormMenu()

	if menuForm.State == huh.StateCompleted && menuForm.GetString("operation") != "none" {

		operation := menuForm.GetString("operation")

		// create-targetgroup
		if operation == "create-targetgroup" {
			config.targetGroup.New = true
			config.targetGroup.Filepath = selectTargetGroupJSON(info)

			if config.targetGroup.Filepath != "" {
				tgConf := viper.New()
				tgConf.SetConfigFile(config.targetGroup.Filepath)
				err := tgConf.ReadInConfig()
				if err != nil {
					logger.Fatal("Could not read file:", err)
				}

				config.targetGroup.Name = tgConf.GetString("targetGroupName")

				config.targetGroupDescription = generateTargetGroupDescription(config.targetGroup.Name, config.targetGroup.Filepath)
			}
		}

		// select-targetgroup
		if operation == "select-targetgroup" {
			targetGroupForm := runFormTargetgroup()
			if targetGroupForm.State == huh.StateCompleted {
				tg := targetGroupForm.Get("targetgroup").(aws.TargetGroup)
				if tg.Arn != "" {
					config.targetGroup.Arn = tg.Arn
					config.targetGroup.Name = tg.Name
					config.targetGroupDescription = generateTargetGroupDescription(tg.Name, tg.Arn)
				}
			}
		}
		if config.targetGroupDescription == "" {
			config.targetGroupDescription = "❌"
		}
		info = generateInfo()

		// create rules
		if operation == "create-targetgroup" || operation == "select-targetgroup" {
			rulesForm := runFormRules()
			if rulesForm.State == huh.StateCompleted && rulesForm.GetBool("confirm") {
				var searchDir string
				var maxRules = 5
				for len(config.rules) < 5 {
					title := fmt.Sprintf("Pick a rule (.json) (%d/%d):", len(config.rules), maxRules)
					file := selectRuleJSON(info, title, searchDir)
					if len(file) > 0 {
						if slices.Contains(config.rules, file) {
							break
						} else {
							config.rules = append(config.rules, file)
							searchDir = filepath.Dir(file)
							info = generateInfo()
						}
					} else {
						break
					}
				}
			}
		}
		if len(config.rules) == 0 {
			config.rulesDescription = "❌"
			info = generateInfo()
		}

		// create service
		if operation == "create-targetgroup" || operation == "select-targetgroup" {
			serviceForm := runFormService()
			if serviceForm.State == huh.StateCompleted && serviceForm.GetBool("confirm") {
				config.service = selectServiceJSON(info)
			}
		} else if operation == "create-service" {
			config.service = selectServiceJSON(info)
		}
		if config.service == "" {
			config.serviceDescription = "❌"
		}
		info = generateInfo()

		fmt.Println(info)

		// @TODO processing
		process(logger)
	}

	fmt.Println("Done")
}
