package infoapp

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/demingongo/my-ecs-helper/model/filepickermodel"
	formmmodel "github.com/demingongo/my-ecs-helper/model/formmodel"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

/*
const (
	// In real life situations we'd adjust the document to fit the width we've
	// detected. In the case of this example we're hardcoding the width, and
	// later using the detected width only to truncate in order to avoid jaggy
	// wrapping.
	width = 96

	columnWidth = 30
)

// Style definitions.
var (

	// General.

	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}

	divider = lipgloss.NewStyle().
		SetString("•").
		Padding(0, 1).
		Foreground(subtle).
		String()

	url = lipgloss.NewStyle().Foreground(special).Render

	// Tabs.

	activeTabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      " ",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┘",
		BottomRight: "└",
	}

	tabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┴",
		BottomRight: "┴",
	}

	tab = lipgloss.NewStyle().
		Border(tabBorder, true).
		BorderForeground(highlight).
		Padding(0, 1)

	activeTab = tab.Copy().Border(activeTabBorder, true)

	tabGap = tab.Copy().
		BorderTop(false).
		BorderLeft(false).
		BorderRight(false)

	// Title.

	titleStyle = lipgloss.NewStyle().
			MarginLeft(1).
			MarginRight(5).
			Padding(0, 1).
			Italic(true).
			Foreground(lipgloss.Color("#FFF7DB")).
			SetString("Lip Gloss")

	descStyle = lipgloss.NewStyle().MarginTop(1)

	infoStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderTop(true).
			BorderForeground(subtle)

	// Dialog.

	dialogBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 0).
			BorderTop(true).
			BorderLeft(true).
			BorderRight(true).
			BorderBottom(true)

	buttonStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFF7DB")).
			Background(lipgloss.Color("#888B7E")).
			Padding(0, 3).
			MarginTop(1)

	activeButtonStyle = buttonStyle.Copy().
				Foreground(lipgloss.Color("#FFF7DB")).
				Background(lipgloss.Color("#F25D94")).
				MarginRight(2).
				Underline(true)

	// List.

	list = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, true, false, false).
		BorderForeground(subtle).
		MarginRight(2).
		Height(8).
		Width(columnWidth + 1)

	listHeader = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(subtle).
			MarginRight(2).
			Render

	listItem = lipgloss.NewStyle().PaddingLeft(2).Render

	checkMark = lipgloss.NewStyle().SetString("✓").
			Foreground(special).
			PaddingRight(1).
			String()

	listDone = func(s string) string {
		return checkMark + lipgloss.NewStyle().
			Strikethrough(true).
			Foreground(lipgloss.AdaptiveColor{Light: "#969B86", Dark: "#696969"}).
			Render(s)
	}

	// Paragraphs/History.

	historyStyle = lipgloss.NewStyle().
			Align(lipgloss.Left).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(highlight).
			Margin(1, 3, 0, 0).
			Padding(1, 2).
			Height(19).
			Width(columnWidth)

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

func Run() {
	physicalWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))
	doc := strings.Builder{}

	// Tabs
	{
		row := lipgloss.JoinHorizontal(
			lipgloss.Top,
			activeTab.Render("Lip Gloss"),
			tab.Render("Blush"),
			tab.Render("Eye Shadow"),
			tab.Render("Mascara"),
			tab.Render("Foundation"),
		)
		gap := tabGap.Render(strings.Repeat(" ", max(0, width-lipgloss.Width(row)-2)))
		row = lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)
		doc.WriteString(row + "\n\n")
	}

	// Title
	{
		var (
			colors = colorGrid(1, 5)
			title  strings.Builder
		)

		for i, v := range colors {
			const offset = 2
			c := lipgloss.Color(v[0])
			fmt.Fprint(&title, titleStyle.Copy().MarginLeft(i*offset).Background(c))
			if i < len(colors)-1 {
				title.WriteRune('\n')
			}
		}

		desc := lipgloss.JoinVertical(lipgloss.Left,
			descStyle.Render("Style Definitions for Nice Terminal Layouts"),
			infoStyle.Render("From Charm"+divider+url("https://github.com/charmbracelet/lipgloss")),
		)

		row := lipgloss.JoinHorizontal(lipgloss.Top, title.String(), desc)
		doc.WriteString(row + "\n\n")
	}

	// Dialog
	{
		okButton := activeButtonStyle.Render("Yes")
		cancelButton := buttonStyle.Render("Maybe")

		question := lipgloss.NewStyle().Width(50).Align(lipgloss.Center).Render("Are you sure you want to eat marmalade?")
		buttons := lipgloss.JoinHorizontal(lipgloss.Top, okButton, cancelButton)
		ui := lipgloss.JoinVertical(lipgloss.Center, question, buttons)

		dialog := lipgloss.Place(width, 9,
			lipgloss.Center, lipgloss.Center,
			dialogBoxStyle.Render(ui),
			lipgloss.WithWhitespaceChars("猫咪"),
			lipgloss.WithWhitespaceForeground(subtle),
		)

		doc.WriteString(dialog + "\n\n")
	}

	// Color grid
	colors := func() string {
		colors := colorGrid(14, 8)

		b := strings.Builder{}
		for _, x := range colors {
			for _, y := range x {
				s := lipgloss.NewStyle().SetString("  ").Background(lipgloss.Color(y))
				b.WriteString(s.String())
			}
			b.WriteRune('\n')
		}

		return b.String()
	}()

	lists := lipgloss.JoinHorizontal(lipgloss.Top,
		list.Render(
			lipgloss.JoinVertical(lipgloss.Left,
				listHeader("Citrus Fruits to Try"),
				listDone("Grapefruit"),
				listDone("Yuzu"),
				listItem("Citron"),
				listItem("Kumquat"),
				listItem("Pomelo"),
			),
		),
		list.Copy().Width(columnWidth).Render(
			lipgloss.JoinVertical(lipgloss.Left,
				listHeader("Actual Lip Gloss Vendors"),
				listItem("Glossier"),
				listItem("Claire‘s Boutique"),
				listDone("Nyx"),
				listItem("Mac"),
				listDone("Milk"),
			),
		),
	)

	doc.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, lists, colors))

	// Marmalade history
	{
		const (
			historyA = "The Romans learned from the Greeks that quinces slowly cooked with honey would “set” when cool. The Apicius gives a recipe for preserving whole quinces, stems and leaves attached, in a bath of honey diluted with defrutum: Roman marmalade. Preserves of quince and lemon appear (along with rose, apple, plum and pear) in the Book of ceremonies of the Byzantine Emperor Constantine VII Porphyrogennetos."
			historyB = "Medieval quince preserves, which went by the French name cotignac, produced in a clear version and a fruit pulp version, began to lose their medieval seasoning of spices in the 16th century. In the 17th century, La Varenne provided recipes for both thick and clear cotignac."
			historyC = "In 1524, Henry VIII, King of England, received a “box of marmalade” from Mr. Hull of Exeter. This was probably marmelada, a solid quince paste from Portugal, still made and sold in southern Europe today. It became a favourite treat of Anne Boleyn and her ladies in waiting."
		)

		doc.WriteString(lipgloss.JoinHorizontal(
			lipgloss.Top,
			historyStyle.Copy().Align(lipgloss.Right).Render(historyA),
			historyStyle.Copy().Align(lipgloss.Center).Render(historyB),
			historyStyle.Copy().MarginRight(0).Render(historyC),
		))

		doc.WriteString("\n\n")
	}

	// Status bar
	{
		w := lipgloss.Width

		statusKey := statusStyle.Render("STATUS")
		encoding := encodingStyle.Render("UTF-8")
		fishCake := fishCakeStyle.Render("🍥 Fish Cake")
		statusVal := statusText.Copy().
			Width(width - w(statusKey) - w(encoding) - w(fishCake)).
			Render("Ravishing")

		bar := lipgloss.JoinHorizontal(lipgloss.Top,
			statusKey,
			statusVal,
			encoding,
			fishCake,
		)

		doc.WriteString(statusBarStyle.Width(width).Render(bar))
	}

	if physicalWidth > 0 {
		docStyle = docStyle.MaxWidth(physicalWidth)
	}

	// Okay, let's print it
	fmt.Println(docStyle.Render(doc.String()))
}

func colorGrid(xSteps, ySteps int) [][]string {
	x0y0, _ := colorful.Hex("#F25D94")
	x1y0, _ := colorful.Hex("#EDFF82")
	x0y1, _ := colorful.Hex("#643AFF")
	x1y1, _ := colorful.Hex("#14F9D5")

	x0 := make([]colorful.Color, ySteps)
	for i := range x0 {
		x0[i] = x0y0.BlendLuv(x0y1, float64(i)/float64(ySteps))
	}

	x1 := make([]colorful.Color, ySteps)
	for i := range x1 {
		x1[i] = x1y0.BlendLuv(x1y1, float64(i)/float64(ySteps))
	}

	grid := make([][]string, ySteps)
	for x := 0; x < ySteps; x++ {
		y0 := x0[x]
		grid[x] = make([]string, xSteps)
		for y := 0; y < xSteps; y++ {
			grid[x][y] = y0.BlendLuv(x1[x], float64(y)/float64(xSteps)).Hex()
		}
	}

	return grid
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
*/

