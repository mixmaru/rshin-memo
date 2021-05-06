package usecases

type CreateNoteUseCaseInterface interface {
	Handle(noteName string) error
}
