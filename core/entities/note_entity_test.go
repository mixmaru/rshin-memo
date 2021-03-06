package entities

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNoteEntity_Name(t *testing.T) {
	// 準備
	entity := NewNoteEntity(
		"note_name",
		`内容text
内容text2
内容text3
`,
	)
	// 実行
	name := entity.Name()
	// 検証
	assert.Equal(t, "note_name", name)
}

func TestNoteEntity_Text(t *testing.T) {
	// 準備
	entity := NewNoteEntity(
		"note_name",
		`内容text
内容text2
内容text3
`,
	)
	// 実行
	text := entity.Text()
	// 検証
	assert.Equal(t, `内容text
内容text2
内容text3
`, text)
}
