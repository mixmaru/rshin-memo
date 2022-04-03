package entities

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDailyDataEntity_Date(t *testing.T) {
	entity := NewDailyDataEntity(
		time.Date(2021, 1, 1, 12, 0, 0, 0, time.Local),
		[]string{
			"Note1",
			"Note2",
			"Note3",
		},
	)
	assert.Equal(t, time.Date(2021, 1, 1, 12, 0, 0, 0, time.Local), entity.Date())
}

func TestDailyDataEntity_NoteNames(t *testing.T) {
	entity := NewDailyDataEntity(
		time.Date(2021, 1, 1, 12, 0, 0, 0, time.Local),
		[]string{
			"Note1",
			"Note2",
			"Note3",
		},
	)
	assert.Equal(t, []string{
		"Note1",
		"Note2",
		"Note3",
	}, entity.NoteNames())
}

func TestNewDailyDataEntityByLoadedData(t *testing.T) {
	entity, err := NewDailyDataEntityByLoadedData("2021-01-01", []string{
		"note1",
		"note2",
	})
	assert.NoError(t, err)
	assert.Equal(t, time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local), entity.Date())
	assert.EqualValues(t, []string{
		"note1",
		"note2",
	}, entity.NoteNames())
}

func TestDailyDataEntity_DateStr(t *testing.T) {
	entity := NewDailyDataEntity(
		time.Date(2021, 1, 1, 12, 0, 0, 0, time.Local),
		[]string{
			"Note1",
			"Note2",
			"Note3",
		},
	)
	assert.Equal(t, "2021-01-01", entity.DateStr())
}

func TestDailyDataEntity_InsertNoteName(t *testing.T) {
	t.Run("先頭", func(t *testing.T) {
		// 準備
		entity := NewDailyDataEntity(
			time.Date(2021, 1, 1, 12, 0, 0, 0, time.Local),
			[]string{
				"Note1",
				"Note2",
				"Note3",
			},
		)

		// 実行
		entity.InsertNoteName("insertNote", 0)

		// 検証
		expected := NewDailyDataEntity(
			time.Date(2021, 1, 1, 12, 0, 0, 0, time.Local),
			[]string{
				"insertNote",
				"Note1",
				"Note2",
				"Note3",
			},
		)
		assert.Equal(t, expected, entity)
	})

	t.Run("中", func(t *testing.T) {
		// 準備
		entity := NewDailyDataEntity(
			time.Date(2021, 1, 1, 12, 0, 0, 0, time.Local),
			[]string{
				"Note1",
				"Note2",
				"Note3",
			},
		)

		// 実行
		entity.InsertNoteName("insertNote", 1)

		// 検証
		expected := NewDailyDataEntity(
			time.Date(2021, 1, 1, 12, 0, 0, 0, time.Local),
			[]string{
				"Note1",
				"insertNote",
				"Note2",
				"Note3",
			},
		)
		assert.Equal(t, expected, entity)
	})

	t.Run("末", func(t *testing.T) {
		// 準備
		entity := NewDailyDataEntity(
			time.Date(2021, 1, 1, 12, 0, 0, 0, time.Local),
			[]string{
				"Note1",
				"Note2",
				"Note3",
			},
		)

		// 実行
		entity.InsertNoteName("insertNote", 3)

		// 検証
		expected := NewDailyDataEntity(
			time.Date(2021, 1, 1, 12, 0, 0, 0, time.Local),
			[]string{
				"Note1",
				"Note2",
				"Note3",
				"insertNote",
			},
		)
		assert.Equal(t, expected, entity)
	})

	t.Run("末以降", func(t *testing.T) {
		// 準備
		entity := NewDailyDataEntity(
			time.Date(2021, 1, 1, 12, 0, 0, 0, time.Local),
			[]string{
				"Note1",
				"Note2",
				"Note3",
			},
		)

		// 実行
		entity.InsertNoteName("insertNote", 10)

		// 検証
		expected := NewDailyDataEntity(
			time.Date(2021, 1, 1, 12, 0, 0, 0, time.Local),
			[]string{
				"Note1",
				"Note2",
				"Note3",
				"insertNote",
			},
		)
		assert.Equal(t, expected, entity)
	})
}
