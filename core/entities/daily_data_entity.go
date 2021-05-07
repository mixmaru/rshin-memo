package entities

import "time"

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

func (d DailyDataEntity) Date() time.Time {
	return d.date
}

func (d DailyDataEntity) NoteNames() []string {
	return d.noteNames
}
