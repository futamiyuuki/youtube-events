package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/service/sqs"
)

type event struct {
	EventType     string `json:"event_type"`
	VideoID       string `json:"video_id"`
	VideoCategory string `json:"video_category"`
	ChannelID     string `json:"channel_id"`
	IsSubscribed  bool   `json:"is_subscirbed"`
	SearchTerm    string `json:"search_term"`
}

type categoryCount struct {
	Category string `json:"category"`
	Count    int    `json:"count"`
}

// Channel represents Youtube channel category views
type channel struct {
	ChannelID      string          `json:"_id"`
	CategoryCounts []categoryCount `json:"categories"`
}

func processEvents(msg *sqs.Message) {
	fmt.Println(*msg.Body)
	var m event
	if err := json.Unmarshal([]byte(*msg.Body), &m); err != nil {
		log.Fatal(err)
	}

	if m.EventType == "video_click" {
		fmt.Println("video clicked")
	}
}
