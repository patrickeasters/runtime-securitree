package queue

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"

	"fmt"
)

type Queue struct {
	sqsSvc   *sqs.SQS
	QueueURL string
	Messages chan *sqs.Message
}

func NewQueue(sess *session.Session, url string) Queue {
	return Queue{
		sqsSvc:   sqs.New(sess),
		QueueURL: url,
		Messages: make(chan *sqs.Message),
	}
}

func (q Queue) Poll(handler func(string) error) {

	output, err := q.sqsSvc.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(q.QueueURL),
		MaxNumberOfMessages: aws.Int64(10),
		WaitTimeSeconds:     aws.Int64(15),
	})

	if err != nil {
		fmt.Printf("Failed to fetch sqs message %s", err)
	}

	for _, msg := range output.Messages {
		err := handler(*msg.Body)
		if err != nil {
			fmt.Printf("Failed to handle sqs message %s", err)
			continue
		}
		q.deleteMessage(msg)
	}

}

func (q Queue) deleteMessage(msg *sqs.Message) {
	q.sqsSvc.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      aws.String(q.QueueURL),
		ReceiptHandle: msg.ReceiptHandle,
	})
}
