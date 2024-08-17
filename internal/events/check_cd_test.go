package events

import (
	"errors"
	"github.com/manomartins/bitbird/internal/mocks"
	"testing"

	"github.com/andygrunwald/go-jira"
	"github.com/manomartins/bitbird/internal/work"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test case when everything works as expected
func TestCheckCD_Execute_Success(t *testing.T) {
	// Arrange
	notifier := new(mocks.Notifier)
	deploymentQueue := new(mocks.DeploymentQueueInterface)

	issueService := new(mocks.IssueService)
	jiraIssue := &jira.Issue{
		Key: "PROJ-123",
		Fields: &jira.IssueFields{
			Summary:  "Deployment 6fc1d3e8e0d91c23df7f2a193f68f8616f7461de",
			Assignee: &jira.User{DisplayName: "John Doe"},
		},
	}

	checkCD := NewCheckCD(notifier, issueService, deploymentQueue)

	// Mock methods
	issueService.On("GetFirstIssueByCodeBase", work.Mobile).Return(jiraIssue)
	issueService.On("GetFirstIssueByCodeBase", work.Backend).Return(nil)
	issueService.On("GetFirstIssueByCodeBase", work.Frontend).Return(nil)
	deploymentQueue.On("GetByCardKey", "PROJ-123").Return(nil, nil)
	notifier.On("SendNotification", mock.Anything).Return("message-id", nil)
	deploymentQueue.On("Create", mock.Anything).Return(nil)

	// Act
	err := checkCD.Execute()

	// Assert
	assert.NoError(t, err)
	notifier.AssertExpectations(t)
	deploymentQueue.AssertExpectations(t)
	issueService.AssertExpectations(t)
}

func TestCheckCD_Execute_NoIssuesFound(t *testing.T) {
	// Arrange
	notifier := new(mocks.Notifier)
	deploymentQueue := new(mocks.DeploymentQueueInterface)

	issueService := new(mocks.IssueService)

	checkCD := NewCheckCD(notifier, issueService, deploymentQueue)

	// Mock methods
	issueService.On("GetFirstIssueByCodeBase", work.Mobile).Return(nil)
	issueService.On("GetFirstIssueByCodeBase", work.Backend).Return(nil)
	issueService.On("GetFirstIssueByCodeBase", work.Frontend).Return(nil)

	// Act
	err := checkCD.Execute()

	// Assert
	assert.NoError(t, err)
	notifier.AssertNotCalled(t, "SendNotification")
	deploymentQueue.AssertNotCalled(t, "GetByCardKey")
	deploymentQueue.AssertNotCalled(t, "Create")
	issueService.AssertExpectations(t)
}

func TestCheckCD_Execute_IssueAssigneeNil(t *testing.T) {
	// Arrange
	notifier := new(mocks.Notifier)
	deploymentQueue := new(mocks.DeploymentQueueInterface)

	issueService := new(mocks.IssueService)
	jiraIssue := &jira.Issue{
		Key: "PROJ-123",
		Fields: &jira.IssueFields{
			Summary: "Deployment 6fc1d3e8e0d91c23df7f2a193f68f8616f7461de",
		},
	}

	checkCD := NewCheckCD(notifier, issueService, deploymentQueue)

	// Mock methods
	issueService.On("GetFirstIssueByCodeBase", work.Mobile).Return(jiraIssue)
	issueService.On("GetFirstIssueByCodeBase", work.Backend).Return(nil)
	issueService.On("GetFirstIssueByCodeBase", work.Frontend).Return(nil)
	deploymentQueue.On("GetByCardKey", "PROJ-123").Return(nil, nil)
	notifier.On("SendNotification", mock.Anything).Return("message-id", nil)
	deploymentQueue.On("Create", mock.Anything).Return(nil)

	// Act
	err := checkCD.Execute()

	// Assert
	assert.NoError(t, err)
	notifier.AssertNotCalled(t, "SendNotification")
	deploymentQueue.AssertNotCalled(t, "GetByCardKey")
	deploymentQueue.AssertNotCalled(t, "Create")
	issueService.AssertExpectations(t)
}

// Test case when there is an error extracting the hash from the issue summary
func TestCheckCD_Execute_ExtractHashError(t *testing.T) {
	// Arrange
	notifier := new(mocks.Notifier)
	deploymentQueue := new(mocks.DeploymentQueueInterface)

	issueService := new(mocks.IssueService)
	jiraIssue := &jira.Issue{
		Key: "PROJ-123",
		Fields: &jira.IssueFields{
			Summary:  "Invalid summary without hash",
			Assignee: &jira.User{DisplayName: "John Doe"},
		},
	}

	checkCD := NewCheckCD(notifier, issueService, deploymentQueue)

	// Mock methods
	issueService.On("GetFirstIssueByCodeBase", work.Mobile).Return(jiraIssue)
	issueService.On("GetFirstIssueByCodeBase", work.Backend).Return(nil)
	issueService.On("GetFirstIssueByCodeBase", work.Frontend).Return(nil)
	deploymentQueue.On("GetByCardKey", "PROJ-123").Return(nil, nil)

	// Act
	err := checkCD.Execute()

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "hash not found")
	notifier.AssertNotCalled(t, "SendNotification", mock.Anything)
	deploymentQueue.AssertNotCalled(t, "Create", mock.Anything)
	issueService.AssertExpectations(t)
}

