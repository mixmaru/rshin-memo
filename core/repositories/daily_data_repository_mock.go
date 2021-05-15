package repositories

import (
	"github.com/mixmaru/rshin-memo/core/entities"
	"time"
)

type DailyDataRepositoryMock struct{}

func (d *DailyDataRepositoryMock) Get() ([]*entities.DailyDataEntity, error) {
	retEntities := []*entities.DailyDataEntity{
		entities.NewDailyDataEntity(
			time.Date(2021, 1, 2, 0, 0, 0, 0, time.Local),
			[]string{
				"noteC",
				"noteD",
			},
		),
		entities.NewDailyDataEntity(
			time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local),
			[]string{
				"noteA",
				"noteB",
			},
		),
	}
	return retEntities, nil
}

func (d *DailyDataRepositoryMock) Save(entity *entities.DailyDataEntity) error {
	return nil
}
