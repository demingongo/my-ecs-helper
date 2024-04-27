package mainapp

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	fpm "github.com/demingongo/my-ecs-helper/model/filepickermodel"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

func readJSONFile[T comparable](filepath string) (map[T]interface{}, error) {
	var res map[T]interface{}
	content, err := os.ReadFile(filepath)
	if err != nil {
		return res, err
	}

	json.Unmarshal([]byte(content), &res)

	return res, err
}

func selectJSONFile(title string, currentDirectory string) string {
	tm, _ := fpm.NewFilepickerModelProgram(fpm.FilepickerModelConfig{
		AllowedTypes:     []string{".json"},
		CurrentDirectory: currentDirectory,
		EnableFastSelect: true,
		Title:            title,
	}).Run()
	mm := tm.(fpm.FilepickerModel)

	return mm.SelectedFile
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
	).WithTheme(huh.ThemeDracula())

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
	).WithTheme(huh.ThemeDracula())

	return form
}

func confirmForm(title string, value *bool) *huh.Form {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(title).
				Value(value),
		),
	).WithTheme(huh.ThemeDracula())

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

func askRules(logger *log.Logger) []string {
	var ok bool
	var files []string
	err := confirmForm("Create rules?", &ok).Run()
	if err != nil {
		logger.Fatal(err)
	}

	if ok {
		var searchDir string
		for i := 0; i < 10; i++ {
			file := selectJSONFile(fmt.Sprintf("Pick a file (%d/10):", len(files)), searchDir)
			if len(file) > 0 {
				if slices.Contains(files, file) {
					break
				} else {
					files = append(files, file)
					searchDir = filepath.Dir(file)
				}
			} else {
				break
			}
		}
	}

	return files
}

func askService(logger *log.Logger) string {
	var ok bool = true
	var value string
	err := confirmForm("Create a service?", &ok).Run()
	if err != nil {
		logger.Fatal(err)
	}

	if ok {
		value = selectJSONFile("", "")
	}

	return value
}

func Run() {
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

	var (
		targetGroup string
		rulesFiles  []string
		serviceFile string
	)

	targetGroup = askTargetGroup(logger)

	if len(targetGroup) > 0 {
		rulesFiles = askRules(logger)
	}

	serviceFile = askService(logger)

	logger.Info(fmt.Sprintf("Target group: %s", targetGroup))
	logger.Info(fmt.Sprintf("Rules: %s", rulesFiles))
	logger.Info(fmt.Sprintf("Service: %s", serviceFile))

	if len(serviceFile) > 0 {
		content, err := readJSONFile[string](serviceFile)
		if err != nil {
			logger.Fatal("Could not load file:", err)
		}

		logger.Info(content["listener"])

	}

	fmt.Println("Done!")
}
