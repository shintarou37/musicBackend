package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	// "reflect"
)
const (
	StatusInternalServerError = 500
)

type Music struct {
    gorm.Model
	Name    string `json:"name"`
	Artist  string `json:"artist"`
	Reason  string `json:"reason"`
}

// グローバルスコープとして定義することで、本ファイルのどの関数でも引数の受け渡しなしに使用可能にする。
var db *gorm.DB
var db_err error

func main() {
    fmt.Println("Start!");
    dsn := "root:@tcp(127.0.0.1:3306)/music?charset=utf8mb4&parseTime=True&loc=Local"
    db, db_err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if db_err != nil {
		panic(db_err)
	}

    http.HandleFunc("/", top);
    http.HandleFunc("/detail", detail);
    http.HandleFunc("/register", register);
    http.ListenAndServe(":8080", nil)
}
/* 
    Top画面 
*/
func top(w http.ResponseWriter, r *http.Request){
    fmt.Println("パス（\"/\"）でGOが呼び出された")

    // ヘッダーをセットする（エラー処理後にセットするとCROSエラーになる）
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Headers", "*")
    w.Header().Set("Content-Type", "application/json")

    // 全レコードを取得する
    ret, orm_err := ReadMulti()

    // jsonエンコード
    outputJson, err := json.Marshal(ret)

    // エラー処理
    if err != nil || !orm_err{
        fmt.Println("error happen!")
        w.WriteHeader(http.StatusInternalServerError)
    }

    // jsonデータを返却する（エラーが発生した場合は空のオブジェクトを返却する）
    fmt.Fprint(w, string(outputJson))
}

/* 
    詳細画面
*/
func detail(w http.ResponseWriter, r *http.Request){
    fmt.Println("パス（\"/detail\"）でGOが呼び出された")

    // ヘッダーをセットする
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Content-Type", "application/json")

    // クエリパラメータ「id」を取得する
    var id string = r.URL.Query().Get("id")

    // React側で画面をリロードするとクエリパラメータがundefinedで送付される
    // その場合は"false"という文字列がパラメーターとして送信されてsqlは発行しない
	if id == "false" {
		panic("no query params")
        // これ以降の処理は行われない
	}
    
    ret, orm_err := Read(id)

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
func register(w http.ResponseWriter, r *http.Request){
    fmt.Println("パス（\"/register\"）でGOが呼び出された")
    // 登録機能時にOPTIONSリクエストが送付される
	fmt.Println(r.Method)
	if r.Method != http.MethodPost {
        return
    }
    // クエリパラメータに含まれた値を使用して構造体を初期化する。
    var create = Music{Name: r.URL.Query().Get("name"), Reason: r.URL.Query().Get("reason"), Artist: r.URL.Query().Get("artist")}

    // レコードの作成
    if orm_err := db.Debug().Create(&create).Error; orm_err != nil {
        fmt.Println("error happen!")
		w.WriteHeader(http.StatusInternalServerError)
	}

		// jsonエンコード
		outputJson, err := json.Marshal(create)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

	fmt.Println(string(outputJson))

    // // ヘッダーをセットする
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Headers", "*")

    // jsonデータを返却する
    fmt.Fprint(w, string(outputJson))
}
/* 
    パス：top
*/
func ReadMulti()([]Music, bool){
    var music_arr []Music
    // return music_arr, false
    if err := db.Debug().Find(&music_arr).Error; err != nil {
		return music_arr, false
	}
    return music_arr, true
}

/* 
    パス：detail
*/
func Read(id string) (Music, bool){
    var music Music
    // return music, false
    // ポインタを引数にしない場合はエラーになる
    if err := db.Debug().First(&music, id).Error; err != nil {
        fmt.Println("error happen!")
		return music, true
	}
    return music, true
}