package main

import (
	"testing"
	"fmt"
    "net/http"
    // "reflect"
    // "io/ioutil"
)

/*
   Top画面
*/
func TestTop(t *testing.T) {
    resp, err := http.Get("http://localhost:8080")
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

/*
   Top画面（検索あり）
*/
func TestSearch(t *testing.T) {
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