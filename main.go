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
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	// "unicode/utf8"
	// "reflect"
	// "time"
)

const (
	StatusBadRequest          = 400
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

	// main関数終了後にデーターベースへの接続を閉じる
	dbClose, _ := db.DB()
	defer dbClose.Close()

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

	// サーバーを起動する
	http.ListenAndServe(port, nil)
}

/*
   Top画面
*/
func top(w http.ResponseWriter, r *http.Request) {

	// クエリパラメーターを取得する
	var search string = r.URL.Query().Get("search")
	// 何も検索しない場合は空文字で送信される
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

	// ヘッダーをセットする
	w.Header().Set("Access-Control-Allow-Origin", os.Getenv("ORIGIN"))
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Content-Type", "application/json")

	// クエリパラメータを取得する
	var id string = r.URL.Query().Get("id")

	// ブラウザをリロードした際にクエリパラメータがundefindで送付される場合がある
	if r.Method != http.MethodGet || id == "undefined" {
		return
	}

	ret, situation, orm_err := models.Read(db, id)
	// エラー処理
	if !orm_err {
		log.Println("orm_error happen!")
		log.Println(orm_err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	// 結果を1つのJSONにまとめる
	var res = unify.ResponseDetail{Mst_situation: situation, Music: ret}
	outputJson, _ := json.Marshal(res)

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
	var userIDString string = r.URL.Query().Get("userID")

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

	// int型に変換する
	var situationID int
	var s string = r.URL.Query().Get("situation")
	situationID, _ = strconv.Atoi(s)
	var userIDInt int
	userIDInt, _ = strconv.Atoi(userIDString)

	// クエリパラメータに含まれた値を使用して構造体を初期化する。
	var create = unify.Music{Name: name, Reason: reason, Artist: artist, Mst_situationID: situationID, UserID: userIDInt}

	// レコードの作成
	ret := models.Register(db, &create)

	if !ret {
		log.Println("登録エラー")
		w.WriteHeader(http.StatusInternalServerError)
	}

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
	
	cookie, cookieErr := r.Cookie("token")

	// Cookieが送付されていない場合
	if cookieErr != nil {
		log.Println("Cookieが送付されていない")
		log.Println(cookieErr)
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, false)
		return
	}

	// Cookieに保存しているJWTをパースする
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
		// JWTが不正である場合は401を返却する
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

	if !retValidate {
		// 文字数が不正である場合は500エラーを返却する
		log.Println("validate_error happen!")
		log.Println(retValidate)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, false)
		return
	}

	// 入力した名前が既に存在しているか確認する
	existName := models.FindName(db, name)
	if existName {
		log.Println("名前が既に存在しています")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, false)
		return
	}

	// パスワードをハッシュ化する
	passwordByte := []byte(password)
	hashed, _ := bcrypt.GenerateFromPassword(passwordByte, 10)

	// クエリパラメータに含まれた値を使用して構造体を初期化する。
	var create = unify.User{Name: name, Password: string(hashed)}

	// // レコードの作成
	ret := models.SignUP(db, &create)

	if !ret {
		log.Println("登録エラー")
		w.WriteHeader(http.StatusInternalServerError)
	}

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

	// 名前がDBに存在しない場合は401を返却する
	if !orm_err {
		log.Println("入力した名前がDBに存在しません")
		log.Println(name)
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

	fmt.Fprint(w, string(outputJson))
}
