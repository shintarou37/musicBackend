package main

import (
	"fmt"
	"net/http"
	"testing"
	// "reflect"
	"backend/unify"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os/exec"
	"time"
	// "net/http/cookiejar"
	// "net/http/httptest"
	// "strings"
	// "net/url"
)

/*
   Server Start
*/
func TestServer(t *testing.T) {
	fmt.Println("server start")
	err := exec.Command("go", "run", "main.go").Start()
	if err != nil {
		t.Error("[SERVER ERROR]", "want nil : ", err)
	}
	cmd := exec.Command("go", "run", "main.go")
	cmd.Start()
	time.Sleep(30 * time.Second)
}

/*
   登録機能（最小値 正常系 ログインなし）
*/
func TestRegisterLeast(t *testing.T) {
	var byte []byte
	resp, err := http.Post("http://127.0.0.1:8080/register?name=a&artist=a&reason=a&situation=1", "application/json", bytes.NewBuffer(byte))
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
		t.Error("[HTTP REQUEST ERROR]", "want nil : ", err)
	}
	defer resp.Body.Close()

	// ステータスコードを確認する
	if resp.Status != "200 OK" {
		t.Error("[STATUS CODE ERROR]", "want 200 OK : ", resp.Status)
	}
}

/*
   登録機能（最小値 異常系）
*/
func TestRegisterLeastFailure(t *testing.T) {
	var byte []byte
	resp, err := http.Post("http://127.0.0.1:8080/register?name=&artist=&reason=&situation=1", "application/json", bytes.NewBuffer(byte))
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
		t.Error("[HTTP REQUEST ERROR]", "want nil", err)
	}
	defer resp.Body.Close()
	fmt.Println(resp.Status)

	// ステータスコードを確認する
	if resp.Status != "400 Bad Request" {
		t.Error("[STATUS CODE ERROR]", "want 400 Bad Request : ", resp.Status)
	}
}

/*
   登録機能（最大値 正常系）
*/
func TestRegisterLargest(t *testing.T) {
	var byte []byte
	resp, err := http.Post("http://127.0.0.1:8080/register?name=ABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUV&artist=ABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUV&reason=ABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKL&situation=1", "application/json", bytes.NewBuffer(byte))
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
		t.Error("[HTTP REQUEST ERROR]", "want nil", err)
	}
	defer resp.Body.Close()

	// ステータスコードを確認する
	if resp.Status != "200 OK" {
		t.Error("[STATUS CODE ERROR]", "want 200 OK : ", resp.Status)
	}
}

/*
   登録機能（最大値 異常系）
*/
func TestRegisterLargestFailure(t *testing.T) {
	var byte []byte
	resp, err := http.Post("http://127.0.0.1:8080/register?name=ABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUV1&artist=ABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUV1&reason=ABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKL1&situation=1", "application/json", bytes.NewBuffer(byte))
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
		t.Error("[HTTP REQUEST ERROR]", "want nil", err)
	}
	defer resp.Body.Close()

	// ステータスコードを確認する
	if resp.Status != "400 Bad Request" {
		t.Error("[STATUS CODE ERROR]", "want 400 Bad Request : ", resp.Status)
	}
}

/*
   Top画面
*/
func TestTop(t *testing.T) {
	resp, err := http.Get("http://127.0.0.1:8080")
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
		t.Error("[HTTP REQUEST ERROR]", "want nil", err)
	}
	defer resp.Body.Close()

	// リクエストBodyを取得する
	body, err := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(body))

	// JSONを変換する
	var result_struct unify.ResponseTop
	if err := json.Unmarshal(body, &result_struct); err != nil {
		t.Error("[JSON UNMARSHAL EROOR]", "want nil : ", err)
	}

	// 複数の投稿がなされていること
	if len(result_struct.Music) < 2 {
		t.Error("[RESULT Music Length ERROR]", "want 2 OR MORE : ", len(result_struct.Music))
	}
	// IDカラムの降順で値を取得できていること
	if result_struct.Music[0].ID < result_struct.Music[1].ID {
		t.Error("[RESULT DATA ID DESC ERROR]", "want TRUE : ", result_struct.Music[0].ID > result_struct.Music[1].ID)
	}

	if resp.Status != "200 OK" {
		t.Error("[STATUS CODE ERROR]", "want 200 OK : ", resp.Status)
	}
}

/*
   Top画面（検索あり）
*/
func TestTopSearch(t *testing.T) {
	resp, err := http.Get("http://127.0.0.1:8080?search=1")
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
		t.Error("[HTTP REQUEST ERROR]", "want nil", err)
	}
	defer resp.Body.Close()

	// リクエストBodyを取得する
	body, err := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(body))

	// JSONを変換する
	var result_struct unify.ResponseTop
	if err := json.Unmarshal(body, &result_struct); err != nil {
		t.Error("[JSON UNMARSHAL EROOR]", "want nil : ", err)
	}

	// 複数の投稿がなされていること（以前のテストで1のシチュエーションに複数の投稿を行っている）
	if len(result_struct.Music) < 2 {
		t.Error("[RESULT Music Length ERROR]", "want 2 OR MORE : ", len(result_struct.Music))
	}
	// IDカラムの降順で値を取得できていること
	if result_struct.Music[0].ID < result_struct.Music[1].ID {
		t.Error("[RESULT DATA ID DESC ERROR]", "want TRUE : ", result_struct.Music[0].ID > result_struct.Music[1].ID)
	}

	if resp.Status != "200 OK" {
		t.Error("[STATUS CODE ERROR]", "want 200 OK", resp.Status)
	}
}