const (
	formWidth    = 60
	summaryWidth = 38

	width = 100
)

var (

	// General.

	subtle  = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	special = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}

	subtleText = lipgloss.NewStyle().Foreground(subtle).Render

	// Titles.

	subtitleStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderTop(true).
			BorderForeground(subtle).
			Foreground(special)

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
		).WithTheme(huh.ThemeBase()).
			WithWidth(formWidth),
	}
}

func (m Model) Init() tea.Cmd {
	return m.form.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// user did what?
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
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

	physicalWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))

	doc := strings.Builder{}

	// Form
	{
		var formView string
		var summaryView string

		// Form / Form
		{
			var style = lipgloss.NewStyle()
			formView = style.Render(m.form.View())
		}

		// Form / Summary
		{
			var style = lipgloss.NewStyle().
				Padding(0, 1).
				Background(lipgloss.Color("#bd93f9")).
				Foreground(lipgloss.Color("#d1cbcb"))

			var (
				targetGroup string
				rules       []string
				service     string
			)

			if len(m.form.GetString("class")) > 0 {
				targetGroup = m.form.GetString("class")
			} else {
				targetGroup = subtleText("-")
			}

			summary := lipgloss.JoinVertical(lipgloss.Left,
				style.Render("SUMMARY"),
				subtitleStyle.Render("Target group"),
				targetGroup,
				subtitleStyle.Render("Rules"),
				subtleText(strings.Join(rules, ", ")),
				subtitleStyle.Render("Service"),
				subtleText(service),
			)

			summaryView = summaryStyle.Render(summary)
		}

		doc.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, formView, summaryView))
		doc.WriteString("\n\n")
	}

	// Status bar
	{
		w := lipgloss.Width

		statusKey := statusStyle.Render("STATUS")
		encoding := encodingStyle.Render("INFO")
		fishCake := fishCakeStyle.Render("💾 My ECS helper")
		statusVal := statusText.Copy().
			Width(width - w(statusKey) - w(encoding) - w(fishCake)).
			Render("Normal")

		bar := lipgloss.JoinHorizontal(lipgloss.Top,
			statusKey,
			statusVal,
			encoding,
			fishCake,
		)

		doc.WriteString(statusBarStyle.Width(width).Render(bar))
	}

	if physicalWidth > 0 {
		docStyle = docStyle.MaxWidth(physicalWidth)
	}

	return docStyle.Render(doc.String())
}

