package chat

import (
	"context"
	"fmt"
	"time"

	db "github.com/hf-chow/glloop/internal/database"
	"github.com/google/uuid"
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

func (m *Model) createHistoryFromLastMessage(modelResp ChatResponse, lastMsg ChatMessage) error {
	q := m.ModelState.DB
	historyArgs := db.CreateHistoryParams{
		ID:        uuid.New(),
		UserID:    m.CurrentUserID,
		CreatedAt: time.Now(),
		Prompt:    lastMsg.Content,
		Reply:     modelResp.Message.Content,
	}
	_, err := q.CreateHistory(context.Background(), historyArgs)

	if err != nil {
		fmt.Printf("error creating chat history: %s\n", err)
		return err
	}
	return nil
}
