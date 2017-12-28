package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/joho/godotenv"
	newrelic "github.com/newrelic/go-agent"
	mgo "gopkg.in/mgo.v2"
)

var dbc *mgo.Collection
var dbe *mgo.Collection

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	nrConfig := newrelic.NewConfig("youtube-events", os.Getenv("NR_KEY"))
	nr, err := newrelic.NewApplication(nrConfig)
	if err != nil {
		log.Fatal("Error initializing New Relic APM")
	}

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

	fmt.Println("Start polling...")
	poll(svc, nr)
}
