package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"backend/unify"
	"backend/models"
	"backend/validates"
	"github.com/joho/godotenv"
	"os"
	"log"
	"io"
	// "unicode/utf8"
	// "reflect"
	"github.com/golang-jwt/jwt/v4"
	// "time"
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
	log.Println("Server Start!")

	// 環境変数の読み込みを行う
	err := godotenv.Load(fmt.Sprintf("env/%s.env", os.Getenv("GO_ENV")))
	if err != nil {
		log.Println("環境変数の読み込みに失敗")
		log.Println(err)
	}
	dsn := os.Getenv("DB_SET")
	port := os.Getenv("PORT")

	// データーベースに接続する
	db, db_err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if db_err != nil {
		log.Println("データーベースの接続に失敗")
		log.Println(db_err)
	}

	// ログファイルの設定をする
	log.SetOutput(io.MultiWriter(logfile, os.Stdout))
	log.SetFlags(log.Ldate | log.Ltime)

	// HTTP handler
	http.HandleFunc("/", top)
	http.HandleFunc("/token", token)
	http.HandleFunc("/te", te)
	http.HandleFunc("/detail", detail)
	http.HandleFunc("/register", register)
	// サーバー起動を起動する
	http.ListenAndServe(port, nil)
}

func token(w http.ResponseWriter, r *http.Request) {
	// トークン生成
	token := jwt.New(jwt.SigningMethodHS256)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	// トークンに電子署名を追加する
	tokenString, _ := token.SignedString([]byte(os.Getenv("SIGNINGKEY")))

	// Cookieに追加する
	cookie := &http.Cookie{
		Name: "hoge",
		Value: tokenString,
		MaxAge: 30 * 10,
	 }
	http.SetCookie(w, cookie)

	// JWTを返却
	w.Write([]byte(tokenString))
}
func te(w http.ResponseWriter, r *http.Request) {
	fmt.Println("----teに来た")
	cookie, _ := r.Cookie("hoge")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	// tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.jezpHBmixG797D1iZt3ihjOD4p01Bignvv7sUxZP4xo"
 
	// パースする
	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("SIGNINGKEY")), nil
	})
	
	// 検証
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println("検証成功")
		fmt.Println(claims)
	} else {
		fmt.Println("検証失敗")
		fmt.Println(err)
	}
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
		log.Println("ReadMulti error happen!")
		log.Println(orm_err)
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
	outputJson, _ := json.Marshal(ret)

	// エラー処理
	if !orm_err {
		log.Println("orm_error happen!")
		log.Println(orm_err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	// jsonデータを返却する
	fmt.Fprint(w, string(outputJson))
}

/*
   登録機能
*/
func register(w http.ResponseWriter, r *http.Request) {

	// 登録機能時にOPTIONSリクエストが送付される
	fmt.Println(r.Method)
	if r.Method != http.MethodPost {
		return
	}

	// ヘッダーをセットする
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	// クエリパラメータを受け取る
	var name string = r.URL.Query().Get("name")
	var artist string = r.URL.Query().Get("artist")
	var reason string = r.URL.Query().Get("reason")

	// 文字数チェック
	retValidate := validates.Register(name, artist, reason)

	if !retValidate {
		// 文字数が不正である場合は400エラーを返却する
		log.Println("validate_error happen!")
		log.Println(retValidate)
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
		log.Println("登録エラー")
		w.WriteHeader(http.StatusInternalServerError)
	}

	// データを返却する
	fmt.Fprint(w, true)
}
