package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"

	tea "github.com/charmbracelet/bubbletea"
	chat "github.com/hf-chow/glloop/internal/chat"
	comp "github.com/hf-chow/glloop/internal/components"
	config "github.com/hf-chow/glloop/internal/config"
	db "github.com/hf-chow/glloop/internal/database"
	llm "github.com/hf-chow/glloop/internal/llm"
)

func main() {
	go func() {
		err := llm.ServeModel()
		if err != nil {
			fmt.Printf("error when serving: %v\n", err)
		}
	}()

	cfg, err := config.ReadConfig()
	if err != nil {
		fmt.Printf("error when reading config: %v\n", err)
	}
	state := &chat.State{Config: &cfg}

	dbtx, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatal("fail to connect to DB")
	}
	dbQueries := db.New(dbtx)
	state.DB = dbQueries

	userID := login(*state.DB)

	p := tea.NewProgram(comp.HistoryModel{})
	m, err := p.Run()
	if err != nil {
		fmt.Printf("error building HistoryModel: %v", err)
		os.Exit(1)
	}

	historyModel, ok := m.(comp.HistoryModel)
	if !ok {
		fmt.Println("Error: returned model is not of type HistoryModel")
		os.Exit(1)
	}
	if historyModel.choice == "Yes" {
		historyModel.clearHistory(*state.DB, userID)
	}
	p = tea.NewProgram(chat.InitModel(userID, state))
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
}

func login(q db.Queries) uuid.UUID {
	fmt.Println("Enter your username: ")
	var username string
	fmt.Scanln(&username)

	userArgs := db.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	}

	q.CreateUser(context.Background(), userArgs)
	userID, err := q.GetIDByUsername(context.Background(), username)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}

	return userID
}
