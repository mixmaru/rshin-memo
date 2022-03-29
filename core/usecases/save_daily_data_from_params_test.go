package usecases

import (
	"github.com/mixmaru/rshin-memo/core/repositories"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateNoteFromParamsUseCaseInteractor_Handle(t *testing.T) {
	t.Run("存在しないNoteNameが1つ含まれている場合正常終了する", func(t *testing.T) {
		// create test data
		interactor := NewSaveDailyDataFromParamsUseCase(
			&repositories.NoteRepositoryMock{},
			&repositories.DailyDataRepositoryMock{},
		)

		// 実行
		err := interactor.Handle(
			time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local),
			"base_memo_name",
			time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local),
			"new_memo_name",
			"new_memo_内容",
		)
		// 検証,意図通りにmockを叩いているか？
		assert.NoError(t, err)

		//// 準備
		//interactor := NewSaveDailyDataUseCase(
		//	&repositories.NoteRepositoryMock{},
		//	&repositories.DailyDataRepositoryMock{},
		//)
		//dailyData := DailyData{
		//	Date: "2021-02-02",
		//	Notes: []string{
		//		repositories.EXIST_NOTE_NAME,
		//		repositories.NOT_EXIST_NOTE_NAME,
		//	},
		//}
		//// 実行
		//err := interactor.Handle(dailyData)
		//// 検証
		//assert.NoError(t, err)
	})

	//t.Run("存在しないNoteNameがない場合正常終了する", func(t *testing.T) {
	//	// 準備
	//	interactor := NewSaveDailyDataUseCase(
	//		&repositories.NoteRepositoryMock{},
	//		&repositories.DailyDataRepositoryMock{},
	//	)
	//	dailyData := DailyData{
	//		Date: "2021-02-02",
	//		Notes: []string{
	//			repositories.EXIST_NOTE_NAME,
	//		},
	//	}
	//	// 実行
	//	err := interactor.Handle(dailyData)
	//	// 検証
	//	assert.NoError(t, err)
	//})
}
