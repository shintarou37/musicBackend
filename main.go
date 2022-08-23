package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	// "reflect"
	"strconv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const StatusInternalServerError int = 500

type Mst_situation struct {
	gorm.Model
	Name   string
	Musics []Music
}

type Music struct {
	gorm.Model
	Name            string
	Artist          string
	Reason          string
	Mst_situationID int
}

type ResponseTop struct {
	Mst_situation string
	Music         string
}

type ResultMusic struct {
	gorm.Model
	Name            	string
	Artist          	string
	Reason          	string
	Mst_situationID		int
	Mst_situationName   string
}

// グローバルスコープとして定義することで、本ファイルのどの関数でも引数の受け渡しなしに使用可能にする。
var db *gorm.DB
var db_err error

func main() {
	fmt.Println("Start!")
	dsn := "root:@tcp(127.0.0.1:3306)/music?charset=utf8mb4&parseTime=True&loc=Local"
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
	music, situation, orm_err := ReadMulti()
	// fmt.Println(reflect.TypeOf(music))
	// fmt.Println(reflect.TypeOf(situation))
	// jsonエンコード
	situationJson, err := json.Marshal(situation)
	musicJson, err := json.Marshal(music)
	// fmt.Println(reflect.TypeOf(situationJson))
	// fmt.Println(reflect.TypeOf(musicJson))
	// arr := [2]string{string(situationJson), string(musicJson)}
	// outputJson, err := json.Marshal(arr)
	var res = ResponseTop{Mst_situation: string(situationJson), Music: string(musicJson)}
	outputJson, err := json.Marshal(res)
	// エラー処理
	if err != nil || !orm_err {
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
	// ヘッダーをセットする
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

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

	// Mst_situationIDをint型に変換する
	var situationID int
	var s string = r.URL.Query().Get("situation")
	situationID, _ = strconv.Atoi(s)

	// クエリパラメータに含まれた値を使用して構造体を初期化する。
	var create = Music{Name: r.URL.Query().Get("name"), Reason: r.URL.Query().Get("reason"), Artist: r.URL.Query().Get("artist"), Mst_situationID: situationID}

	// レコードの作成
	if orm_err := db.Debug().Create(&create).Error; orm_err != nil {
		fmt.Println("error happen!")
		w.WriteHeader(http.StatusInternalServerError)
	}

	// // ヘッダーをセットする
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	// データを返却する
	fmt.Fprint(w, true)
}

/*
   パス：top
*/
func ReadMulti() ([]ResultMusic, []Mst_situation, bool) {
	var music_result []ResultMusic
	var situation_arr []Mst_situation

	if err := db.Table("musics").Debug().Select("musics.id, musics.name, musics.artist, musics.reason, musics.mst_situation_id, `mst_situations`.name AS Mst_situationName").Joins("INNER JOIN mst_situations AS `mst_situations` ON `musics`.mst_situation_id = `mst_situations`.id").Find(&music_result).Error; err != nil {
	    fmt.Println(err)
		return music_result, situation_arr, false
	}

	if err := db.Debug().Find(&situation_arr).Error; err != nil {
		return music_result, situation_arr, false
	}
	return music_result, situation_arr, true
}

/*
   パス：detail
*/
func Read(id string) (ResultMusic, bool) {
	var music_result ResultMusic
	// return music, false
	// ポインタを引数にしない場合はエラーになる
	// if err := db.Debug().First(&music_result, id).Error; err != nil {
	// 	fmt.Println("error happen!")
	// 	return music_result, true
	// }
	if err := db.Table("musics").Debug().Select("musics.id, musics.name, musics.artist, musics.reason, musics.mst_situation_id, `mst_situations`.name AS Mst_situationName").Joins("INNER JOIN mst_situations AS `mst_situations` ON `musics`.mst_situation_id = `mst_situations`.id").Find(&music_result).Error; err != nil {
	    fmt.Println(err)
		return music_result, false
	}
	fmt.Println(json.Marshal(music_result))
	
	return music_result, true
}
