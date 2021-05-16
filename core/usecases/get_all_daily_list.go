package usecases

import "github.com/mixmaru/rshin-memo/core/repositories"

type GetAllDailyListUsecase struct {
	dailyDataRepository repositories.DailyDataRepositoryInterface
}

func NewGetAllDailyListUsecase(dailyDataRepository repositories.DailyDataRepositoryInterface) *GetAllDailyListUsecase {
	return &GetAllDailyListUsecase{dailyDataRepository: dailyDataRepository}
}

func (i *GetAllDailyListUsecase) Handle() ([]DailyData, error) {
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
