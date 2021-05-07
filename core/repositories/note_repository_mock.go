package repositories

import (
	"github.com/mixmaru/rshin-memo/core/entities"
)

type NoteRepositoryMock struct{}

func (n *NoteRepositoryMock) Save(entity *entities.NoteEntity) error {
	return nil
}

const EXIST_NOTE_NAME = "exist_note_name"
const NOT_EXIST_NOTE_NAME = "not_exist_note_name"

func (n *NoteRepositoryMock) GetByNoteName(noteName string) (*entities.NoteEntity, error) {
	switch noteName {
	case EXIST_NOTE_NAME:
		entity := entities.NewNoteEntity(
			"なんかしらのNote",
			`なんかしらの内容。
なんかしらの内容2行目`,
		)
		return entity, nil
	case NOT_EXIST_NOTE_NAME:
		return nil, nil
	default:
		return nil, nil
	}
}
