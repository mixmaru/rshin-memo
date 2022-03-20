package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
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
