package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/junaidmd0980/gator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("you must provide username")
	}

	name := cmd.Args[0]

	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("Error: user '%s' does not exist in the database\n", name)
			os.Exit(1)
		}
		return fmt.Errorf("database error: %w", err)
	}

	err = s.cfgPtr.SetUser(name)
	if err != nil {
		return err
	}
	fmt.Printf("User '%s' has been set..\n", cmd.Args[0])
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("you must provide username")
	}

	name := cmd.Args[0]

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	})
	if err != nil {
		return fmt.Errorf("failed to create user %s: %w", name, err)
	}

	err = s.cfgPtr.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("failed to set current user: %w", err)
	}

	fmt.Printf("User was created successfully:\n%+v\n", user)

	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to reset users data: %w", err)
	}

	fmt.Println("Data reset successfully!")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to fetch users data: %w", err)
	}

	if len(users) == 0 {
		fmt.Println("Users list is empty")
		os.Exit(1)
	}

	for _, user := range users {
		if user.Name == s.cfgPtr.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Println("*", user.Name)
		}
	}

	return nil
}
