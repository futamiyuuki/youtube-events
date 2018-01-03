package main

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	statsd "gopkg.in/alexcesaro/statsd.v2"
)

const (
	maxNumberOfMessages int64 = 10
	waitTimeSeconds     int64 = 20
)

func poll(svc *sqs.SQS, sc *statsd.Client) {
	qURL := os.Getenv("AWS_QUEUE_URL")

	go func() {
		for {
			resp, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
				// AttributeNames: []*string{
				// 		aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
				// },
				MessageAttributeNames: []*string{
					aws.String(sqs.QueueAttributeNameAll),
				},
				QueueUrl:            &qURL,
				MaxNumberOfMessages: aws.Int64(maxNumberOfMessages),
				WaitTimeSeconds:     aws.Int64(waitTimeSeconds),
			})
			if err != nil {
				log.Println("Error receiving message from SQS")
				continue
			}
			mlen := len(resp.Messages)
			if mlen > 0 {
				mcnt := 0
				dmreqs := []*sqs.DeleteMessageBatchRequestEntry{}
				for _, msg := range resp.Messages {
					t := sc.NewTiming()
					// processEvents(msg)
					JobQueue <- Job{msg}
					t.Send("events.response_time")
					dmreqs = append(dmreqs, &sqs.DeleteMessageBatchRequestEntry{
						Id:            msg.MessageId,
						ReceiptHandle: msg.ReceiptHandle,
					})
					mcnt++
					if mcnt == mlen {
						svc.DeleteMessageBatch(&sqs.DeleteMessageBatchInput{
							QueueUrl: &qURL,
							Entries:  dmreqs,
						})
					}
				}
			}
		}
	}()
}
