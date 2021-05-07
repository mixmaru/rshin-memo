package usecases

import (
	"github.com/mixmaru/rshin-memo/core/repositories"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetNoteUseCaseInteractor_Handle(t *testing.T) {
	t.Run("対象データがある時、noteTextが帰ってくる", func(t *testing.T) {
		// 準備
		interactor := NewGetNoteUseCaseInteractor(&repositories.NoteRepositoryMock{})
		// 実行
		noteText, notExist, err := interactor.Handle(repositories.EXIST_NOTE_NAME)
		assert.NoError(t, err)
		// 検証
		assert.False(t, notExist)
		assert.Equal(t, `なんかしらの内容。
なんかしらの内容2行目`, noteText)
	})

	t.Run("対象データがない時、noteTextはカラ文字、notExistにtrueが入る", func(t *testing.T) {
		// 準備
		interactor := NewGetNoteUseCaseInteractor(&repositories.NoteRepositoryMock{})
		// 実行
		noteText, notExist, err := interactor.Handle(repositories.NOT_EXIST_NOTE_NAME)
		assert.NoError(t, err)
		// 検証
		assert.True(t, notExist)
		assert.Equal(t, "", noteText)
	})
}