func selectJSONFile(title string, currentDirectory string, info string) string {
	tm, _ := filepickermodel.NewFilepickerModelProgram(filepickermodel.FilepickerModelConfig{
		AllowedTypes:     []string{".json"},
		CurrentDirectory: currentDirectory,
		EnableFastSelect: true,
		Title:            title,
		InfoBubble:       info,
	}).Run()
	mm := tm.(filepickermodel.FilepickerModel)

	return mm.SelectedFile
}

func selectTargetGroupJSON(info string) string {
	value := selectJSONFile("Pick a target group (.json):", "", info)
	return value
}

func generateInfo(targetGroup string, rules []string, service string) string {
	var style = lipgloss.NewStyle().
		Padding(0, 1).
		Background(lipgloss.Color("#bd93f9")).
		Foreground(lipgloss.Color("#d1cbcb"))
	if len(targetGroup) == 0 {
		targetGroup = subtleText("-❌")
	}

	content := lipgloss.JoinVertical(lipgloss.Left,
		style.Render("SUMMARY"),
		subtitleStyle.Render("Target group"),
		targetGroup,
		subtitleStyle.Render("Rules"),
		subtleText(strings.Join(rules, ", ")),
		subtitleStyle.Render("Service"),
		subtleText(service),
	)

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#bd93f9")).
		BorderTop(true).
		BorderLeft(true).
		BorderRight(true).
		BorderBottom(true).
		Width(38).Render(content)
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

	return logger
}

