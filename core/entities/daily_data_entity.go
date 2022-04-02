package entities

import (
	"github.com/pkg/errors"
	"time"
)

type DailyDataEntity struct {
	date      time.Time
	noteNames []string
}

func NewDailyDataEntity(date time.Time, noteNames []string) *DailyDataEntity {
	return &DailyDataEntity{
		date:      date,
		noteNames: noteNames,
	}
}

func NewDailyDataEntityByLoadedData(dateStr string, noteNames []string) (*DailyDataEntity, error) {
	date, err := time.ParseInLocation("2006-01-02", dateStr, time.Local)
	if err != nil {
		return nil, errors.Wrapf(err, "日付パース失敗. %v", dateStr)
	}
	return &DailyDataEntity{
		date:      date,
		noteNames: noteNames,
	}, nil
}

func (d *DailyDataEntity) Date() time.Time {
	return d.date
}

func (d *DailyDataEntity) DateStr() string {
	return d.date.Format("2006-01-02")
}

func (d *DailyDataEntity) NoteNames() []string {
	return d.noteNames
}

func (d *DailyDataEntity) InsertNoteName(noteName string, index int) {
	// 末挿入の場合
	if len(d.noteNames) <= index {
		d.noteNames = append(d.noteNames, noteName)
	} else {
		// 中に挿入の場合
		d.noteNames = append(d.noteNames[:index+1], d.noteNames[index:]...)
		d.noteNames[index] = noteName
	}
}
