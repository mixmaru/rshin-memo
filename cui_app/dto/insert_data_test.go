package dto

import (
	"github.com/mixmaru/rshin-memo/core/usecases"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInsertData_GenerateNewDailyData(t *testing.T) {
	t.Run("間に突っ込む", func(t *testing.T) {
		insertData := InsertData{}
		insertData.TargetDailyData = []usecases.DailyData{
			{
				Date: "2021-01-01",
				Notes: []string{
					"note1",
					"note2",
					"note3",
				},
			},
		}
		insertData.DateStr = "2021-01-01"
		insertData.InsertNum = 1
		insertData.NoteName = "NewNote"
		got, err := insertData.GenerateNewDailyData()
		assert.NoError(t, err)
		expect := usecases.DailyData{
			Date: "2021-01-01",
			Notes: []string{
				"note1",
				"NewNote",
				"note2",
				"note3",
			},
		}

		assert.Equal(t, expect, got)
	})

	t.Run("先頭に突っ込む", func(t *testing.T) {
		insertData := InsertData{}
		insertData.TargetDailyData = []usecases.DailyData{
			{
				Date: "2021-01-01",
				Notes: []string{
					"note1",
					"note2",
					"note3",
				},
			},
		}
		insertData.DateStr = "2021-01-01"
		insertData.InsertNum = 0
		insertData.NoteName = "NewNote"
		got, err := insertData.GenerateNewDailyData()
		assert.NoError(t, err)
		expect := usecases.DailyData{
			Date: "2021-01-01",
			Notes: []string{
				"NewNote",
				"note1",
				"note2",
				"note3",
			},
		}

		assert.Equal(t, expect, got)
	})

	t.Run("末に突っ込む", func(t *testing.T) {
		insertData := InsertData{}
		insertData.TargetDailyData = []usecases.DailyData{
			{
				Date: "2021-01-01",
				Notes: []string{
					"note1",
					"note2",
					"note3",
				},
			},
		}
		insertData.DateStr = "2021-01-01"
		insertData.InsertNum = 3
		insertData.NoteName = "NewNote"
		got, err := insertData.GenerateNewDailyData()
		assert.NoError(t, err)
		expect := usecases.DailyData{
			Date: "2021-01-01",
			Notes: []string{
				"note1",
				"note2",
				"note3",
				"NewNote",
			},
		}
		assert.Equal(t, expect, got)
	})

	t.Run("範囲外は末に追加される", func(t *testing.T) {
		insertData := InsertData{}
		insertData.TargetDailyData = []usecases.DailyData{
			{
				Date: "2021-01-01",
				Notes: []string{
					"note1",
					"note2",
					"note3",
				},
			},
		}
		insertData.DateStr = "2021-01-01"
		insertData.InsertNum = 100
		insertData.NoteName = "NewNote"
		got, err := insertData.GenerateNewDailyData()
		assert.NoError(t, err)
		expect := usecases.DailyData{
			Date: "2021-01-01",
			Notes: []string{
				"note1",
				"note2",
				"note3",
				"NewNote",
			},
		}
		assert.Equal(t, expect, got)
	})
}

func createTestDailyList() []usecases.DailyData {
	dailyList := []usecases.DailyData{
		{
			Date: "2021-03-30",
			Notes: []string{
				"a",
				"b",
			},
		},
		{
			Date: "2021-03-29",
			Notes: []string{
				"a",
				"b",
			},
		},
		{
			Date: "2021-03-27",
			Notes: []string{
				"a",
				"b",
			},
		},
	}
	return dailyList
}

func Test_generateNewDailyData(t *testing.T) {
	t.Run("先頭に新規追加", func(t *testing.T) {
		dailyList := createTestDailyList()
		result, err := generateNewDailyData(dailyList, "newNote", "2021-04-01", 0)
		assert.NoError(t, err)
		expected := usecases.DailyData{
			Date: "2021-04-01",
			Notes: []string{
				"newNote",
			},
		}
		assert.Equal(t, expected, result)
	})

	t.Run("日付の先頭位置への挿入", func(t *testing.T) {
		dailyList := createTestDailyList()
		result, err := generateNewDailyData(dailyList, "newNote", "2021-03-30", 0)
		assert.NoError(t, err)
		expected := usecases.DailyData{
			Date: "2021-03-30",
			Notes: []string{
				"newNote",
				"a",
				"b",
			},
		}
		assert.Equal(t, expected, result)

		dailyList = createTestDailyList()
		result, err = generateNewDailyData(dailyList, "newNote", "2021-03-29", 2)
		assert.NoError(t, err)
		expected = usecases.DailyData{
			Date: "2021-03-29",
			Notes: []string{
				"newNote",
				"a",
				"b",
			},
		}
		assert.Equal(t, expected, result)
	})

	t.Run("日付の中間位置への挿入", func(t *testing.T) {
		dailyList := createTestDailyList()
		result, err := generateNewDailyData(dailyList, "newNote", "2021-03-30", 1)
		assert.NoError(t, err)
		expected := usecases.DailyData{
			Date: "2021-03-30",
			Notes: []string{
				"a",
				"newNote",
				"b",
			},
		}
		assert.Equal(t, expected, result)

		dailyList = createTestDailyList()
		result, err = generateNewDailyData(dailyList, "newNote", "2021-03-29", 3)
		assert.NoError(t, err)
		expected = usecases.DailyData{
			Date: "2021-03-29",
			Notes: []string{
				"a",
				"newNote",
				"b",
			},
		}
		assert.Equal(t, expected, result)
	})

	t.Run("日付の末尾位置への挿入", func(t *testing.T) {
		dailyList := createTestDailyList()
		result, err := generateNewDailyData(dailyList, "newNote", "2021-03-30", 2)
		assert.NoError(t, err)
		expected := usecases.DailyData{
			Date: "2021-03-30",
			Notes: []string{
				"a",
				"b",
				"newNote",
			},
		}
		assert.Equal(t, expected, result)

		dailyList = createTestDailyList()
		result, err = generateNewDailyData(dailyList, "newNote", "2021-03-29", 4)
		assert.NoError(t, err)
		expected = usecases.DailyData{
			Date: "2021-03-29",
			Notes: []string{
				"a",
				"b",
				"newNote",
			},
		}
		assert.Equal(t, expected, result)
	})

	t.Run("日付の先頭位置に別日の挿入", func(t *testing.T) {
		dailyList := createTestDailyList()
		result, err := generateNewDailyData(dailyList, "newNote", "2021-04-01", 0)
		assert.NoError(t, err)
		expected := usecases.DailyData{
			Date: "2021-04-01",
			Notes: []string{
				"newNote",
			},
		}
		assert.Equal(t, expected, result)

		dailyList = createTestDailyList()
		result, err = generateNewDailyData(dailyList, "newNote", "2021-03-28", 4)
		assert.NoError(t, err)
		expected = usecases.DailyData{
			Date: "2021-03-28",
			Notes: []string{
				"newNote",
			},
		}
		assert.Equal(t, expected, result)
	})

	t.Run("日付の末尾位置に別日の挿入", func(t *testing.T) {
		dailyList := createTestDailyList()
		result, err := generateNewDailyData(dailyList, "newNote", "2021-03-26", 6)
		assert.NoError(t, err)
		expected := usecases.DailyData{
			Date: "2021-03-26",
			Notes: []string{
				"newNote",
			},
		}
		assert.Equal(t, expected, result)
	})

	t.Run("異なる日付の位置に挿入しようとするとエラー", func(t *testing.T) {
		dailyList := createTestDailyList()
		_, err := generateNewDailyData(dailyList, "newNote", "2021-01-27", 1)
		assert.Error(t, err)
	})
}
