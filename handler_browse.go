package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/junaidmd0980/gator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.Args) == 1 {
		parsedLimit, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("invalid limit: %v", err)
		}
		limit = parsedLimit
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("failed getting posts for user: %w", err)
	}

	fmt.Printf("Found %v posts for user %s:\n", len(posts), user.Name)
	for _, post := range posts {
		fmt.Printf("------------------------------------\n")
		fmt.Printf("Title: %s\n", post.Title)
		fmt.Printf("URL: %s\n", post.Url)
		if post.PublishedAt.Valid {
			fmt.Printf("Published:   %s\n", post.PublishedAt.Time.Format("Jan 2, 2006 at 3:04pm"))
		}
		if post.Description.Valid {
			fmt.Printf("Desc: %s\n", post.Description.String)
		}
	}

	return nil
}
