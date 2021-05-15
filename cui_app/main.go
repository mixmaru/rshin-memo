package main

import (
	"github.com/mixmaru/rshin-memo/core/repositories"
	"github.com/mixmaru/rshin-memo/core/usecases"
	"log"
	"os"
	"path/filepath"
)

func main() {
	getNoteUseCaseInteractor := usecases.NewGetNoteUseCaseInteractor(&repositories.NoteRepositoryMock{})

	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Panicf("homeDir取得失敗. %v", err)
	}
	rshinMemoBaseDir := filepath.Join(homedir, "rshin_memo")
	saveDailyDataUseCaseInteractor := usecases.NewSaveDailyDataUseCaseInteractor(
		repositories.NewNoteRepository(rshinMemoBaseDir),
		repositories.NewDailyDataRepository(filepath.Join(rshinMemoBaseDir, "daily_data.json")),
	)
	rshinMemo := NewRshinMemo(
		&GetAllDailyListUsecaseMock{},
		getNoteUseCaseInteractor,
		saveDailyDataUseCaseInteractor,
	)
	defer rshinMemo.Close()

	err = rshinMemo.Run()
	if err != nil {
		log.Panicf("%+v", err)
	}
}

type GetAllDailyListUsecaseMock struct{}

func (u *GetAllDailyListUsecaseMock) Handle() ([]usecases.DailyData, error) {
	retList := []usecases.DailyData{
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
	}
	return retList, nil
}
