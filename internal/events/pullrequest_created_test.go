package events_test

import (
	"context"
	"github.com/manomartins/bitbird/internal/events"
	mocks "github.com/manomartins/bitbird/internal/mocks"
	"github.com/manomartins/bitbird/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
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
