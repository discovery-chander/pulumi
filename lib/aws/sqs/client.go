package sqs

import (
	"github.com/EurosportDigital/global-transcoding-platform/lib/errors"
	"github.com/EurosportDigital/global-transcoding-platform/lib/logger"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

const (
	visibilityTimeoutInSeconds = 120
	maxNumberOfMessages        = 10
	waitTimeSeconds            = 0
)

// clientImpl provides the ability to connect to an SQS queue.
type clientImpl struct {
	queueURL   string
	sqsService sqsiface.SQSAPI
}

// Client defines the functions that interact with an SQS queue.
type Client interface {
	// ReceiveMessages retrieves multiple messages from an SQS queue.
	ReceiveMessages() ([]Message, error)

	// DeleteMessages deletes the messages from the queue.
	DeleteMessages(messages []Message) error

	// SendMessage queues a message.
	SendMessage(string) (string, error)
}

// Stored as a variable so it can be overridden in tests.
var (
	newAwsSession = session.NewSession
	newSqsService = sqs.New
)

// NewClient returns a client that can send and receive messages from SQS.
func NewClient(awsRegion string, queueURL string) (Client, error) {
	sess, err := newAwsSession(&aws.Config{
		Region: aws.String(awsRegion),
	})
	if err != nil {
		return nil, errors.Wrap(err, "error initializing aws session")
	}
	sqsService := newSqsService(sess, &aws.Config{})

	logger.Debugf("Initializing SQS client with queue url: %v", queueURL)
	return &clientImpl{queueURL, sqsService}, nil
}

// ReceiveMessages returns up to 10 messages from an SQS queue.
func (client *clientImpl) ReceiveMessages() ([]Message, error) {
	result, err := client.sqsService.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            &(client.queueURL),
		MaxNumberOfMessages: aws.Int64(maxNumberOfMessages),
		VisibilityTimeout:   aws.Int64(visibilityTimeoutInSeconds),
		WaitTimeSeconds:     aws.Int64(waitTimeSeconds),
	})

	if err != nil {
		return make([]Message, 0), errors.Wrap(err, "error receiving sqs messages")
	}

	if len(result.Messages) == 0 {
		return make([]Message, 0), nil
	}

	return unmarshalMessages(result.Messages)
}

// DeleteMessages allows you to delete messages that have been processed from the queue.
func (client *clientImpl) DeleteMessages(messages []Message) error {
	var errorCount int = 0
	var lastError error = nil
	for i := range messages {
		_, err := client.sqsService.DeleteMessage(&sqs.DeleteMessageInput{
			QueueUrl:      &client.queueURL,
			ReceiptHandle: &messages[i].ReceiptHandle,
		})
		if err != nil {
			errorCount++
			lastError = errors.Wrapf(err, "error deleting %v messages, captured last error", errorCount)
		}
	}

	if lastError != nil {
		return lastError
	}
	return nil
}

// SendMessage queues a message.
func (client *clientImpl) SendMessage(message string) (string, error) {
	output, err := client.sqsService.SendMessage(&sqs.SendMessageInput{
		QueueUrl:    &client.queueURL,
		MessageBody: &message,
	})
	if err != nil {
		return "", errors.Wrapf(err, "unable to send message to %v", client.queueURL)
	}

	if output.MessageId == nil {
		return "", errors.Errorf("queue %v did not return a message ID", client.queueURL)
	}

	return *output.MessageId, nil
}

func unmarshalMessages(sqsMessages []*sqs.Message) ([]Message, error) {
	var messages []Message
	for _, msg := range sqsMessages {
		messages = append(messages, Message{*(msg.ReceiptHandle), *(msg.Body)})
	}
	return messages, nil
}
