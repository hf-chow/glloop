package chat

import (
	"context"
	"fmt"
)


func (m *Model) createMessagesFromHistory(lastPrompt string) ([]Message, error) {
	q := m.ModelState.DB

	username, err := q.GetUsernameByID(context.Background(), m.CurrentUserID)
	if err != nil {
		return []Message{}, err
	}

	histories, err := q.GetAllHistoryByUserID(context.Background(), m.CurrentUserID)
	if err != nil {
		return []Message{}, err
	}

	msgs := []Message{}

	for _, history := range histories {
		msg := Message{
			Role: 		username,
			Content: 	history.Prompt,
		}
		msgs = append(msgs, msg)

		replyMsg := Message {
			Role: 		"assistant",
			Content: 	history.Reply,
		}
		msgs = append(msgs, replyMsg)

		lastMsg := Message {
			Role: 		username,
			Content: 	lastPrompt,
		}
		msgs = append(msgs, lastMsg)
	}

	return msgs, nil
}

func (m *Model) historyExist() bool {
	q := m.ModelState.DB
	histories, err := q.GetAllHistoryByUserID(context.Background(), m.CurrentUserID)
	if err != nil {
		fmt.Printf("error getting histories by id: %s", err)
		return false
	}
	if len(histories) > 0 {
		return true
	} else {
		return false
	}
}
