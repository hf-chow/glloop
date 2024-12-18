package components

import (
	"context"
	"strings"

	"github.com/google/uuid"

	tea "github.com/charmbracelet/bubbletea"
	db "github.com/hf-chow/glloop/internal/database"
)

var historyChoices = []string{"Yes", "No"}

type HistoryModel struct {
	cursor int
	choice string
}

func (m HistoryModel) Init() tea.Cmd {
	return nil
}

func (m HistoryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "enter":
			m.choice = historyChoices[m.cursor]
			return m, tea.Quit

		case "down", "j":
			m.cursor++
			if m.cursor >= len(historyChoices) {
				m.cursor = 0
			}

		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(historyChoices) - 1
			}
		}
	}
	return m, nil
}

func (m HistoryModel) View() string {
	s := strings.Builder{}
	s.WriteString("Would you like to continue where you left off?\n\n")

	for i := 0; i < len(historyChoices); i ++ {
		if m.cursor == i {
			s.WriteString("(•)")
		} else {
			s.WriteString("( )")
		}
		s.WriteString(historyChoices[i])
		s.WriteString("\n")
	}
	s.WriteString("\n(press q to quit)\n")

	return s.String()
}

func (m HistoryModel) clearHistory(q db.Queries, userID uuid.UUID) {
	q.DeleteAllHistoryByUserID(context.Background(), userID)
}

