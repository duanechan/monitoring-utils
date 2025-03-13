package model

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	email "github.com/duanechan/monitoring-utils/email/internal"
)

type (
	EmailModel struct {
		width            int
		height           int
		cursor           int
		selectedTemplate int
		sendResults      []string
		templates        []templates
		err              error
		mode             mode
		parseResult      email.ParseResult
		config           email.EmailConfig
		input            textinput.Model
		table            table.Model
		progressBar      progress.Model
		progressChan     chan float64
	}

	templates struct {
		name     string
		template email.Template
	}

	mode struct {
		Quit   bool
		Help   bool
		Parser bool
		Editor bool
		Send   bool
	}

	sendEmails       struct{}
	readInputMessage struct{}
	initializeEditor struct{}
	progressMsg      struct {
		progress float64
		results  []string
	}
)

func InitializeModel() EmailModel {
	config, err := email.LoadConfig()
	if err != nil {
		fmt.Printf("error loading config: %s", err)
	}

	return EmailModel{
		input:        initParser(),
		progressChan: make(chan float64),
		templates: []templates{
			{name: "CRED", template: email.Credentials},
			{name: "LATE", template: email.Late},
			{name: "ABST", template: email.Credentials},
		},
		mode: mode{
			Quit:   false,
			Help:   false,
			Parser: true,
			Editor: false,
			Send:   false,
		},
		config: config,
	}
}

func (e EmailModel) Init() tea.Cmd {
	return tea.Batch(tea.ClearScreen, tea.EnterAltScreen, textinput.Blink)
}

func (e EmailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "?":
			e.mode.Help = !e.mode.Help
			return e, nil
		case "esc":
			if e.mode.Send && e.progressBar.Percent() == 1.0 {
				e.progressBar.SetPercent(0.0)
				e.sendResults = nil
				e.mode.Parser = true
				e.mode.Editor = false
				e.mode.Send = false
				e.mode.Help = false
				return e, textinput.Blink
			} else if e.progressBar.Percent() < 1.0 {
				return e, nil
			}
			e.mode.Quit = !e.mode.Quit
		case "ctrl+c":
			return e, tea.Quit
		case "enter":
			e.err = nil
			switch {
			case e.mode.Quit && e.cursor == 1:
				return e, tea.Quit
			case e.mode.Quit && e.cursor == 0:
				e.mode.Quit = false
				return e, nil
			case e.mode.Parser:
				return e, e.handleInput
			case e.mode.Editor && e.cursor == 1:
				return e, e.startSending
			case e.mode.Editor && e.cursor == 0:
				e.mode.Editor = false
				e.mode.Parser = true
				e.input.Focus()
				return e, textinput.Blink
			}
		case "shift+tab":
			if (!e.mode.Quit || !e.mode.Editor || !e.mode.Send) && e.selectedTemplate > 0 {
				e.selectedTemplate--
			}
		case "tab":
			if (!e.mode.Quit || !e.mode.Editor || !e.mode.Send) && e.selectedTemplate < len(e.templates)-1 {
				e.selectedTemplate++
			}
		case "left":
			if (e.mode.Quit || e.mode.Editor) && e.cursor > 0 {
				e.cursor--
			}
		case "right":
			if (e.mode.Quit || e.mode.Editor) && e.cursor < 1 {
				e.cursor++
			}
		}

	case tea.WindowSizeMsg:
		e.width = msg.Width
		e.height = msg.Height

	case readInputMessage:
		input := strings.ReplaceAll(e.input.Value(), "\"", "")
		records, err := email.ParseData(input)
		if err != nil {
			e.err = err
			return e, nil
		}
		e.mode.Parser = false
		e.parseResult = email.ValidateRecords(records)

		e.input.Reset()
		e.input.Blur()
		return e, e.initializeEditor

	case initializeEditor:
		e.mode.Editor = true
		e.table = initEditor(e.parseResult)

	case sendEmails:
		e.mode.Send = true
		return e, e.SendEmails()

	case progressMsg:
		e.sendResults = msg.results
		cmds = append(cmds, e.progressBar.SetPercent(float64(msg.progress)))
	}

	e.input, cmd = e.input.Update(msg)
	cmds = append(cmds, cmd)

	e.table, cmd = e.table.Update(msg)
	cmds = append(cmds, cmd)

	model, cmd := e.progressBar.Update(msg)
	e.progressBar = model.(progress.Model)
	cmds = append(cmds, cmd)

	return e, tea.Batch(cmds...)
}

