package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/service/sqs"
	"gopkg.in/mgo.v2/bson"
)

func processEvents(msg *sqs.Message) {
	fmt.Println(*msg.Body)
	var m event
	if err := json.Unmarshal([]byte(*msg.Body), &m); err != nil {
		log.Fatal(err)
	}

	if m.EventType == "video_click" {
		fmt.Printf("message: %+v", m)
		fmt.Println("video clicked")
		var ch ChannelModel
		fmt.Println("id:", m.ChannelID)
		if err := dbc.Find(bson.M{"_id": m.ChannelID}).One(&ch); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("channel: %+v", ch)
	}
}
