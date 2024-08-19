package events

import (
	"errors"
	"github.com/manomartins/bitbird/internal/mocks"
	"github.com/manomartins/bitbird/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPullRequestChangesRequest_Execute_Success(t *testing.T) {
	// Arrange
	notifier := new(mocks.Notifier)
	messagesStorage := new(mocks.PullRequestMessagesInterface)

	event := PullRequestEvent{PullRequest: PullRequest{ID: 123}}
	prMessage := &model.PullRequestMessageModel{PrID: "123", MessageID: "message-id"}

	messagesStorage.On("GetById", "123").Return(prMessage, nil)
	notifier.On("AddChangeRequestEmoji", "message-id").Return(nil)

	sut := NewPullRequestChangesRequest(notifier, messagesStorage)

	// Act
	err := sut.Execute(event)

	// Assert
	assert.NoError(t, err)
	messagesStorage.AssertExpectations(t)
	notifier.AssertExpectations(t)
}

func TestPullRequestChangesRequest_Execute_GetPullRequestMessageError(t *testing.T) {
	// Arrange
	notifier := new(mocks.Notifier)
	messagesStorage := new(mocks.PullRequestMessagesInterface)

	expectedError := errors.New("storage error")

	event := PullRequestEvent{PullRequest: PullRequest{ID: 123}}

	messagesStorage.On("GetById", "123").Return(nil, expectedError)

	sut := NewPullRequestChangesRequest(notifier, messagesStorage)

	// Act
	err := sut.Execute(event)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	messagesStorage.AssertExpectations(t)
	notifier.AssertNotCalled(t, "AddChangeRequestEmoji", mock.Anything)
}

func TestPullRequestChangesRequest_Execute_AddChangeRequestEmojiError(t *testing.T) {
	// Arrange
	notifier := new(mocks.Notifier)
	messagesStorage := new(mocks.PullRequestMessagesInterface)

	expectedError := errors.New("notifier error")

	event := PullRequestEvent{PullRequest: PullRequest{ID: 123}}
	prMessage := &model.PullRequestMessageModel{PrID: "123", MessageID: "message-id"}

	messagesStorage.On("GetById", "123").Return(prMessage, nil)
	notifier.On("AddChangeRequestEmoji", "message-id").Return(expectedError)

	sut := NewPullRequestChangesRequest(notifier, messagesStorage)

	// Act
	err := sut.Execute(event)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	messagesStorage.AssertExpectations(t)
	notifier.AssertExpectations(t)
}
