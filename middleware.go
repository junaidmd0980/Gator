package main

import (
	"context"
	"fmt"

	"github.com/junaidmd0980/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfgPtr.CurrentUserName)
		if err != nil {
			return fmt.Errorf("User must be logged in: %w", err)
		}

		return handler(s, cmd, user)
	}
}
