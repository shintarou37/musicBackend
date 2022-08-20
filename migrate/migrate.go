package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
  "reflect"
)

type Music struct {
  gorm.Model
	Name    string
	Reason  string
}

func main() {
    fmt.Println("Start migrate!");
    dsn := "root:@tcp(127.0.0.1:3306)/music?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    fmt.Println(reflect.TypeOf(db))
    if err != nil {
		panic("failed to connect database")
		} else {
      down01(dsn, db);
      up01(dsn, db)
	}

    fmt.Println("End migrate!");
}

func up01(dsn string, db *gorm.DB) {
    fmt.Println("Start up01!");
    // charsetをutf8mb4にしないと、ORMをDBに接続した際のcharsetと合わずに文字列を登録すると「?」になる
    db.Set("gorm:table_options", "charset=utf8mb4").AutoMigrate(Music{})
    fmt.Println("End up01!");

}
func down01(dsn string, db *gorm.DB) {
    fmt.Println("Start down01!");
    // テーブル削除
    db.Migrator().DropTable(&Music{})
    fmt.Println("End down01!");
}
