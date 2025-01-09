package components

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hf-chow/glloop/internal/config"
	db "github.com/hf-chow/glloop/internal/database"
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

type LoginModelState struct {
	Config 	*config.Config
	DB		*db.Queries
}

type LoginModel struct {
	inputs 		[]textinput.Model
	focused 	int
	err			error
	UserID		uuid.UUID
	Username	string
	State		*LoginModelState
}

func InitLoginModel() LoginModel {
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
				m.Username = m.inputs[username].Value()
				return m, tea.Quit
			}
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
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


func (m LoginModel) Login(q db.Queries, username string) (uuid.UUID, error) {
	if username == "" {
		return uuid.UUID{}, fmt.Errorf("username cannot be empty")
	}

	exists, err := q.UsernameExists(context.Background(), username)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("error when checking if username exists: %s", err)
	}

	var userID uuid.UUID
	if exists {
		userID, err = q.GetIDByUsername(context.Background(), username)
		if err != nil {
			return uuid.UUID{}, fmt.Errorf("error when retrieving userID: %s", err)
		}
	} else {
		userID = uuid.New()
		userArgs := db.CreateUserParams{
			ID:        userID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      username,
		}
		_, err := q.CreateUser(context.Background(), userArgs)
		if err != nil {
			return uuid.UUID{}, fmt.Errorf("error when creating new user: %s", err)
		}
	}
	return userID, nil
}
