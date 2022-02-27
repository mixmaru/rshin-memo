package main

import (
	"github.com/labstack/echo/v4"
	"github.com/mixmaru/rshin-memo/core/repositories"
	"github.com/mixmaru/rshin-memo/core/usecases"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

type WebApp struct {
	port        string
	dataDirPath string
}

func NewWebApp(port string, dataDirPath string) *WebApp {
	return &WebApp{dataDirPath: dataDirPath, port: port}
}

func (w *WebApp) Run() {
	e := w.initRouter()
	e.Logger.Fatal(e.Start(":" + w.port))
	//http.HandleFunc("/", w.list)

	//log.Printf("Server listening on port %s", w.port)
	//log.Print(http.ListenAndServe(":"+w.port, nil))
}

func (w *WebApp) initRouter() *echo.Echo {
	e := echo.New()
	e.GET("/", w.list)
	return e
}

func (w *WebApp) list(c echo.Context) error {
	return c.String(http.StatusOK, "こんにちは!")
}

func (w *WebApp) list_old(writer http.ResponseWriter, request *http.Request) {
	// メモ一覧データ取得
	rep := repositories.NewDailyDataRepository(filepath.Join(w.dataDirPath, "daily_data.json"))
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
