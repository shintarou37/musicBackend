package unify

import (
	"gorm.io/gorm"
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
	Mst_situationID int
}

type ResponseTop struct {
	Mst_situation []Mst_situation
	Music         []ResultMusic
}

type ResultMusic struct {
	gorm.Model
	Name            	string
	Artist          	string
	Reason          	string
	Mst_situationID		int
	Mst_situationName   string
}

type User struct {
  gorm.Model
  Name            string
  Password        []byte
  Mst_situationID int
}

type Like struct {
  gorm.Model
  MusicID    int
  UserID     int
}