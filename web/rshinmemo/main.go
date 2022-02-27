package main

import (
	"flag"
	"github.com/mixmaru/rshin-memo/utils"
	"log"
	"path/filepath"
)

func main() {
	dataDirPath, err := utils.GetRshinMamoBaseDirPath()
	if err != nil {
		log.Fatalf("fail: %v", err)
	}
	// コマンド引数定義
	port := flag.String("port", "8080", "利用ポート番号")
	dataDirPathFlag := flag.String("datadir", dataDirPath, "利用ポート番号")
	flag.Parse()

	// データディレクトリパス
	dataDirAbsPath, err := filepath.Abs(*dataDirPathFlag)
	if err != nil {
		log.Fatalf("fail getting data dir abs path: %v,  %v", dataDirAbsPath, err)
	}
	log.Printf("dataDir %s", dataDirAbsPath)

	app := NewWebApp(*port, dataDirAbsPath)
	app.Run()
}
