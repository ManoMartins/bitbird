package interfaces

import "github.com/manomartins/bitbird/internal/model"

type PullRequestMessagesInterface interface {
	GetById(prID string) (*model.PullRequestMessageModel, error)
	FindAll() ([]model.PullRequestMessageModel, error)
	Save()
	Update(prID string, channelID string, messageID string) error
	Create(prID string, channelID string, messageID string) error
}
