package usecases

import (
	"github.com/mixmaru/rshin-memo/core/repositories"
)

type CreateNoteUseCaseInterface interface {
	Handle(dailyData DailyData) error
}

type CreateNoteUseCaseInteractor struct {
	noteRepository repositories.NoteRepositoryInterface
}

func NewCreateNoteUseCaseInteractor(noteRepository repositories.NoteRepositoryInterface) *CreateNoteUseCaseInteractor {
	return &CreateNoteUseCaseInteractor{noteRepository: noteRepository}
}

func (c *CreateNoteUseCaseInteractor) Handle(dailyData DailyData) error {
	// dailyListの作成
	// Noteの作成（存在しないNoteであれば作成する）
	//noteEntity := entities.NewNoteEntity(
	//	noteName,
	//	"",
	//)
	//err := c.noteRepository.Save(noteEntity)
	//if err != nil {
	//	return err
	//}
	return nil
}