func generateTargetGroupDescription(name string, filepath string) string {

	r := name

	filepathMaxSize := 60

	if filepath != "" {
		if len(filepath) > filepathMaxSize {
			r += " (..." + filepath[len(filepath)-filepathMaxSize:] + ")"
		} else {
			r += " (" + filepath + ")"
		}
	}

	return r
}

func generateFirstSelectForm() *huh.Form {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("What do we do:").
				Key("targetgroup").
				Options(
					huh.NewOption("Create a target group", "create"),
					huh.NewOption("Select a target group", "select"),
					huh.NewOption("Other", "other"),
				),

			huh.NewConfirm().
				Key("confirm").
				Title("Are you sure?"),
		),
	).WithTheme(huh.ThemeDracula())

	return form
}

/*
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
*/

func Run() {

	var (
		targetGroup            string
		targetGroupDescription string
		rules                  []string
		service                string
	)

	logger := createLogger()

	info := generateInfo(targetGroupDescription, rules, service)

	firstSelectForm := generateFirstSelectForm()
	firstSelect := formmmodel.NewModel(formmmodel.ModelConfig{
		Form:       firstSelectForm,
		InfoBubble: info,
	})

	tea.NewProgram(&firstSelect).Run()

	if firstSelectForm.State == huh.StateCompleted && firstSelectForm.GetString("targetgroup") == "create" {
		/*
			m := confirmmodel.NewModel(confirmmodel.ConfirmModelConfig{
				InfoBubble: info,
				Title:      "Create a target group?",
				Key:        "confirm-create-targetgroup",
			})
			tm, _ := tea.NewProgram(&m).Run()

			mm := tm.(confirmmodel.Model)

			if mm.State == huh.StateCompleted && mm.Confirmed {
				info = generateInfo(targetGroupDescription, rules, service)
				targetGroup = selectTargetGroupJSON(info)
			}
		*/

		info = generateInfo(targetGroupDescription, rules, service)
		targetGroup = selectTargetGroupJSON(info)

		if targetGroup != "" {
			tgConf := viper.New()
			tgConf.SetConfigFile(targetGroup)
			err := tgConf.ReadInConfig()
			if err != nil {
				logger.Fatal("Could not read file:", err)
			}

			targetGroupDescription = generateTargetGroupDescription(tgConf.GetString("targetGroupName"), targetGroup)
		}
	}

	/*
		mm := tm.(Model)

		if mm.form.State == huh.StateCompleted {
			fmt.Printf("So yea basically, you selected: %s, Lvl. %d\n", mm.form.GetString("class"), mm.form.GetInt("level"))
		}
	*/

	info = generateInfo(targetGroupDescription, rules, service)

	fmt.Println(info)

	fmt.Println("Done")
}
