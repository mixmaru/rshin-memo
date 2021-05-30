package dto

import (
	"github.com/mixmaru/rshin-memo/core/usecases"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInsertData_GenerateNewDailyData(t *testing.T) {
	t.Run("間に突っ込む", func(t *testing.T) {
		insertData := InsertData{}
		insertData.SetTargetDailyData(usecases.DailyData{
			Date: "2021-01-01",
			Notes: []string{
				"note1",
				"note2",
				"note3",
			},
		})
		err := insertData.SetInsertNum(1)
		assert.NoError(t, err)
		insertData.SetNoteName("NewNote")
		got := insertData.GenerateNewDailyData()
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
		insertData.SetTargetDailyData(usecases.DailyData{
			Date: "2021-01-01",
			Notes: []string{
				"note1",
				"note2",
				"note3",
			},
		})
		err := insertData.SetInsertNum(0)
		assert.NoError(t, err)
		insertData.SetNoteName("NewNote")
		got := insertData.GenerateNewDailyData()
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
		insertData.SetTargetDailyData(usecases.DailyData{
			Date: "2021-01-01",
			Notes: []string{
				"note1",
				"note2",
				"note3",
			},
		})
		err := insertData.SetInsertNum(3)
		assert.NoError(t, err)
		insertData.SetNoteName("NewNote")
		got := insertData.GenerateNewDailyData()
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

	t.Run("範囲外は指定できない", func(t *testing.T) {
		insertData := InsertData{}
		insertData.SetTargetDailyData(usecases.DailyData{
			Date: "2021-01-01",
			Notes: []string{
				"note1",
				"note2",
				"note3",
			},
		})
		err := insertData.SetInsertNum(100)
		assert.Error(t, err)
	})
}
