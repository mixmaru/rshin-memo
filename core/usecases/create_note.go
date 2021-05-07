package usecases

import (
	"github.com/mixmaru/rshin-memo/core/entities"
	"github.com/mixmaru/rshin-memo/core/repositories"
	"time"
)

type CreateNoteUseCaseInterface interface {
	Handle(noteName string, date time.Time) error
}

type CreateNoteUseCaseInteractor struct {
	noteRepository repositories.NoteRepositoryInterface
}

func NewCreateNoteUseCaseInteractor(noteRepository repositories.NoteRepositoryInterface) *CreateNoteUseCaseInteractor {
	return &CreateNoteUseCaseInteractor{noteRepository: noteRepository}
}

func (c *CreateNoteUseCaseInteractor) Handle(noteName string, date time.Time) error {
	// dailyListの作成
	// Noteの作成
	noteEntity := entities.NewNoteEntity(
		noteName,
		"",
	)
	err := c.noteRepository.Save(noteEntity)
	if err != nil {
		return err
	}
	return nil
}
