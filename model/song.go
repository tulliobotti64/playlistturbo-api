package model

import (
	"time"

	"github.com/gofrs/uuid"
)

type Song struct {
	ID              uuid.UUID `gorm:"not null;type:uuid;primary_key;default:uuid_generate_v4()"`
	GenreID         uint      `gorm:"type:uint;index;column:genre"`
	Genre           Genre     `gorm:"foreignKey:genre"`
	GenreTag        string    `gorm:"not null;type:varchar(100)"`
	Artist          string    `gorm:"not null;type:varchar(200)"`
	Album           string    `gorm:"not null;type:varchar(200)"`
	Title           string    `gorm:"not null;type:varchar(200)"`
	Lenght          string    `gorm:"lenght;type:varchar(20)"`
	AlbumDate       string    `gorm:"not null;type:varchar(200)"`
	FilePath        string    `gorm:"not null;type:varchar(500);uniqueIndex"`
	TwonkyLink      string    `gorm:"not null;type:varchar(500)"`
	Favorite        bool      `gorm:"not null;type:bool;default:false"`
	TrackNumber     int       `gorm:"type:uint"`
	Format          string    `gorm:"not null;type:varchar(200)"`
	SampleFrequency int       `gorm:"type:uint"`
	Bitrate         int       `gorm:"type:uint"`
	AlbumArtURI     string    `gorm:"not null;type:varchar(200)"`
	ListenQty       uint16    `gorm:"type:uint"`
	UpdatedAt       time.Time
}
