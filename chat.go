package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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


func (m *Model) setAndGo() {
	m.viewport.SetContent(strings.Join(m.messages, "\n"))
	m.textarea.Reset()
	m.viewport.GotoBottom()
}

func (m *Model) Send(v string) (string, error) {
	m.messages = append(m.messages, m.senderStyle.Render("You: ") + v)
	m.setAndGo()
	return v, nil
}

func (m *Model) Reply(msg string) error {
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
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var modelResp Response
	err = json.Unmarshal(body, &modelResp)
	if err != nil {
		return err
	}
	m.messages = append(m.messages, m.responderStyle.Render("Bot: ")+modelResp.Response)
	m.setAndGo()

	return nil
}
