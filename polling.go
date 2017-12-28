package main

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

const (
	maxNumberOfMessages int64 = 10
	waitTimeSeconds     int64 = 20
)

func poll(svc *sqs.SQS) {
	qURL := os.Getenv("AWS_QUEUE_URL")

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
				go func(msg *sqs.Message) {
					// st := time.Now()
					processEvents(msg)
					dmreqs = append(dmreqs, &sqs.DeleteMessageBatchRequestEntry{
						Id:            msg.MessageId,
						ReceiptHandle: msg.ReceiptHandle,
					})
					mcnt++
					if mcnt == mlen {
						// fmt.Println("Batch deleting...")
						svc.DeleteMessageBatch(&sqs.DeleteMessageBatchInput{
							QueueUrl: &qURL,
							Entries:  dmreqs,
						})
						// fmt.Printf("\nFinished deleting message %s\n", time.Now().Sub(st).String())
					}
				}(msg)
			}
		}
	}
}
