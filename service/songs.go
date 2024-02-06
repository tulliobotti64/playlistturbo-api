package service

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/dhowden/tag"
	"github.com/go-flac/flacvorbis"
	"github.com/go-flac/go-flac"
	"github.com/gofrs/uuid"
	"playlistturbo.com/config"
	"playlistturbo.com/dto"
	"playlistturbo.com/model"
	"playlistturbo.com/plterror"
	"playlistturbo.com/utils"
)

const (
	EXT_MP3  = "mp3"
	EXT_FLAC = "flac"
)

type SongsService interface {
	AddSong(model.Song) error
	GetMainList() ([]model.Song, error)
	ImportSongs(importSongs dto.ImportSongs) (dto.ImportedSongs, error)
	GetSongsByTitle(title string, limit int) ([]dto.Songs, error)
	MoveSongs(moveSongs dto.MoveSongs) error
	UpdateTwonkyLinks() ([]model.Song, error)
	RemoveSong(importSongs dto.ImportSongs) error
	SetFavoriteSong(id uuid.UUID) error
	GetGenres() ([]model.Genre, error)
	GetArtistByGenre(genreID int) ([]dto.List, error)
	GetAlbumByArtist(artist string) ([]dto.List, error)
	GetSongsByAlbum(album string, limit int) ([]dto.Songs, error)
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

func (svc *PLTService) MoveSongs(moveInfo dto.MoveSongs) error {
	songs, err := svc.DB.GetSongsByPath(moveInfo.OldPath)
	if err != nil {
		return err
	}
	for _, song := range songs {
		if _, err := os.Stat(song.FilePath); errors.Is(err, os.ErrNotExist) {
			err = svc.MoveSong(song, moveInfo)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (svc *PLTService) MoveSong(song model.Song, moveInfo dto.MoveSongs) error {
	// extractFilename(song.FilePath)
	fileArray := strings.Split(song.FilePath, "/")
	newFilePath := ""
	for _, file := range fileArray {
		if strings.Contains(file, ".mp3") {
			newFilePath = moveInfo.NewPath + "/" + file
			if _, err := os.Stat(newFilePath); errors.Is(err, os.ErrNotExist) {
				return err
			} else {
				err := svc.DB.UpdateFilePath(song.ID, newFilePath)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (svc *PLTService) ImportSongs(importSongs dto.ImportSongs) (dto.ImportedSongs, error) {
	var msg dto.ImportedSongs
	msg.Message = "Quantidade de musicas importadas:"
	msg.StartTime = time.Now()

	songExtraTable, err := GetSongExtraTable(importSongs)

	fileList := make([]string, 0)
	var fList = make([]string, 0)
	extension := "*." + importSongs.SongExtension
	fList, err = svc.WalkMatch(importSongs.Path, extension, importSongs.Recursive)
	if err != nil {
		return msg, err
	}
	fileList = append(fileList, fList...)

	for _, songPath := range fileList {
		song := model.Song{}
		var err error
		if strings.Contains(songPath, "mp3") {
			// songTitle := strings.Split(songPath, "/")
			song, err = svc.processMp3(songPath, songExtraTable, importSongs.GenreFromPath)
			if err != nil {
				return msg, err
			}
		}

		if strings.Contains(songPath, "flac") {
			song, err = svc.processFlac(songPath)
			if err != nil {
				return msg, err
			}
		}

		exist, err := svc.DB.SearchFilePath(songPath)
		if err != nil {
			return msg, err
		}

		if !exist {
			_, err = svc.DB.AddSong(song)
			if err != nil {
				fmt.Printf("error adding song: %s", song.Title)
				fmt.Printf("error : %v\n", err)
				return msg, err
			}
			msg.SongQty++
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

func GetSongExtraTable(importSongs dto.ImportSongs) ([]dto.SongExtraTable, error) {
	var msg dto.ImportedSongs
	songET := []dto.SongExtraTable{}
	// verify path
	info, err := os.Stat(importSongs.Path)
	if err != nil || !info.IsDir() {
		return songET, plterror.InvalidSongPath
	}

	arrPath := strings.Split(importSongs.Path, "/")
	if len(arrPath) < 5 {
		return songET, plterror.InvalidSongPath
	}

	genre := arrPath[4]
	// pegar lista dos Diretorios de Generos do Nas
	nasGenreURL, err := GetTwonkyGenre(genre)
	if err != nil {
		return songET, err
	}

	artist := ""
	if len(arrPath)-1 > 4 {
		artist = arrPath[5]
	}
	nasArtistURL, err := GetTwonkyArtistFolder(nasGenreURL, artist, &msg)
	if err != nil {
		return songET, err
	}

	var albumLinkList []string
	album := ""
	if len(arrPath)-1 > 5 {
		album = arrPath[6]
	}
	albumLinkList, err = GetArtistAlbumList(nasArtistURL, album, &msg)
	if err != nil {
		return songET, err
	}

	albumCDList, err := GetArtistAlbumCDList(albumLinkList)
	if err != nil {
		return songET, err
	}

	songET, err = GetAlbumSongList(albumCDList, &msg, importSongs.SongExtension)
	if err != nil {
		return songET, err
	}
	return songET, nil
}

func (svc *PLTService) processMp3(songPath string, songExtraTable []dto.SongExtraTable, genreFromPath bool) (model.Song, error) {
	var songMp3 model.Song
	f, err := os.Open(songPath)
	if err != nil {
		log.Fatal("Error while opening mp3 file: ", err)
	}
	defer f.Close()

	genreTagAux := ""
	albumAux := "unknown"
	artistAux := "unknown"
	titleAux := ""
	trackAux := 0
	var yearAux uint = 0
	var genreID uint = 1
	pathSplit := strings.Split(songPath, "/")
	fname := extractFilename(songPath)

	mp3Tag, err := tag.ReadFrom(f)
	if err != nil {
		fmt.Printf("error reading file: %v\n", err)
		fmt.Println("file:", songPath)
	} else {
		genreTagAux = mp3Tag.Genre()
		yearAux = uint(mp3Tag.Year())
		trackAux, _ = mp3Tag.Track()
		titleAux = utils.RemoveAccent(mp3Tag.Title())

		if len(mp3Tag.Artist()) > 0 {
			artistAux = utils.RemoveAccent(mp3Tag.Artist()) // tratar campo com acento
		} else {
			artistAux = pathSplit[5]
		}

		if len(mp3Tag.Album()) > 0 {
			albumAux = utils.RemoveAccent(mp3Tag.Album()) //*
		} else {
			albumAux = pathSplit[6]
		}
	}
	if genreFromPath || genreTagAux == "" {
		genreTagAux = pathSplit[4]
	}
	genreID, err = svc.DB.SearchGenre(genreTagAux)
	if err != nil {
		return songMp3, plterror.ErrServerError
	}
	if genreID == 0 {
		genreID = 1
	}

	if len(titleAux) == 0 {
		titleAux = pathSplit[len(pathSplit)-1]
	}

	songMp3.Album = albumAux
	songMp3.AlbumDate = yearAux
	songMp3.Artist = artistAux
	songMp3.GenreID = genreID
	songMp3.GenreTag = genreTagAux
	songMp3.Title = titleAux
	songMp3.UpdatedAt = time.Now()
	songMp3.FilePath = utils.RemoveAccent(songPath) //*
	songMp3.TrackNumber = uint(trackAux)
	songMp3.Format = "mp3"

	for _, s := range songExtraTable {
		if s.Filename == fname {
			songMp3.TwonkyLink = s.URL
			songMp3.SampleFrequency = uint(s.Frequency)
			songMp3.Bitrate = uint(s.Bitrate)
			songMp3.AlbumArtURI = s.AlbumArtURI
			songMp3.Lenght = s.Duration
			break
		}
	}

	return songMp3, nil
}
func (svc *PLTService) processFlac(songPath string) (model.Song, error) {
	var songFlac model.Song
	f, err := flac.ParseFile(songPath)
	if err != nil {
		panic(err)
	}

	data, err := f.GetStreamInfo()
	if err != nil {
		panic(err)
	}

	flacLenght := data.SampleCount / int64(data.SampleRate)
	var genrex, albumx, titlex, artistx, trackx string
	var yearx int = 0
	var cmt *flacvorbis.MetaDataBlockVorbisComment

	for idx, meta := range f.Meta {
		fmt.Println("indice:", idx)
		if meta.Type == flac.VorbisComment {
			cmt, err = flacvorbis.ParseFromMetaDataBlock(*meta)
			if err != nil {
				panic(err)
			}
			title, _ := cmt.Get("TITLE")
			titlex = title[0]
			album, _ := cmt.Get("ALBUM")
			albumx = album[0]
			track, _ := cmt.Get("TRACKNUMBER")
			trackx = track[0]
			artist, _ := cmt.Get("ARTIST")
			artistx = artist[0]
			genre, _ := cmt.Get("GENRE")
			genrex = genre[0]
			year, _ := cmt.Get("DATE")
			yearx, err = strconv.Atoi(year[0])
			if err != nil {
				yearx = 0
			}
		}
	}
	genreID, err := svc.DB.SearchGenre(genrex)
	if err != nil {
		return songFlac, plterror.ErrServerError
	}
	if genreID == 0 {
		genreID = 1
	}
	songFlac.Album = albumx
	songFlac.AlbumDate = uint(yearx)
	songFlac.Artist = artistx
	songFlac.GenreID = genreID
	songFlac.Lenght = utils.SecondsToMinutes(int(flacLenght))
	songFlac.Title = titlex
	songFlac.UpdatedAt = time.Now()
	songFlac.FilePath = songPath
	tracki, _ := strconv.Atoi(trackx)
	songFlac.TrackNumber = uint(tracki)
	return songFlac, nil
}

func (svc *PLTService) WalkMatch(root, pattern string, recursive bool) ([]string, error) {
	matches := make([]string, 0)
	qtd := 0
	paradeler := false
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if qtd > 0 && !recursive {
				paradeler = true
			}
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched && !paradeler {
			matches = append(matches, path)
			qtd++
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}

func GetTwonkyGenre(genre string) (string, error) {

	url := ""

	req, err := http.NewRequest(http.MethodGet, config.Config.DlnaGenreUrl, nil)
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

func GetTwonkyArtistFolder(urlArtist, artist string, msg *dto.ImportedSongs) ([]string, error) {
	var url []string
	takeAll := false
	if len(artist) == 0 {
		takeAll = true
	}
	req, err := http.NewRequest(http.MethodGet, urlArtist, nil)
	if err != nil {
		return url, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("resp:", resp)
		return url, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var artistFolder dto.TwonkyArtistAlbumFolder

	if err := xml.Unmarshal([]byte(body), &artistFolder); err != nil {
		log.Println("unmarshal:", err)
		return url, err
	}

	for _, item := range artistFolder.Channel.Item {
		if takeAll || artist == item.Title {
			url = append(url, item.Enclosure.Url)
			msg.ArtistQty++
		}
	}
	if len(url) == 0 {
		return url, plterror.Tabelavazia
	}

	return url, nil
}

func GetArtistAlbumList(nasArtistAlbumURL []string, oneAlbum string, msg *dto.ImportedSongs) ([]string, error) {
	var albumList []string
	takeAll := false
	if len(oneAlbum) == 0 {
		takeAll = true
	}
	msg.AlbumQty = len(nasArtistAlbumURL)

	for _, album := range nasArtistAlbumURL {
		req, err := http.NewRequest(http.MethodGet, album, nil)
		if err != nil {
			return albumList, err
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println("resp:", resp)
			return albumList, err
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)

		var artistAlbumList dto.TwonkyArtistAlbumFolder

		if err := xml.Unmarshal([]byte(body), &artistAlbumList); err != nil {
			log.Println("unmarshal:", err)
			return albumList, err
		}

		for _, item := range artistAlbumList.Channel.Item {
			if takeAll || oneAlbum == item.Title {
				albumList = append(albumList, item.Enclosure.Url)
			}
		}
	}

	return albumList, nil
}

func GetArtistAlbumCDList(nasArtistAlbumCDURL []string) ([]string, error) {
	var albumList []string

	for _, album := range nasArtistAlbumCDURL {
		req, err := http.NewRequest(http.MethodGet, album, nil)
		if err != nil {
			return albumList, err
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println("resp:", resp)
			return albumList, err
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)

		var artistAlbumList dto.TwonkyArtistAlbumFolder

		if err := xml.Unmarshal([]byte(body), &artistAlbumList); err != nil {
			log.Println("unmarshal:", err)
			return albumList, err
		}

		for _, item := range artistAlbumList.Channel.Item {
			albumList = append(albumList, item.Enclosure.Url)
		}
	}

	return albumList, nil
}

func GetAlbumSongList(albumUrlList []string, msg *dto.ImportedSongs, ext string) ([]dto.SongExtraTable, error) {
	songList := make([]dto.SongExtraTable, 0)

	for _, albumUrl := range albumUrlList {

		req, err := http.NewRequest(http.MethodGet, albumUrl, nil)
		if err != nil {
			return songList, err
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println("resp:", resp)
			return songList, err
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)

		var albumSongList dto.TwonkySongList

		if err := xml.Unmarshal([]byte(body), &albumSongList); err != nil {
			log.Println("unmarshal:", err)
			return songList, err
		}

		for _, item := range albumSongList.Channel.Item {
			entry := dto.SongExtraTable{
				Filename:    item.Title,
				URL:         item.Meta.Res.CharData,
				Frequency:   item.Meta.Res.SampleFrequency,
				Bitrate:     item.Meta.Res.Bitrate,
				Duration:    item.Meta.Res.Duration,
				AlbumArtURI: item.Meta.AlbumArtURI.CharData,
			}

			if item.Meta.Extension == ext {
				songList = append(songList, entry)
			}
		}
	}

	return songList, nil
}

func (svc *PLTService) GetSongsByTitle(title string, limit int) ([]dto.Songs, error) {
	var songs []dto.Songs
	songsDB, err := svc.DB.GetSongsByTitle(title, limit)
	if err != nil {
		return songs, err
	}

	for _, songDB := range songsDB {
		var song dto.Songs
		song = dto.ToDtoSongs(songDB, song)
		songs = append(songs, song)
	}
	return songs, nil
}

func (svc *PLTService) UpdateTwonkyLinks() ([]model.Song, error) {
	songs, err := svc.DB.GetEmptyTwonkyLinks()
	if err != nil {
		return songs, err
	}

	if len(songs) == 0 {
		return songs, nil
	}

	var param dto.ImportSongs
	param.Path = extractPath(songs[0].FilePath)
	param.Recursive = true
	if strings.Contains(songs[0].FilePath, EXT_MP3) {
		param.SongExtension = EXT_MP3
	}
	if strings.Contains(songs[0].FilePath, EXT_FLAC) {
		param.SongExtension = EXT_FLAC
	}

	var processedPath string
	var songsET []dto.SongExtraTable
	for _, song := range songs {
		if extractPath(song.FilePath) != processedPath {
			processedPath = extractPath(song.FilePath)
			songsET, err = GetSongExtraTable(param)
			if err != nil {
				return songs, err
			}
		}
		// get townky info from ET and update db
		for _, songET := range songsET {
			if extractFilename(song.FilePath) == songET.Filename {
				svc.DB.UpdateTwonkyLinks(song.ID, songET.URL, songET.AlbumArtURI)
				break
			}
		}
	}

	return songs, nil
}

func (svc *PLTService) RemoveSong(importSongs dto.ImportSongs) error {
	split := strings.Split(importSongs.Path, "/")

	id := len(split)
	path := ""
	for x := 1; x < id; x++ {
		path += "/" + split[x]
		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			break
		}
	}

	err := svc.DB.RemoveSong(path)
	if err != nil {
		return err
	}
	return nil
}

func (svc *PLTService) SetFavoriteSong(id uuid.UUID) error {
	err := svc.DB.SetFavoriteSong(id)
	if err != nil {
		return err
	}
	return nil
}
func extractPath(path string) string {
	temp := strings.Split(path, "/")
	var newPath string
	for x := 1; x < len(temp)-1; x++ {
		piece := "/" + temp[x]
		newPath += piece
	}
	return newPath
}

func extractFilename(path string) string {
	pathSplit := strings.Split(path, "/")
	pIndex := len(pathSplit) - 1
	fname := pathSplit[pIndex]
	return fname
}

func (svc *PLTService) GetGenres() ([]model.Genre, error) {
	var genres []model.Genre
	songs, err := svc.DB.GetGenres()
	if err != nil {
		return genres, err
	}
	for _, song := range songs {
		var genre model.Genre
		genre.ID = song.GenreID
		genre.Name = song.Genre.Name
		genres = append(genres, genre)
	}
	return genres, nil
}

func (svc *PLTService) GetArtistByGenre(genreID int) ([]dto.List, error) {
	var artists []dto.List
	songs, err := svc.DB.GetArtistByGenre(genreID)
	if err != nil {
		return artists, err
	}
	for x, song := range songs {
		var artist dto.List
		artist.ID = uint(x + 1)
		artist.Name = song.Artist
		artists = append(artists, artist)
	}
	return artists, nil
}
func (svc *PLTService) GetAlbumByArtist(artist string) ([]dto.List, error) {
	var albums []dto.List
	songs, err := svc.DB.GetAlbumByArtist(artist)
	if err != nil {
		return albums, err
	}
	for x, song := range songs {
		var album dto.List
		album.ID = uint(x + 1)
		album.Name = song.Album
		albums = append(albums, album)
	}
	return albums, nil
}

func (svc *PLTService) GetSongsByAlbum(album string, limit int) ([]dto.Songs, error) {
	var songs []model.Song
	var songsDto []dto.Songs
	albumDecoded, err := url.QueryUnescape(album)
	songs, err = svc.DB.GetSongsByAlbum(albumDecoded, limit)
	if err != nil {
		return songsDto, err
	}

	for _, songDB := range songs {
		var songDto dto.Songs
		songDto = dto.ToDtoSongs(songDB, songDto)
		songsDto = append(songsDto, songDto)
	}

	return songsDto, nil
}
