package dto

import (
	"time"

	"github.com/gofrs/uuid"
	"playlistturbo.com/model"
)

type ImportSongs struct {
	Path             string `json:"path"`
	Recursive        bool   `json:"recursive"`
	SongExtension    string `json:"songExtension"`
	GenreFromPath    bool   `json:"genreFromPath"`
	GenreArtistAlbum string `json:"genreArtistAlbum"`
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
	ID          uuid.UUID `json:"id"`
	Genre       string    `json:"genre"`
	Artist      string    `json:"artist"`
	Album       string    `json:"album"`
	Title       string    `json:"title"`
	TwonkyLink  string    `json:"songUrl"`
	AlbumArtURI string    `json:"albumArtUri"`
	Favorite    bool      `json:"favorite"`
	AlbumDate   uint      `json:"albumDate"`
}

func ToDtoSongs(i model.Song, songDto Songs) Songs {
	dto := Songs{
		ID:          i.ID,
		Genre:       i.GenreTag,
		Artist:      i.Artist,
		Album:       i.Album,
		Title:       i.Title,
		TwonkyLink:  i.TwonkyLink,
		AlbumArtURI: i.AlbumArtURI,
		Favorite:    i.Favorite,
		AlbumDate:   i.AlbumDate,
	}
	return dto
}

type List struct {
	ID   uint   `json:"ID"`
	Name string `json:"Name"`
}

type Favorites struct {
	Genre  string `json:"genre"`
	Artist string `json:"artist"`
}

type SongsByArtist struct {
	Artist string `json:"artist"`
	Option string `json:"option"`
	Limit  int    `json:"limit"`
}
