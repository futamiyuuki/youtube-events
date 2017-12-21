package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/joho/godotenv"
)

// var pm = make(chan string)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	sess, err := session.NewSession()
	if err != nil {
		log.Fatal("Error loading AWS session")
	}

	svc := sqs.New(sess)

	fmt.Println("Start polling...")
	poll(svc)
}
