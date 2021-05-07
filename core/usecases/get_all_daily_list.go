package usecases

type GetAllDailyListUsecaseInterface interface {
	Handle() ([]DailyData, error)
}

type GetAllDailyListUsecaseInteractor struct{}

func (i *GetAllDailyListUsecaseInteractor) Handle() ([]DailyData, error) {
	// 実際にファイルからデータを取得する処理を書く。TDDで。
	return nil, nil
}
