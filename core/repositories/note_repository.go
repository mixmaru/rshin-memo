package repositories

import (
	"github.com/mixmaru/rshin-memo/core/entities"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type NoteRepositoryInterface interface {
	GetByNoteName(noteName string) (*entities.NoteEntity, error)
	GetAllNotesOnlyName() ([]*entities.NoteEntity, error)
	Save(entity *entities.NoteEntity) error
	GetBySearchText(text string) ([]*entities.NoteEntity, error)
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

// メモリ軽減のため、内容は取得しない
func (n *NoteRepository) GetAllNotesOnlyName() ([]*entities.NoteEntity, error) {
	retEntities := []*entities.NoteEntity{}

	files, err := ioutil.ReadDir(n.dirPath)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	r := regexp.MustCompile(`\.txt$`)
	for _, file := range files {
		if !file.IsDir() {
			if r.MatchString(file.Name()) {
				noteName := strings.Replace(file.Name(), ".txt", "", 1)
				retEntities = append(retEntities, entities.NewNoteEntity(noteName, ""))
			}
		}
	}

	return retEntities, nil
}

func (n *NoteRepository) Save(entity *entities.NoteEntity) error {
	// NoteFileが存在しなければ新規作成、あれば上書きする
	err := ioutil.WriteFile(filepath.Join(n.dirPath, entity.Name()+".txt"), []byte(entity.Text()), 0644)
	if err != nil {
		return errors.Wrapf(err, "NoteFile create error. entity: %+v", entity)
	}
	return nil
}
