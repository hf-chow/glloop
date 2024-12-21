package components


import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	db "github.com/hf-chow/glloop/internal/database"
)

func Login(q db.Queries) uuid.UUID {
	fmt.Println("Enter your username: ")
	var username string
	fmt.Scanln(&username)
	userID, err := q.GetIDByUsername(context.Background(), username)
	if err != nil {
		userArgs := db.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      username,
		}
		q.CreateUser(context.Background(), userArgs)
		userID = userArgs.ID
	}

	return userID
}
