package usecases

import (
	"github.com/mixmaru/rshin-memo/core/entities"
	"github.com/mixmaru/rshin-memo/core/repositories"
)

type SaveDailyDataUseCaseInterface interface {
	Handle(dailyData DailyData) error
}

type SaveDailyDataUseCaseInteractor struct {
	noteRepository repositories.NoteRepositoryInterface
}

func NewSaveDailyDataUseCaseInteractor(noteRepository repositories.NoteRepositoryInterface) *SaveDailyDataUseCaseInteractor {
	return &SaveDailyDataUseCaseInteractor{noteRepository: noteRepository}
}

func (c *SaveDailyDataUseCaseInteractor) Handle(dailyData DailyData) error {
	// dailyListの作成

	// Noteの作成（存在しないNoteであれば作成する）
	for _, noteName := range dailyData.Notes {
		noteEntity, err := c.noteRepository.GetByNoteName(noteName)
		if err != nil {
			return err
		}
		if noteEntity == nil {
			newNote := entities.NewNoteEntity(noteName, "")
			err := c.noteRepository.Save(newNote)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
