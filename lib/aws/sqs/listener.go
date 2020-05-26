package sqs

import (
	"github.com/EurosportDigital/global-transcoding-platform/lib/logger"
)

type listenerClient struct {
	sqsClient    Client
	closePolling chan bool
}

// Listener is an interface for ReceiveMessages and DeleteMessages.
type Listener interface {
	Poll(chan<- Message, bool) error
	Close() error
}

// Message is a data model for passing the body of an SQS message and its receipt handle.
type Message struct {
	ReceiptHandle string
	Message       string
}

// NewListener is a function for initializing a new client of interface Listener
func NewListener(sqsClient Client) Listener {
	return &listenerClient{sqsClient, make(chan bool)}
}

// Poll infinitely checks for messages and sends them to a channel. It should be used as a goroutine to listen to a Sqs queue.
func (lc *listenerClient) Poll(channel chan<- Message, autocleanup bool) error {
	for {
		select {
		case <-lc.closePolling:
			close(channel)
			return nil
		default:
		}

		messages, err := lc.sqsClient.ReceiveMessages()
		if err != nil {
			logger.Error(err, "Error on reading SQS messages")
			continue
		}
		for i := range messages {
			channel <- messages[i]
		}
		if autocleanup {
			err := lc.sqsClient.DeleteMessages(messages)
			if err != nil {
				logger.Error(err, "Error on deleting SQS messages")
			}
		}
	}
}

// Close gives the user the ability to stop polling after calling Poll
func (lc *listenerClient) Close() error {
	lc.closePolling <- true
	return nil
}
