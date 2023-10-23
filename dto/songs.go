package dto

type ImportSongs struct {
	Path          string `json:"path"`
	Recursive     bool   `json:"recursive"`
	SongExtension string `json:"songExtension"`
}

type SongMetadata struct {
	Artist    string
	Album     string
	Title     string
	AlbumYear uint
	Lenght    string
	LenghtSec float64
	GenreID   uint
}
