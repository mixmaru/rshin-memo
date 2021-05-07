package repositories

import "github.com/mixmaru/rshin-memo/core/entities"

type DailyDataRepositoryInterface interface {
	Save(entity *entities.DailyDataEntity) error
}