// Test case when there is an error sending the notification
func TestCheckCD_Execute_SendNotificationError(t *testing.T) {
	// Arrange
	notifier := new(mocks.Notifier)
	deploymentQueue := new(mocks.DeploymentQueueInterface)

	issueService := new(mocks.IssueService)
	jiraIssue := &jira.Issue{
		Key: "PROJ-123",
		Fields: &jira.IssueFields{
			Summary:  "Deployment 6fc1d3e8e0d91c23df7f2a193f68f8616f7461de",
			Assignee: &jira.User{DisplayName: "John Doe"},
		},
	}

	checkCD := NewCheckCD(notifier, issueService, deploymentQueue)

	// Mock methods
	issueService.On("GetFirstIssueByCodeBase", work.Mobile).Return(jiraIssue)
	issueService.On("GetFirstIssueByCodeBase", work.Backend).Return(nil)
	issueService.On("GetFirstIssueByCodeBase", work.Frontend).Return(nil)
	deploymentQueue.On("GetByCardKey", "PROJ-123").Return(nil, nil)
	notifier.On("SendNotification", mock.Anything).Return("", errors.New("notification error"))

	// Act
	err := checkCD.Execute()

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "notification error")
	deploymentQueue.AssertNotCalled(t, "Create", mock.Anything)
	issueService.AssertExpectations(t)
	notifier.AssertExpectations(t)
}

// Test case when there is an error creating the deployment queue entry
func TestCheckCD_Execute_CreateDeploymentQueueError(t *testing.T) {
	// Arrange
	notifier := new(mocks.Notifier)
	deploymentQueue := new(mocks.DeploymentQueueInterface)

	issueService := new(mocks.IssueService)
	jiraIssue := &jira.Issue{
		Key: "PROJ-123",
		Fields: &jira.IssueFields{
			Summary:  "Deployment 6fc1d3e8e0d91c23df7f2a193f68f8616f7461de",
			Assignee: &jira.User{DisplayName: "John Doe"},
		},
	}

	checkCD := NewCheckCD(notifier, issueService, deploymentQueue)

	// Mock methods
	issueService.On("GetFirstIssueByCodeBase", work.Mobile).Return(jiraIssue)
	issueService.On("GetFirstIssueByCodeBase", work.Backend).Return(nil)
	issueService.On("GetFirstIssueByCodeBase", work.Frontend).Return(nil)
	deploymentQueue.On("GetByCardKey", "PROJ-123").Return(nil, nil)
	notifier.On("SendNotification", mock.Anything).Return("message-id", nil)
	deploymentQueue.On("Create", mock.Anything).Return(errors.New("create queue error"))

	// Act
	err := checkCD.Execute()

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "create queue error")
	issueService.AssertExpectations(t)
	notifier.AssertExpectations(t)
	deploymentQueue.AssertExpectations(t)
}
