package main

import (
	"backend/unify"
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	// "reflect"
)

type Mst_situation struct {
	gorm.Model
	Name   string
	Musics []Music
}

type Music struct {
	gorm.Model
	Name            string
	Artist          string
	Reason          string
	Mst_situationID int `gorm:"not null"`
	UserID          int
}

type User struct {
	gorm.Model
	Name     string
	Password string
}

type Like struct {
	gorm.Model
	MusicID int `gorm:"not null"`
	UserID  int `gorm:"not null"`
}

func main() {
	fmt.Println("Start migrate!")
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
		// down02(dsn, db);
		// down03(dsn, db);
		up01(dsn, db)
		up02(dsn, db)
		up03(dsn, db)
		insert01(dsn, db)
	}

	fmt.Println("End migrate!")
}

// migrate up関数
func up01(dsn string, db *gorm.DB) {
	fmt.Println("Start up01!")
	// charsetをutf8mb4にしないと、ORMをDBに接続した際のcharsetと合わずに文字列を登録すると「?」になる
	db.Set("gorm:table_options", "charset=utf8mb4").AutoMigrate(Mst_situation{})
	db.Set("gorm:table_options", "charset=utf8mb4").AutoMigrate(Music{})
	fmt.Println("End up01!")
}
func up02(dsn string, db *gorm.DB) {
	fmt.Println("Start up02!")
	db.Set("gorm:table_options", "charset=utf8mb4").AutoMigrate(User{})
	db.Set("gorm:table_options", "charset=utf8mb4").AutoMigrate(Like{})
	fmt.Println("End up02!")
}
func up03(dsn string, db *gorm.DB) {
	fmt.Println("Start up03!")
	db.Migrator().AddColumn(&Music{}, "UserID")
	fmt.Println("End up03!")
}

// migrate down関数
func down01(dsn string, db *gorm.DB) {
	fmt.Println("Start down01!")
	// テーブル削除
	db.Migrator().DropTable(&Music{})
	db.Migrator().DropTable(&Mst_situation{})
	fmt.Println("End down01!")
}
func down02(dsn string, db *gorm.DB) {
	fmt.Println("Start down02!")
	// テーブル削除
	db.Migrator().DropTable(&User{})
	db.Migrator().DropTable(&Like{})
	fmt.Println("End down02!")
}
func down03(dsn string, db *gorm.DB) {
	fmt.Println("Start down02!")
	db.Migrator().DropColumn(&Music{}, "UserID")
	fmt.Println("End down02!")
}

func insert01(dsn string, db *gorm.DB) {
	fmt.Println("Start insert01!")
	var create = unify.Mst_situation{Name: "name"}
	if orm_err := db.Debug().Create(&create).Error; orm_err != nil {
		fmt.Println("登録関数エラー")
		fmt.Println(orm_err)
	}
}
