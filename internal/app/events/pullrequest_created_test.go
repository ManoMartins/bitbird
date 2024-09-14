package events_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/manomartins/bitbird/internal/app/events"
	"github.com/manomartins/bitbird/internal/app/interfaces"
	"github.com/manomartins/bitbird/internal/app/utils"
	mocks "github.com/manomartins/bitbird/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestPullRequestCreated_Execute(t *testing.T) {
	mockNotifier := new(mocks.MockNotifier)
	mockMessagesStorage := new(mocks.MockPullRequestMessagesInterface)

	pullRequestCreated := events.NewPullRequestCreated(mockNotifier, mockMessagesStorage)

	ctx := context.Background()
	prEvent := events.PullRequestEvent{
		Actor:      events.Actor{DisplayName: "manoel_martins"},
		Repository: events.Repository{Name: "example-repo"},
		PullRequest: events.PullRequest{
			ID:          1,
			Title:       "Fix bug in login",
			State:       "open",
			Destination: events.Destination{Branch: events.Branch{Name: "master"}},
			Reviewers:   []events.Reviewer{{DisplayName: "jean_paes_rabello"}},
			Links:       events.PullRequestLinks{HTML: events.HTML{Href: "http://example.com/pr/1"}},
		},
	}

	events.DiscordUsers[utils.ToSnakeCase(prEvent.Actor.DisplayName)] = "667184274428002345"

	mockNotifier.On("GetUserAvatarURL", mock.Anything, "667184274428002345").Return("http://example.com/avatar/1", nil)
	mockNotifier.On("SendNotificationEmbed", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("interfaces.EmbedData")).Return("message-id-1", nil)
	mockMessagesStorage.On("Create", mock.Anything, "1", mock.AnythingOfType("string"), "message-id-1").Return(nil)

	err := pullRequestCreated.Execute(ctx, prEvent)

	assert.NoError(t, err)
	mockNotifier.AssertExpectations(t)
	mockMessagesStorage.AssertExpectations(t)
}

func TestPullRequestCreated_GetUserAvatarURLNotFound(t *testing.T) {
	mockNotifier := new(mocks.MockNotifier)
	mockMessagesStorage := new(mocks.MockPullRequestMessagesInterface)

	pullRequestCreated := events.NewPullRequestCreated(mockNotifier, mockMessagesStorage)

	ctx := context.Background()
	prEvent := events.PullRequestEvent{
		Actor:      events.Actor{DisplayName: "user-not-exist"},
		Repository: events.Repository{Name: "example-repo"},
		PullRequest: events.PullRequest{
			ID:          1,
			Title:       "Fix bug in login",
			State:       "open",
			Destination: events.Destination{Branch: events.Branch{Name: "master"}},
			Links:       events.PullRequestLinks{HTML: events.HTML{Href: "http://example.com/pr/1"}},
		},
	}

	mockNotifier.On("SendNotificationEmbed", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("interfaces.EmbedData")).Return("message-id-1", nil)
	mockMessagesStorage.On("Create", mock.Anything, "1", mock.AnythingOfType("string"), "message-id-1").Return(nil)

	err := pullRequestCreated.Execute(ctx, prEvent)

	assert.NoError(t, err)
	mockNotifier.AssertNotCalled(t, "GetUserAvatarURL")
	mockMessagesStorage.AssertCalled(t, "Create", mock.Anything, "1", mock.AnythingOfType("string"), "message-id-1")
}

func TestPullRequestCreated_GetUserAvatarURLFailure(t *testing.T) {
	expectedError := errors.New("error to get avatar url")

	mockNotifier := new(mocks.MockNotifier)
	mockMessagesStorage := new(mocks.MockPullRequestMessagesInterface)

	pullRequestCreated := events.NewPullRequestCreated(mockNotifier, mockMessagesStorage)

	ctx := context.Background()
	prEvent := events.PullRequestEvent{
		Actor:      events.Actor{DisplayName: "manoel_martins"},
		Repository: events.Repository{Name: "example-repo"},
		PullRequest: events.PullRequest{
			ID:          1,
			Title:       "Fix bug in login",
			State:       "open",
			Destination: events.Destination{Branch: events.Branch{Name: "master"}},
			Links:       events.PullRequestLinks{HTML: events.HTML{Href: "http://example.com/pr/1"}},
		},
	}

	events.DiscordUsers[utils.ToSnakeCase(prEvent.Actor.DisplayName)] = "667184274428002345"

	mockNotifier.On("GetUserAvatarURL", mock.Anything, "667184274428002345").Return("", expectedError)
	mockNotifier.On("SendNotificationEmbed", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("interfaces.EmbedData")).Return("message-id-1", nil)
	mockMessagesStorage.On("Create", mock.Anything, "1", mock.AnythingOfType("string"), "message-id-1").Return(nil)

	err := pullRequestCreated.Execute(ctx, prEvent)

	assert.EqualError(t, err, expectedError.Error())
	mockNotifier.AssertCalled(t, "GetUserAvatarURL", mock.Anything, "667184274428002345")
	mockNotifier.AssertNotCalled(t, "SendNotificationEmbed")
	mockMessagesStorage.AssertNotCalled(t, "Create")
}

