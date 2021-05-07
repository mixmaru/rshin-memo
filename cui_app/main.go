package main

import (
	"github.com/mixmaru/rshin-memo/core/repositories"
	"github.com/mixmaru/rshin-memo/core/usecases"
	"log"
)

func main() {
	getNoteUseCaseInteractor := usecases.NewGetNoteUseCaseInteractor(&repositories.NoteRepositoryMock{})
	rshinMemo := NewRshinMemo(
		&GetAllDailyListUsecaseMock{},
		getNoteUseCaseInteractor,
		&CreateNoteUseCaseMock{},
	)
	defer rshinMemo.Close()

	err := rshinMemo.Run()
	if err != nil {
		log.Panicf("%+v", err)
	}
}

type GetAllDailyListUsecaseMock struct{}

func (u *GetAllDailyListUsecaseMock) Handle() (usecases.GetAllDailyListUsecaseResponse, error) {
	response := usecases.GetAllDailyListUsecaseResponse{
		DailyList: []usecases.OneDayList{
			{
				Date: "2021-04-30",
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
				Date: "2021-04-29",
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

type CreateNoteUseCaseMock struct{}

func (c CreateNoteUseCaseMock) Handle(noteName string) error {
	return nil
}
