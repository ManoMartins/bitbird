package events

import (
	"context"
	"errors"
	"github.com/manomartins/bitbird/internal/app/model"
	"github.com/manomartins/bitbird/internal/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test case when everything works as expected
func TestPullRequestCommentCreated_Execute_Success(t *testing.T) {
	// Arrange
	notifier := new(interfaces.MockNotifier)
	messagesStorageMock := new(interfaces.MockPullRequestMessagesInterface)

	event := PullRequestEvent{
		PullRequest: PullRequest{ID: 123},
		Comment:     Comment{Content: CommentContent{Raw: "This is a test comment"}},
	}
	prMessage := &model.PullRequestMessageModel{PrID: "123", MessageID: "message-id", ChannelID: "channel-id"}

	messagesStorageMock.On("GetById", "123").Return(prMessage, nil)
	notifier.On("SendCommentNotification", "channel-id", "message-id", "This is a test comment").Return(nil)

	sut := NewPullRequestCommentCreated(notifier, messagesStorageMock)

	// Act
	err := sut.Execute(context.Background(), event)

	// Assert
	assert.NoError(t, err)
	messagesStorageMock.AssertExpectations(t)
	notifier.AssertExpectations(t)
}

// Test case when there is an error retrieving the pull request message
func TestPullRequestCommentCreated_Execute_GetPullRequestMessageError(t *testing.T) {
	// Arrange
	notifier := new(interfaces.MockNotifier)
	messagesStorageMock := new(interfaces.MockPullRequestMessagesInterface)

	expectedError := errors.New("storage error")

	event := PullRequestEvent{PullRequest: PullRequest{ID: 123}}

	messagesStorageMock.On("GetById", "123").Return(nil, expectedError)

	sut := NewPullRequestCommentCreated(notifier, messagesStorageMock)

	// Act
	err := sut.Execute(context.Background(), event)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	messagesStorageMock.AssertExpectations(t)
	notifier.AssertNotCalled(t, "SendCommentNotification", mock.Anything, mock.Anything)
}

// Test case when there is an error sending the comment notification
func TestPullRequestCommentCreated_Execute_SendCommentNotificationError(t *testing.T) {
	// Arrange
	notifier := new(interfaces.MockNotifier)
	messagesStorageMock := new(interfaces.MockPullRequestMessagesInterface)

	expectedError := errors.New("notifier error")
	event := PullRequestEvent{
		PullRequest: PullRequest{ID: 123},
		Comment:     Comment{Content: CommentContent{Raw: "This is a test comment"}},
	}
	prMessage := &model.PullRequestMessageModel{PrID: "123", MessageID: "message-id", ChannelID: "channel-id"}

	messagesStorageMock.On("GetById", "123").Return(prMessage, nil)
	notifier.On("SendCommentNotification", "channel-id", "message-id", "This is a test comment").Return(expectedError)

	sut := NewPullRequestCommentCreated(notifier, messagesStorageMock)

	// Act
	err := sut.Execute(context.Background(), event)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	messagesStorageMock.AssertExpectations(t)
	notifier.AssertExpectations(t)
}