/*
   詳細画面（正常系）
*/
func TestDetail(t *testing.T) {
	resp, err := http.Get("http://127.0.0.1:8080/detail?id=1")
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
		t.Error("[HTTP REQUEST ERROR]", "want nil", err)
	}
	defer resp.Body.Close()

	// リクエストBodyを取得する
	body, err := ioutil.ReadAll(resp.Body)

	// JSONを変換する
	var result_struct unify.ResponseDetail
	if err := json.Unmarshal(body, &result_struct); err != nil {
		t.Error("[JSON UNMARSHAL EROOR]", "want nil : ", err)
	}

	// 取得するIDカラムを確認する
	if result_struct.Music.ID != 1 {
		t.Error("[RESULT DATA ID ERROR]", "want 1 : ", result_struct.Music.ID)
	}

	// ステータスコードを確認する
	if resp.Status != "200 OK" {
		t.Error("[STATUS CODE ERROR]", "want 200 OK : ", resp.Status)
	}
}

/*
   詳細画面（パラメータ undefined）
*/
func TestDetailUndefind(t *testing.T) {
	resp, err := http.Get("http://127.0.0.1:8080/detail?id=undefined")
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
		t.Error("[HTTP REQUEST ERROR]", "want nil", err)
	}
	defer resp.Body.Close()

	// リクエストBodyを取得する
	body, _ := ioutil.ReadAll(resp.Body)

	// Bodyが空であることをを確認する
	if len(body) != 0 {
		t.Error("[RESULT BODY ERROR]", "want 0 : ", len(body))
	}
	// ステータスコードを確認する
	if resp.Status != "200 OK" {
		t.Error("[STATUS CODE ERROR]", "want 200 OK : ", resp.Status)
	}
}

/*
   詳細画面（不正なパラメーター）
*/
func TestDetailError(t *testing.T) {
	resp, err := http.Get("http://127.0.0.1:8080/detail?id=ERROR")
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
		t.Error("[HTTP REQUEST ERROR]", "want nil", err)
	}
	defer resp.Body.Close()

	// ステータスコードを確認する
	if resp.Status != "500 Internal Server Error" {
		t.Error("[STATUS CODE ERROR]", "want 500 Internal Server Error : ", resp.Status)
	}
}
/*
   更新機能（異常系 Cookieなし）
*/
func TestUpdateCookieFailure(t *testing.T) {
	var byte []byte
	resp, err := http.Post("http://127.0.0.1:8080/update?id=1&name=a&artist=a&reason=a&situation=1", "application/json", bytes.NewBuffer(byte))
	if err != nil {
		t.Error("[HTTP REQUEST ERROR]", "want nil", err)
	}
	defer resp.Body.Close()

	// Cookieがないので400エラーが返却される
	if resp.Status != "400 Bad Request" {
		t.Error("[STATUS CODE ERROR]", "want 400 Bad Request : ", resp.Status)
	}
}
/*
   更新機能（異常系 Cookieあり 最小値）
*/
func TestUpdateLeastFailure(t *testing.T) {
	var client http.Client
	cookie := &http.Cookie{
        Name:   "token",
        Value:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.jezpHBmixG797D1iZt3ihjOD4p01Bignvv7sUxZP4xo",
    }

    req, _ := http.NewRequest("POST", "http://127.0.0.1:8080/update?id=1&name=&artist=&reason=&situation=1", nil)
    req.AddCookie(cookie)

    resp, err := client.Do(req)
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
		t.Error("[HTTP REQUEST ERROR]", "want nil", err)
	}
    defer resp.Body.Close()

	if resp.Status != "400 Bad Request" {
		t.Error("[STATUS CODE ERROR]", "want 400 Bad Request : ", resp.Status)
	}
}
/*
   更新機能（正常系 Cookieあり最小値）
*/
func TestUpdateCookie(t *testing.T) {
	var client http.Client
	cookie := &http.Cookie{
        Name:   "token",
        Value:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.jezpHBmixG797D1iZt3ihjOD4p01Bignvv7sUxZP4xo",
    }

    req, _ := http.NewRequest("POST", "http://127.0.0.1:8080/update?id=1&name=a&artist=a&reason=a&situation=1", nil)
    req.AddCookie(cookie)
    resp, err := client.Do(req)
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
		t.Error("[HTTP REQUEST ERROR]", "want nil", err)
	}
    defer resp.Body.Close()

    if resp.Status != "200 OK" {
		t.Error("[STATUS CODE ERROR]", "200 OK : ", resp.Status)
	}
}
/*
   更新機能（異常系 Cookieあり 最大値）
*/
func TestUpdateLargestFailure(t *testing.T) {
	var client http.Client
	cookie := &http.Cookie{
        Name:   "token",
        Value:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.jezpHBmixG797D1iZt3ihjOD4p01Bignvv7sUxZP4xo",
    }

    req, _ := http.NewRequest("POST", "http://127.0.0.1:8080/update?id=1&name=ABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUV1&artist=ABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUV1&reason=ABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKL1&situation=1", nil)
    req.AddCookie(cookie)

    resp, err := client.Do(req)
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
		t.Error("[HTTP REQUEST ERROR]", "want nil", err)
	}
    defer resp.Body.Close()

	if resp.Status != "400 Bad Request" {
		t.Error("[STATUS CODE ERROR]", "want 400 Bad Request : ", resp.Status)
	}
}
