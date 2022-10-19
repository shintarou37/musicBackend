package models

import (
	"backend/unify"
	"gorm.io/gorm"
	"log"
	"strconv"
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

	// クエリパラメータにsearchがある場合
	if search != "" {
		var Mst_situationID, _ = strconv.Atoi(search)
		musicSearch.Mst_situationID = Mst_situationID
	}

	// Musicテーブルのレコードを取得する
	if err := db.Table("musics").Debug().Order("id desc").Select("musics.id, musics.name, musics.artist, musics.reason, musics.user_id, musics.mst_situation_id, `mst_situations`.name AS Mst_situationName").Joins("INNER JOIN mst_situations AS `mst_situations` ON `musics`.mst_situation_id = `mst_situations`.id").Order("musics.id asc").Find(&music, musicSearch).Error; err != nil {
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

func Read(db *gorm.DB, id string) (unify.ResultMusic, []unify.Mst_situation, bool) {
	var music unify.ResultMusic
	var situation_arr []unify.Mst_situation

	// テーブル名を指定しないと構造体の名称「ResultMusic」をテーブル名をみなす
	if err := db.Debug().Table("musics").Select("musics.*, `mst_situations`.name AS Mst_situationName, `users`.name AS UserName").Joins("INNER JOIN mst_situations AS `mst_situations` ON `musics`.mst_situation_id = `mst_situations`.id").Joins("LEFT OUTER JOIN users AS `users` ON `musics`.user_id = `users`.id").First(&music, id).Error; err != nil {
		log.Println("Read関数のmusicsテーブルのデータ取得時にエラー")
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
   パス：register
*/

func Register(db *gorm.DB, create *unify.Music) bool {
	if orm_err := db.Debug().Create(&create).Error; orm_err != nil {
		log.Println("Register関数エラー")
		log.Println(orm_err)
		return false
	}

	return true
}

/*
   パス：update
*/
func Update(db *gorm.DB, id, name, reason, artist string, situationID int) bool {
	if orm_err := db.Debug().Model(&unify.Music{}).Where("id = ?", id).Updates(unify.Music{Name: name, Reason: reason, Artist: artist, Mst_situationID: situationID}).Error; orm_err != nil {
		log.Println("Update関数エラー")
		log.Println(orm_err)
		return false
	}

	return true
}

/*
   パス：signup
*/
func SignUP(db *gorm.DB, create *unify.User) bool {
	if orm_err := db.Debug().Create(&create).Error; orm_err != nil {
		log.Println("SignUP関数エラー")
		log.Println(orm_err)
		return false
	}

	return true
}

/*
   パス：signup
*/
func FindName(db *gorm.DB, name string) bool {
	var user unify.User

	if err := db.Debug().Table("users").Select("users.id").First(&user, "name = ?", name).Error; err != nil {
		log.Println("FindName関数エラー")
		log.Println(err)
		return false
	}
	return true
}

/*
   パス：signin
*/
func FindUser(db *gorm.DB, name string) (unify.SignInRet, bool) {
	var user unify.SignInRet
	// return user, false
	if err := db.Debug().Table("users").Select("users.*").First(&user, "name = ?", name).Error; err != nil {
		log.Println("FindUser関数エラー")
		log.Println(err)
		return user, false
	}

	return user, true
}
