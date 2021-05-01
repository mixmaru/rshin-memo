package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRshinMemo_New(t *testing.T) {
	t.Run("インスタンス化できる", func(t *testing.T) {
        // 実行
        memo := NewRshinMemo()
        // 検証
		assert.IsType(t, &RshinMemo{},  memo)
    })
}
