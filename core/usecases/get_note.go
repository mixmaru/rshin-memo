package usecases

import "github.com/mixmaru/rshin-memo/core/repositories"

type GetNoteUseCaseInterface interface {
	Handle(noteName string) (text string, notExist bool, err error)
}

type GetNoteUseCaseInteractor struct {
	noteRepository repositories.NoteRepositoryInterface
}

func NewGetNoteUseCaseInteractor(noteRepository repositories.NoteRepositoryInterface) *GetNoteUseCaseInteractor {
	return &GetNoteUseCaseInteractor{noteRepository: noteRepository}
}

func (g *GetNoteUseCaseInteractor) Handle(noteName string) (text string, notExist bool, err error) {
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
