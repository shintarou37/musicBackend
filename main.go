package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	// "reflect"
	"strconv"
	// "unicode/utf8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"backend/unify"
	"backend/models"
	"backend/validates"
	"github.com/joho/godotenv"
	"os"
	"log"
	"io"
)

const (
	StatusBadRequest          = 400
	StatusInternalServerError = 500
)

// グローバルスコープとして定義することで、本ファイルのどの関数でも引数の受け渡しなしに使用可能にする。
var db *gorm.DB
var db_err error

var logfile, _ = os.OpenFile("./request.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

func main() {
	fmt.Println("Start!")
	err := godotenv.Load(fmt.Sprintf("env/%s.env", os.Getenv("GO_ENV")))
	if err != nil {
			fmt.Println(err)
	}
	dsn := os.Getenv("DB_SET")
	port := os.Getenv("PORT")
	db, db_err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if db_err != nil {
		fmt.Println("gorm Open err")
		panic(db_err)
	}
	log.SetOutput(io.MultiWriter(logfile, os.Stdout))

	log.SetFlags(log.Ldate | log.Ltime)
	http.HandleFunc("/", top)
	http.HandleFunc("/detail", detail)
	http.HandleFunc("/register", register)
	http.ListenAndServe(port, nil)
}

/*
   Top画面
*/
func top(w http.ResponseWriter, r *http.Request) {
	// クエリパラメーター"search"を取得する
	var search string = r.URL.Query().Get("search")
	if search == "" {
		log.Println("params「search」が空文字列です")
	}

	// ヘッダーをセットする（エラー処理後にセットするとCROSエラーになる）
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Content-Type", "application/json")

	// 全レコードを取得する
	music, situation, orm_err := models.ReadMulti(db, search)
	if !orm_err {
		fmt.Println("ReadMulti error happen!")
		fmt.Println(orm_err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, false)
		return
	}

	var res = unify.ResponseTop{Mst_situation: situation, Music: music}

	// jsonエンコード
	outputJson, _ := json.Marshal(res)

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
	if r.Method != http.MethodGet || id == "undefined" {
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
	var artist string = r.URL.Query().Get("artist")
	var reason string = r.URL.Query().Get("reason")

	// 文字数チェック
	retVatidate := validates.Register(name, artist, reason)

	if !retVatidate {
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
		fmt.Println("登録エラー")
		w.WriteHeader(http.StatusInternalServerError)
	}

	// データを返却する
	fmt.Fprint(w, true)
}
