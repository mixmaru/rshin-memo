package main

import (
	"github.com/mixmaru/rshin-memo/core/repositories"
	"github.com/mixmaru/rshin-memo/utils"
	"path/filepath"
)

func main() {
	baseDirPath, err := utils.GetRshinMamoBaseDirPath()
	if err != nil {
		panic(err)
	}
	dailyDataJsonFilePath, err := getRshinMamoDailyDataJsonFilePath()
	if err != nil {
		panic(err)
	}
	rshinMemo := NewRshinMemo(
		baseDirPath,
		repositories.NewDailyDataRepository(dailyDataJsonFilePath),
		repositories.NewNoteRepository(baseDirPath),
	)
	err = rshinMemo.Run()
	if err != nil {
		panic(err)
	}
}

func getRshinMamoDailyDataJsonFilePath() (string, error) {
	baseDirPath, err := utils.GetRshinMamoBaseDirPath()
	if err != nil {
		return "", err
	}

	return filepath.Join(baseDirPath, "daily_data.json"), nil
}
