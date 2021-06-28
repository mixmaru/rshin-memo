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
	t.Run("指定のNoteファイルが存在しない場合、nullが返る", func(t *testing.T) {
		////// 準備
		thisDir, err := os.Getwd()
		if err != nil {
			t.Fatalf("%+v", err)
		}
		targetDirPath := filepath.Join(thisDir, "test")

		// 既存ファイルがあれば事前に削除
		if _, err := os.Stat(filepath.Join(targetDirPath, "not_exist_note")); err == nil {
			err = os.Remove(filepath.Join(targetDirPath, "not_exist_note"))
			if err != nil {
				t.Fatalf("%+v", err)
			}
		}

		////// 実行
		rep := NewNoteRepository(targetDirPath)
		note, err := rep.GetByNoteName("not_exist_note")
		assert.NoError(t, err)

		////// 検証
		assert.Nil(t, note)
	})

	t.Run("指定のNoteが存在する場合、noteデータが返る", func(t *testing.T) {
		////// 準備
		thisDir, err := os.Getwd()
		if err != nil {
			t.Fatalf("%+v", err)
		}
		targetDirPath := filepath.Join(thisDir, "test")
		rep := NewNoteRepository(targetDirPath)

		// 事前に登録する
		existNote := entities.NewNoteEntity("exist_note", "存在するノート")
		err = rep.Save(existNote)
		assert.NoError(t, err)

		////// 実行
		note, err := rep.GetByNoteName("exist_note")
		assert.NoError(t, err)

		////// 検証
		assert.EqualValues(t, existNote, note)
	})
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
		newNote := entities.NewNoteEntity("TestNoteRepository_Save", "TestNoteRepository_Saveの内容")
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
		oldNoteEntity := entities.NewNoteEntity("TestNoteRepository_Save2", "先に存在する内容")
		err = rep.Save(oldNoteEntity)
		assert.NoError(t, err)

		////// 上書き実行
		newNote := entities.NewNoteEntity("TestNoteRepository_Save2", "上書きされた内容")
		err = rep.Save(newNote)

		////// テスト
		assert.NoError(t, err)
		bytes, err := ioutil.ReadFile(filepath.Join(targetDirPath, "TestNoteRepository_Save2.txt"))
		assert.NoError(t, err)
		assert.Equal(t, "上書きされた内容", string(bytes))
	})
}

func TestNoteRepository_GetAllNotes(t *testing.T) {
	thisDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	targetDirPath := filepath.Join(thisDir, "test", "TestNoteRepository_GetAllNotes")

	t.Run("あるとき", func(t *testing.T) {
		////// 準備
		if _, err := os.Stat(targetDirPath); err == nil {
			// あれば一旦削除して作り直す
			err = os.RemoveAll(targetDirPath)
		}
		err = os.Mkdir(targetDirPath, 0777)
		if err != nil {
			t.Fatalf("%+v", err)
		}

		repository := NewNoteRepository(targetDirPath)
		// なんこかNoteを追加しておく
		newNoteEntityA := entities.NewNoteEntity("noteA", "noteAの内容")
		newNoteEntityB := entities.NewNoteEntity("noteB", "noteBの内容")
		err = repository.Save(newNoteEntityA)
		if err != nil {
			assert.Failf(t, "%+v", err.Error())
		}
		err = repository.Save(newNoteEntityB)
		if err != nil {
			assert.Failf(t, "%+v", err.Error())
		}

		result, err := repository.GetAllNotesOnlyName()
		assert.NoError(t, err)
		expected := []*entities.NoteEntity{
			entities.NewNoteEntity("noteA", ""),
			entities.NewNoteEntity("noteB", ""),
		}

		assert.EqualValues(t, expected, result)
	})

	t.Run("ないとき", func(t *testing.T) {
		////// 準備
		if _, err := os.Stat(targetDirPath); err == nil {
			// あれば一旦削除して作り直す
			err = os.RemoveAll(targetDirPath)
		}
		err = os.Mkdir(targetDirPath, 0777)
		if err != nil {
			t.Fatalf("%+v", err)
		}

		repository := NewNoteRepository(targetDirPath)
		result, err := repository.GetAllNotesOnlyName()
		assert.NoError(t, err)

		assert.Len(t, result, 0)
	})
}

func TestNoteRepository_GetBySearchText(t *testing.T) {
	thisDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	targetDirPath := filepath.Join(thisDir, "test", "TestNoteRepository_GetBySearchText")

	t.Run("検索対象があったとき", func(t *testing.T) {
		////// 準備
		if _, err := os.Stat(targetDirPath); err == nil {
			// あれば一旦削除して作り直す
			err = os.RemoveAll(targetDirPath)
		}
		err = os.Mkdir(targetDirPath, 0777)
		if err != nil {
			t.Fatalf("%+v", err)
		}

		repository := NewNoteRepository(targetDirPath)
		// なんこかNoteを追加しておく
		newNoteEntityA := entities.NewNoteEntity("noteA", "noteAの内容")
		newNoteEntityB := entities.NewNoteEntity("noteB", "noteBの内容")
		newNoteEntityC := entities.NewNoteEntity("ヒットしないやつ", "noteCの内容")
		err = repository.Save(newNoteEntityA)
		if err != nil {
			assert.Failf(t, "%+v", err.Error())
		}
		err = repository.Save(newNoteEntityB)
		if err != nil {
			assert.Failf(t, "%+v", err.Error())
		}
		err = repository.Save(newNoteEntityC)
		if err != nil {
			assert.Failf(t, "%+v", err.Error())
		}

		////// 実行
		results, err := repository.GetBySearchText("note")
		assert.NoError(t, err)

		////// 検証
		assert.Len(t, results, 2)
		assert.EqualValues(t, newNoteEntityA.Name(), results[0].Name())
		assert.EqualValues(t, newNoteEntityB.Name(), results[1].Name())
	})

	t.Run("検索対象がなかったとき", func(t *testing.T) {
		////// 準備
		if _, err := os.Stat(targetDirPath); err == nil {
			// あれば一旦削除して作り直す
			err = os.RemoveAll(targetDirPath)
		}
		err = os.Mkdir(targetDirPath, 0777)
		if err != nil {
			t.Fatalf("%+v", err)
		}

		repository := NewNoteRepository(targetDirPath)
		// なんこかNoteを追加しておく
		newNoteEntityA := entities.NewNoteEntity("noteA", "noteAの内容")
		newNoteEntityB := entities.NewNoteEntity("noteB", "noteBの内容")
		err = repository.Save(newNoteEntityA)
		if err != nil {
			assert.Failf(t, "%+v", err.Error())
		}
		err = repository.Save(newNoteEntityB)
		if err != nil {
			assert.Failf(t, "%+v", err.Error())
		}

		////// 実行
		results, err := repository.GetBySearchText("ヒットせず")
		assert.NoError(t, err)

		////// 検証
		assert.Len(t, results, 0)
	})
}
