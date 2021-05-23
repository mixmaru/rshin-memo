package usecases

import (
	"github.com/mixmaru/rshin-memo/core/repositories"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetAllNotesUseCase_Handle(t *testing.T) {
	t.Run("", func(t *testing.T) {
		useCase := NewGetAllNotesUseCase(&repositories.NoteRepositoryMock{})
		result, err := useCase.Handle()
		assert.NoError(t, err)
		assert.Equal(t, []string{"NoteA", "NoteB"}, result)
	})
}
