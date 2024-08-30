package events_test

import (
	"context"
	"github.com/andygrunwald/go-jira"
	"github.com/manomartins/bitbird/internal/events"
	interfaces "github.com/manomartins/bitbird/internal/mocks"
	"github.com/manomartins/bitbird/internal/work"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestCheckCD_Execute(t *testing.T) {
	mockNotifier := new(interfaces.MockNotifier)
	mockIssueService := new(interfaces.MockIssueService)
	mockDeploymentQueueStorage := new(interfaces.MockDeploymentQueueInterface)

	checkCD := events.NewCheckCD(mockNotifier, mockIssueService, mockDeploymentQueueStorage)

	ctx := context.Background()
	mockIssue := &jira.Issue{
		Key: "TEST-123",
		Fields: &jira.IssueFields{
			Summary: "This is a commit hash: a5b9f6a9a5b9f6a9a5b9f6a9a5b9f6a9a5b9f6a9",
			Assignee: &jira.User{
				DisplayName: "John Doe",
			},
		},
	}

	// Mock behaviors
	mockIssueService.On("GetFirstIssueByCodeBase", work.Frontend).Return(mockIssue)
	mockIssueService.On("GetFirstIssueByCodeBase", work.Backend).Return(nil)
	mockIssueService.On("GetFirstIssueByCodeBase", work.Mobile).Return(nil)

	mockDeploymentQueueStorage.On("GetByCardKey", "TEST-123").Return(nil, nil)
	mockNotifier.On("SendNotificationEmbed", ctx, mock.Anything, mock.Anything).Return("message-id", nil)
	mockDeploymentQueueStorage.On("Create", mock.Anything).Return(nil)

	// Execute the function
	err := checkCD.Execute(ctx)

	// Assert expectations
	assert.NoError(t, err)
	mockIssueService.AssertExpectations(t)
	mockDeploymentQueueStorage.AssertExpectations(t)
	mockNotifier.AssertExpectations(t)
}

func TestCheckCD_Execute_HashNotFound(t *testing.T) {
	mockNotifier := new(interfaces.MockNotifier)
	mockIssueService := new(interfaces.MockIssueService)
	mockDeploymentQueueStorage := new(interfaces.MockDeploymentQueueInterface)

	checkCD := events.NewCheckCD(mockNotifier, mockIssueService, mockDeploymentQueueStorage)

	ctx := context.Background()
	mockIssue := &jira.Issue{
		Key: "TEST-123",
		Fields: &jira.IssueFields{
			Summary: "This is an issue without a commit hash",
			Assignee: &jira.User{
				DisplayName: "John Doe",
			},
		},
	}

	// Mock behaviors
	mockIssueService.On("GetFirstIssueByCodeBase", work.Frontend).Return(mockIssue)
	mockIssueService.On("GetFirstIssueByCodeBase", work.Backend).Return(nil)
	mockIssueService.On("GetFirstIssueByCodeBase", work.Mobile).Return(nil)
	mockDeploymentQueueStorage.On("GetByCardKey", "TEST-123").Return(nil, nil)

	// Execute the function
	err := checkCD.Execute(ctx)

	// Assert expectations
	assert.Error(t, err)
	assert.EqualError(t, err, "hash not found in the provided string")
}
