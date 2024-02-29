package model

type Genre struct {
	ID   uint   `gorm:"type:uint;primaryKey;autoIncrement"`
	Name string `gorm:"type:varchar(50)"`
}

//https://mutagen-specs.readthedocs.io/en/latest/id3/id3v2.4.0-frames.html
