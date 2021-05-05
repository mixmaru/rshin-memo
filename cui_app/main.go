package main

import (
	"github.com/mixmaru/rshin-memo/core/usecases"
	"log"
	"time"
)

func main() {
	rshinMemo := NewRshinMemo(&GetAllDailyListUsecaseMock{})
	defer rshinMemo.Close()

	err := rshinMemo.Run()
	if err != nil {
		log.Panicf("%+v", err)
	}
}

type GetAllDailyListUsecaseMock struct {}

func (u *GetAllDailyListUsecaseMock) Handle() (usecases.GetAllDailyListUsecaseResponse, error) {
	local, _ := time.LoadLocation("Local")
	response := usecases.GetAllDailyListUsecaseResponse{
		DailyList: []usecases.OneDayList{
			{
				Date: time.Date(2021, 04, 30, 0,0,0,0, local),
				Notes: []string{
					"なんかしらのNote1",
					"なんかしらのNote2",
					"なんかしらのNote3",
					"なんかしらのNote4",
					"なんかしらのNote5",
					"なんかしらのNote6",
					"なんかしらのNote7",
					"なんかしらのNote8",
				},
			},
			{
				Date: time.Date(2021, 04, 29, 0,0,0,0, local),
				Notes: []string{
					"なんかしらのNote1",
					"なんかしらのNote2",
					"なんかしらのNote3",
					"なんかしらのNote4",
					"なんかしらのNote5",
					"なんかしらのNote6",
					"なんかしらのNote7",
					"なんかしらのNote8",
				},
			},
		},
	}
	return response, nil
}
