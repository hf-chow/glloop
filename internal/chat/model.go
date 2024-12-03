package chat

import (
	"fmt"

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
	SystemStyle		lipgloss.Style
	requestCh		chan string
	responseCh		chan string
	username		string
	err				error
}


func InitModel(username string) Model {
	requestCh := make(chan string, 1)
	responseCh := make(chan string, 1)

	ta := textarea.New()
	ta.Placeholder = "Type a messge..."
	ta.Focus()

	ta.Prompt = "| "
	ta.SetWidth(100)
	ta.SetHeight(5)

	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.ShowLineNumbers = false
	vp := viewport.New(100, 10)
	vp.SetContent(`You are in the chat room. Type a message and press Enter to send.`)
	ta.KeyMap.InsertNewline.SetEnabled(false)
	return Model{
		textarea:		ta,
		messages:		[]string{},
		viewport:		vp,
		senderStyle: 	lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		responderStyle:	lipgloss.NewStyle().Foreground(lipgloss.Color("1")),
		SystemStyle:	lipgloss.NewStyle().Foreground(lipgloss.Color("2")),
		requestCh: 		requestCh,
		responseCh: 	responseCh,
		username: 		username,
		err:			nil,
	}
}

func (m Model) Init() tea.Cmd {
	return textarea.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
//	fmt.Printf("Received message type: %T\n", msg)
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
			m.Send(v)
			m.textarea.Reset()
			return m, m.WaitForResponse()
		default:
			var cmd tea.Cmd
			m.textarea, cmd = m.textarea.Update(msg)
			return m, cmd
		}
	case cursor.BlinkMsg:
		var cmd tea.Cmd
		go m.fetchReply()
		m.textarea, cmd = m.textarea.Update(msg)
		return m, cmd
	case BotResponseMsg:
		m.BotReply(msg)
		return m, nil
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
