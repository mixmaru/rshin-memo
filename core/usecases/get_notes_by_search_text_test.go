package usecases

import (
	"github.com/mixmaru/rshin-memo/core/repositories"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetNotesBySearchTextUseCase_Handle(t *testing.T) {
	////// 準備
	useCase := NewGetNotesBySearchTextUseCase(&repositories.NoteRepositoryMock{})

	////// 実行
	result, err := useCase.Handle("検索ワード")
	assert.NoError(t, err)

	////// 検証
	expected := []string{
		"検索されたNote_A",
		"検索されたNote_B",
	}
	assert.EqualValues(t, expected, result)
}
