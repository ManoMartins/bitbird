package events

import (
	"errors"
	"github.com/manomartins/bitbird/internal/mocks"
	"github.com/manomartins/bitbird/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test case when everything works as expected
func TestPullRequestCommentCreated_Execute_Success(t *testing.T) {
	// Arrange
	notifier := new(mocks.Notifier)
	messagesStorageMock := new(mocks.PullRequestMessagesInterface)

	event := PullRequestEvent{
		PullRequest: PullRequest{ID: 123},
		Comment:     Comment{Content: CommentContent{Raw: "This is a test comment"}},
	}
	prMessage := &model.PullRequestMessageModel{PrID: "123", MessageID: "message-id"}

	messagesStorageMock.On("GetPullRequestMessage", "123").Return(prMessage, nil)
	notifier.On("SendCommentNotification", "message-id", "This is a test comment").Return(nil)

	sut := NewPullRequestCommentCreated(notifier, messagesStorageMock)

	// Act
	err := sut.Execute(event)

	// Assert
	assert.NoError(t, err)
	messagesStorageMock.AssertExpectations(t)
	notifier.AssertExpectations(t)
}

// Test case when there is an error retrieving the pull request message
func TestPullRequestCommentCreated_Execute_GetPullRequestMessageError(t *testing.T) {
	// Arrange
	notifier := new(mocks.Notifier)
	messagesStorage := new(mocks.PullRequestMessagesInterface)

	expectedError := errors.New("storage error")

	event := PullRequestEvent{PullRequest: PullRequest{ID: 123}}

	messagesStorage.On("GetPullRequestMessage", "123").Return(nil, expectedError)

	sut := NewPullRequestCommentCreated(notifier, messagesStorage)

	// Act
	err := sut.Execute(event)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	messagesStorage.AssertExpectations(t)
	notifier.AssertNotCalled(t, "SendCommentNotification", mock.Anything, mock.Anything)
}

// Test case when there is an error sending the comment notification
func TestPullRequestCommentCreated_Execute_SendCommentNotificationError(t *testing.T) {
	// Arrange
	notifier := new(mocks.Notifier)
	messagesStorage := new(mocks.PullRequestMessagesInterface)

	expectedError := errors.New("notifier error")
	event := PullRequestEvent{
		PullRequest: PullRequest{ID: 123},
		Comment:     Comment{Content: CommentContent{Raw: "This is a test comment"}},
	}
	prMessage := &model.PullRequestMessageModel{PrID: "123", MessageID: "message-id"}

	messagesStorage.On("GetPullRequestMessage", "123").Return(prMessage, nil)
	notifier.On("SendCommentNotification", "message-id", "This is a test comment").Return(expectedError)

	sut := NewPullRequestCommentCreated(notifier, messagesStorage)

	// Act
	err := sut.Execute(event)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	messagesStorage.AssertExpectations(t)
	notifier.AssertExpectations(t)
}
