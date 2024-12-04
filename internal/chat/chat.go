package chat

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
	db "github.com/hf-chow/glloop/internal/database"
)

type Request struct {
	Model 	string `json:"model"`
	Prompt 	string `json:"prompt"`
	Stream	bool 	`json:"stream"`
}

type Response struct {
	Model  					string	`json:"model"`
	Created_at				string	`json:"created_at"`
	Response				string	`json:"response"`
	Done					bool	`json:"done"`
	Context					[]int	`json:"context"`
	Total_duration			int		`json:"total_duration"`
	Load_duration			int		`json:"load_duration"`
	Prompt_eval_count		int		`json:"prompt_eval_count"`
	Prompt_eval_duration	int		`json:"prompt_eval_duration"`
	Eval_count				int		`json:"eval_count"`
	Eval_duration			int		`json:"eval_duration"`
}

type BotResponseMsg string

func (m *Model) setAndGo() {
	m.viewport.SetContent(strings.Join(m.messages, "\n"))
	m.textarea.Reset()
	m.viewport.GotoBottom()
}

func (m *Model) Send(v string) {
	m.messages = append(
		m.messages, m.senderStyle.Render(fmt.Sprintf(m.CurrentUser) +": ") + v,
	)
	go m.sendRequest(v)
	m.setAndGo()
}

func (m *Model) SysReply() {
	m.messages = append(m.messages, m.SystemStyle.Render("Standby..."))
	m.setAndGo()
}

func (m *Model) Blink() {
	// For testing only
	m.messages = append(m.messages, m.SystemStyle.Render("Blink..."))
	m.setAndGo()
}

func (m *Model) sendRequest(v string) {
	m.requestCh <- v
}

func (m *Model) fetchReply() {
	msg := <- m.requestCh
	postBody, err := json.Marshal(Request{
		Model: "llama3.2",
		Prompt: msg,
		Stream: false,
	})
	if err != nil {
		fmt.Printf("Error marshalling request: %v\n", err)
	}
	buf := bytes.NewBuffer(postBody)
	resp, err := http.Post(
		"http://localhost:11434/api/generate", "application/json", buf,
	)
	if err != nil {
		fmt.Printf("Error marshalling request: %v\n", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
	}

	var modelResp Response
	err = json.Unmarshal(body, &modelResp)
	if err != nil {
		fmt.Printf("Error unmarshalling response: %v\n", err)
	}

	userID := db.GetIDByUsername(context.Background(), m.CurrentUser)

	historyArgs := db.CreateHistoryParams{
		ID: 			uuid.New(),
		userID:			userID,
		CreatedAt:		time.Now(),
		Prompt:			msg,
		Reply:			modelResp.Response,
	}
	_, err = db.CreateHisory(context.Background(), historyArgs)
	if err != nil {
		fmt.Printf("error creating chat history: %s\n", err)
	}

	m.responseCh <- modelResp.Response
}

func (m *Model) reply() {
	response :=  <- m.responseCh
// 	fmt.Printf("Got reply: %s", response)
	m.messages = append(m.messages, m.responderStyle.Render("Bot: ")+response)
	m.setAndGo()
}

func (m *Model) BotReply(msg BotResponseMsg) {
	m.messages = append(m.messages, m.responderStyle.Render("Bot: ") + string(msg))
	m.viewport.SetContent(strings.Join(m.messages, "\n"))
	m.textarea.Reset()
	m.viewport.GotoBottom()
}

func (m *Model) WaitForResponse() tea.Cmd {
	return func() tea.Msg {
		response := <- m.responseCh
//		fmt.Println("Got response:", response)
		return BotResponseMsg(response)
	}
}
