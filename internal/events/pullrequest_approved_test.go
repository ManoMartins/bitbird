package events

import (
	"errors"
	"github.com/manomartins/bitbird/internal/mocks"
	"github.com/manomartins/bitbird/internal/model"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPullRequestApproved_Execute_Success(t *testing.T) {
	// Create the mocks
	notifier := new(mocks.Notifier)
	messagesStorage := new(mocks.PullRequestMessagesInterface)

	// Set up the expected interactions
	expectedPrID := "1"
	expectedMessageID := "1"
	expectedPrMessage := &model.PullRequestMessageModel{PrID: "1", MessageID: "1"}

	messagesStorage.On("GetPullRequestMessage", expectedPrID).Return(expectedPrMessage, nil)
	notifier.On("AddApprovalEmoji", expectedMessageID).Return(nil)

	// Create the PullRequestApproved instance
	prApproved := NewPullRequestApproved(notifier, messagesStorage)

	// Define the event to test
	event := PullRequestEvent{
		PullRequest: PullRequest{
			ID: 1,
		},
	}

	// Call the method
	err := prApproved.Execute(event)

	// Assert that there were no errors
	assert.NoError(t, err)

	// Assert that the expected methods were called
	messagesStorage.AssertCalled(t, "GetPullRequestMessage", expectedPrID)
	notifier.AssertCalled(t, "AddApprovalEmoji", expectedMessageID)
}

func TestPullRequestApproved_Execute_GetMessageError(t *testing.T) {
	// Create the mocks
	notifier := new(mocks.Notifier)
	messagesStorage := new(mocks.PullRequestMessagesInterface)

	// Set up the expected interactions
	expectedPrID := "1"
	expectedError := errors.New("error getting PR message")

	messagesStorage.On("GetPullRequestMessage", expectedPrID).Return(nil, expectedError)

	// Create the PullRequestApproved instance
	prApproved := NewPullRequestApproved(notifier, messagesStorage)

	// Define the event to test
	event := PullRequestEvent{
		PullRequest: PullRequest{
			ID: 1,
		},
	}

	// Call the method
	err := prApproved.Execute(event)

	// Assert that the error is the one expected
	assert.Equal(t, expectedError, err)

	// Assert that AddApprovalEmoji was not called
	notifier.AssertNotCalled(t, "AddApprovalEmoji", mock.Anything)
}

func TestPullRequestApproved_Execute_AddEmojiError(t *testing.T) {
	// Create the mocks
	notifier := new(mocks.Notifier)
	messagesStorage := new(mocks.PullRequestMessagesInterface)

	// Set up the expected interactions
	expectedPrID := "1"
	expectedMessageID := "1"
	expectedPrMessage := &model.PullRequestMessageModel{PrID: "1", MessageID: "1"}
	expectedError := errors.New("error adding emoji")

	messagesStorage.On("GetPullRequestMessage", expectedPrID).Return(expectedPrMessage, nil)
	notifier.On("AddApprovalEmoji", expectedMessageID).Return(expectedError)

	// Create the PullRequestApproved instance
	prApproved := NewPullRequestApproved(notifier, messagesStorage)

	// Define the event to test
	event := PullRequestEvent{
		PullRequest: PullRequest{
			ID: 1,
		},
	}

	// Call the method
	err := prApproved.Execute(event)

	// Assert that the error is the one expected
	assert.Equal(t, expectedError, err)

	// Assert that the expected methods were called
	messagesStorage.AssertCalled(t, "GetPullRequestMessage", expectedPrID)
	notifier.AssertCalled(t, "AddApprovalEmoji", expectedMessageID)
}
