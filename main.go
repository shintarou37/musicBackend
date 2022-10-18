package main

import (
	"backend/models"
	"backend/unify"
	"backend/validates"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	// "unicode/utf8"
	// "reflect"
	"github.com/golang-jwt/jwt/v4"
	// "time"
	"golang.org/x/crypto/bcrypt"
)

const (
	StatusBadRequest		  = 400
	StatusNotAcceptable       = 406
	StatusUnauthorized        = 401
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
	http.HandleFunc("/detail", detail)
	http.HandleFunc("/update", update)
	http.HandleFunc("/register", register)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/signin", signin)
	// サーバー起動を起動する
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
	w.Header().Set("Access-Control-Allow-Origin", os.Getenv("ORIGIN"))
	w.Header().Set("Access-Control-Allow-Credentials", "true")
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
	w.Header().Set("Access-Control-Allow-Origin", os.Getenv("ORIGIN"))
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Content-Type", "application/json")

	// クエリパラメータ「id」を取得する
	var id string = r.URL.Query().Get("id")

	// ブラウザをリロードした際にクエリパラメータがundefindで送付される場合がある
	if r.Method != http.MethodGet || id == "undefined" {
		return
	}

	ret, situation, orm_err := models.Read(db, id)
	var res = unify.ResponseDetail{Mst_situation: situation, Music: ret}
	// jsonエンコード
	outputJson, _ := json.Marshal(res)

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
	w.Header().Set("Access-Control-Allow-Origin", os.Getenv("ORIGIN"))
	w.Header().Set("Access-Control-Allow-Credentials", "true")

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

/*
   編集機能
*/
func update(w http.ResponseWriter, r *http.Request) {
	// ヘッダーをセットする
	w.Header().Set("Access-Control-Allow-Origin", os.Getenv("ORIGIN"))
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	// 更新機能時にOPTIONSリクエストが送付される
	fmt.Println(r.Method)
	if r.Method != http.MethodPost {
		return
	}

	// Cookieに保存しているJWTをパースする
	cookie, _ := r.Cookie("token")
	token, _ := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SIGNINGKEY")), nil
	})
	
	// JWT検証
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		log.Println("編集機能 JWT検証成功")
	} else {
		log.Println("編集機能 JWT検証失敗")
		log.Println(claims)
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, false)
		return
	}

	// クエリパラメータを受け取る
	var id string = r.URL.Query().Get("id")
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

	// レコードの作成
	ret := models.Update(db, id, name, reason, artist, situationID)

	if !ret {
		log.Println("更新エラー")
		w.WriteHeader(http.StatusInternalServerError)
	}

	// データを返却する
	fmt.Fprint(w, true)
}
/*
   利用者登録機能
*/
func signup(w http.ResponseWriter, r *http.Request) {

	// ヘッダーをセットする
	w.Header().Set("Access-Control-Allow-Origin", os.Getenv("ORIGIN"))
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	// クエリパラメータを受け取る
	var name string = r.URL.Query().Get("name")
	var password string = r.URL.Query().Get("password")

	// 文字数チェック
	retValidate := validates.SignUp(name, password)

	passwordByte := []byte(password)
	hashed, _ := bcrypt.GenerateFromPassword(passwordByte, 10)

	if !retValidate {
		// 文字数が不正である場合は400エラーを返却する
		log.Println("validate_error happen!")
		log.Println(retValidate)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, false)
		return
	}

	// クエリパラメータに含まれた値を使用して構造体を初期化する。
	var create = unify.User{Name: name, Password: string(hashed)}

	// // レコードの作成
	ret := models.SignUP(db, &create)

	if !ret {
		log.Println("登録エラー")
		w.WriteHeader(http.StatusInternalServerError)
	}

	// データを返却する
	fmt.Fprint(w, true)
}

/*
   ログイン機能
*/
func signin(w http.ResponseWriter, r *http.Request) {
	// ヘッダーをセットする
	w.Header().Set("Access-Control-Allow-Origin", os.Getenv("ORIGIN"))
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	// クエリパラメータを受け取る
	var name string = r.URL.Query().Get("name")

	// 入力した名前をDBから取得する
	ret, orm_err := models.FindUser(db, name)

	// 名前がDBに存在しない場合
	if !orm_err {
		log.Println("名前が違います!")
		log.Println(orm_err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var password string = r.URL.Query().Get("password")
	passwordByte := []byte(password)

	// 第一引数にDBに保存しているカラムの値、第２引数に入力したパスワードをbyte型に変更して確認する
	err := bcrypt.CompareHashAndPassword([]byte(ret.Password), passwordByte)
	if err != nil {
		log.Println("パスワードが違います")
		log.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	// トークン生成
	token := jwt.New(jwt.SigningMethodHS256)
	// トークンに電子署名を追加する
	tokenString, _ := token.SignedString([]byte(os.Getenv("SIGNINGKEY")))

	ret.Token = tokenString
	outputJson, _ := json.Marshal(ret)
	// データを返却する
	fmt.Fprint(w, string(outputJson))
}