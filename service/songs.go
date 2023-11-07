package service

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"playlistturbo.com/dto"
	"playlistturbo.com/model"
	"playlistturbo.com/plterror"
)

const (
	EXT_MP3  = "mp3"
	EXT_FLAC = "flac"
)

type SongsService interface {
	AddSong(model.Song) error
	GetMainList() ([]model.Song, error)
	ImportSongs(importSongs dto.ImportSongs) (dto.ImportedSongs, error)
}

func (svc *PLTService) AddSong(Song model.Song) error {
	ct := time.Now()
	fmt.Println(ct.Format("15:04:05"))
	_, err := svc.DB.AddSong(Song)
	if err != nil {
		return err
	}
	return nil
}
func (svc *PLTService) GetMainList() ([]model.Song, error) {
	music, err := svc.DB.GetMainList()

	return music, err
}

// Essa funçao de Importacao de Musicas vai receber um Path pra ler os arquivos recursivamente ou nao
// Cada um dos arquivos, se for uma extensao suportada (pegamos do arquivo de config), abrimos ele com
// o id3v2. Se as infos de tag mp3 forem corretas, adicionamos ele no DB
func (svc *PLTService) ImportSongs(importSongs dto.ImportSongs) (dto.ImportedSongs, error) {
	var msg dto.ImportedSongs
	var err error
	msg.Message = "Quantidade de musicas importadas:"
	msg.StartTime = time.Now()
	//pegar lista dos Diretorios de Generos do Nas
	nasGenreURL, err := GetTwonkyGenre(importSongs.Genre)
	if err != nil {
		return msg, err
	}

	artists, err := GetTwonkyArtistFolder(nasGenreURL, importSongs.Artist, &msg)
	if err != nil {
		return msg, err
	}

	for _, artist := range artists {

		err := svc.GetArtistAlbum(artist, importSongs.Genre, importSongs.Album, importSongs.Extension, &msg)
		if err != nil {
			return msg, err
		}
	}

	msg.EndTime = time.Now()
	dur := msg.EndTime.Sub(msg.StartTime).Seconds()

	if dur <= 60 {
		msg.Duration = fmt.Sprintf("%.2f", dur) + " seconds"
	} else {
		msg.Duration = fmt.Sprintf("%.2f", msg.EndTime.Sub(msg.StartTime).Minutes()) + " minutes"
	}

	return msg, nil
}

func GetTwonkyGenre(genre string) (string, error) {

	URLGenre := "http://192."
	url := ""

	req, err := http.NewRequest(http.MethodGet, URLGenre, nil)
	if err != nil {
		return url, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("resp:", resp)
		return url, plterror.ErrDLNAAccess
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var nasGenreList dto.TwonkyGengeList

	json.Unmarshal(body, &nasGenreList)
	if err != nil {
		log.Println("Unmarshal error:", err)
		return url, err
	}

	//lopar item pra pegar o  folder de acordo com o Genero
	for _, myFolder := range nasGenreList.Item {
		if genre == myFolder.Title {
			url = myFolder.Enclosure.URL
		}
	}

	if len(url) == 0 {
		return url, plterror.InvalidGenre
	}

	return url, nil
}

func GetTwonkyArtistFolder(urlArtist, artist string, msg *dto.ImportedSongs) ([]dto.ImportAlbums, error) {
	var albums []dto.ImportAlbums
	takeAll := false
	if len(artist) == 0 {
		takeAll = true
	}
	req, err := http.NewRequest(http.MethodGet, urlArtist, nil)
	if err != nil {
		return albums, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("resp:", resp)
		return albums, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var artistFolder dto.TwonkyArtistAlbumFolder

	if err := xml.Unmarshal([]byte(body), &artistFolder); err != nil {
		fmt.Println("fudeu", err)
		return albums, err
	}

	for _, item := range artistFolder.Channel.Item {
		if takeAll || artist == item.Title {
			tbj := dto.ImportAlbums{
				AlbumURL:   item.Enclosure.Url,
				AlbumTitle: item.Title,
			}
			msg.ArtistQty++
			albums = append(albums, tbj)
		}
	}
	if len(albums) == 0 {
		return albums, plterror.Tabelavazia
	}

	return albums, nil
}

func (svc *PLTService) GetArtistAlbum(artist dto.ImportAlbums, genre, oneAlbum, extension string, msg *dto.ImportedSongs) error {
	takeAll := false
	if len(oneAlbum) == 0 {
		takeAll = true
	}

	req, err := http.NewRequest(http.MethodGet, artist.AlbumURL, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("resp:", resp)
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var artistAlbumList dto.TwonkyArtistAlbumFolder

	if err := xml.Unmarshal([]byte(body), &artistAlbumList); err != nil {
		fmt.Println("fudeu", err)
		return err
	}

	for _, item := range artistAlbumList.Channel.Item {
		if takeAll || oneAlbum == item.Title {
			msg.AlbumQty++
			var songsQty int
			songsQty, err = svc.GetAlbumSongList(item.Enclosure.Url, oneAlbum, genre, extension)
			msg.SongQty = msg.SongQty + songsQty
		}
	}

	return nil
}

func (svc *PLTService) GetAlbumSongList(albumUrl string, album, genre, extension string) (int, error) {
	songs := 0

	req, err := http.NewRequest(http.MethodGet, albumUrl, nil)
	if err != nil {
		return 0, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("resp:", resp)
		return 0, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var albumSongList dto.TwonkySongList

	if err := xml.Unmarshal([]byte(body), &albumSongList); err != nil {
		fmt.Println("fudeu", err)
		return 0, err
	}

	for _, item := range albumSongList.Channel.Item {
		if item.Meta.Extension == extension {
			ok, err := svc.SaveSong(item, genre, item.Title)

			if err != nil {
				return 0, err
			}
			if ok {
				songs++
			}
		}
	}

	return songs, nil
}

func (svc *PLTService) SaveSong(item dto.TwonkySongListItem, genre, album string) (bool, error) {
	var song model.Song

	genreTag := item.Meta.Genre
	if genreTag == "Unknown" {
		genreTag = genre
	}
	genreID, err := svc.DB.SearchGenre(strings.ToLower(genreTag))
	if err != nil {
		return false, plterror.ErrServerError
	}
	if genreID == 0 {
		genreID = 1
	}
	song.GenreTag = genreTag
	song.GenreID = genreID

	song.Album = item.Meta.Album
	song.AlbumDate = item.Meta.Date
	song.Artist = item.Meta.Artist
	song.Lenght = item.Meta.Duration
	song.Title = item.Meta.Title
	song.UpdatedAt = time.Now()
	song.TwonkyLink = item.Meta.Res.CharData
	songPath := ""
	if strings.Contains(song.Title, ".mp3") {
		songPath = "/mnt/Ironman/Musicas-MP3"
	} else {
		songPath = "/mnt/Ironman/Musicas-HQ"
	}
	filePath := songPath + "/" + genre + "/" + album + "/" + song.Title
	song.FilePath = filePath
	song.TrackNumber = item.Meta.OriginalTrackNumber
	song.Format = item.Meta.Format
	song.SampleFrequency = item.Meta.Res.SampleFrequency
	song.Bitrate = item.Meta.Res.Bitrate
	song.AlbumArtURI = item.Meta.AlbumArtURI.CharData

	exist, err := svc.DB.SearchFilePath(filePath)
	if err != nil {
		return false, err
	}

	if !exist {
		_, err = svc.DB.AddSong(song)
		if err != nil {
			fmt.Printf("error adding song: %s", song.Title)
			fmt.Printf("error : %v\n", err)
			return false, nil
		}
	}

	return true, nil
}
