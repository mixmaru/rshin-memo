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
	req.Header.Set("Content-Type", "text/html") //formからの入力ということを指定してるっぽい
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	////// 検証
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestWebApp_list_old(t *testing.T) {
	////// 準備
	app := NewWebApp("8080", "./testdata/")

	////// 実行
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	app.list_old(rec, req)

	////// 検証
	assert.Equal(t, rec.Code, 200)
}
