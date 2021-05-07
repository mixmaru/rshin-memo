package entities

import "time"

type NoteEntity struct {
	date time.Time
	name string
	text string
}

func NewNoteEntity(date time.Time, name string, text string) *NoteEntity {
	return &NoteEntity{date: date, name: name, text: text}
}

func (n *NoteEntity) Date() time.Time {
	return n.date
}

func (n *NoteEntity) Name() string {
	return n.name
}

func (n *NoteEntity) Text() string {
	return n.text
}
