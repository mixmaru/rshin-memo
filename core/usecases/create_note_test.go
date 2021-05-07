package usecases

import (
	"github.com/mixmaru/rshin-memo/core/repositories"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateNoteUseCaseInteractor_Handle(t *testing.T) {
	// 準備
	interactor := NewCreateNoteUseCaseInteractor(&repositories.NoteRepositoryMock{})
	dailyData := DailyData{
		Date: "2021-02-02",
		Notes: []string{
			"Note1",
			"Note2",
		},
	}
	// 実行
	err := interactor.Handle(dailyData)
	// 検証
	assert.NoError(t, err)
}
