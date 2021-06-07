package main

import (
	"github.com/mixmaru/rshin-memo/core/repositories"
	"log"
	"os"
	"path/filepath"
)

func main() {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Panicf("homeDir取得失敗. %v", err)
	}
	rshinMemoBaseDir := filepath.Join(homedir, "rshin_memo")
	rshinMemo := NewRshinMemo(
		rshinMemoBaseDir,
		repositories.NewDailyDataRepository(filepath.Join(rshinMemoBaseDir, "daily_data.json")),
		repositories.NewNoteRepository(rshinMemoBaseDir),
	)
	defer rshinMemo.Close()

	err = rshinMemo.Run()
	if err != nil {
		log.Panicf("%+v", err)
	}
}
