package commands

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/znataniel/gator/internal/database"
)

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("error: wrong number of arguments provided")
	}

	if _, err := s.Db.GetUserByName(context.Background(),
		sql.NullString{
			String: cmd.Args[0],
			Valid:  true,
		}); err == sql.ErrNoRows {
		return fmt.Errorf("error: user is not registered, register and try again")
	}

	err := s.Cfg.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}
	println("user", cmd.Args[0], "has logged in")
	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("error: wrong number of arguments provided")
	}

	if _, err := s.Db.GetUserByName(context.Background(),
		sql.NullString{
			String: cmd.Args[0],
			Valid:  true,
		}); err != sql.ErrNoRows {
		return fmt.Errorf("error: user already exists")
	}

	createUserParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      sql.NullString{String: cmd.Args[0], Valid: true},
	}
	s.Db.CreateUser(context.Background(), createUserParams)

	err := s.Cfg.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}

	fmt.Println("User was created")
	fmt.Println("id:", createUserParams.ID)
	fmt.Println("created_at:", createUserParams.CreatedAt)
	fmt.Println("updated_at:", createUserParams.UpdatedAt)
	fmt.Println("name:", cmd.Args[0])
	return nil
}
