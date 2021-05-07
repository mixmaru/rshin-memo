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
