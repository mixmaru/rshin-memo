package usecases

type GetNoteSummaryUsecaseInterface interface {
    Handle(noteName string) (string, error)
}
