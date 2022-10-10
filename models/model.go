package models

import (
	"strconv"
	"gorm.io/gorm"
	"backend/unify"
	"log"
	// "reflect"
	// "fmt"
	// "encoding/json"
)

/*
   パス：top
*/
func ReadMulti(db *gorm.DB, search string) ([]unify.ResultMusic, []unify.Mst_situation, bool) {
	// 構造体を初期化する
	var music []unify.ResultMusic
	var musicSearch = unify.ResultMusic{}
	var situation_arr []unify.Mst_situation

	// クエリパラメータにsearchがある場合"Mst_situationID"を検索する
	if search != ""{
		var Mst_situationID, _ = strconv.Atoi(search)
		musicSearch.Mst_situationID = Mst_situationID
	}

	// return music, situation_arr, false
	// Musicテーブルのレコードを取得する
	if err := db.Table("musics").Debug().Select("musics.id, musics.name, musics.artist, musics.reason, musics.mst_situation_id, `mst_situations`.name AS Mst_situationName").Joins("INNER JOIN mst_situations AS `mst_situations` ON `musics`.mst_situation_id = `mst_situations`.id").Order("musics.id asc").Find(&music, musicSearch).Error; err != nil {
	  log.Println("ReadMult関数のmusicsテーブルのデータ取得時にエラー")
	  log.Println(err)
		return music, situation_arr, false
	}

	// Mst_situationテーブルのレコードを取得する
	if err := db.Debug().Find(&situation_arr).Error; err != nil {
		log.Println("ReadMult関数のmst_situationテーブルのデータ取得時にエラー")
	  log.Println(err)
		return music, situation_arr, false
	}

	return music, situation_arr, true
}

/*
   パス：detail
*/

func Read(db *gorm.DB, id string) (unify.ResultMusic, bool) {
	var music unify.ResultMusic
	// return music, false
	// テーブル名を指定しないと構造体の名称「ResultMusic」をテーブル名をみなす
	if err := db.Debug().Table("musics").Select("musics.*, `mst_situations`.name AS Mst_situationName").Joins("INNER JOIN mst_situations AS `mst_situations` ON `musics`.mst_situation_id = `mst_situations`.id").First(&music, id).Error; err != nil {
		log.Println("Read関数のmusicsテーブルのデータ取得時にエラー")
	  log.Println(err)
		return music, false
	}
	
	return music, true
}

/*
   パス：register
*/

func Register(db *gorm.DB, create *unify.Music) (bool) {
	if orm_err := db.Debug().Create(&create).Error; orm_err != nil {
		log.Println("Register関数エラー")
	  log.Println(orm_err)
		return false
	}

	return true
}