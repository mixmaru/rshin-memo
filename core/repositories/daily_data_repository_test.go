package repositories

import (
	"github.com/mixmaru/rshin-memo/core/entities"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func getJsonFilePathForTest() (string, error) {
	thisDirPath, err := os.Getwd()
	if err != nil {
		return "", err
	}
	jsonFilePath := filepath.Join(thisDirPath, "test", "daily_data.json")
	return jsonFilePath, nil
}

func TestDailyDataRepository_Save(t *testing.T) {
	filePath, err := getJsonFilePathForTest()
	if err != nil {
		t.Fatalf("エラー %v", err)
	}
	rep := NewDailyDataRepository(filePath)

	t.Run("すべて新規Noteだった場合、全体が新規追加される", func(t *testing.T) {
		// 一旦削除
		if _, err := os.Stat(filePath); err == nil {
			err := os.Remove(filePath)
			if err != nil {
				assert.Fail(t, err.Error())
			}
		}

		// 準備
		newEntity := entities.NewDailyDataEntity(
			time.Date(2021, 1, 2, 0, 0, 0, 0, time.Local),
			[]string{
				"Note1-1",
				"Note1-2",
				"Note1-3",
				"Note1-4",
				"Note1-5",
			},
		)

		// 実行
		err = rep.Save(newEntity)
		// 検証
		assert.NoError(t, err)
		loadedEntities, err := rep.Get()
		assert.NoError(t, err)
		assert.Len(t, loadedEntities, 1)
		assert.EqualValues(t, newEntity, loadedEntities[0])
	})

	t.Run("一部が新規Noteだった場合、一部のみが新規追加される", func(t *testing.T) {
		////// 準備
		// 一旦削除
		if _, err := os.Stat(filePath); err == nil {
			err := os.Remove(filePath)
			if err != nil {
				assert.Fail(t, err.Error())
			}
		}
		preEntity := entities.NewDailyDataEntity(
			time.Date(2021, 1, 2, 0, 0, 0, 0, time.Local),
			[]string{
				"Note2-1",
				"Note2-2",
			},
		)
		err = rep.Save(preEntity)
		assert.NoError(t, err)
		// 一部のみ新規Noteのentityを用意
		newEntity := entities.NewDailyDataEntity(
			time.Date(2021, 1, 2, 0, 0, 0, 0, time.Local),
			[]string{
				"Note2-1",
				"Note2-3",
				"Note2-2",
			},
		)

		////// 実行
		err = rep.Save(newEntity)

		////// 検証
		assert.NoError(t, err)
		loadedEntities, err := rep.Get()
		assert.NoError(t, err)
		assert.Len(t, loadedEntities, 1)
		assert.EqualValues(t, newEntity, loadedEntities[0])
	})
}

func TestDailyDataRepository_Get(t *testing.T) {
	t.Run("既存データがある場合", func(t *testing.T) {
		////// 準備
		// 既存jsonファイル削除
		jsonfilepath, err := getJsonFilePathForTest()
		if err != nil {
			assert.Fail(t, err.Error())
		}
		if _, err := os.Stat(jsonfilepath); err == nil {
			err = os.Remove(jsonfilepath)
			if err != nil {
				assert.Fail(t, err.Error())
			}
		}
		rep := NewDailyDataRepository(jsonfilepath)
		// 事前データ登録
		entity1 := entities.NewDailyDataEntity(
			time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local),
			[]string{
				"note_A",
				"note_B",
				"note_C",
			},
		)
		err = rep.Save(entity1)
		if err != nil {
			assert.Fail(t, err.Error())
		}
		entity2 := entities.NewDailyDataEntity(
			time.Date(2021, 1, 2, 0, 0, 0, 0, time.Local),
			[]string{
				"note_A",
				"note_B",
				"note_C",
			},
		)
		err = rep.Save(entity2)
		assert.NoError(t, err)

		////// 実行
		dailyDataEntities, err := rep.Get()
		assert.NoError(t, err)

		////// 検証
		// 日付は降順で取得される
		assert.EqualValues(t, entity2, dailyDataEntities[0])
		assert.EqualValues(t, entity1, dailyDataEntities[1])

	})

	t.Run("既存データがない場合", func(t *testing.T) {
		////// 準備
		// 既存jsonファイル削除
		jsonfilepath, err := getJsonFilePathForTest()
		if err != nil {
			assert.Fail(t, err.Error())
		}
		if _, err := os.Stat(jsonfilepath); err == nil {
			err = os.Remove(jsonfilepath)
			if err != nil {
				assert.Fail(t, err.Error())
			}
		}

		////// 実行
		rep := NewDailyDataRepository(jsonfilepath)
		dailyDataEntities, err := rep.Get()
		assert.NoError(t, err)

		////// 検証
		// 日付は降順で取得される
		assert.Len(t, dailyDataEntities, 0)
	})
}
