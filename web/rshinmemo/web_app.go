package main

import (
	"github.com/labstack/echo/v4"
	"github.com/mixmaru/rshin-memo/core/repositories"
	"github.com/mixmaru/rshin-memo/core/usecases"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
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
	t := &Template{
		templates: template.Must(template.ParseGlob("template/*.html")),
	}
	e.Renderer = t
	e.GET("/", w.list)
	e.GET("/:memo", w.memo)
	return e
}

func (w *WebApp) list(c echo.Context) error {
	// メモ一覧データ取得
	rep := repositories.NewDailyDataRepository(filepath.Join(w.dataDirPath, "daily_data.json"))
	useCase := usecases.NewGetAllDailyListUsecase(rep)
	dailyData, err := useCase.Handle()
	if err != nil {
		log.Fatalf("fail getting data: %v", err)
	}

	// 出力
	return c.Render(http.StatusOK, "index.html", map[string]interface{}{
		"Title":     "一覧",
		"DailyData": dailyData,
	})
}
func (w *WebApp) memo(c echo.Context) error {
	noteName, err := url.PathUnescape(c.Param("memo"))
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}
	useCase := usecases.NewGetNoteUseCase(repositories.NewNoteRepository(w.dataDirPath))
	note, notExist, err := useCase.Handle(noteName)
	if err != nil {
		log.Fatalf("fail getting data: %v", err)
	}
	if notExist {
		return c.NoContent(http.StatusNotFound)
	}

	// 出力
	return c.Render(http.StatusOK, "memo.html", map[string]interface{}{
		"Title":       noteName,
		"memoTitle":   noteName,
		"memoContent": note,
	})
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
