package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go/service/sqs"
	"gopkg.in/mgo.v2/bson"
)

func processEvents(msg *sqs.Message) {
	// fmt.Println(*msg.Body)
	var m event
	if err := json.Unmarshal([]byte(*msg.Body), &m); err != nil {
		log.Fatal(err)
	}

	if m.EventType == "video_click" {
		// fmt.Printf("video clicked, message: %+v\n", m)
		var ch channelModel
		if err := dbc.FindId(m.ChannelID).One(&ch); err != nil {
			// fmt.Println("Channel not found, creating new channel")
			cms := []categoryModel{}
			cms = append(cms, categoryModel{
				m.VideoCategory,
				1,
			})
			ch.ID = m.ChannelID
			ch.Categories = cms
			if err := dbc.Insert(&ch); err != nil {
				log.Fatal(err)
			}
		} else {
			// fmt.Println("Channel found, updating channel", ch.ID)
			for i, cat := range ch.Categories {
				if cat.Category == m.VideoCategory {
					ch.Categories[i].Count++
				}
			}
			dbc.UpdateId(ch.ID, bson.M{"$set": ch})
		}
		// fmt.Printf("channel: %+v\n", ch)
	}
}
