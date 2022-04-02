package usecases

import (
	"github.com/golang/mock/gomock"
	"github.com/mixmaru/rshin-memo/core/entities"
	mock_repositories "github.com/mixmaru/rshin-memo/core/repositories/mock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateNoteFromParamsUseCaseInteractor_Handle(t *testing.T) {
	t.Run("存在しないNoteNameが1つ含まれている場合正常終了する", func(t *testing.T) {
		retData := []*entities.DailyDataEntity{
			entities.NewDailyDataEntity(
				time.Date(2021, 1, 10, 0, 0, 0, 0, time.Local),
				[]string{
					"noteC",
					"noteB",
					"noteA",
				},
			),
			entities.NewDailyDataEntity(
				time.Date(2021, 1, 5, 0, 0, 0, 0, time.Local),
				[]string{
					"noteC",
					"noteB",
					"noteA",
				},
			),
		}

		saveingEntity := entities.NewDailyDataEntity(
			time.Date(2021, 1, 8, 0, 0, 0, 0, time.Local),
			[]string{
				"newMemoName",
			},
		)

		noteEntity := entities.NewNoteEntity("newMemoName", "new_memo_内容")

		// repository mock
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		dailyDataRep := mock_repositories.NewMockDailyDataRepositoryInterface(ctrl)
		dailyDataRep.EXPECT().Get().Return(retData, nil)
		dailyDataRep.EXPECT().Save(saveingEntity)

		noteRep := mock_repositories.NewMockNoteRepositoryInterface(ctrl)
		noteRep.EXPECT().Save(noteEntity).Return(nil)

		// create test data
		interactor := NewSaveDailyDataFromParamsUseCase(
			noteRep,
			dailyDataRep,
		)

		// 実行
		err := interactor.Handle(
			time.Date(2021, 1, 5, 0, 0, 0, 0, time.Local),
			"noteC",
			time.Date(2021, 1, 8, 0, 0, 0, 0, time.Local),
			"newMemoName",
			"new_memo_内容",
			INSERT_NEWER_MODE,
		)
		// 検証,意図通りにmockを叩いているか？
		assert.NoError(t, err)
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
