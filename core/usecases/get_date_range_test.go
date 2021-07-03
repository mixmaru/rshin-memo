package usecases

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetDateRangeUseCase_Handle(t *testing.T) {
	////// 準備
	useCase := NewGetDateRangeUseCase()
	now := time.Date(2021, 2, 1, 0, 0, 0, 0, time.Local)
	////// 実行
	from, to := useCase.Handle(now, 2, 1)
	////// 検証
	expectedFrom := time.Date(2021, 1, 30, 0, 0, 0, 0, time.Local)
	expectedTo := time.Date(2021, 2, 2, 0, 0, 0, 0, time.Local)
	assert.Equal(t, expectedFrom, from)
	assert.Equal(t, expectedTo, to)
}
