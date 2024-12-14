package chat

import (
	"context"
)

func (m *Model) createMessagesFromHistory(lastPrompt string) ([]ChatMessage, error) {
	q := m.ModelState.DB

	//username, err := q.GetUsernameByID(context.Background(), m.CurrentUserID)
	//if err != nil {
	//	return []Message{}, err
	//}
	histories, err := q.GetAllHistoryByUserID(context.Background(), m.CurrentUserID)
	if err != nil {
		return []ChatMessage{}, err
	}

	msgs := []ChatMessage{}
	for _, history := range histories {
		msg := ChatMessage{
			Role: 		"user",
			Content: 	history.Prompt,
		}
		msgs = append(msgs, msg)
		if len(history.Reply) > 0 {
			replyMsg := ChatMessage {
				Role: 		"assistant",
				Content: 	history.Reply,
			}
			msgs = append(msgs, replyMsg)
		}
	}
	lastMsg := ChatMessage {
		Role: 		"user",
		Content: 	lastPrompt,
	}
	msgs = append(msgs, lastMsg)
	return msgs, nil
}

func (m *Model) historyExist() bool {
	q := m.ModelState.DB
	histories, err := q.GetAllHistoryByUserID(context.Background(), m.CurrentUserID)
	if err != nil {
		return false
	}
	if len(histories) > 0 {
		return true
	} else {
		return false
	}
}
