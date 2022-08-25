package models

import (
	// "encoding/json"
	"fmt"
	"gorm.io/gorm"
	"backend/unify"
)

/*
   パス：top
*/
func ReadMulti(db *gorm.DB) ([]unify.ResultMusic, []unify.Mst_situation, bool) {
	var music []unify.ResultMusic
	var situation_arr []unify.Mst_situation
	// return music, situation_arr, false
	// Musicテーブルのレコードを取得する
	if err := db.Table("musics").Debug().Select("musics.id, musics.name, musics.artist, musics.reason, musics.mst_situation_id, `mst_situations`.name AS Mst_situationName").Joins("INNER JOIN mst_situations AS `mst_situations` ON `musics`.mst_situation_id = `mst_situations`.id").Find(&music).Error; err != nil {
	    fmt.Println(err)
		return music, situation_arr, false
	}

	// Mst_situationテーブルのレコードを取得する
	if err := db.Debug().Find(&situation_arr).Error; err != nil {
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
		fmt.Println(err)
		return music, false
	}
	
	return music, true
}
