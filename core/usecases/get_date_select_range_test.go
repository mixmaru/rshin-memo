package usecases

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetDateRangeUseCase_Handle(t *testing.T) {
	////// 準備
	useCase := NewGetDateSelectRangeUseCase()
	now := time.Date(2021, 2, 1, 0, 0, 0, 0, time.Local)

	////// 検証1
	dates := useCase.Handle(now, 2, 1)
	expected := []time.Time{
		time.Date(2021, 1, 30, 0, 0, 0, 0, time.Local),
		time.Date(2021, 1, 31, 0, 0, 0, 0, time.Local),
		time.Date(2021, 2, 1, 0, 0, 0, 0, time.Local),
		time.Date(2021, 2, 2, 0, 0, 0, 0, time.Local),
	}
	assert.Equal(t, expected, dates)

	////// 検証2
	dates = useCase.Handle(now, 0, 0)
	expected = []time.Time{
		time.Date(2021, 2, 1, 0, 0, 0, 0, time.Local),
	}
	assert.Equal(t, expected, dates)
}