func (e EmailModel) View() string {
	if e.mode.Quit {
		return e.quitView()
	}

	sections := []string{}
	sections = append(sections, e.headerView())

	if e.mode.Parser {
		sections = append(
			sections, lipgloss.NewStyle().
				Padding(1, 1).
				Background(lipgloss.Color("42")).
				Render("Email Template: "+e.templates[e.selectedTemplate].name))
		sections = append(
			sections, lipgloss.NewStyle().
				Padding(1, 0).
				Render(e.input.View()))
	} else if e.mode.Editor {
		sections = append(sections, e.table.View())
		if e.mode.Send {
			if !e.progressBar.IsAnimating() && e.progressBar.Percent() == 1.0 {
				sections = append(
					sections, lipgloss.NewStyle().
						Bold(true).
						Foreground(lipgloss.Color("46")).
						Render("\nDone!\n"))
				if len(e.sendResults) > 0 {
					sections = append(sections, e.sendResults...)
				}
				sections = append(sections, "Press ESC to go back")
			} else {
				sections = append(sections, "\nSending emails...\n")
				sections = append(sections, e.progressBar.View())
			}
		} else {
			sections = append(sections, e.resultView())
		}
	}

	if e.err != nil {
		sections = append(sections, "\n"+e.errorView())
	} else {
		sections = append(sections, "\n\n\n")
	}

	sections = append(sections, e.helpView())

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (e EmailModel) startSending() tea.Msg { return sendEmails{} }

func (e EmailModel) handleInput() tea.Msg { return readInputMessage{} }

func (e EmailModel) initializeEditor() tea.Msg { return initializeEditor{} }

func (e *EmailModel) SendEmails() tea.Cmd {
	e.progressBar = progress.New(progress.WithGradient("#005DAD", "#6796BF"))
	total := float64(len(e.parseResult.Recipients))
	progress := 0.0
	e.input.Focus()

	return func() tea.Msg {
		results := []string{}
		progressChan := make(chan float64)
		resultChan := make(chan string)

		go func() {
			for p := range progressChan {
				e.progressBar.SetPercent(p)
			}
		}()

		for _, r := range e.parseResult.Recipients {

			em := email.Email{
				Body:   e.templates[e.selectedTemplate].template,
				To:     email.User{Name: r.Name, Email: r.Email},
				Config: e.config,
			}

			if err := em.Send(); err != nil {
				results = append(
					results,
					lipgloss.NewStyle().
						Foreground(lipgloss.Color(Red)).
						Render(fmt.Sprintf("✖ Failed: %s\n", r.Email)))
			} else {
				results = append(
					results, lipgloss.NewStyle().
						Foreground(lipgloss.Color("46")).
						Render(fmt.Sprintf("✔ Sent: %s\n", r.Email)))
			}

			progress++
			progressChan <- progress / total
		}

		close(progressChan)
		close(resultChan)
		return progressMsg{
			progress: 1.0,
			results:  results,
		}
	}
}

func (e EmailModel) headerView() string {
	header := " ______                 _ _   _    _      _\n" +
		"|  ____|               (_) | | |  | |    | |\n" +
		"| |__   _ __ ___   __ _ _| | | |__| | ___| |_ __   ___ _ __ \n" +
		"|  __| | '_ ' _ \\ / _' | | | |  __  |/ _ \\ | '_ \\ / _ \\ '__|\n" +
		"| |____| | | | | | (_| | | | | |  | |  __/ | |_) |  __/ |\n" +
		"|______|_| |_| |_|\\__,_|_|_| |_|  |_|\\___|_| .__/ \\___|_|\n" +
		"										   | |\n" +
		"										   |_|\n"

	return Header.Render(header)
}

func (e EmailModel) resultView() string {
	resultStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00"))
	validStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("46"))

	var result string

	if badEmails := len(e.parseResult.BadEmails); badEmails > 1 {
		result = resultStyle.Render(fmt.Sprintf("There are %d bad emails detected in the file.", badEmails))
	} else if badEmails == 1 {
		result = resultStyle.Render(fmt.Sprintf("There is %d bad email detected in the file.", badEmails))
	} else {
		result = validStyle.Render("✔ All emails are valid!")
	}

	var cancel, confirm string
	if e.cursor == 0 {
		cancel = QuitButtonContainer.Render(QuitSelectedStyle.Render("Cancel"))
		confirm = QuitButtonContainer.Render(QuitUnselectedStyle.Render("Confirm"))
	} else {
		cancel = QuitButtonContainer.Render(QuitUnselectedStyle.Render("Cancel"))
		confirm = QuitButtonContainer.Render(QuitSelectedStyle.Render("Confirm"))
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		"\n",
		result,
		"\n",
		lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(2, 3).
			Render(
				lipgloss.JoinVertical(
					lipgloss.Center,
					"Do you want to send the emails?",
					lipgloss.JoinHorizontal(lipgloss.Center, cancel, confirm),
				),
			),
	)
}

