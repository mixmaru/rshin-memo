package usecases

import (
	"github.com/mixmaru/rshin-memo/core/entities"
	"github.com/mixmaru/rshin-memo/core/repositories"
	"github.com/pkg/errors"
	"time"
)

type SaveDailyDataUseCase struct {
	noteRepository      repositories.NoteRepositoryInterface
	dailyDataRepository repositories.DailyDataRepositoryInterface
}

func NewSaveDailyDataUseCase(noteRepository repositories.NoteRepositoryInterface, dailyDataRepository repositories.DailyDataRepositoryInterface) *SaveDailyDataUseCase {
	return &SaveDailyDataUseCase{noteRepository: noteRepository, dailyDataRepository: dailyDataRepository}
}

func (c *SaveDailyDataUseCase) Handle(dailyData DailyData) error {
	// dailyListの作成
	date, err := time.ParseInLocation("2006-01-02", dailyData.Date, time.Local)
	if err != nil {
		return errors.Wrapf(err, "日付Parseの失敗。%+v", dailyData)
	}
	newDailyDataEntity := entities.NewDailyDataEntity(
		date,
		dailyData.Notes,
	)
	err = c.dailyDataRepository.Save(newDailyDataEntity)
	if err != nil {
		return err
	}

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
