package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/google/uuid"

	tea "github.com/charmbracelet/bubbletea"
	db "github.com/hf-chow/glloop/internal/database"
	llm "github.com/hf-chow/glloop/internal/llm"
	chat "github.com/hf-chow/glloop/internal/chat"
	config "github.com/hf-chow/glloop/internal/config"
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
	clearHistory(*state.DB, userID)

	p := tea.NewProgram(chat.InitModel(userID, state))
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
}

func login(q db.Queries) uuid.UUID {
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
	userID, err := q.GetIDByUsername(context.Background(), username)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}

	return userID
}

func clearHistory(q db.Queries, userID uuid.UUID) {
	fmt.Println("Do you want to continue where you left off? ([Y]/N)")
	var resp string
	fmt.Scanln(&resp)
	if strings.ToLower(resp) == "n" {
		q.DeleteAllHistoryByUserID(context.Background(), userID)
	}
}
