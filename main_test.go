package main

import (
	"testing"
	"fmt"
    "net/http"
    "reflect"
    "io/ioutil"
    // "encoding/json"
    // "backend/unify"
)

/*
   Top画面
*/
func TestTop(t *testing.T) {
    fmt.Println("Top画面")
    resp, err := http.Get("http://localhost:8080")
    if err != nil {
        fmt.Println("err")
        fmt.Println(err)
        t.Error("[HTTP REQUEST ERROR]", "want nil", err)
    }
    defer resp.Body.Close()

    if resp.Status != "200 OK"{
        t.Error("[STATUS CODE ERROR]", "want 200 OK : ", resp.Status)
    }
}

/*
   Top画面（検索あり）
*/
func TestTopSearch(t *testing.T) {
    fmt.Println("Top画面（検索あり）")
    resp, err := http.Get("http://localhost:8080?search=1")
    if err != nil {
        fmt.Println("err")
        fmt.Println(err)
        t.Error("[HTTP REQUEST ERROR]", "want nil", err)
    }
    defer resp.Body.Close()

    if resp.Status != "200 OK"{
        t.Error("[STATUS CODE ERROR]", "want 200 OK", resp.Status)
    }
}

// /*
//    詳細画面（正常系）
// */
func TestDetail(t *testing.T) {
    fmt.Println("詳細画面（正常系）")
    resp, err := http.Get("http://localhost:8080/detail?id=1")
    if err != nil {
        fmt.Println("err")
        fmt.Println(err)
        t.Error("[HTTP REQUEST ERROR]", "want nil", err)
    }
    defer resp.Body.Close()
    b, err := ioutil.ReadAll(resp.Body)
    fmt.Println(reflect.TypeOf(b))

    if resp.Status != "200 OK"{
        t.Error("[STATUS CODE ERROR]", "want 200 OK : ", resp.Status)
    }
}
/*
   詳細画面（パラメータ undefined）
*/
func TestDetailUndefind(t *testing.T) {
    fmt.Println("詳細画面（異常系）")
    resp, err := http.Get("http://localhost:8080/detail?id=undefined")
    if err != nil {
        fmt.Println("err")
        fmt.Println(err)
        t.Error("[HTTP REQUEST ERROR]", "want nil", err)
    }
    defer resp.Body.Close()

    if resp.Status != "200 OK"{
        t.Error("[STATUS CODE ERROR]", "want 200 OK : ", resp.Status)
    }
}

/*
   詳細画面（不正なパラメーター）
*/
func TestDetailError(t *testing.T) {
    fmt.Println("詳細画面（不正なパラメーター）")
    resp, err := http.Get("http://localhost:8080/detail?id=ERROR")
    if err != nil {
        fmt.Println("err")
        fmt.Println(err)
        t.Error("[HTTP REQUEST ERROR]", "want nil", err)
    }
    defer resp.Body.Close()

    if resp.Status != "500 Internal Server Error"{
        t.Error("[STATUS CODE ERROR]", "want 500 Internal Server Error : ", resp.Status)
    }
}