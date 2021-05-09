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
		err := rep.Save(newEntity)
		// 検証
		assert.NoError(t, err)
	})

	t.Run("一部が新規Noteだった場合、一部のみが新規追加される", func(t *testing.T) {
		////// 準備
		preEntity := entities.NewDailyDataEntity(
			time.Date(2021, 1, 2, 0, 0, 0, 0, time.Local),
			[]string{
				"Note2-1",
				"Note2-2",
			},
		)
		err := rep.Save(preEntity)
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
	})
}
