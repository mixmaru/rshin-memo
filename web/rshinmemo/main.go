package main

import (
	"flag"
	"github.com/mixmaru/rshin-memo/core/repositories"
	"github.com/mixmaru/rshin-memo/core/usecases"
	"github.com/mixmaru/rshin-memo/utils"
	"html/template"
	"log"
	"net/http"
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

	http.HandleFunc("/", list)
	log.Printf("Server listening on port %s", *port)
	log.Print(http.ListenAndServe(":"+*port, nil))
}

func list(writer http.ResponseWriter, request *http.Request) {
	// メモ一覧データ取得
	rep := repositories.NewDailyDataRepository("testdata/daily_data.json")
	useCase := usecases.NewGetAllDailyListUsecase(rep)
	dailyData, err := useCase.Handle()
	if err != nil {
		log.Fatalf("fail getting data: %v", err)
	}

	// 出力
	t, err := template.ParseFiles("template/index.html")
	if err != nil {
		log.Fatalf("template error: %v", err)
	}
	if err := t.Execute(writer, struct {
		Title     string
		DailyData []usecases.DailyData
	}{
		Title:     "一覧",
		DailyData: dailyData,
	}); err != nil {
		log.Printf("failed to execute template: %v", err)
	}
}
