package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
  "github.com/joho/godotenv"
	"os"
  "backend/unify"
  // "reflect"
)

type Mst_situation struct {
  gorm.Model
	Name    string
  Musics  []Music
}

type Music struct {
  gorm.Model
	Name    string
  Artist  string
	Reason  string
  Mst_situationID int `gorm:"not null"`
}

func main() {
    fmt.Println("Start migrate!");
    // データーベースに接続する
    err := godotenv.Load(fmt.Sprintf("env/%s.env", os.Getenv("GO_ENV")))
    if err != nil {
        fmt.Println(err)
    }
    dsn := os.Getenv("DB_SET")
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
		  panic("failed to connect database")
    // エラーが発生しなかった場合にmigrateを実行する
		} else {
      // down01(dsn, db);
      up01(dsn, db)
      insert01(dsn, db)
	}

    fmt.Println("End migrate!");
}

// migrate up関数
func up01(dsn string, db *gorm.DB) {
    fmt.Println("Start up01!");
    // charsetをutf8mb4にしないと、ORMをDBに接続した際のcharsetと合わずに文字列を登録すると「?」になる
    db.Set("gorm:table_options", "charset=utf8mb4").AutoMigrate(Mst_situation{})
    db.Set("gorm:table_options", "charset=utf8mb4").AutoMigrate(Music{})
    fmt.Println("End up01!");
}
// migrate down関数
func down01(dsn string, db *gorm.DB) {
    fmt.Println("Start down01!");
    // テーブル削除
    db.Migrator().DropTable(&Music{})
    db.Migrator().DropTable(&Mst_situation{})
    fmt.Println("End down01!");
}

func insert01(dsn string, db *gorm.DB) {
  fmt.Println("Start insert01!");
  var create = unify.Mst_situation{Name: "name"}
	if orm_err := db.Debug().Create(&create).Error; orm_err != nil {
		fmt.Println("登録関数エラー")
		fmt.Println(orm_err)
	}
}