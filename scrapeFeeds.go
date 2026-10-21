package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/junaidmd0980/gator/internal/database"
)

func scrapeFeeds(s *state) {
	ctx := context.Background()
	nextFeed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No feed found to fetch.")
			return
		}
		log.Printf("error getting next feed: %v\n", err)
		return
	}

	_, err = s.db.MarkFeedFetched(ctx, nextFeed.ID)
	if err != nil {
		log.Printf("error marking feed %s fetched", nextFeed.Name)
		return
	}

	rssFeed, err := fetchFeed(ctx, nextFeed.Url)
	if err != nil {
		log.Printf("error fetching feed: %v", err)
		return
	}

	fmt.Printf("--- Feeds: %s ---\n", rssFeed.Channel.Title)
	for _, item := range rssFeed.Channel.Item {
		pubAt := parsePubDate(item.PubDate)

		_, err = s.db.CreatePost(ctx, database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title:     item.Title,
			Url:       item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid:  item.Description != "",
			},
			PublishedAt: pubAt,
			FeedID:      nextFeed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "violates unique constraint") {
				continue
			}
			log.Printf("Error creating post %s: %v\n", item.Title, err)
		}
	}

	fmt.Printf("Feed %s collected, %v posts found.\n", nextFeed.Name, len(rssFeed.Channel.Item))
}

func parsePubDate(dateStr string) sql.NullTime {
	layouts := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC3339,
		"Mon, 02 Jan 2006 15:04:05 MST",
		"2006-01-02T15:04:05Z07:00",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, dateStr); err == nil {
			return sql.NullTime{Time: t, Valid: true}
		}
	}

	return sql.NullTime{Valid: false}
}
