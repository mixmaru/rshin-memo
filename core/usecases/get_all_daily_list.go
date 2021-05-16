package usecases

import "github.com/mixmaru/rshin-memo/core/repositories"

type GetAllDailyListUsecaseInterface interface {
	Handle() ([]DailyData, error)
}

type GetAllDailyListUsecaseInteractor struct {
	dailyDataRepository repositories.DailyDataRepositoryInterface
}

func NewGetAllDailyListUsecaseInteractor(dailyDataRepository repositories.DailyDataRepositoryInterface) *GetAllDailyListUsecaseInteractor {
	return &GetAllDailyListUsecaseInteractor{dailyDataRepository: dailyDataRepository}
}

func (i *GetAllDailyListUsecaseInteractor) Handle() ([]DailyData, error) {
	entities, err := i.dailyDataRepository.Get()
	if err != nil {
		return nil, err
	}

	retData := []DailyData{}
	for _, entity := range entities {
		retData = append(retData, DailyData{
			Date:  entity.DateStr(),
			Notes: entity.NoteNames(),
		})
	}
	return retData, nil
}
