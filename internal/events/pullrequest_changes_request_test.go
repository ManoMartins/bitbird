package events

import (
	"context"
	"errors"
	interfaces "github.com/manomartins/bitbird/internal/mocks"
	"github.com/manomartins/bitbird/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPullRequestChangesRequest_Execute_Success(t *testing.T) {
	// Arrange
	notifier := new(interfaces.MockNotifier)
	messagesStorage := new(interfaces.MockPullRequestMessagesInterface)

	event := PullRequestEvent{PullRequest: PullRequest{ID: 123}}
	prMessage := &model.PullRequestMessageModel{PrID: "123", MessageID: "message-id"}

	messagesStorage.On("GetById", "123").Return(prMessage, nil)
	notifier.On("AddChangeRequestEmoji", mock.AnythingOfType("string"), "message-id").Return(nil)

	sut := NewPullRequestChangesRequest(notifier, messagesStorage)

	// Act
	err := sut.Execute(context.Background(), event)

	// Assert
	assert.NoError(t, err)
	messagesStorage.AssertExpectations(t)
	notifier.AssertExpectations(t)
}

func TestPullRequestChangesRequest_Execute_GetPullRequestMessageError(t *testing.T) {
	// Arrange
	notifier := new(interfaces.MockNotifier)
	messagesStorage := new(interfaces.MockPullRequestMessagesInterface)

	expectedError := errors.New("storage error")

	event := PullRequestEvent{PullRequest: PullRequest{ID: 123}}

	messagesStorage.On("GetById", "123").Return(nil, expectedError)

	sut := NewPullRequestChangesRequest(notifier, messagesStorage)

	// Act
	err := sut.Execute(context.Background(), event)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	messagesStorage.AssertExpectations(t)
	notifier.AssertNotCalled(t, "AddChangeRequestEmoji", mock.Anything)
}

func TestPullRequestChangesRequest_Execute_AddChangeRequestEmojiError(t *testing.T) {
	// Arrange
	notifier := new(interfaces.MockNotifier)
	messagesStorage := new(interfaces.MockPullRequestMessagesInterface)

	expectedError := errors.New("notifier error")

	event := PullRequestEvent{PullRequest: PullRequest{ID: 123}}
	prMessage := &model.PullRequestMessageModel{PrID: "123", MessageID: "message-id"}

	messagesStorage.On("GetById", "123").Return(prMessage, nil)
	notifier.On("AddChangeRequestEmoji", mock.AnythingOfType("string"), "message-id").Return(expectedError)

	sut := NewPullRequestChangesRequest(notifier, messagesStorage)

	// Act
	err := sut.Execute(context.Background(), event)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	messagesStorage.AssertExpectations(t)
	notifier.AssertExpectations(t)
}
