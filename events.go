package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go/service/sqs"
	"gopkg.in/mgo.v2/bson"
)

type Job struct {
	Msg *sqs.Message
}

var JobQueue chan Job

type Worker struct {
	WorkerPool chan chan Job
	JobChannel chan Job
	quit       chan bool
}

func NewWorker(wPool chan chan Job) Worker {
	return Worker{
		WorkerPool: wPool,
		JobChannel: make(chan Job),
		quit:       make(chan bool),
	}
}

func (w Worker) Start() {
	go func() {
		for {
			w.WorkerPool <- w.JobChannel
			select {
			case j := <-w.JobChannel:
				j.ProcessEvents()
			case <-w.quit:
				return
			}
		}
	}()
}

func (w Worker) Stop() {
	w.quit <- true
}

func (j Job) ProcessEvents() {
	// st := time.Now()
	var m event
	if err := json.Unmarshal([]byte(*j.Msg.Body), &m); err != nil {
		log.Fatal(err)
	}

	go func(m event) {
		e := eventModel{m.EventType, m.VideoID, m.VideoCategory, m.ChannelID, m.IsSubscribed, m.SearchTerm}
		if err := dbe.Insert(&e); err != nil {
			log.Fatal(err)
		}
		// fmt.Printf("\nFinished inserting event in %s\n", time.Now().Sub(st).String())
	}(m)

	if m.EventType == "video_click" {
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
			for i, cat := range ch.Categories {
				if cat.Category == m.VideoCategory {
					ch.Categories[i].Count++
				}
			}
			dbc.UpdateId(ch.ID, bson.M{"$set": ch})
		}
		// fmt.Printf("channel: %+v\n", ch)
		// fmt.Printf("\nFinished processing click event in %s\n", time.Now().Sub(st).String())
	}
}
