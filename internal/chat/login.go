package chat


import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	db "github.com/hf-chow/glloop/internal/database"
)

func (m *Model) UserLogin() {
	q := m.ModelState.DB
	name := m.ModelState.Config.CurrentUsername

	if name == "" {
		fmt.Print("username cannot be empty")
	}

	exists, err := q.UsernameExists(context.Background(), name)
	if err != nil {
		fmt.Printf("error when checking if username exists: %s", err)
	}

	var userID uuid.UUID
	if exists {
		userID, err = q.GetIDByUsername(context.Background(), name)
		if err != nil {
			fmt.Printf("error when retrieving userID: %s", err)
		}
	} else {
		userID = uuid.New()
		userArgs := db.CreateUserParams{
			ID:        userID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      name,
		}
		_, err := q.CreateUser(context.Background(), userArgs)
		if err != nil {
			fmt.Printf("error when creating new user: %s", err)
		}
	}
}
