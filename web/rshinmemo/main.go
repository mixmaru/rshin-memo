package main

import (
	"github.com/mixmaru/rshin-memo/core/repositories"
	"github.com/mixmaru/rshin-memo/core/usecases"
	"html/template"
	"log"
	"net/http"
)

func main() {
	port := "8080"
	http.HandleFunc("/", list)
	log.Printf("Server listening on port %s", port)
	log.Print(http.ListenAndServe(":"+port, nil))
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
