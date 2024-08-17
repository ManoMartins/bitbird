package interfaces

import "github.com/manomartins/bitbird/internal/model"

type PullRequestMessagesInterface interface {
	GetPullRequestMessage(prID string) (*model.PullRequestMessageModel, error)
	FindAllPullRequestMessages() ([]model.PullRequestMessageModel, error)
	SavePullRequestMessage()
	UpdatePullRequestMessage(prID string, messageID string) error
}
