package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/junaidmd0980/gator/internal/database"
)

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("Usage: cli <command> <args>")
	}

	name := cmd.Args[0]
	url := cmd.Args[1]

	currentUsername := s.cfgPtr.CurrentUserName

	user, err := s.db.GetUser(context.Background(), currentUsername)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user must be logged in: %w", err)
		}
		return fmt.Errorf("database error: %w", err)
	}

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to create feed: %w", err)
	}

	fmt.Printf("Feed created successfully:\n%+v\n", feed)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.ListFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("failed fetch feeds: %w", err)
	}

	if len(feeds) == 0 {
		return fmt.Errorf("feeds list is empty")
	}

	for _, feed := range feeds {
		fmt.Println("Name of the feed:", feed.Name)
		fmt.Println("Url of the feed:", feed.Url)
		fmt.Println("Publisher of the feed:", feed.Username)
	}

	return nil
}
