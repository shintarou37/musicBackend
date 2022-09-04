package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
  "reflect"
  "github.com/joho/godotenv"
	"os"
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
    db_set := os.Getenv("DB_SET")
    dsn := db_set
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    fmt.Println(reflect.TypeOf(db))
    if err != nil {
		panic("failed to connect database")
		} else {
      // down01(dsn, db);
      up01(dsn, db)
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
