package usecases

import "github.com/mixmaru/rshin-memo/core/repositories"

type GetNotesBySearchTextUseCase struct {
	noteRepository repositories.NoteRepositoryInterface
}

func NewGetNotesBySearchTextUseCase(noteRepository repositories.NoteRepositoryInterface) *GetNotesBySearchTextUseCase {
	return &GetNotesBySearchTextUseCase{
		noteRepository: noteRepository,
	}
}

func (g *GetNotesBySearchTextUseCase) Handle(searchText string) (notes []string, err error) {
	// repositoryに問い合わせ
	entities, err := g.noteRepository.GetBySearchText(searchText)
	if err != nil {
		return nil, err
	}
	retNotes := []string{}
	for _, entity := range entities {
		retNotes = append(retNotes, entity.Name())
	}
	return retNotes, nil
}
