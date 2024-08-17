package events

import (
	"errors"
	"github.com/manomartins/bitbird/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestPullRequestCreated_Execute_Success(t *testing.T) {
	// Create the mocks
	notifier := new(mocks.Notifier)
	messagesStorage := new(mocks.PullRequestMessagesInterface)

	// Set up the expected interactions
	expectedPrID := "123"
	expectedMessageID := "123"

	expectedMessage := "**🚀 *Detalhes do Pull Request:* **\n\n" +
		"**Título:** `Update README`\n" +
		"**Status:** `OPEN`\n" +
		"**Autor:** manoel martins\n" +
		"**Destino:** `develop`\n" +
		"**🌟 Repositorio:** `example-repo`\n" +
		"**Link do PR:** [Abrir PR](https://example.com/pr/123)\n\n" +
		"*Nenhum revisor atribuído.*\n"

	messagesStorage.On("UpdatePullRequestMessage", expectedPrID, expectedMessageID).Return(nil)
	notifier.On("SendNotification", expectedMessage).Return(expectedMessageID, nil)

	// Create the PullRequestApproved instance
	prCreated := NewPullRequestCreated(notifier, messagesStorage)

	// Define the event to test
	event := PullRequestEvent{
		Repository: Repository{Name: "example-repo"},
		Actor:      Actor{DisplayName: "manoel martins"},
		PullRequest: PullRequest{
			ID:     123,
			Title:  "Update README",
			State:  "OPEN",
			Source: PullRequestSource{Branch: Branch{Name: "develop"}},
			Links:  PullRequestLinks{HTML: HTML{Href: "https://example.com/pr/123"}},
		},
	}

	// Call the method
	err := prCreated.Execute(event)

	// Assert that there were no errors
	assert.NoError(t, err)

	notifier.AssertNumberOfCalls(t, "SendNotification", 1)
	messagesStorage.AssertCalled(t, "UpdatePullRequestMessage", expectedPrID, expectedMessageID)
}

func TestPullRequestCreated_Execute_SuccessWithReviewers(t *testing.T) {
	// Create the mocks
	notifier := new(mocks.Notifier)
	messagesStorage := new(mocks.PullRequestMessagesInterface)

	// Set up the expected interactions
	expectedPrID := "123"
	expectedMessageID := "123"

	expectedMessage := "**🚀 *Detalhes do Pull Request:* **\n\n" +
		"**Título:** `Update README`\n" +
		"**Status:** `OPEN`\n" +
		"**Autor:** manoel martins\n" +
		"**Destino:** `develop`\n" +
		"**🌟 Repositorio:** `example-repo`\n" +
		"**Link do PR:** [Abrir PR](https://example.com/pr/123)\n\n" +
		"**📝 Revisores:**\n- tassyo monteiro\n- manoel martins\n"

	messagesStorage.On("UpdatePullRequestMessage", expectedPrID, expectedMessageID).Return(nil)
	notifier.On("SendNotification", expectedMessage).Return(expectedMessageID, nil)

	// Create the PullRequestApproved instance
	prCreated := NewPullRequestCreated(notifier, messagesStorage)

	// Define the event to test
	event := PullRequestEvent{
		Repository: Repository{Name: "example-repo"},
		Actor:      Actor{DisplayName: "manoel martins"},
		PullRequest: PullRequest{
			ID:     123,
			Title:  "Update README",
			State:  "OPEN",
			Source: PullRequestSource{Branch: Branch{Name: "develop"}},
			Links:  PullRequestLinks{HTML: HTML{Href: "https://example.com/pr/123"}},
			Reviewers: []Reviewer{{
				User: struct {
					DisplayName string `json:"display_name"`
				}{DisplayName: "tassyo monteiro"},
			}, {
				User: struct {
					DisplayName string `json:"display_name"`
				}{DisplayName: "manoel martins"},
			}},
		},
	}

	// Call the method
	err := prCreated.Execute(event)

	// Assert that there were no errors
	assert.NoError(t, err)

	notifier.AssertExpectations(t)
	messagesStorage.AssertExpectations(t)
}

func TestPullRequestCreated_Execute_SendNotificationError(t *testing.T) {
	// Create the mocks
	notifier := new(mocks.Notifier)
	messagesStorage := new(mocks.PullRequestMessagesInterface)

	// Set up the expected interactions
	expectedError := errors.New("error sending notification")

	notifier.On("SendNotification", mock.Anything).Return("", expectedError)

	// Create the PullRequestApproved instance
	prCreated := NewPullRequestCreated(notifier, messagesStorage)

	// Define the event to test
	event := PullRequestEvent{
		Repository: Repository{Name: "example-repo"},
		Actor:      Actor{DisplayName: "manoel martins"},
		PullRequest: PullRequest{
			ID:     123,
			Title:  "Update README",
			State:  "OPEN",
			Source: PullRequestSource{Branch: Branch{Name: "develop"}},
			Links:  PullRequestLinks{HTML: HTML{Href: "https://example.com/pr/123"}},
		},
	}

	// Call the method
	err := prCreated.Execute(event)

	// Assert that the error is the one expected
	assert.Equal(t, expectedError, err)

	notifier.AssertExpectations(t)
	messagesStorage.AssertNotCalled(t, "UpdatePullRequestMessage", mock.Anything)
}

func TestPullRequestCreated_Execute_UpdatePullRequestMessageError(t *testing.T) {
	// Create the mocks
	notifier := new(mocks.Notifier)
	messagesStorage := new(mocks.PullRequestMessagesInterface)

	// Set up the expected interactions
	expectedPrID := "123"
	expectedMessageID := "123"
	expectedError := errors.New("error to update")

	notifier.On("SendNotification", mock.Anything).Return(expectedMessageID, nil)
	messagesStorage.On("UpdatePullRequestMessage", expectedPrID, expectedMessageID).Return(expectedError)

	// Create the PullRequestApproved instance
	prCreated := NewPullRequestCreated(notifier, messagesStorage)

	// Define the event to test
	event := PullRequestEvent{
		Repository: Repository{Name: "example-repo"},
		Actor:      Actor{DisplayName: "manoel martins"},
		PullRequest: PullRequest{
			ID:     123,
			Title:  "Update README",
			State:  "OPEN",
			Source: PullRequestSource{Branch: Branch{Name: "develop"}},
			Links:  PullRequestLinks{HTML: HTML{Href: "https://example.com/pr/123"}},
		},
	}

	// Call the method
	err := prCreated.Execute(event)

	// Assert that the error is the one expected
	assert.Equal(t, expectedError, err)

	notifier.AssertExpectations(t)
	messagesStorage.AssertExpectations(t)
}
