package model

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	email "github.com/duanechan/monitoring-utils/email/internal"
)

type (
	EmailModel struct {
		helpMode    bool
		quitMode    bool
		quitChoice  int
		width       int
		height      int
		err         error
		parseResult email.ParseResult
		parser      ParserModel
		editor      EditorModel
		sender      SenderModel
	}

	initializeEditor struct{}
	toggleHelp       struct{}
	toggleQuit       struct{}
)

func InitializeModel() EmailModel {
	return EmailModel{
		parser: initParser(),
		sender: initSender(),
	}
}

func (e EmailModel) Init() tea.Cmd {
	return tea.Batch(tea.ClearScreen, tea.EnterAltScreen)
}

func (e EmailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "?":
			return e, e.showHelp
		case "esc":
			return e, e.quit
		case "ctrl+c":
			return e, tea.Quit
		case "enter":
			if e.quitMode && e.quitChoice == 1 {
				return e, tea.Quit
			} else if e.quitMode && e.quitChoice == 0 {
				return e, e.quit
			}
		case "left":
			if e.quitMode && e.quitChoice > 0 {
				e.quitChoice--
			}
		case "right":
			if e.quitMode && e.quitChoice < 1 {
				e.quitChoice++
			}
		}

	case tea.WindowSizeMsg:
		e.width = msg.Width
		e.height = msg.Height

	case toggleHelp:
		e.helpMode = !e.helpMode

	case toggleQuit:
		e.quitMode = !e.quitMode

	case parserMessage:
		e.parseResult, e.err = msg.result, msg.err
		return e, e.initializeEditor

	case initializeEditor:
		e.editor = initEditor(e.parseResult.Raw)
		return e, nil
	}

	var parser tea.Model
	parser, cmd = e.parser.Update(msg)
	e.parser = parser.(ParserModel)
	return e, cmd
}

func (e EmailModel) View() string {
	if e.quitMode {
		return e.quitView()
	}

	sections := []string{}
	// sections = append(sections, e.headerView())

	if e.parser.result.IsEmpty() {
		sections = append(sections, e.parser.View())
	} else {
		sections = append(sections, e.editor.View())
	}

	if e.err != nil {
		sections = append(sections, "\n"+e.errorView())
	} else {
		sections = append(sections, "\n\n\n")
	}

	sections = append(sections, e.helpView())

	return lipgloss.Place(
		e.width,
		e.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			Render(lipgloss.JoinVertical(
				lipgloss.Center,
				e.headerView(),
				lipgloss.JoinVertical(lipgloss.Left, sections...),
			),
			),
	)
}

func (e *EmailModel) initializeEditor() tea.Msg { return initializeEditor{} }

func (e EmailModel) quit() tea.Msg { return toggleQuit{} }

func (e EmailModel) showHelp() tea.Msg { return toggleHelp{} }

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

func (e EmailModel) helpView() string {
	if !e.helpMode {
		return Help.Render(lipgloss.NewStyle().Padding(1, 2).Render("? / toggle help"))
	}

	commands := lipgloss.JoinHorizontal(
		lipgloss.Center,
		lipgloss.NewStyle().Padding(1, 2).Render("? / toggle help"),
		lipgloss.NewStyle().Padding(1, 2).Render("ctrl+c, esc / exit program"),
	)

	return Help.Render(commands)
}

func (e EmailModel) errorView() string {
	if e.err != nil {
		return Error.Render(fmt.Sprint("Error: ", e.err))
	}
	return ""
}

func (e EmailModel) quitView() string {
	var cancel, confirm string
	if e.quitChoice == 0 {
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
