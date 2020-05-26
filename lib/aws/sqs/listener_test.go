package sqs

import (
	"testing"
	"time"

	"github.com/EurosportDigital/global-transcoding-platform/lib/aws/mocks"
	"github.com/aws/aws-sdk-go/aws"
	awsclient "github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var expectedMessage = `{"field1":"value1","field2":"value2","field3":"value3"}`
var expectedReceiptHandle = "abcdef123456"
var mockQueueUrl = "aQueueUrl"
var mockReceiveMessageOutput = sqs.ReceiveMessageOutput{
	Messages: []*sqs.Message{
		{
			Body:          aws.String(expectedMessage),
			ReceiptHandle: aws.String(expectedReceiptHandle),
		},
	},
}
var mockDeleteMessageInput = sqs.DeleteMessageInput{
	QueueUrl:      &mockQueueUrl,
	ReceiptHandle: &expectedReceiptHandle,
}
var maxTimeout = 3 * time.Second

func TestNewListener(t *testing.T) {
	previousNewAwsSession := newAwsSession
	previousNewSqsService := newSqsService
	defer func() {
		newAwsSession = previousNewAwsSession
		newSqsService = previousNewSqsService
	}()
	newAwsSession = func(...*aws.Config) (*session.Session, error) {
		return nil, nil
	}
	newSqsService = func(awsclient.ConfigProvider, ...*aws.Config) *sqs.SQS {
		return nil
	}
	listener := NewListener(nil)
	require.NotNil(t, listener)
}

func TestPoll(t *testing.T) {
	mockSqsService := mocks.SQSAPI{}
	mockSqsService.On("ReceiveMessage", mock.AnythingOfType("*sqs.ReceiveMessageInput")).Return(&mockReceiveMessageOutput, nil)
	mockSqsService.On("DeleteMessage", &mockDeleteMessageInput).Return(&sqs.DeleteMessageOutput{}, nil)
	testClient := listenerClient{sqsClient: &clientImpl{queueURL: mockQueueUrl, sqsService: &mockSqsService}, closePolling: make(chan bool)}
	channel := make(chan Message)
	done := make(chan bool)

	go func() {
		testClient.Poll(channel, true)
		done <- true
	}()
	go func() {
		for {
			<-channel
		}
	}()

	select {
	case msg := <-channel:
		assert.Equal(t, expectedMessage, msg.Message)
		assert.Equal(t, expectedReceiptHandle, msg.ReceiptHandle)
	case <-time.After(maxTimeout):
		assert.Fail(t, "Test timed out without receiving message from channel")
	}

	testClient.closePolling <- true
	<-done
	mockSqsService.AssertExpectations(t)
}

func TestClose(t *testing.T) {
	mockSqsService := mocks.SQSAPI{}
	testClient := listenerClient{sqsClient: &clientImpl{queueURL: mockQueueUrl, sqsService: &mockSqsService}, closePolling: make(chan bool)}

	// In goroutine so that sending to closePolling doesn't block before there's a listener
	go func() {
		testClient.Close()
	}()

	select {
	case <-testClient.closePolling:
	// Succeed test and exit before timeout is reached.
	case <-time.After(maxTimeout):
		assert.Fail(t, "Test timed out without receiving message from channel")
	}
	mockSqsService.AssertExpectations(t)
}
