package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
  "reflect"
  "github.com/joho/godotenv"
	"os"
  "backend/unify"
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
    err := godotenv.Load(fmt.Sprintf("env/%s.env", os.Getenv("GO_ENV")))
    if err != nil {
        fmt.Println(err)
    }
    dsn := os.Getenv("DB_SET")
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    fmt.Println(reflect.TypeOf(db))
    if err != nil {
		panic("failed to connect database")
		} else {
      // down01(dsn, db);
      up01(dsn, db)
      insert01(dsn, db)
	}

    fmt.Println("End migrate!");
}

func up01(dsn string, db *gorm.DB) {
    fmt.Println("Start up01!");
    // charsetをutf8mb4にしないと、ORMをDBに接続した際のcharsetと合わずに文字列を登録すると「?」になる
    db.Set("gorm:table_options", "charset=utf8mb4").AutoMigrate(Mst_situation{})
    db.Set("gorm:table_options", "charset=utf8mb4").AutoMigrate(Music{})
    fmt.Println("End up01!");
}
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