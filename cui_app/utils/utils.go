package utils

import (
	"github.com/mattn/go-runewidth"
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