func (e EmailModel) helpView() string {
	if !e.mode.Help {
		return lipgloss.NewStyle().
			Foreground(Gray).
			Padding(1, 2).
			Render("? / toggle help")
	}

	commands := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.JoinHorizontal(
			lipgloss.Center,
			lipgloss.NewStyle().Padding(1, 2).Render("? / toggle help"),
			lipgloss.NewStyle().Padding(1, 2).Render("esc / quit"),
			lipgloss.NewStyle().Padding(1, 2).Render("ctrl+c / force quit"),
		),
		lipgloss.JoinHorizontal(
			lipgloss.Center,
			lipgloss.NewStyle().Padding(1, 2).Render("tab / next template"),
			lipgloss.NewStyle().Padding(1, 2).Render("shift+tab / previous template"),
		),
	)

	return lipgloss.NewStyle().
		Foreground(Gray).
		Render(commands)
}

func (e EmailModel) errorView() string {
	if e.err != nil {
		return Error.Render(fmt.Sprint("Error: ", e.err))
	}
	return ""
}

func (e EmailModel) quitView() string {
	var cancel, confirm string
	if e.cursor == 0 {
		cancel = QuitButtonContainer.Render(QuitSelectedStyle.Render("Cancel"))
		confirm = QuitButtonContainer.Render(QuitUnselectedStyle.Render("Confirm"))
	} else {
		cancel = QuitButtonContainer.Render(QuitUnselectedStyle.Render("Cancel"))
		confirm = QuitButtonContainer.Render(QuitSelectedStyle.Render("Confirm"))
	}

	prompt := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(6, 12).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Center,
				"Are you sure you want to exit?\n",
				lipgloss.NewStyle().Foreground(Gray).Render("ctrl+c to force exit\n\n"),
				lipgloss.JoinHorizontal(lipgloss.Center, cancel, confirm),
			),
		)

	return lipgloss.Place(
		e.width, e.height,
		lipgloss.Center, lipgloss.Center,
		prompt,
	)
}

func initParser() textinput.Model {
	ti := textinput.New()
	ti.Prompt = lipgloss.NewStyle().Foreground(Primary).Render("Input the filepath: ")
	ti.PromptStyle = lipgloss.NewStyle().
		AlignHorizontal(lipgloss.Center).
		Bold(true)
	ti.Placeholder = "C:/path/to/file"
	ti.Width = 90
	ti.CharLimit = 150
	ti.Focus()

	return ti
}

func initEditor(result email.ParseResult) table.Model {
	rows := []table.Row{}

	redStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	greenStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("46"))

	for i, r := range result.Raw {
		var status string
		if e, exists := result.BadEmails[i+1]; exists {
			if strings.Contains(e, "Duplicate") {
				status = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00")).Render(fmt.Sprintf("━ %s", e))
			} else {
				status = redStyle.Render(fmt.Sprintf("✖ %s", e))
			}
		} else {
			status = greenStyle.Render("✔")
		}

		rows = append(rows, table.Row{fmt.Sprintf("%d", i+1), r[0], r[1], status})
	}

	cols := []table.Column{
		{Title: "Row", Width: 5},
		{Title: "Name", Width: 25},
		{Title: "Email", Width: 25},
		{Title: "Valid", Width: 100},
	}

	t := table.New(
		table.WithHeight(len(result.Raw)+1),
		table.WithColumns(cols),
		table.WithRows(rows),
		table.WithStyles(table.Styles{
			Header:   lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("210")),
			Selected: lipgloss.NewStyle().Background(lipgloss.Color("236")).Foreground(lipgloss.Color("15")),
		}),
	)

	return t
}
