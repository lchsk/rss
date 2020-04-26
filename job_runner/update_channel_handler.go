package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/lchsk/rss/comms"
	"github.com/mmcdole/gofeed"
)

func refreshChannelHandler(message *comms.Message, body []byte) {
	data := comms.RefreshChannel{}
	err := json.Unmarshal(body, &data)

	if err != nil {
		log.Printf("Error decoding refresh channel message: %s\n", err)
		return
	}

	var url string

	debug := os.Getenv("DEBUG")

	if debug == "false" {
		url = data.Url
	} else if debug == "true" {
		// Host some channel locally
		url = "https://lchsk.com/posts.xml"
	}

	if url == "" {
		log.Printf("Url not provided for refresh channel message, debug: %s\n", debug)
		return
	}

	log.Printf("Refreshing channel %s\n", url)

	fp := gofeed.NewParser()

	// TODO: Download manually with timeout and size limit
	feed, err := fp.ParseURL(url)

	if err != nil {
		log.Printf("Error getting channel data for %s: %s\n", url, err)
		return
	}

	// Call UpdateChannel(feed)
	fmt.Println(feed.Title)
}
