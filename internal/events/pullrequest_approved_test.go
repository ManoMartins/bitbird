package events_test

import (
	"context"
	mocks "github.com/manomartins/bitbird/internal/mocks"
	"testing"

	"github.com/manomartins/bitbird/internal/events"
	"github.com/manomartins/bitbird/internal/model"
	"github.com/manomartins/bitbird/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPullRequestApproved_Execute_AddApprovalEmoji(t *testing.T) {
	mockNotifier := new(mocks.MockNotifier)
	mockMessagesStorage := new(mocks.MockPullRequestMessagesInterface)

	pullRequestApproved := events.NewPullRequestApproved(mockNotifier, mockMessagesStorage)

	ctx := context.Background()
	prID := 123
	pr := &model.PullRequestMessageModel{
		ChannelID: "channel-id",
		MessageID: "message-id",
	}

	mockMessagesStorage.On("GetById", "123").Return(pr, nil)
	mockNotifier.On("AddApprovalEmoji", "channel-id", "message-id").Return(nil)

	event := events.PullRequestEvent{
		PullRequest: events.PullRequest{
			ID:    prID,
			Title: "Test PR",
			Links: events.PullRequestLinks{
				HTML: events.HTML{Href: "http://example.com/pr"},
			},
		},
		Actor: events.Actor{
			DisplayName: "Manoel Martins",
		},
	}

	err := pullRequestApproved.Execute(ctx, event)
	assert.NoError(t, err)

	mockMessagesStorage.AssertExpectations(t)
	mockNotifier.AssertExpectations(t)
}

func TestPullRequestApproved_Execute_SendDirectMessage(t *testing.T) {
	mockNotifier := new(mocks.MockNotifier)
	mockMessagesStorage := new(mocks.MockPullRequestMessagesInterface)

	pullRequestApproved := events.NewPullRequestApproved(mockNotifier, mockMessagesStorage)

	ctx := context.Background()
	prID := 123
	pr := &model.PullRequestMessageModel{
		ChannelID: "channel-id",
		MessageID: "message-id",
	}

	mockMessagesStorage.On("GetById", "123").Return(pr, nil)
	mockNotifier.On("AddApprovalEmoji", "channel-id", "message-id").Return(nil)

	event := events.PullRequestEvent{
		PullRequest: events.PullRequest{
			ID:    prID,
			Title: "Test PR",
			Links: events.PullRequestLinks{
				HTML: events.HTML{Href: "http://example.com/pr"},
			},
			Actor: events.Actor{
				DisplayName: "Manoel Martins",
			},
		},
	}

	// Simula o usuário no DiscordUsers
	events.DiscordUsers[utils.ToSnakeCase(event.PullRequest.Actor.DisplayName)] = "123456789"

	expectedMessage := "Manoel Martins o pull request **Test PR** foi aprovado! [**Clique aqui para ver o PR**](http://example.com/pr). 🎉"

	mockNotifier.On("SendDirectMessage", ctx, "123456789", expectedMessage).Return(nil)

	err := pullRequestApproved.Execute(ctx, event)
	assert.NoError(t, err)

	mockMessagesStorage.AssertExpectations(t)
	mockNotifier.AssertExpectations(t)
}

func TestPullRequestApproved_Execute_NoDirectMessageForNonAcceptedUsers(t *testing.T) {
	mockNotifier := new(mocks.MockNotifier)
	mockMessagesStorage := new(mocks.MockPullRequestMessagesInterface)

	pullRequestApproved := events.NewPullRequestApproved(mockNotifier, mockMessagesStorage)

	ctx := context.Background()
	prID := 123
	pr := &model.PullRequestMessageModel{
		ChannelID: "channel-id",
		MessageID: "message-id",
	}

	mockMessagesStorage.On("GetById", "123").Return(pr, nil)
	mockNotifier.On("AddApprovalEmoji", "channel-id", "message-id").Return(nil)

	event := events.PullRequestEvent{
		PullRequest: events.PullRequest{
			ID:    prID,
			Title: "Test PR",
			Links: events.PullRequestLinks{
				HTML: events.HTML{Href: "http://example.com/pr"},
			},
		},
		Actor: events.Actor{
			DisplayName: "Manoel Martins",
		},
	}

	err := pullRequestApproved.Execute(ctx, event)
	assert.NoError(t, err)

	mockMessagesStorage.AssertExpectations(t)
	mockNotifier.AssertExpectations(t)

	// Verifica que SendDirectMessage não foi chamado
	mockNotifier.AssertNotCalled(t, "SendDirectMessage", mock.Anything, mock.Anything, mock.Anything)
}
