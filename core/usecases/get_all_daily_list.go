package usecases

import "time"

type GetAllDailyListUsecaseInterface interface {
    Handle() (GetAllDailyListUsecaseResponse, error)
}

type GetAllDailyListUsecaseResponse struct {
    DailyList []OneDayList
}

type OneDayList struct {
    Date time.Time
    Notes []string
}

type GetAllDailyListUsecaseInteractor struct {}

func (i *GetAllDailyListUsecaseInteractor) Handle() (GetAllDailyListUsecaseResponse, error) {
    // 実際にファイルからデータを取得する処理を書く。TDDで。
    return GetAllDailyListUsecaseResponse{}, nil
}
