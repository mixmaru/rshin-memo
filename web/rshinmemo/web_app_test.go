package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestWebApp_list(t *testing.T) {
	////// 準備
	app := NewWebApp("8080", "./testdata/")
	router := app.initRouter()

	////// 実行
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Content-Type", "text/html")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	////// 検証
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestWebApp_memo(t *testing.T) {
	t.Run("200", func(t *testing.T) {
		////// 準備
		app := NewWebApp("8080", "./testdata/")
		router := app.initRouter()

		////// 実行
		req := httptest.NewRequest("GET", "/rshin_memo構築", nil)
		req.Header.Set("Content-Type", "text/html")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		////// 検証
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("404", func(t *testing.T) {
		////// 準備
		app := NewWebApp("8080", "./testdata/")
		router := app.initRouter()

		////// 実行
		req := httptest.NewRequest("GET", "/notFound", nil)
		req.Header.Set("Content-Type", "text/html")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		////// 検証
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestWebApp_noteNew(t *testing.T) {
	t.Run("200", func(t *testing.T) {
		////// 準備
		app := NewWebApp("8080", "./testdata/")
		router := app.initRouter()

		////// 実行
		req := httptest.NewRequest("GET", "/note/new?base=rshin_memo構築&date=2022-02-19&to=newer", nil)
		req.Header.Set("Content-Type", "text/html")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		////// 検証
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestWebApp_addNewNote(t *testing.T) {
	t.Run("302", func(t *testing.T) {
		//// リクエストパラメータ作成
		//body := url.Values{}
		//body.Set("type", "aaaa")
		//body.Set("name", "個人　太郎")
		//
		//// リクエスト実行
		//req := httptest.NewRequest("POST", "/users/", strings.NewReader(body.Encode()))
		//req.Header.Set("Content-Type", "application/x-www-form-urlencoded") //formからの入力ということを指定してるっぽい
		//rec := httptest.NewRecorder()
		//router.ServeHTTP(rec, req)
		////// 準備
		app := NewWebApp("8080", "./testdata/")
		router := app.initRouter()
		// リクエストパラメータ作成
		body := url.Values{}
		body.Set("base_memo_name", "rshin_memo構築")
		body.Set("base_memo_date", "2022-02-19")
		body.Set("new_memo_name", "新規memo")
		body.Set("new_memo_date", "2022-02-19")
		body.Set("to", "older")

		////// 実行
		req := httptest.NewRequest("POST", "/note/new", strings.NewReader(body.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		////// 検証
		assert.Equal(t, http.StatusFound, rec.Code)
		assert.Equal(t, "/", rec.Header().Get("Location"))
	})
}
