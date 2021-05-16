package usecases

import "github.com/mixmaru/rshin-memo/core/repositories"

type GetNoteUseCase struct {
	noteRepository repositories.NoteRepositoryInterface
}

func NewGetNoteUseCase(noteRepository repositories.NoteRepositoryInterface) *GetNoteUseCase {
	return &GetNoteUseCase{noteRepository: noteRepository}
}

func (g *GetNoteUseCase) Handle(noteName string) (text string, notExist bool, err error) {
	// repositoryに問い合わせ
	note, err := g.noteRepository.GetByNoteName(noteName)
	if err != nil {
		return "", false, err
	}
	if note == nil {
		return "", true, nil
	} else {
		return note.Text(), false, nil
	}
}
