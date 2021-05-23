package repositories

import (
	"github.com/mixmaru/rshin-memo/core/entities"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

type NoteRepositoryInterface interface {
	GetByNoteName(noteName string) (*entities.NoteEntity, error)
	GetAllNotes() ([]*entities.NoteEntity, error)
	Save(entity *entities.NoteEntity) error
}

type NoteRepository struct {
	dirPath string // Note保存先DirPath
}

func NewNoteRepository(dirPath string) *NoteRepository {
	return &NoteRepository{dirPath: dirPath}
}

func (n *NoteRepository) GetByNoteName(noteName string) (*entities.NoteEntity, error) {
	noteFilePath := filepath.Join(n.dirPath, noteName+".txt")
	_, err := os.Stat(noteFilePath)
	if err != nil {
		return nil, nil
	}

	bytes, err := ioutil.ReadFile(noteFilePath)
	if err != nil {
		return nil, errors.Wrapf(err, "noteファイル作成失敗. %v", noteFilePath)
	}
	retEntity := entities.NewNoteEntity(noteName, string(bytes))
	return retEntity, nil
}

func (n *NoteRepository) Save(entity *entities.NoteEntity) error {
	// NoteFileが存在しなければ新規作成、あれば上書きする
	err := ioutil.WriteFile(filepath.Join(n.dirPath, entity.Name()+".txt"), []byte(entity.Text()), 0644)
	if err != nil {
		return errors.Wrapf(err, "NoteFile create error. entity: %+v", entity)
	}
	return nil
}
