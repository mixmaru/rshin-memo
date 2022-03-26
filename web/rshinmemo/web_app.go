package main

import (
	"github.com/labstack/echo/v4"
	"github.com/mixmaru/rshin-memo/core/repositories"
	"github.com/mixmaru/rshin-memo/core/usecases"
	"github.com/pkg/errors"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"regexp"
	"time"
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
}

func (w *WebApp) initRouter() *echo.Echo {
	e := echo.New()
	t := &Template{
		templates: template.Must(template.ParseGlob("template/*.html")),
	}

	e.Renderer = t
	e.GET("/", w.list)
	e.GET("/:memo", w.memo)
	e.GET("/note/new", w.noteNew)
	e.POST("/note/new", w.addNewNote)
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

	reg := "\r\n|\n"
	noteLines := regexp.MustCompile(reg).Split(note, -1)
	// 出力
	return c.Render(http.StatusOK, "memo.html", map[string]interface{}{
		"Title":            noteName,
		"memoTitle":        noteName,
		"memoContentLines": noteLines,
	})
}

func (w *WebApp) noteNew(c echo.Context) error {
	/**
	候補日付取得
	選択の一つまえの日付と、選択日付と、選択の一つあとの日付が必要
	*/
	// 候補日付取得
	memoName := c.QueryParam("base")
	memoDate, err := time.ParseInLocation("2006-01-02T15:04:05.000000Z", c.QueryParam("date")+"T00:00:00.000000Z", time.Local)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	dateList, err := w.getDateList(memoName, memoDate, c.QueryParam("to"))
	if err != nil {
		log.Fatalf("fail getting data: %v", err)
	}
	return w.renderNewNoteForm(
		c,
		c.QueryParam("base"),
		c.QueryParam("date"),
		c.QueryParam("to"),
		dateList,
	)
}

func (w *WebApp) getDateList(memoName string, memoDate time.Time, to string) ([]time.Time, error) {
	now := time.Now()
	rep := repositories.NewDailyDataRepository(filepath.Join(w.dataDirPath, "daily_data.json"))
	usecase := usecases.NewGetDateSelectRangeUseCase(now, rep)
	var mode usecases.InsertMode
	switch to {
	case "newer":
		mode = usecases.INSERT_NEWER_MODE
	case "older":
		mode = usecases.INSERT_OLDER_MODE
	default:
		return nil, errors.Errorf("toがおかしい to: %v", to)
	}
	return usecase.Handle(memoName, memoDate, mode)
}

func (w *WebApp) renderNewNoteForm(c echo.Context, base, date, to string, dateList []time.Time) error {
	return c.Render(http.StatusOK, "new_form.html", map[string]interface{}{
		"Title":    "note追加",
		"dateList": dateList,
		"Base":     base,
		"Date":     date,
		"To":       to,
	})
}

func (w *WebApp) addNewNote(c echo.Context) error {
	// パラメータ取得
	baseMemoDate, err := time.ParseInLocation("2006-01-02T15:04:05.000000Z", c.FormValue("base_memo_date")+"T00:00:00.000000Z", time.Local)
	if err != nil {
		// todo: バリデーションエラーにする
		log.Fatalf("fail date parse: %v", errors.WithStack(err))
	}
	baseMemoName := c.FormValue("base_memo_name")

	newMemoDate, err := time.ParseInLocation("2006-01-02T15:04:05.000000Z", c.FormValue("new_memo_date")+"T00:00:00.000000Z", time.Local)
	if err != nil {
		// todo: バリデーションエラーにする
		log.Fatalf("fail date parse: %v", errors.WithStack(err))
	}
	newMemoName := c.FormValue("new_memo_name")
	to := c.FormValue("to")
	dummyUseCase(baseMemoDate, baseMemoName, newMemoDate, newMemoName, to)
	/*

	 */
	// todo: メモ編集画面へリダイレクトさせる
	return c.Redirect(http.StatusFound, "/")
}

func dummyUseCase(date time.Time, name string, date2 time.Time, name2 string, to string) {

}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
