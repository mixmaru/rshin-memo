package usecases

import (
	"github.com/mixmaru/rshin-memo/core/entities"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetNoteUseCaseInteractor_Handle(t *testing.T) {
	t.Run("対象データがある時、noteTextが帰ってくる", func(t *testing.T) {
		// 準備
		interactor := NewGetNoteUseCaseInteractor(&NoteRepositoryMock_DataExist{})
		// 実行
		noteText, notExist, err := interactor.Handle("なにかしらのNote")
		assert.NoError(t, err)
		// 検証
		assert.False(t, notExist)
		assert.Equal(t, `なんかしらの内容。
なんかしらの内容2行目`, noteText)
	})

	t.Run("対象データがない時、noteTextはカラ文字、notExistにtrueが入る", func(t *testing.T) {
		// 準備
		interactor := NewGetNoteUseCaseInteractor(&NoteRepositoryMock_DataNotExist{})
		// 実行
		noteText, notExist, err := interactor.Handle("なにかしらのNote")
		assert.NoError(t, err)
		// 検証
		assert.True(t, notExist)
		assert.Equal(t, "", noteText)
	})
}

type NoteRepositoryMock_DataExist struct{}

func (n *NoteRepositoryMock_DataExist) GetByNoteName(noteName string) (*entities.NoteEntity, error) {
	entity := entities.NewNoteEntity(
		time.Date(2021, 2, 1, 0, 0, 0, 0, time.Local),
		"なんかしらのNote",
		`なんかしらの内容。
なんかしらの内容2行目`,
	)
	return entity, nil
}

type NoteRepositoryMock_DataNotExist struct{}

func (n *NoteRepositoryMock_DataNotExist) GetByNoteName(noteName string) (*entities.NoteEntity, error) {
	return nil, nil
}
