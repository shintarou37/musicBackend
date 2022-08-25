package validates

import (
	// "reflect"
	// "encoding/json"
	"fmt"
	"unicode/utf8"
)

func Register(name, artist, reason string) (bool) {

	var nameln int = utf8.RuneCountInString(name)
	var artistln int = utf8.RuneCountInString(artist)
	var reasonln int = utf8.RuneCountInString(reason)

	// 文字数を確認する
	if nameln == 0 || nameln >= 101 || artistln == 0 || artistln >= 101 || reasonln >= 1001{
		fmt.Println("文字数エラー")
		return false
	}

	return true
}