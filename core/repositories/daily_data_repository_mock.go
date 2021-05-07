package repositories

import "github.com/mixmaru/rshin-memo/core/entities"

type DailyDataRepositoryMock struct{}

func (d *DailyDataRepositoryMock) Save(entity *entities.DailyDataEntity) error {
	return nil
}
