package repositories

import (
	"github.com/mixmaru/rshin-memo/core/entities"
)

type DailyDataRepositoryMock struct {
	getFunc func() ([]*entities.DailyDataEntity, error)
}

func (d *DailyDataRepositoryMock) Get() ([]*entities.DailyDataEntity, error) {
	return d.getFunc()
}

func (d *DailyDataRepositoryMock) SetGetFunc(function func() ([]*entities.DailyDataEntity, error)) {
	d.getFunc = function
}

func (d *DailyDataRepositoryMock) Save(entity *entities.DailyDataEntity) error {
	return nil
}
