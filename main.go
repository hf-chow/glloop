package main

import _ "github.com/lib/pq"
import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
	"github.com/hf-chow/glloop/internal/database"
	db "github.com/hf-chow/glloop/internal/database"
	llm "github.com/hf-chow/glloop/internal/llm"
	chat "github.com/hf-chow/glloop/internal/chat"

)

type Model struct {
	viewport 		viewport.Model
	messages		[]string
	textarea		textarea.Model
	senderStyle		lipgloss.Style
	responderStyle	lipgloss.Style
	SystemStyle		lipgloss.Style
	requestCh		chan string
	responseCh		chan string
	username		string
	err				error
}

type State struct {
	Config	*Config
	DB 		*database.Queries
}

func main() {
	go func() {
		err := llm.ServeModel()
		if err != nil {
			fmt.Printf("error when serving: %v\n", err)
		}
	}()

	cfg, err := ReadConfig()
	if err != nil {
		fmt.Printf("error when reading config: %v\n", err)
	}
	state := &State{Config: &cfg}
	fmt.Printf("%s\n", cfg.DBURL)

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatal("fail to connect to DB")
	}
	dbQueries := database.New(db)
	state.DB = dbQueries

	username := login(*state.DB)

	p := tea.NewProgram(chat.initModel(username))
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
}

func login(q db.Queries) string{
	fmt.Println("Enter your username: ")
	var username string
	fmt.Scanln(&username)

	userArgs := db.CreateUserParams{
		ID: 		uuid.New(),
		CreatedAt:	time.Now(),
		UpdatedAt:	time.Now(),
		Name:		username,
	}

	q.CreateUser(context.Background(), userArgs)

	return username
}
