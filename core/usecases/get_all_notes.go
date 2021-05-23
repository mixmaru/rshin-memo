package usecases

import "github.com/mixmaru/rshin-memo/core/repositories"

type GetAllNotesUseCase struct {
	noteRepository repositories.NoteRepositoryInterface
}

func NewGetAllNotesUseCase(noteRepository repositories.NoteRepositoryInterface) *GetAllNotesUseCase {
	return &GetAllNotesUseCase{noteRepository: noteRepository}
}

func (g *GetAllNotesUseCase) Handle() ([]string, error) {
	// repositoryに問い合わせ
	notes, err := g.noteRepository.GetAllNotesOnlyName()
	if err != nil {
		return nil, err
	}
	retNoteNames := []string{}
	for _, note := range notes {
		retNoteNames = append(retNoteNames, note.Name())
	}
	return retNoteNames, nil
}
