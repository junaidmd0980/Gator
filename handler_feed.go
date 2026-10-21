package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/junaidmd0980/gator/internal/database"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("Usage: cli <command> <args>")
	}

	name := cmd.Args[0]
	url := cmd.Args[1]
	ctx := context.Background()

	// 	currentUsername := s.cfgPtr.CurrentUserName
	//
	// 	user, err := s.db.GetUser(ctx, currentUsername)
	// 	if err != nil {
	// 		if err == sql.ErrNoRows {
	// 			return fmt.Errorf("user must be logged in: %w", err)
	// 		}
	// 		return fmt.Errorf("failed to fetch user: %w", err)
	// 	}

	feed, err := s.db.CreateFeed(ctx, database.CreateFeedParams{
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

	_, err = s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to follow the feed: %w", err)
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
