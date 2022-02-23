package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_list(t *testing.T) {
	////// 実行
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	list(rec, req)

	////// 検証
	assert.Equal(t, rec.Code, 200)
}
