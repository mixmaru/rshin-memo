package repositories

import (
	"github.com/mixmaru/rshin-memo/core/entities"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestNoteRepository_GetByNoteName(t *testing.T) {
	t.Skip()
}

func TestNoteRepository_Save(t *testing.T) {
	t.Run("Noteが新規追加だったとき、fileが新規追加される", func(t *testing.T) {
		////// 準備
		thisDir, err := os.Getwd()
		if err != nil {
			t.Fatalf("%+v", err)
		}
		targetDirPath := filepath.Join(thisDir, "test")

		// 既存ファイルがあれば事前に削除
		if _, err := os.Stat(filepath.Join(targetDirPath, "TestNoteRepository_Save.txt")); err == nil {
			err = os.Remove(filepath.Join(targetDirPath, "TestNoteRepository_Save.txt"))
			if err != nil {
				t.Fatalf("%+v", err)
			}
		}

		rep := NewNoteRepository(targetDirPath)

		////// 実行
		newNote := entities.NewNoteEntity("TestNoteRepository_Save.txt", "TestNoteRepository_Saveの内容")
		err = rep.Save(newNote)

		////// テスト
		assert.NoError(t, err)
		bytes, err := ioutil.ReadFile(filepath.Join(targetDirPath, "TestNoteRepository_Save.txt"))
		assert.NoError(t, err)
		assert.Equal(t, "TestNoteRepository_Saveの内容", string(bytes))
	})

	t.Run("Noteが既に存在するとき、上書き新規追加される", func(t *testing.T) {
		////// 準備
		thisDir, err := os.Getwd()
		if err != nil {
			t.Fatalf("%+v", err)
		}
		targetDirPath := filepath.Join(thisDir, "test")
		rep := NewNoteRepository(targetDirPath)
		// 先に既存ファイルを作っておく
		oldNoteEntity := entities.NewNoteEntity("TestNoteRepository_Save2.txt", "先に存在する内容")
		err = rep.Save(oldNoteEntity)
		assert.NoError(t, err)

		////// 上書き実行
		newNote := entities.NewNoteEntity("TestNoteRepository_Save2.txt", "上書きされた内容")
		err = rep.Save(newNote)

		////// テスト
		assert.NoError(t, err)
		bytes, err := ioutil.ReadFile(filepath.Join(targetDirPath, "TestNoteRepository_Save2.txt"))
		assert.NoError(t, err)
		assert.Equal(t, "上書きされた内容", string(bytes))
	})
}
