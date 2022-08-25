package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	// "reflect"
	"strconv"
	"unicode/utf8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"backend/unify"
	"backend/models"
)

const (
	StatusBadRequest			= 400
	StatusInternalServerError	= 500
)

// グローバルスコープとして定義することで、本ファイルのどの関数でも引数の受け渡しなしに使用可能にする。
var db *gorm.DB
var db_err error

func main() {
	fmt.Println("Start!")
	dsn := unify.DBSet
	db, db_err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if db_err != nil {
		panic(db_err)
	}

	http.HandleFunc("/", top)
	http.HandleFunc("/detail", detail)
	http.HandleFunc("/register", register)
	http.ListenAndServe(":8080", nil)
}

/*
   Top画面
*/
func top(w http.ResponseWriter, r *http.Request) {
	fmt.Println("パス（\"/\"）でGOが呼び出された")

	// ヘッダーをセットする（エラー処理後にセットするとCROSエラーになる）
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Content-Type", "application/json")

	// 全レコードを取得する
	music, situation, orm_err := models.ReadMulti(db)
	if !orm_err {
		fmt.Println("error happen!")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, false)
		return
	}

	// jsonエンコード
	situationJson, errSitu := json.Marshal(situation)
	musicJson, errMusic := json.Marshal(music)
	if errSitu != nil || errMusic != nil{
		fmt.Println("error happen!")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, false)
		return
	}

	var res = unify.ResponseTop{Mst_situation: string(situationJson), Music: string(musicJson)}
	outputJson, err := json.Marshal(res)

	// エラー処理
	if err != nil {
		fmt.Println("error happen!")
		w.WriteHeader(http.StatusInternalServerError)
	}

	fmt.Fprint(w, string(outputJson))
}

/*
   詳細画面
*/
func detail(w http.ResponseWriter, r *http.Request) {
	fmt.Println("パス（\"/detail\"）でGOが呼び出された")
	
	// ヘッダーをセットする
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Content-Type", "application/json")

	// クエリパラメータ「id」を取得する
	var id string = r.URL.Query().Get("id")

	// ブラウザをリロードした際にクエリパラメータがundefindで送付される場合がある
	if r.Method != http.MethodGet || id == "undefined"{
		return
	}

	ret, orm_err := models.Read(db, id)

	// jsonエンコード
	outputJson, err := json.Marshal(ret)

	// エラー処理
	if err != nil || !orm_err {
		fmt.Println("error happen!")
		w.WriteHeader(http.StatusInternalServerError)
	}

	// jsonデータを返却する
	fmt.Fprint(w, string(outputJson))
}

/*
   登録機能
*/
func register(w http.ResponseWriter, r *http.Request) {
	fmt.Println("パス（\"/register\"）でGOが呼び出された")
	// 登録機能時にOPTIONSリクエストが送付される
	fmt.Println(r.Method)
	if r.Method != http.MethodPost {
		return
	}

	// // ヘッダーをセットする
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	// クエリパラメータを受け取る
	var name string = r.URL.Query().Get("name")
	var nameln int = utf8.RuneCountInString(name)
	// fmt.Println(nameln)
	var artist string = r.URL.Query().Get("artist")
	var artistln int = utf8.RuneCountInString(artist)
	// fmt.Println(artistln)
	var reason string = r.URL.Query().Get("reason")
	var reasonln int = utf8.RuneCountInString(reason)
	// fmt.Println(reasonln)

	// 文字数チェック
	if nameln == 0 || nameln >= 101 || artistln == 0 || artistln >= 101 || reasonln >= 1001{
		// 文字数が不正である場合は400エラーを返却する
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, false)
		return 
	}

	// Mst_situationIDをint型に変換する
	var situationID int
	var s string = r.URL.Query().Get("situation")
	situationID, _ = strconv.Atoi(s)

	// クエリパラメータに含まれた値を使用して構造体を初期化する。
	var create = unify.Music{Name: name, Reason: reason, Artist: artist, Mst_situationID: situationID}

	// レコードの作成
	ret := models.Register(db, &create)

	if !ret {
		fmt.Println("error happen!")
		w.WriteHeader(http.StatusInternalServerError)
	}

	// データを返却する
	fmt.Fprint(w, true)
}

