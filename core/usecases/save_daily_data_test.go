package usecases

import (
	"github.com/mixmaru/rshin-memo/core/repositories"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateNoteUseCaseInteractor_Handle(t *testing.T) {
	t.Run("存在しないNoteNameが1つ含まれている場合正常終了する", func(t *testing.T) {
		// 準備
		interactor := NewSaveDailyDataUseCase(
			&repositories.NoteRepositoryMock{},
			&repositories.DailyDataRepositoryMock{},
		)
		dailyData := DailyData{
			Date: "2021-02-02",
			Notes: []string{
				repositories.EXIST_NOTE_NAME,
				repositories.NOT_EXIST_NOTE_NAME,
			},
		}
		// 実行
		err := interactor.Handle(dailyData)
		// 検証
		assert.NoError(t, err)
	})

	t.Run("存在しないNoteNameがない場合正常終了する", func(t *testing.T) {
		// 準備
		interactor := NewSaveDailyDataUseCase(
			&repositories.NoteRepositoryMock{},
			&repositories.DailyDataRepositoryMock{},
		)
		dailyData := DailyData{
			Date: "2021-02-02",
			Notes: []string{
				repositories.EXIST_NOTE_NAME,
			},
		}
		// 実行
		err := interactor.Handle(dailyData)
		// 検証
		assert.NoError(t, err)
	})
}
