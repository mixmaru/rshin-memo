package usecases

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetDateRangeUseCase_Handle(t *testing.T) {
	t.Run("カーソル上下にデータがある場合のテスト", func(t *testing.T) {
		t.Run("INSERT_UNDER_DATE_MODE", func(t *testing.T) {
			////// 準備
			now := time.Date(2021, 5, 1, 0, 0, 0, 0, time.Local)
			useCase := NewGetDateSelectRangeUseCase(now)

			overCurrentDate := time.Date(2021, 1, 15, 0, 0, 0, 0, time.Local)
			currentDate := time.Date(2021, 1, 10, 0, 0, 0, 0, time.Local)
			underCurrentDate := time.Date(2021, 1, 5, 0, 0, 0, 0, time.Local)

			////// 検証1
			dates, err := useCase.Handle(overCurrentDate, currentDate, underCurrentDate, INSERT_UNDER_DATE_MODE)
			assert.NoError(t, err)
			expected := []time.Time{
				time.Date(2021, 1, 10, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 9, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 8, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 7, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 6, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 5, 0, 0, 0, 0, time.Local),
			}
			assert.Equal(t, expected, dates)
		})

		t.Run("INSERT_UNDER_DATE_MODE 間がだいぶ空いている場合", func(t *testing.T) {
			////// 準備
			now := time.Date(2021, 5, 1, 0, 0, 0, 0, time.Local)
			useCase := NewGetDateSelectRangeUseCase(now)

			overCurrentDate := time.Date(2021, 1, 15, 0, 0, 0, 0, time.Local)
			currentDate := time.Date(2021, 1, 10, 0, 0, 0, 0, time.Local)
			underCurrentDate := time.Date(2019, 10, 1, 0, 0, 0, 0, time.Local)

			////// 検証1
			dates, err := useCase.Handle(overCurrentDate, currentDate, underCurrentDate, INSERT_UNDER_DATE_MODE)
			assert.NoError(t, err)
			expected := []time.Time{
				time.Date(2021, 1, 10, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 9, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 8, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 7, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 6, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 5, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 4, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 3, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 2, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local),
				time.Date(2020, 12, 31, 0, 0, 0, 0, time.Local),
				time.Date(2020, 12, 30, 0, 0, 0, 0, time.Local),
				time.Date(2020, 12, 29, 0, 0, 0, 0, time.Local),
				time.Date(2020, 12, 28, 0, 0, 0, 0, time.Local),
				time.Date(2020, 12, 27, 0, 0, 0, 0, time.Local),
				time.Date(2020, 12, 26, 0, 0, 0, 0, time.Local),
				time.Date(2020, 12, 25, 0, 0, 0, 0, time.Local),
				time.Date(2020, 12, 24, 0, 0, 0, 0, time.Local),
				time.Date(2020, 12, 23, 0, 0, 0, 0, time.Local),
				time.Date(2020, 12, 22, 0, 0, 0, 0, time.Local),
				time.Date(2020, 12, 21, 0, 0, 0, 0, time.Local),
				time.Date(2020, 12, 20, 0, 0, 0, 0, time.Local),
				time.Date(2020, 12, 19, 0, 0, 0, 0, time.Local),
				time.Date(2020, 12, 18, 0, 0, 0, 0, time.Local),
				time.Date(2020, 12, 17, 0, 0, 0, 0, time.Local),
				time.Date(2020, 12, 16, 0, 0, 0, 0, time.Local),
				time.Date(2020, 12, 15, 0, 0, 0, 0, time.Local),
				time.Date(2020, 12, 14, 0, 0, 0, 0, time.Local),
				time.Date(2020, 12, 13, 0, 0, 0, 0, time.Local),
				time.Date(2020, 12, 12, 0, 0, 0, 0, time.Local),
			}
			assert.Equal(t, expected, dates)
		})

		t.Run("INSERT_OVER_DATE_MODE", func(t *testing.T) {
			////// 準備
			now := time.Date(2021, 5, 1, 0, 0, 0, 0, time.Local)
			useCase := NewGetDateSelectRangeUseCase(now)

			overCurrentDate := time.Date(2021, 1, 15, 0, 0, 0, 0, time.Local)
			currentDate := time.Date(2021, 1, 10, 0, 0, 0, 0, time.Local)
			underCurrentDate := time.Date(2021, 1, 5, 0, 0, 0, 0, time.Local)

			////// 検証1
			dates, err := useCase.Handle(overCurrentDate, currentDate, underCurrentDate, INSERT_OVER_DATE_MODE)
			assert.NoError(t, err)
			expected := []time.Time{
				time.Date(2021, 1, 10, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 11, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 12, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 13, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 14, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 15, 0, 0, 0, 0, time.Local),
			}
			assert.Equal(t, expected, dates)
		})
	})

	t.Run("カーソルの下行にデータがない場合のテスト", func(t *testing.T) {
		t.Run("INSERT_UNDER_DATE_MODE", func(t *testing.T) {
			////// 準備
			now := time.Date(2021, 5, 1, 0, 0, 0, 0, time.Local)
			useCase := NewGetDateSelectRangeUseCase(now)

			overCurrentDate := time.Date(2021, 1, 15, 0, 0, 0, 0, time.Local)
			currentDate := time.Date(2021, 1, 10, 0, 0, 0, 0, time.Local)
			underCurrentDate := time.Time{}

			////// 検証1
			dates, err := useCase.Handle(overCurrentDate, currentDate, underCurrentDate, INSERT_UNDER_DATE_MODE)
			assert.NoError(t, err)
			expected := []time.Time{
				time.Date(2021, 1, 10, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 9, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 8, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 7, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 6, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 5, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 4, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 3, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 2, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local),
				time.Date(2020, 12, 31, 0, 0, 0, 0, time.Local),
				time.Date(2020, 12, 30, 0, 0, 0, 0, time.Local),
				time.Date(2020, 12, 29, 0, 0, 0, 0, time.Local),
				time.Date(2020, 12, 28, 0, 0, 0, 0, time.Local),
				time.Date(2020, 12, 27, 0, 0, 0, 0, time.Local),
				time.Date(2020, 12, 26, 0, 0, 0, 0, time.Local),
				time.Date(2020, 12, 25, 0, 0, 0, 0, time.Local),
				time.Date(2020, 12, 24, 0, 0, 0, 0, time.Local),
				time.Date(2020, 12, 23, 0, 0, 0, 0, time.Local),
				time.Date(2020, 12, 22, 0, 0, 0, 0, time.Local),
				time.Date(2020, 12, 21, 0, 0, 0, 0, time.Local),
				time.Date(2020, 12, 20, 0, 0, 0, 0, time.Local),
				time.Date(2020, 12, 19, 0, 0, 0, 0, time.Local),
				time.Date(2020, 12, 18, 0, 0, 0, 0, time.Local),
				time.Date(2020, 12, 17, 0, 0, 0, 0, time.Local),
				time.Date(2020, 12, 16, 0, 0, 0, 0, time.Local),
				time.Date(2020, 12, 15, 0, 0, 0, 0, time.Local),
				time.Date(2020, 12, 14, 0, 0, 0, 0, time.Local),
				time.Date(2020, 12, 13, 0, 0, 0, 0, time.Local),
				time.Date(2020, 12, 12, 0, 0, 0, 0, time.Local),
			}
			assert.Equal(t, expected, dates)
		})

		t.Run("INSERT_OVER_DATE_MODE", func(t *testing.T) {
			////// 準備
			now := time.Date(2021, 5, 1, 0, 0, 0, 0, time.Local)
			useCase := NewGetDateSelectRangeUseCase(now)

			overCurrentDate := time.Date(2021, 1, 15, 0, 0, 0, 0, time.Local)
			currentDate := time.Date(2021, 1, 10, 0, 0, 0, 0, time.Local)
			underCurrentDate := time.Time{}

			////// 検証1
			dates, err := useCase.Handle(overCurrentDate, currentDate, underCurrentDate, INSERT_OVER_DATE_MODE)
			assert.NoError(t, err)
			expected := []time.Time{
				time.Date(2021, 1, 10, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 11, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 12, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 13, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 14, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 15, 0, 0, 0, 0, time.Local),
			}
			assert.Equal(t, expected, dates)
		})
	})

	t.Run("カーソルの上行にデータがない場合のテスト", func(t *testing.T) {
		t.Run("INSERT_UNDER_DATE_MODE", func(t *testing.T) {
			////// 準備
			now := time.Date(2021, 5, 1, 0, 0, 0, 0, time.Local)
			useCase := NewGetDateSelectRangeUseCase(now)

			overCurrentDate := time.Time{}
			currentDate := time.Date(2021, 1, 10, 0, 0, 0, 0, time.Local)
			underCurrentDate := time.Date(2021, 1, 5, 0, 0, 0, 0, time.Local)

			////// 検証1
			dates, err := useCase.Handle(overCurrentDate, currentDate, underCurrentDate, INSERT_UNDER_DATE_MODE)
			assert.NoError(t, err)
			expected := []time.Time{
				time.Date(2021, 1, 10, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 9, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 8, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 7, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 6, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 5, 0, 0, 0, 0, time.Local),
			}
			assert.Equal(t, expected, dates)
		})

		t.Run("INSERT_OVER_DATE_MODE", func(t *testing.T) {
			////// 準備
			now := time.Date(2021, 5, 1, 0, 0, 0, 0, time.Local)
			useCase := NewGetDateSelectRangeUseCase(now)

			overCurrentDate := time.Time{}
			currentDate := time.Date(2021, 1, 10, 0, 0, 0, 0, time.Local)
			underCurrentDate := time.Date(2021, 1, 5, 0, 0, 0, 0, time.Local)

			////// 検証1
			dates, err := useCase.Handle(overCurrentDate, currentDate, underCurrentDate, INSERT_OVER_DATE_MODE)
			assert.NoError(t, err)
			expected := []time.Time{
				time.Date(2021, 1, 10, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 11, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 12, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 13, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 14, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 15, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 16, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 17, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 18, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 19, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 20, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 21, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 22, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 23, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 24, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 25, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 26, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 27, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 28, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 29, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 30, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 31, 0, 0, 0, 0, time.Local),
				time.Date(2021, 2, 1, 0, 0, 0, 0, time.Local),
				time.Date(2021, 2, 2, 0, 0, 0, 0, time.Local),
				time.Date(2021, 2, 3, 0, 0, 0, 0, time.Local),
				time.Date(2021, 2, 4, 0, 0, 0, 0, time.Local),
				time.Date(2021, 2, 5, 0, 0, 0, 0, time.Local),
				time.Date(2021, 2, 6, 0, 0, 0, 0, time.Local),
				time.Date(2021, 2, 7, 0, 0, 0, 0, time.Local),
				time.Date(2021, 2, 8, 0, 0, 0, 0, time.Local),
			}
			assert.Equal(t, expected, dates)
		})
	})

	t.Run("最初まったくデータがないときにつかうやつ", func(t *testing.T) {
		t.Run("INSERT_UNDER_DATE_MODE", func(t *testing.T) {
			////// 準備
			now := time.Date(2021, 1, 10, 0, 0, 0, 0, time.Local)
			useCase := NewGetDateSelectRangeUseCase(now)

			overCurrentDate := time.Time{}
			currentDate := time.Time{}
			underCurrentDate := time.Time{}

			////// 検証1
			dates, err := useCase.Handle(overCurrentDate, currentDate, underCurrentDate, INSERT_UNDER_DATE_MODE)
			assert.NoError(t, err)
			expected := []time.Time{
				time.Date(2021, 1, 24, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 23, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 22, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 21, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 20, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 19, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 18, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 17, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 16, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 15, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 14, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 13, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 12, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 11, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 10, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 9, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 8, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 7, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 6, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 5, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 4, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 3, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 2, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local),
				time.Date(2020, 12, 31, 0, 0, 0, 0, time.Local),
				time.Date(2020, 12, 30, 0, 0, 0, 0, time.Local),
				time.Date(2020, 12, 29, 0, 0, 0, 0, time.Local),
				time.Date(2020, 12, 28, 0, 0, 0, 0, time.Local),
				time.Date(2020, 12, 27, 0, 0, 0, 0, time.Local),
				time.Date(2020, 12, 26, 0, 0, 0, 0, time.Local),
			}
			assert.Equal(t, expected, dates)
		})

		t.Run("INSERT_OVER_DATE_MODE", func(t *testing.T) {
			////// 準備
			now := time.Date(2021, 1, 10, 0, 0, 0, 0, time.Local)
			useCase := NewGetDateSelectRangeUseCase(now)

			overCurrentDate := time.Time{}
			currentDate := time.Time{}
			underCurrentDate := time.Time{}

			////// 検証1
			dates, err := useCase.Handle(overCurrentDate, currentDate, underCurrentDate, INSERT_OVER_DATE_MODE)
			assert.NoError(t, err)
			expected := []time.Time{
				time.Date(2020, 12, 26, 0, 0, 0, 0, time.Local),
				time.Date(2020, 12, 27, 0, 0, 0, 0, time.Local),
				time.Date(2020, 12, 28, 0, 0, 0, 0, time.Local),
				time.Date(2020, 12, 29, 0, 0, 0, 0, time.Local),
				time.Date(2020, 12, 30, 0, 0, 0, 0, time.Local),
				time.Date(2020, 12, 31, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 2, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 3, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 4, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 5, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 6, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 7, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 8, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 9, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 10, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 11, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 12, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 13, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 14, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 15, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 16, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 17, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 18, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 19, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 20, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 21, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 22, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 23, 0, 0, 0, 0, time.Local),
				time.Date(2021, 1, 24, 0, 0, 0, 0, time.Local),
			}
			assert.Equal(t, expected, dates)
		})
	})
}
