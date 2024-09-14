package interfaces

import "github.com/manomartins/bitbird/internal/app/model"

type DeploymentQueueInterface interface {
	Create(data model.DeploymentQueueModel) error
	GetByCardKey(key string) (*model.DeploymentQueueModel, error)
}
