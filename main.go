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

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"

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

	username := login(*state.DB)
	clearHistory(*state.DB, username)

	p := tea.NewProgram(chat.InitModel(username, state))
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
}

func login(q db.Queries) string {
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

func clearHistory(q db.Queries, username string) {
	fmt.Println("Do you want to continue where you left off? ([Y]/N)")
	var resp string
	fmt.Scanln(&resp)
	if strings.ToLower(resp) != "y" {
		userID, err := q.GetIDByUsername(context.Background(), username)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		q.DeleteAllHistoryByUserID(context.Background(), userID)
	}
}
