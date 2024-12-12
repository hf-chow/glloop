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

type ChatMessage struct {
	Role 	string 	`json:"role"`
	Content	string	`json:"content"`
}

type GenerateRequest struct {
	Model 	string 	`json:"model"`
	Prompt 	string 	`json:"prompt"`
	Stream	bool   	`json:"stream"`
}

type ChatRequest struct {
	Model 		string 			`json:"model"`
	Messages 	[]ChatMessage	`json:"messages"`	
	Stream		bool 			`json:"stream"`
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
	q := m.ModelState.DB

	username, err := q.GetUsernameByID(context.Background(), m.CurrentUserID)

	if err != nil {
		fmt.Printf("error in retrieving username: %v\n", err)
	}
	m.messages = append(
		m.messages, m.senderStyle.Render(fmt.Sprintf(username) +": ") + v,
	)
	go m.sendRequest(v)
	m.setAndGo()
}

func (m *Model) sendRequest(v string) {
	m.requestCh <- v
}

func (m *Model) fetchSingleReply() {
	q := m.ModelState.DB

	msg := <- m.requestCh
	dat, err := json.Marshal(GenerateRequest{
		Model: "llama3.2",
		Prompt: msg,
		Stream: false,
	})
	if err != nil {
		fmt.Printf("Error marshalling request: %v\n", err)
	}
	buf := bytes.NewBuffer(dat)
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

	historyArgs := db.CreateHistoryParams{
		ID: 			uuid.New(),
		UserID:			m.CurrentUserID,
		CreatedAt:		time.Now(),
		Prompt:			msg,
		Reply:			modelResp.Response,
	}

	_, err = q.CreateHistory(context.Background(), historyArgs)
	if err != nil {
		fmt.Printf("error creating chat history: %s\n", err)
	}

	m.responseCh <- modelResp.Response
}

func (m *Model) fetchReplyWithHistory() {
	q := m.ModelState.DB

	lastPrompt := <- m.requestCh

	msgs, err := m.createMessagesFromHistory(lastPrompt)
	if err != nil {
		fmt.Printf("error creating messages from history: %s\n", err)
	}

	dat, err := json.Marshal(ChatRequest{
		Model: 		m.CurrentModel,
		Messages: 	msgs,
		Stream: 	false,
	})
	if err != nil {
		fmt.Printf("Error marshalling request: %v\n", err)
	}

	buf := bytes.NewBuffer(dat)
	resp, err := http.Post(
		"http://localhost:11434/api/chat", "application/json", buf,
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

	lastMsg := msgs[len(msgs)-1].Content
	historyArgs := db.CreateHistoryParams{
		ID: 			uuid.New(),
		UserID:			m.CurrentUserID,
		CreatedAt:		time.Now(),
		Prompt:			lastMsg,
		Reply:			modelResp.Response,
	}
	_, err = q.CreateHistory(context.Background(), historyArgs)

	if err != nil {
		fmt.Printf("error creating chat history: %s\n", err)
	} 
	m.responseCh <- modelResp.Response
}

func (m *Model) reply() {
	response :=  <- m.responseCh
 	fmt.Printf("Got reply: %s", response)
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
