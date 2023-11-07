package dto

import "time"

type ImportSongs struct {
	Genre     string `json:"genre"`
	Artist    string `json:"artist"`
	Album     string `json:"album"`
	Extension string `json:"extension"`
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

type ImportedSongs struct {
	Message   string
	ArtistQty int
	AlbumQty  int
	SongQty   int
	StartTime time.Time
	EndTime   time.Time
	Duration  string
}

type ImportAlbums struct {
	AlbumURL   string
	AlbumTitle string
}
