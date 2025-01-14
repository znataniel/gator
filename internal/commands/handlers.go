package commands

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/znataniel/gator/internal/database"
	"github.com/znataniel/gator/internal/rss"
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

func HandlerReset(s *State, cmd Command) error {
	err := s.Db.DeleteAllUsers(context.Background())
	return err
}

func HandlerUsers(s *State, cmd Command) error {
	users, err := s.Db.GetUsers(context.Background())
	if err == sql.ErrNoRows {
		return fmt.Errorf("error: no users found")
	}

	for _, u := range users {
		if s.Cfg.CurrentUser == u.Name.String {
			fmt.Println("\t*", u.Name.String, "(current)")
			continue
		}
		fmt.Println("\t*", u.Name.String)
	}
	return nil
}

func HandlerAgg(s *State, cmd Command) error {
	feed, err := rss.FetchFeed(context.Background(), "url here")
	if err != nil {
		return err
	}

	fmt.Println(feed)
	return nil
}

func HandlerAddFeed(s *State, cmd Command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("error: wrong number of arguments provided")
	}

	ctx := context.Background()

	currentUserRow, err := s.Db.GetUserByName(ctx, sql.NullString{
		String: s.Cfg.CurrentUser,
		Valid:  true,
	})
	if err != nil {
		return fmt.Errorf("error: could not retrieve current user data")
	}

	feed, err := s.Db.CreateFeed(ctx, database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    currentUserRow.ID,
	})
	if err != nil {
		return fmt.Errorf("error: could not create feed")
	}

	fmt.Println("New feed created!")
	fmt.Println(feed)
	return nil
}

func HandlerFeeds(s *State, cmd Command) error {
	feedData, err := s.Db.GetFeedsToPrint(context.Background())
	if err != nil {
		return err
	}

	for _, f := range feedData {
		fmt.Println("*", f.Name)
		fmt.Println("\turl:", f.Url)
		fmt.Println("\tadded by:", f.Name_2.String)
	}
	return nil
}
