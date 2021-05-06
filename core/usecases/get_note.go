package usecases

type GetNoteUseCaseInterface interface {
	Handle(noteName string) (text string, notExist bool, err error)
}
