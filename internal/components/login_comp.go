package components

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"

	db "github.com/hf-chow/glloop/internal/database"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	errMsg error
)

const (
	username = iota 
)

const (
	hotPink = lipgloss.Color("#FF06B7")
	darkGray = lipgloss.Color("#767676")
)

var (
	inputStyle 	  = lipgloss.NewStyle().Foreground(hotPink)
	continueStyle = lipgloss.NewStyle().Foreground(darkGray)
)

type LoginModel struct {
	inputs 		[]textinput.Model
	focused 	int
	err			error
	UserID		uuid.UUID
}

func InitLoginModel() LoginModel {
	// Possibly add a password feature?
	var inputs []textinput.Model = make([]textinput.Model, 1)
	inputs[username] = textinput.New()
	inputs[username].Placeholder = "How would you like to be referred?"
	inputs[username].Focus()
	inputs[username].Prompt = ""
	return LoginModel {
		inputs:		inputs,
		focused:	0,
		err:		nil,
	}
}

func (m LoginModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m LoginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.focused == len(m.inputs)-1 {
				return m, tea.Quit
			}
			m.nextInput()
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyShiftTab, tea.KeyCtrlP:
			m.prevInput()
		case tea.KeyTab, tea.KeyCtrlN:
			m.nextInput()
		}
		for i := range m.inputs {
			m.inputs[i].Blur()
		}
		m.inputs[m.focused].Focus()

	case errMsg:
		m.err = msg
		return m, nil
	}

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}

func (m LoginModel) View() string {
	return fmt.Sprintf(
`
 %s
 %s
 %s
`,
		inputStyle.Width(30).Render("Username"),
		m.inputs[username].View(),
		continueStyle.Render("Press Enter to continue"),
	) + "\n"
}

func (m *LoginModel) nextInput() {
	m.focused = (m.focused + 1) % len(m.inputs)
}

func (m *LoginModel) prevInput() {
	m.focused--
	if m.focused < 0 {
		m.focused = len(m.inputs) - 1
	}
}

func (m *LoginModel) Login(q db.Queries) uuid.UUID {
	fmt.Println("Enter your username: ")
	name := m.inputs[username].View()
	userID, err := q.GetIDByUsername(context.Background(), name)
	if err != nil {
		userArgs := db.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      name,
		}
		q.CreateUser(context.Background(), userArgs)
		userID = userArgs.ID
	}
	return userID
}