func TestPullRequestCreated_FormatMessageWithoutReviewers(t *testing.T) {
	mockNotifier := new(mocks.MockNotifier)
	mockMessagesStorage := new(mocks.MockPullRequestMessagesInterface)

	pullRequestCreated := events.NewPullRequestCreated(mockNotifier, mockMessagesStorage)

	ctx := context.Background()
	prEvent := events.PullRequestEvent{
		Actor:      events.Actor{DisplayName: "john doe"},
		Repository: events.Repository{Name: "example-repo"},
		PullRequest: events.PullRequest{
			ID:          1,
			Title:       "Fix bug in login",
			State:       "open",
			Destination: events.Destination{Branch: events.Branch{Name: "master"}},
			Links:       events.PullRequestLinks{HTML: events.HTML{Href: "http://example.com/pr/1"}},
		},
	}
	expectedEmbedData := interfaces.EmbedData{
		Title:     "🚀 Novo Pull Request",
		CreatedAt: time.Now(),
		Message:   "Título: `Fix bug in login`\nDestino: `master`\nRepositorio: `example-repo`\n\n*Nenhum revisor atribuído.*\n",
		Author:    "john doe",
		AuthorURL: "author-url",
		Content:   "📰 Um novo Pull Request foi aberto ||<@discord-id>||",
		Fields: []*interfaces.EmbedField{
			{
				Name:   " ",
				Value:  fmt.Sprintf("[**🔗 Clique aqui para ver o PR**](%s)", prEvent.PullRequest.Links.HTML.Href),
				Inline: false,
			},
		},
	}

	events.DiscordUsers[utils.ToSnakeCase(prEvent.Actor.DisplayName)] = "discord-id"

	mockNotifier.On("GetUserAvatarURL", mock.Anything, "discord-id").Return("author-url", nil)
	mockNotifier.On("SendNotificationEmbed", mock.Anything, mock.AnythingOfType("string"), mock.MatchedBy(func(data interfaces.EmbedData) bool {
		return data.Title == expectedEmbedData.Title &&
			data.Message == expectedEmbedData.Message &&
			data.Author == expectedEmbedData.Author &&
			data.AuthorURL == expectedEmbedData.AuthorURL &&
			data.Content == expectedEmbedData.Content &&
			len(data.Fields) == len(expectedEmbedData.Fields)
	})).Return("message-id", nil)
	mockMessagesStorage.On("Create", mock.Anything, "1", mock.AnythingOfType("string"), "message-id").Return(nil)

	err := pullRequestCreated.Execute(ctx, prEvent)

	assert.NoError(t, err)
	mockNotifier.AssertExpectations(t)
	mockMessagesStorage.AssertExpectations(t)
}

func TestPullRequestCreated_SendNotificationEmbedFailure(t *testing.T) {
	expectedError := errors.New("failure to send notification embed")

	mockNotifier := new(mocks.MockNotifier)
	mockMessagesStorage := new(mocks.MockPullRequestMessagesInterface)

	pullRequestCreated := events.NewPullRequestCreated(mockNotifier, mockMessagesStorage)

	ctx := context.Background()
	prEvent := events.PullRequestEvent{
		Actor:      events.Actor{DisplayName: "john_doe"},
		Repository: events.Repository{Name: "example-repo"},
		PullRequest: events.PullRequest{
			ID:          1,
			Title:       "Fix bug in login",
			State:       "open",
			Destination: events.Destination{Branch: events.Branch{Name: "master"}},
			Links:       events.PullRequestLinks{HTML: events.HTML{Href: "http://example.com/pr/1"}},
		},
	}

	events.DiscordUsers[utils.ToSnakeCase(prEvent.Actor.DisplayName)] = "discord-id"

	mockNotifier.On("GetUserAvatarURL", mock.Anything, "discord-id").Return("author-url", nil)
	mockNotifier.On("SendNotificationEmbed", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("interfaces.EmbedData")).Return("", expectedError)
	mockMessagesStorage.On("Create", mock.Anything, "1", mock.AnythingOfType("string"), "message-id-1").Return(nil)

	err := pullRequestCreated.Execute(ctx, prEvent)

	assert.Error(t, err)
	assert.EqualError(t, err, expectedError.Error())
	mockNotifier.AssertExpectations(t)
	mockMessagesStorage.AssertNotCalled(t, "Create")
}

func TestPullRequestCreated_StoreMessageFailure(t *testing.T) {
	expectedError := errors.New("failure to store message")

	mockNotifier := new(mocks.MockNotifier)
	mockMessagesStorage := new(mocks.MockPullRequestMessagesInterface)

	pullRequestCreated := events.NewPullRequestCreated(mockNotifier, mockMessagesStorage)

	ctx := context.Background()
	prEvent := events.PullRequestEvent{
		Actor:      events.Actor{DisplayName: "john doe"},
		Repository: events.Repository{Name: "example-repo"},
		PullRequest: events.PullRequest{
			ID:          1,
			Title:       "Fix bug in login",
			State:       "open",
			Destination: events.Destination{Branch: events.Branch{Name: "master"}},
			Links:       events.PullRequestLinks{HTML: events.HTML{Href: "http://example.com/pr/1"}},
		},
	}

	events.DiscordUsers[utils.ToSnakeCase(prEvent.Actor.DisplayName)] = "discord-id"

	mockNotifier.On("GetUserAvatarURL", mock.Anything, "discord-id").Return("author-url", nil)
	mockNotifier.On("SendNotificationEmbed", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("interfaces.EmbedData")).Return("message-id", nil)
	mockMessagesStorage.On("Create", mock.Anything, "1", mock.AnythingOfType("string"), "message-id").Return(expectedError)

	err := pullRequestCreated.Execute(ctx, prEvent)

	assert.Error(t, err)
	assert.EqualError(t, err, expectedError.Error())
	mockNotifier.AssertExpectations(t)
	mockMessagesStorage.AssertExpectations(t)
}
