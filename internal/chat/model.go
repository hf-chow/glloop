package chat

import (
	"fmt"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"

	"github.com/google/uuid"

	config "github.com/hf-chow/glloop/internal/config"
	db "github.com/hf-chow/glloop/internal/database"
	tea "github.com/charmbracelet/bubbletea"

)

type State struct {
	Config	*config.Config
	DB 		*db.Queries
}

type Model struct {
	viewport 		viewport.Model
	messages		[]string
	textarea		textarea.Model
	senderStyle		lipgloss.Style
	responderStyle	lipgloss.Style
	SystemStyle		lipgloss.Style
	requestCh		chan string
	responseCh		chan string
	CurrentUserID	uuid.UUID
	CurrentModel	string
	ModelState		*State
	err				error
}

func InitModel(userID uuid.UUID, s *State) Model {
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
	vp := viewport.New(100, 5)
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
		CurrentUserID: 	userID,
		CurrentModel: 	"llama3.2",
		ModelState: 	s,
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
		if m.historyExist() {
			go m.fetchReplyWithHistory()
		} else {
			go m.fetchSingleReply()
		}
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
