package main

import (
	"strings"
)

func (m *Model) setAndGo() {
	m.viewport.SetContent(strings.Join(m.messages, "\n"))
	m.textarea.Reset()
	m.viewport.GotoBottom()

}

func (m *Model) Send(v string) error {
	m.messages = append(m.messages, m.senderStyle.Render("You: ") + v)
	m.setAndGo()
	return nil
}

func (m *Model) Reply() error {
	last_msg := strings.ReplaceAll(
		m.messages[len(m.messages)-1], "You: ", "Your last message is: ",
	)
	m.messages = append(m.messages, m.responderStyle.Render("Bot: ") + last_msg)
	m.setAndGo()
	return nil
}
