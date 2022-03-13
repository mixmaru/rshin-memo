package usecases

import (
	"github.com/mixmaru/rshin-memo/core/entities"
	"github.com/mixmaru/rshin-memo/core/repositories"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetAllDailyListUsecaseInteractor_Handle(t *testing.T) {
	////// 準備
	rep := &repositories.DailyDataRepositoryMock{}
	rep.SetGetFunc(func() ([]*entities.DailyDataEntity, error) {
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
	})
	useCase := NewGetAllDailyListUsecase(rep)

	////// 実行
	dailyList, err := useCase.Handle()
	assert.NoError(t, err)

	////// 検証
	expected := []DailyData{
		{
			Date: "2021-01-02",
			Notes: []string{
				"noteC",
				"noteD",
			},
		},
		{
			Date: "2021-01-01",
			Notes: []string{
				"noteA",
				"noteB",
			},
		},
	}
	assert.EqualValues(t, expected, dailyList)
}
