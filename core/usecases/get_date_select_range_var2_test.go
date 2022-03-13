package usecases

import (
	"github.com/mixmaru/rshin-memo/core/entities"
	"github.com/mixmaru/rshin-memo/core/repositories"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetDateRangeUseCaseVer2_Handle(t *testing.T) {
	t.Run("カーソル上下にデータがある場合のテスト", func(t *testing.T) {
		////// 準備
		now := time.Date(2021, 5, 1, 0, 0, 0, 0, time.Local)
		repo := &repositories.DailyDataRepositoryMock{}
		repo.SetGetFunc(func() ([]*entities.DailyDataEntity, error) {
			retEntities := []*entities.DailyDataEntity{
				entities.NewDailyDataEntity(
					time.Date(2021, 1, 5, 0, 0, 0, 0, time.Local),
					[]string{
						"noteC",
						"noteB",
						"noteA",
					},
				),
				entities.NewDailyDataEntity(
					time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local),
					[]string{
						"noteC",
						"noteB",
						"noteA",
					},
				),
				entities.NewDailyDataEntity(
					time.Date(2020, 12, 25, 0, 0, 0, 0, time.Local),
					[]string{
						"noteC",
						"noteB",
						"noteA",
					},
				),
			}
			return retEntities, nil
		})
		useCase := NewGetDateSelectRangeVer2UseCase(now, repo)

		//t.Run("一番古い指定で更に古い方に追加", func(t *testing.T) {
		//
		//	////// 検証1
		//	date := time.Date(2021, 2, 1, 0, 0, 0, 0, time.Local)
		//	dates, err := useCase.Handle("noteA", date, INSERT_NEWER_MODE)
		//	assert.NoError(t, err)
		//	expected := []time.Time{
		//		time.Date(2020, 12, 25, 0, 0, 0, 0, time.Local),
		//		time.Date(2020, 12, 24, 0, 0, 0, 0, time.Local),
		//		// ...
		//	}
		//	assert.Equal(t, expected, dates)
		//})

		t.Run("同日のみ", func(t *testing.T) {
			////// 検証
			date := time.Date(2020, 12, 25, 0, 0, 0, 0, time.Local)
			dates, err := useCase.Handle("noteB", date, INSERT_NEWER_MODE)
			assert.NoError(t, err)
			expected := []time.Time{
				time.Date(2020, 12, 25, 0, 0, 0, 0, time.Local),
			}
			assert.Equal(t, expected, dates)
		})

		//t.Run("INSERT_OLDER_MODE", func(t *testing.T) {
		//	////// 準備
		//	now := time.Date(2021, 5, 1, 0, 0, 0, 0, time.Local)
		//	repo := &repositories.DailyDataRepositoryMock{}
		//	repo.SetGetFunc(func() ([]*entities.DailyDataEntity, error) {
		//		retEntities := []*entities.DailyDataEntity{
		//			entities.NewDailyDataEntity(
		//				time.Date(2021, 1, 2, 0, 0, 0, 0, time.Local),
		//				[]string{
		//					"noteC",
		//					"noteD",
		//				},
		//			),
		//			entities.NewDailyDataEntity(
		//				time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local),
		//				[]string{
		//					"noteA",
		//					"noteB",
		//				},
		//			),
		//		}
		//		return retEntities, nil
		//	})
		//	useCase := NewGetDateSelectRangeVer2UseCase(now, repo)
		//
		//	currentDate := time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local)
		//
		//	////// 検証1
		//	//dates, err := useCase.Handle(overCurrentDate, currentDate, underCurrentDate, INSERT_OLDER_MODE)
		//	dates, err := useCase.Handle("noteA", currentDate, INSERT_NEWER_MODE)
		//	assert.NoError(t, err)
		//	expected := []time.Time{
		//		//time.Date(2021, 1, 10, 0, 0, 0, 0, time.Local),
		//		//time.Date(2021, 1, 9, 0, 0, 0, 0, time.Local),
		//		//time.Date(2021, 1, 8, 0, 0, 0, 0, time.Local),
		//		//time.Date(2021, 1, 7, 0, 0, 0, 0, time.Local),
		//		//time.Date(2021, 1, 6, 0, 0, 0, 0, time.Local),
		//		time.Date(2021, 1, 2, 0, 0, 0, 0, time.Local),
		//		time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local),
		//	}
		//	assert.Equal(t, expected, dates)
		//})

		//	t.Run("INSERT_OLDER_MODE 間がだいぶ空いている場合", func(t *testing.T) {
		//		////// 準備
		//		now := time.Date(2021, 5, 1, 0, 0, 0, 0, time.Local)
		//		useCase := NewGetDateSelectRangeVer2UseCase(now)
		//
		//		overCurrentDate := time.Date(2021, 1, 15, 0, 0, 0, 0, time.Local)
		//		currentDate := time.Date(2021, 1, 10, 0, 0, 0, 0, time.Local)
		//		underCurrentDate := time.Date(2019, 10, 1, 0, 0, 0, 0, time.Local)
		//
		//		////// 検証1
		//		dates, err := useCase.Handle(overCurrentDate, currentDate, underCurrentDate, INSERT_OLDER_MODE)
		//		assert.NoError(t, err)
		//		expected := []time.Time{
		//			time.Date(2021, 1, 10, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 9, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 8, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 7, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 6, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 5, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 4, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 3, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 2, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local),
		//			time.Date(2020, 12, 31, 0, 0, 0, 0, time.Local),
		//			time.Date(2020, 12, 30, 0, 0, 0, 0, time.Local),
		//			time.Date(2020, 12, 29, 0, 0, 0, 0, time.Local),
		//			time.Date(2020, 12, 28, 0, 0, 0, 0, time.Local),
		//			time.Date(2020, 12, 27, 0, 0, 0, 0, time.Local),
		//			time.Date(2020, 12, 26, 0, 0, 0, 0, time.Local),
		//			time.Date(2020, 12, 25, 0, 0, 0, 0, time.Local),
		//			time.Date(2020, 12, 24, 0, 0, 0, 0, time.Local),
		//			time.Date(2020, 12, 23, 0, 0, 0, 0, time.Local),
		//			time.Date(2020, 12, 22, 0, 0, 0, 0, time.Local),
		//			time.Date(2020, 12, 21, 0, 0, 0, 0, time.Local),
		//			time.Date(2020, 12, 20, 0, 0, 0, 0, time.Local),
		//			time.Date(2020, 12, 19, 0, 0, 0, 0, time.Local),
		//			time.Date(2020, 12, 18, 0, 0, 0, 0, time.Local),
		//			time.Date(2020, 12, 17, 0, 0, 0, 0, time.Local),
		//			time.Date(2020, 12, 16, 0, 0, 0, 0, time.Local),
		//			time.Date(2020, 12, 15, 0, 0, 0, 0, time.Local),
		//			time.Date(2020, 12, 14, 0, 0, 0, 0, time.Local),
		//			time.Date(2020, 12, 13, 0, 0, 0, 0, time.Local),
		//			time.Date(2020, 12, 12, 0, 0, 0, 0, time.Local),
		//		}
		//		assert.Equal(t, expected, dates)
		//	})
		//
		//	t.Run("INSERT_NEWER_MODE", func(t *testing.T) {
		//		////// 準備
		//		now := time.Date(2021, 5, 1, 0, 0, 0, 0, time.Local)
		//		useCase := NewGetDateSelectRangeVer2UseCase(now)
		//
		//		overCurrentDate := time.Date(2021, 1, 15, 0, 0, 0, 0, time.Local)
		//		currentDate := time.Date(2021, 1, 10, 0, 0, 0, 0, time.Local)
		//		underCurrentDate := time.Date(2021, 1, 5, 0, 0, 0, 0, time.Local)
		//
		//		////// 検証1
		//		dates, err := useCase.Handle(overCurrentDate, currentDate, underCurrentDate, INSERT_NEWER_MODE)
		//		assert.NoError(t, err)
		//		expected := []time.Time{
		//			time.Date(2021, 1, 10, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 11, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 12, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 13, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 14, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 15, 0, 0, 0, 0, time.Local),
		//		}
		//		assert.Equal(t, expected, dates)
		//	})
		//})
		//
		//t.Run("カーソルの下行にデータがない場合のテスト", func(t *testing.T) {
		//	t.Run("INSERT_OLDER_MODE", func(t *testing.T) {
		//		////// 準備
		//		now := time.Date(2021, 5, 1, 0, 0, 0, 0, time.Local)
		//		useCase := NewGetDateSelectRangeVer2UseCase(now)
		//
		//		overCurrentDate := time.Date(2021, 1, 15, 0, 0, 0, 0, time.Local)
		//		currentDate := time.Date(2021, 1, 10, 0, 0, 0, 0, time.Local)
		//		underCurrentDate := time.Time{}
		//
		//		////// 検証1
		//		dates, err := useCase.Handle(overCurrentDate, currentDate, underCurrentDate, INSERT_OLDER_MODE)
		//		assert.NoError(t, err)
		//		expected := []time.Time{
		//			time.Date(2021, 1, 10, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 9, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 8, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 7, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 6, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 5, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 4, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 3, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 2, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local),
		//			time.Date(2020, 12, 31, 0, 0, 0, 0, time.Local),
		//			time.Date(2020, 12, 30, 0, 0, 0, 0, time.Local),
		//			time.Date(2020, 12, 29, 0, 0, 0, 0, time.Local),
		//			time.Date(2020, 12, 28, 0, 0, 0, 0, time.Local),
		//			time.Date(2020, 12, 27, 0, 0, 0, 0, time.Local),
		//			time.Date(2020, 12, 26, 0, 0, 0, 0, time.Local),
		//			time.Date(2020, 12, 25, 0, 0, 0, 0, time.Local),
		//			time.Date(2020, 12, 24, 0, 0, 0, 0, time.Local),
		//			time.Date(2020, 12, 23, 0, 0, 0, 0, time.Local),
		//			time.Date(2020, 12, 22, 0, 0, 0, 0, time.Local),
		//			time.Date(2020, 12, 21, 0, 0, 0, 0, time.Local),
		//			time.Date(2020, 12, 20, 0, 0, 0, 0, time.Local),
		//			time.Date(2020, 12, 19, 0, 0, 0, 0, time.Local),
		//			time.Date(2020, 12, 18, 0, 0, 0, 0, time.Local),
		//			time.Date(2020, 12, 17, 0, 0, 0, 0, time.Local),
		//			time.Date(2020, 12, 16, 0, 0, 0, 0, time.Local),
		//			time.Date(2020, 12, 15, 0, 0, 0, 0, time.Local),
		//			time.Date(2020, 12, 14, 0, 0, 0, 0, time.Local),
		//			time.Date(2020, 12, 13, 0, 0, 0, 0, time.Local),
		//			time.Date(2020, 12, 12, 0, 0, 0, 0, time.Local),
		//		}
		//		assert.Equal(t, expected, dates)
		//	})
		//
		//	t.Run("INSERT_NEWER_MODE", func(t *testing.T) {
		//		////// 準備
		//		now := time.Date(2021, 5, 1, 0, 0, 0, 0, time.Local)
		//		useCase := NewGetDateSelectRangeVer2UseCase(now)
		//
		//		overCurrentDate := time.Date(2021, 1, 15, 0, 0, 0, 0, time.Local)
		//		currentDate := time.Date(2021, 1, 10, 0, 0, 0, 0, time.Local)
		//		underCurrentDate := time.Time{}
		//
		//		////// 検証1
		//		dates, err := useCase.Handle(overCurrentDate, currentDate, underCurrentDate, INSERT_NEWER_MODE)
		//		assert.NoError(t, err)
		//		expected := []time.Time{
		//			time.Date(2021, 1, 10, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 11, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 12, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 13, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 14, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 15, 0, 0, 0, 0, time.Local),
		//		}
		//		assert.Equal(t, expected, dates)
		//	})
		//})
		//
		//t.Run("カーソルの上行にデータがない場合のテスト", func(t *testing.T) {
		//	t.Run("INSERT_OLDER_MODE", func(t *testing.T) {
		//		////// 準備
		//		now := time.Date(2021, 5, 1, 0, 0, 0, 0, time.Local)
		//		useCase := NewGetDateSelectRangeVer2UseCase(now)
		//
		//		overCurrentDate := time.Time{}
		//		currentDate := time.Date(2021, 1, 10, 0, 0, 0, 0, time.Local)
		//		underCurrentDate := time.Date(2021, 1, 5, 0, 0, 0, 0, time.Local)
		//
		//		////// 検証1
		//		dates, err := useCase.Handle(overCurrentDate, currentDate, underCurrentDate, INSERT_OLDER_MODE)
		//		assert.NoError(t, err)
		//		expected := []time.Time{
		//			time.Date(2021, 1, 10, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 9, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 8, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 7, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 6, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 5, 0, 0, 0, 0, time.Local),
		//		}
		//		assert.Equal(t, expected, dates)
		//	})
		//
		//	t.Run("INSERT_NEWER_MODE", func(t *testing.T) {
		//		////// 準備
		//		now := time.Date(2021, 5, 1, 0, 0, 0, 0, time.Local)
		//		useCase := NewGetDateSelectRangeVer2UseCase(now)
		//
		//		overCurrentDate := time.Time{}
		//		currentDate := time.Date(2021, 1, 10, 0, 0, 0, 0, time.Local)
		//		underCurrentDate := time.Date(2021, 1, 5, 0, 0, 0, 0, time.Local)
		//
		//		////// 検証1
		//		dates, err := useCase.Handle(overCurrentDate, currentDate, underCurrentDate, INSERT_NEWER_MODE)
		//		assert.NoError(t, err)
		//		expected := []time.Time{
		//			time.Date(2021, 1, 10, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 11, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 12, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 13, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 14, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 15, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 16, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 17, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 18, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 19, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 20, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 21, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 22, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 23, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 24, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 25, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 26, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 27, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 28, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 29, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 30, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 31, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 2, 1, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 2, 2, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 2, 3, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 2, 4, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 2, 5, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 2, 6, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 2, 7, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 2, 8, 0, 0, 0, 0, time.Local),
		//		}
		//		assert.Equal(t, expected, dates)
		//	})
		//})
		//
		//t.Run("最初まったくデータがないときにつかうやつ", func(t *testing.T) {
		//	t.Run("INSERT_OLDER_MODE", func(t *testing.T) {
		//		////// 準備
		//		now := time.Date(2021, 1, 10, 0, 0, 0, 0, time.Local)
		//		useCase := NewGetDateSelectRangeVer2UseCase(now)
		//
		//		overCurrentDate := time.Time{}
		//		currentDate := time.Time{}
		//		underCurrentDate := time.Time{}
		//
		//		////// 検証1
		//		dates, err := useCase.Handle(overCurrentDate, currentDate, underCurrentDate, INSERT_OLDER_MODE)
		//		assert.NoError(t, err)
		//		expected := []time.Time{
		//			time.Date(2021, 1, 24, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 23, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 22, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 21, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 20, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 19, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 18, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 17, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 16, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 15, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 14, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 13, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 12, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 11, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 10, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 9, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 8, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 7, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 6, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 5, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 4, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 3, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 2, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local),
		//			time.Date(2020, 12, 31, 0, 0, 0, 0, time.Local),
		//			time.Date(2020, 12, 30, 0, 0, 0, 0, time.Local),
		//			time.Date(2020, 12, 29, 0, 0, 0, 0, time.Local),
		//			time.Date(2020, 12, 28, 0, 0, 0, 0, time.Local),
		//			time.Date(2020, 12, 27, 0, 0, 0, 0, time.Local),
		//			time.Date(2020, 12, 26, 0, 0, 0, 0, time.Local),
		//		}
		//		assert.Equal(t, expected, dates)
		//	})
		//
		//	t.Run("INSERT_NEWER_MODE", func(t *testing.T) {
		//		////// 準備
		//		now := time.Date(2021, 1, 10, 0, 0, 0, 0, time.Local)
		//		useCase := NewGetDateSelectRangeVer2UseCase(now)
		//
		//		overCurrentDate := time.Time{}
		//		currentDate := time.Time{}
		//		underCurrentDate := time.Time{}
		//
		//		////// 検証1
		//		dates, err := useCase.Handle(overCurrentDate, currentDate, underCurrentDate, INSERT_NEWER_MODE)
		//		assert.NoError(t, err)
		//		expected := []time.Time{
		//			time.Date(2020, 12, 26, 0, 0, 0, 0, time.Local),
		//			time.Date(2020, 12, 27, 0, 0, 0, 0, time.Local),
		//			time.Date(2020, 12, 28, 0, 0, 0, 0, time.Local),
		//			time.Date(2020, 12, 29, 0, 0, 0, 0, time.Local),
		//			time.Date(2020, 12, 30, 0, 0, 0, 0, time.Local),
		//			time.Date(2020, 12, 31, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 2, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 3, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 4, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 5, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 6, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 7, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 8, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 9, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 10, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 11, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 12, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 13, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 14, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 15, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 16, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 17, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 18, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 19, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 20, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 21, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 22, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 23, 0, 0, 0, 0, time.Local),
		//			time.Date(2021, 1, 24, 0, 0, 0, 0, time.Local),
		//		}
		//		assert.Equal(t, expected, dates)
		//	})
	})
}
