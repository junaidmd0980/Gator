package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/junaidmd0980/gator/internal/database"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("Usage: cli <command> <args>")
	}
	ctx := context.Background()

	// name := s.cfgPtr.CurrentUserName
	url := cmd.Args[0]

	// user, err := s.db.GetUser(ctx, name)
	// if err != nil {
	// 	if err == sql.ErrNoRows {
	// 		fmt.Printf("Error: user '%s' is not registered\n", name)
	// 		os.Exit(1)
	// 	}
	// 	return fmt.Errorf("failed to fetch user: %w", err)
	// }

	feed, err := s.db.GetFeedByURL(ctx, url)
	if err != nil {
		return fmt.Errorf("failed to fetch feed: %w", err)
	}

	ff, err := s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to follow the feed: %w", err)
	}

	fmt.Printf("User %s is now following the feed: %s\n", ff.UserName, ff.FeedName)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	ctx := context.Background()
	// 	name := s.cfgPtr.CurrentUserName
	//
	// 	user, err := s.db.GetUser(ctx, name)
	// 	if err != nil {
	// 		if err == sql.ErrNoRows {
	// 			fmt.Printf("Error: user '%s' is not registered\n", name)
	// 			os.Exit(1)
	// 		}
	// 		return fmt.Errorf("failed to fetch user: %w", err)
	// 	}

	feedFollows, err := s.db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("failed to get following feeds: %w", err)
	}

	fmt.Printf("Feeds that %s are following:\n", user.Name)
	for _, ff := range feedFollows {
		fmt.Printf("* %s\n", ff.FeedName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("Usage: cli <command> <args>")
	}

	url := cmd.Args[0]
	ctx := context.Background()

	feed, err := s.db.GetFeedByURL(ctx, url)
	if err != nil {
		return fmt.Errorf("failed to fetch feed: %w", err)
	}

	err = s.db.DeleteFeedFollow(ctx, database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to unfollow feed: %w", err)
	}

	fmt.Printf("Successfully unfollowed feed: %s (%s)\n", feed.Name, url)
	return nil
}
