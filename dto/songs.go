package dto

import "time"

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

type ImportedSongs struct {
	Message   string
	ArtistQty int
	AlbumQty  int
	SongQty   int
	StartTime time.Time
	EndTime   time.Time
	Duration  string
}

type SongExtraTable struct {
	Filename    string
	URL         string
	Frequency   int
	Bitrate     int
	AlbumArtURI string
	Duration    string
}
