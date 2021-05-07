package repositories

import "github.com/mixmaru/rshin-memo/core/entities"

type NoteRepositoryInterface interface {
	GetByNoteName(noteName string) (*entities.NoteEntity, error)
	Save(entity *entities.NoteEntity) error
}
