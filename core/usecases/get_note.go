package usecases

type GetNoteUseCaseInterface interface {
	Handle(noteName string) (text string, notExist bool, err error)
}

type GetNoteUseCaseInteractor struct {}

func (g GetNoteUseCaseInteractor) Handle(noteName string) (text string, notExist bool, err error) {
	return "", false, nil
}

