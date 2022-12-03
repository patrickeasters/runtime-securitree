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

func (q Queue) Poll(handler func(string) error) error {

	output, err := q.sqsSvc.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(q.QueueURL),
		MaxNumberOfMessages: aws.Int64(10),
		WaitTimeSeconds:     aws.Int64(15),
	})

	if err != nil {
		return fmt.Errorf("failed to fetch sqs message: %w", err)
	}

	for _, msg := range output.Messages {
		err := handler(*msg.Body)
		if err != nil {
			return fmt.Errorf("failed to handle sqs message: %w", err)
		}
		err = q.deleteMessage(msg)
		if err != nil {
			return fmt.Errorf("failed to delete sqs message: %w", err)
		}
	}
	return nil
}

func (q Queue) deleteMessage(msg *sqs.Message) error {
	_, err := q.sqsSvc.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      aws.String(q.QueueURL),
		ReceiptHandle: msg.ReceiptHandle,
	})
	return err
}
