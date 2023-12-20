package dto

import (
	"time"

	"playlistturbo.com/model"
)

type ImportSongs struct {
	Path          string `json:"path"`
	Recursive     bool   `json:"recursive"`
	SongExtension string `json:"songExtension"`
}

type MoveSongs struct {
	NewPath       string `json:"newPath"`
	OldPath       string `json:"oldPath"`
	Recursive     bool   `json:"recursive"`
	SongExtension string `json:"songExtension"`
}

type SongMetadata struct {
	Artist    string
	Album     string
	Title     string
	AlbumYear uint
	Lenght    string
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

type Songs struct {
	Artist      string `json:"artist"`
	Album       string `json:"album"`
	Title       string `json:"title"`
	TwonkyLink  string `json:"songUrl"`
	AlbumArtURI string `json:"albumArtUri"`
}

func ToDtoSongs(i model.Song, songDto Songs) Songs {
	dto := Songs{
		Artist:      i.Artist,
		Album:       i.Album,
		Title:       i.Title,
		TwonkyLink:  i.TwonkyLink,
		AlbumArtURI: i.AlbumArtURI,
	}
	return dto
}
