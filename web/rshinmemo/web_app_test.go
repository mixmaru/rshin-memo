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
