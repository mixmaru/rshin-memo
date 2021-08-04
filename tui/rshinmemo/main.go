package main

import (
	"github.com/mixmaru/rshin-memo/core/repositories"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

func main() {
	baseDirPath, err := getRshinMamoBaseDirPath()
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

func getRshinMamoBaseDirPath() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.WithStack(err)
	}
	return filepath.Join(homedir, "rshinmemo"), nil
}

func getRshinMamoDailyDataJsonFilePath() (string, error) {
	baseDirPath, err := getRshinMamoBaseDirPath()
	if err != nil {
		return "", err
	}

	return filepath.Join(baseDirPath, "daily_data.json"), nil
}
