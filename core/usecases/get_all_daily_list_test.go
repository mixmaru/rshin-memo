package usecases

import (
	"github.com/mixmaru/rshin-memo/core/repositories"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetAllDailyListUsecaseInteractor_Handle(t *testing.T) {
	////// 準備
	useCase := NewGetAllDailyListUsecase(&repositories.DailyDataRepositoryMock{})

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
