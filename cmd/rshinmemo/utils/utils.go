package utils

import (
	"github.com/mattn/go-runewidth"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"strings"
)

// 全角文字だったら後ろに半角スペースを入れる。
// gocuiが全角文字表示に対応しておらず、こうしないとうまく表示されないため
func ConvertStringForView(s string) string {
	runeArr := []rune{}
	for _, r := range s {
		runeArr = append(runeArr, r)
		// if もし全角文字だったら
		if runewidth.StringWidth(string(r)) == 2 {
			runeArr = append(runeArr, ' ')
		}
	}
	return string(runeArr)
}

// ConvertStringForViewでいれた半角スペースを除去する
func ConvertStringForLogic(s string) string {
	return strings.ReplaceAll(s, " ", "")
}

// vimで対象noteを開く
func OpenVim(filePath string) error {
	c := exec.Command("vim", filePath)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err := c.Run()
	if err != nil {
		return errors.Wrap(err, "vim起動エラー")
	}
	return nil
}
