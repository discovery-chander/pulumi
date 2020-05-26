package sqs

import (
	stderrors "errors"
	"fmt"
	"testing"

	"github.com/EurosportDigital/global-transcoding-platform/lib/aws/mocks"
	"github.com/EurosportDigital/global-transcoding-platform/lib/errors"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestReceiveMessagesWithValidMessages(t *testing.T) {
	mockSqsService := mocks.SQSAPI{}
	mockSqsService.On("ReceiveMessage", mock.AnythingOfType("*sqs.ReceiveMessageInput")).Return(&mockReceiveMessageOutput, nil)

	testClient := clientImpl{sqsService: &mockSqsService}
	result, err := testClient.ReceiveMessages()
	require.NoError(t, err)

	assert.Equal(t, expectedMessage, result[0].Message)
	assert.Equal(t, expectedReceiptHandle, result[0].ReceiptHandle)
	mockSqsService.AssertExpectations(t)
}

func TestReceiveMessagesWithNoMessages(t *testing.T) {
	mockEmptyReceiveMessageOutput := sqs.ReceiveMessageOutput{
		Messages: []*sqs.Message{},
	}
	mockSqsService := mocks.SQSAPI{}
	mockSqsService.On("ReceiveMessage", mock.AnythingOfType("*sqs.ReceiveMessageInput")).Return(&mockEmptyReceiveMessageOutput, nil)

	testClient := clientImpl{queueURL: mockQueueUrl, sqsService: &mockSqsService}
	result, err := testClient.ReceiveMessages()

	require.NoError(t, err)
	assert.Empty(t, result)
	mockSqsService.AssertExpectations(t)
}

func TestReceiveMessagesWithSqsError(t *testing.T) {
	expectedErrorMsg := "error receiving sqs messages: error returned in ReceiveMessage"
	mockSqsService := mocks.SQSAPI{}
	mockSqsService.On("ReceiveMessage", mock.AnythingOfType("*sqs.ReceiveMessageInput")).Return(nil, errors.New("error returned in ReceiveMessage"))
	testClient := clientImpl{sqsService: &mockSqsService}
	_, err := testClient.ReceiveMessages()
	assert.EqualError(t, err, expectedErrorMsg)
	mockSqsService.AssertExpectations(t)
}

func TestDeleteMessages(t *testing.T) {
	mockReceivedMessages := []Message{
		{
			ReceiptHandle: expectedReceiptHandle,
			Message:       expectedMessage,
		},
	}
	mockDeleteMessageInputs := getMockDeleteMessageInputs(mockQueueUrl, mockReceivedMessages)
	mockSqsService := mocks.SQSAPI{}
	mockSqsService.On("DeleteMessage", &mockDeleteMessageInputs[0]).Return(&sqs.DeleteMessageOutput{}, nil)
	testClient := clientImpl{queueURL: mockQueueUrl, sqsService: &mockSqsService}
	err := testClient.DeleteMessages(mockReceivedMessages)
	require.NoError(t, err)
	mockSqsService.AssertExpectations(t)
}

func TestDeleteMessagesEmptyMessages(t *testing.T) {
	mockReceivedMessages := []Message{}
	mockSqsService := mocks.SQSAPI{}
	testClient := clientImpl{sqsService: &mockSqsService}
	err := testClient.DeleteMessages(mockReceivedMessages)
	require.NoError(t, err)
	mockSqsService.AssertExpectations(t)
}

func TestDeleteMessagesNilMessages(t *testing.T) {
	mockSqsService := mocks.SQSAPI{}
	testClient := clientImpl{sqsService: &mockSqsService}
	err := testClient.DeleteMessages(nil)
	require.NoError(t, err)
	mockSqsService.AssertExpectations(t)
}

func TestDeleteMessagesWithOneSqsError(t *testing.T) {
	expectedErrorMsg := "error deleting 1 messages, captured last error: inner error message"
	mockReceivedMessages := []Message{
		{
			ReceiptHandle: "123",
			Message:       "I'm going to succeed on deletion",
		},
		{
			ReceiptHandle: "456",
			Message:       "I'm going to fail on deletion",
		},
		{
			ReceiptHandle: "789",
			Message:       "I'm also going to succeed on deletion",
		},
	}
	mockDeleteMessageInputs := getMockDeleteMessageInputs(mockQueueUrl, mockReceivedMessages)
	mockSqsService := mocks.SQSAPI{}
	mockSqsService.On("DeleteMessage", &mockDeleteMessageInputs[0]).Return(&sqs.DeleteMessageOutput{}, nil)
	mockSqsService.On("DeleteMessage", &mockDeleteMessageInputs[1]).Return(nil, errors.New("inner error message"))
	mockSqsService.On("DeleteMessage", &mockDeleteMessageInputs[2]).Return(&sqs.DeleteMessageOutput{}, nil)
	testClient := clientImpl{queueURL: mockQueueUrl, sqsService: &mockSqsService}

	err := testClient.DeleteMessages(mockReceivedMessages)
	assert.Equal(t, expectedErrorMsg, err.Error())
	mockSqsService.AssertExpectations(t)
}

func TestSendMessageReturnsMessageIDOnSuccess(t *testing.T) {
	expectedMessageID := "my message id"
	mockSqsService := &mocks.SQSAPI{}
	testClient := clientImpl{queueURL: mockQueueUrl, sqsService: mockSqsService}

	mockSqsService.On("SendMessage", &sqs.SendMessageInput{
		QueueUrl:    &mockQueueUrl,
		MessageBody: &expectedMessage,
	}).Return(
		&sqs.SendMessageOutput{
			MessageId: &expectedMessageID,
		},
		nil,
	)

	messageID, err := testClient.SendMessage(expectedMessage)
	require.NoError(t, err)
	require.EqualValues(t, expectedMessageID, messageID)
}

func TestSendMessageReturnsErrorFromService(t *testing.T) {
	expectedError := stderrors.New("my error")
	mockSqsService := &mocks.SQSAPI{}
	testClient := clientImpl{queueURL: mockQueueUrl, sqsService: mockSqsService}
	mockSqsService.On("SendMessage", &sqs.SendMessageInput{
		QueueUrl:    &mockQueueUrl,
		MessageBody: &expectedMessage,
	}).Return(nil, expectedError)

	_, err := testClient.SendMessage(expectedMessage)
	require.Error(t, err)
	require.EqualError(t, errors.Cause(err), expectedError.Error())
}

func TestSendMessageErrorsWhenReturnedNilID(t *testing.T) {
	mockSqsService := &mocks.SQSAPI{}
	testClient := clientImpl{queueURL: mockQueueUrl, sqsService: mockSqsService}
	mockSqsService.On("SendMessage", &sqs.SendMessageInput{
		QueueUrl:    &mockQueueUrl,
		MessageBody: &expectedMessage,
	}).Return(&sqs.SendMessageOutput{
		MessageId: nil,
	}, nil)

	_, err := testClient.SendMessage(expectedMessage)
	require.Error(t, err)
	require.EqualError(t, errors.Cause(err), fmt.Sprintf("queue %v did not return a message ID", mockQueueUrl))
}

func getMockDeleteMessageInputs(queueUrl string, messages []Message) []sqs.DeleteMessageInput {
	mockDeleteMessageInputs := make([]sqs.DeleteMessageInput, 0)
	for i := range messages {
		mockDeleteMessageInputs = append(mockDeleteMessageInputs, sqs.DeleteMessageInput{
			QueueUrl:      &queueUrl,
			ReceiptHandle: &messages[i].ReceiptHandle,
		})
	}
	return mockDeleteMessageInputs
}
