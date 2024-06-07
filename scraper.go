package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/Hultyrhik/rssAggregator/internal/database"
)

func startScraping(
	db *database.Queries,
	concurrency int, // how many goroutines will do scraping
	timeBetweenRequest time.Duration, // cooldown between scraping
) {
	log.Printf("Scraping on %v goroutines every %s duration",
		concurrency, timeBetweenRequest)

	ticker := time.NewTicker(timeBetweenRequest)

	// execute body of the loop each time a new value
	// comes accross a ticker channel
	// if -  for range ticker.C - we will wait 1 minute upfront
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(
			context.Background(), // Global context
			int32(concurrency),
		)
		if err != nil {
			log.Println("error fetching feeds: ", err)
			continue // makes loop work forever
		}

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)

			go scrapeFeed(db, wg, feed)
		}
		wg.Wait()
	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("Error marking feed as fetched:", err)
		return
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("Error fetching feed", err)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		log.Println("Found post", item.Title, "on feed", feed.Name)
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))

}
