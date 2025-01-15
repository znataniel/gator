package commands

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/znataniel/gator/internal/database"
)

func MiddlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(s *State, cmd Command) error {
	return func(s *State, cmd Command) error {
		currentUser, err := s.Db.GetUserByName(context.Background(), sql.NullString{
			String: s.Cfg.CurrentUser,
			Valid:  true,
		})
		if err != nil {
			return fmt.Errorf("could not retrieve current user data: %s", err)
		}

		return handler(s, cmd, currentUser)
	}
}
