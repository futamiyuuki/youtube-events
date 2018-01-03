package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/joho/godotenv"
	statsd "gopkg.in/alexcesaro/statsd.v2"
	mgo "gopkg.in/mgo.v2"
)

const (
	maxJobCount    = 12000
	maxWorkerCount = 100
)

var dbc *mgo.Collection
var dbe *mgo.Collection

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	JobQueue = make(chan Job, maxJobCount)
}

func main() {
	stop := make(chan bool)

	sc, err := statsd.New()
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	sess, err := session.NewSession()
	if err != nil {
		log.Fatal("Error loading AWS session")
	}

	svc := sqs.New(sess)

	mURI := os.Getenv("MONGODB_URI")
	mSess, err := mgo.Dial(mURI)
	if err != nil {
		panic(err)
	}
	fmt.Println("Using DB:", mURI)
	dbc = mSess.DB("events").C("channels")
	dbe = mSess.DB("events").C("events")
	defer mSess.Close()

	fmt.Printf("Create %d workers\n", maxWorkerCount)
	wPool := make(chan chan Job, maxWorkerCount)
	for i := 0; i < maxWorkerCount; i++ {
		w := NewWorker(wPool)
		w.Start()
	}

	// dispatch
	go func() {
		for j := range JobQueue {
			go func(j Job) {
				jChan := <-wPool
				jChan <- j
			}(j)
		}
	}()

	fmt.Println("Start polling...")
	for i := 0; i < 38; i++ {
		poll(svc, sc)
	}

	<-stop
}
