package usecases

import (
	"github.com/mixmaru/rshin-memo/core/repositories"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateNoteUseCaseInteractor_Handle(t *testing.T) {
	// 準備
	interactor := NewCreateNoteUseCaseInteractor(&repositories.NoteRepositoryMock{})
	// 実行
	err := interactor.Handle("新規NOTE", time.Date(2021, 2, 1, 0, 0, 0, 0, time.Local))
	// 検証
	assert.NoError(t, err)
}
