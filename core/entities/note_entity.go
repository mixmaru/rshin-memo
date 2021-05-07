package entities

type NoteEntity struct {
	name string
	text string
}

func NewNoteEntity(name string, text string) *NoteEntity {
	return &NoteEntity{name: name, text: text}
}

func (n *NoteEntity) Name() string {
	return n.name
}

func (n *NoteEntity) Text() string {
	return n.text
}
