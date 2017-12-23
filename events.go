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
		fmt.Printf("video clicked, message: %+v\n", m)
		var ch ChannelModel
		fmt.Println("id:", m.ChannelID)
		if err := dbc.FindId(m.ChannelID).One(&ch); err != nil {
			fmt.Println("Channel not found, creating new channel")
			cms := []CategoryModel{}
			cms = append(cms, CategoryModel{
				m.VideoCategory,
				1,
			})
			ch.ID = m.ChannelID
			ch.Categories = cms
			if err := dbc.Insert(&ch); err != nil {
				log.Fatal(err)
			}
		} else {
			fmt.Println("Channel found, updating channel", ch.ID)
			for i, cat := range ch.Categories {
				if cat.Category == m.VideoCategory {
					fmt.Printf("Increment Category: %s %d\n", cat.Category, cat.Count)
					ch.Categories[i].Count++
				}
			}
			fmt.Printf("%+v\n", ch)
			dbc.UpdateId(ch.ID, bson.M{"$set": ch})
		}
		fmt.Printf("channel: %+v\n", ch)
	}
}
