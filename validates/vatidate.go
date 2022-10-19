package validates

import (
	"unicode/utf8"
	"log"
	// "reflect"
	// "encoding/json"
	// "fmt"
)
/*
   音楽投稿機能バリデーション
*/
func Register(name, artist, reason string) (bool) {

	var nameLn int = utf8.RuneCountInString(name)
	var artistln int = utf8.RuneCountInString(artist)
	var reasonLn int = utf8.RuneCountInString(reason)

	// 文字数を確認する
	if nameLn == 0 || nameLn >= 101 || artistln == 0 || artistln >= 101 || reasonLn >= 1001{
		log.Println("文字数エラー")
		return false
	}

	return true
}

/*
   利用者新規登録機能バリデーション
*/
func SignUp(name, password string) (bool) {

	var nameLn int = utf8.RuneCountInString(name)
	var passwordLn int = utf8.RuneCountInString(password)

	// 文字数を確認する
	if nameLn == 0 || nameLn >= 10 || passwordLn < 8 || passwordLn >= 16{
		log.Println("文字数エラー")
		return false
	}

	return true
}