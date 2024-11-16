package main

import (
	"fmt"
	"os"
//	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	viewport 		viewport.Model
	messages		[]string
	textarea		textarea.Model
	senderStyle		lipgloss.Style
	responderStyle	lipgloss.Style
	err				error
}

func main() {
	go func() {
		err := serveModel()
		if err != nil {
			fmt.Printf("Error when serving: %v\n", err)
		}
	}()
	p := tea.NewProgram(initalModel())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
}

func initalModel() Model {
	ta := textarea.New()
	ta.Placeholder = "Type a messge..."
	ta.Focus()

	ta.Prompt = "| "
	ta.CharLimit = 280
	ta.SetWidth(30)
	ta.SetHeight(3)

	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.ShowLineNumbers = false
	vp := viewport.New(100, 5)
	vp.SetContent(`You are in the chat room. Type a message and press Enter to send.`)
	ta.KeyMap.InsertNewline.SetEnabled(false)
	return Model{
		textarea:		ta,
		messages:		[]string{},
		viewport:		vp,
		senderStyle: 	lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		responderStyle:	lipgloss.NewStyle().Foreground(lipgloss.Color("1")),
		err:			nil,
	}
}

func (m Model) Init() tea.Cmd {
	return textarea.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case "enter":
			v := m.textarea.Value()
			if v == "" {
				return m, nil
			}
			v, err := m.Send(v)
			err = m.Reply(v)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			}
			return m, nil
		default:
			var cmd tea.Cmd
			m.textarea, cmd = m.textarea.Update(msg)
			return m, cmd
		}
	case cursor.BlinkMsg:
		var cmd tea.Cmd
		m.textarea, cmd = m.textarea.Update(msg)
		return m, cmd
	default:
		return m, nil

	}
}

func (m Model) View() string {
	return fmt.Sprintf(
		"%s\n\n%s",
		m.viewport.View(),
		m.textarea.View(),
	) + "\n\n"
} 
