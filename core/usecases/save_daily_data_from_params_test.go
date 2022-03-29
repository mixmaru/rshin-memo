package usecases

import (
	"github.com/golang/mock/gomock"
	"github.com/mixmaru/rshin-memo/core/entities"
	"github.com/mixmaru/rshin-memo/core/repositories"
	mock_repositories "github.com/mixmaru/rshin-memo/core/repositories/mock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateNoteFromParamsUseCaseInteractor_Handle(t *testing.T) {
	t.Run("存在しないNoteNameが1つ含まれている場合正常終了する", func(t *testing.T) {
		//ctrl := gomock.NewController(t)
		//
		//// Assert that Bar() is invoked.
		//defer ctrl.Finish()
		//
		//m := NewMockFoo(ctrl)
		//
		//// Asserts that the first and only call to Bar() is passed 99.
		//// Anything else will fail.
		//m.
		//	EXPECT().
		//	Bar(gomock.Eq(99)).
		//	Return(101)
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

		// repository mock
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		m := mock_repositories.NewMockDailyDataRepositoryInterface(ctrl)
		m.EXPECT().Get().Return(retData, nil)

		// create test data
		interactor := NewSaveDailyDataFromParamsUseCase(
			&repositories.NoteRepositoryMock{},
			m,
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
