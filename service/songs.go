package service

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/dhowden/tag"
	"github.com/go-flac/flacvorbis"
	"github.com/go-flac/go-flac"
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
	GetSongsByTitle(title string) ([]dto.Songs, error)
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
	msg.Message = "Quantidade de musicas importadas:"
	msg.StartTime = time.Now()

	//verify path
	info, err := os.Stat(importSongs.Path)
	if err != nil || !info.IsDir() {
		return msg, plterror.InvalidSongPath
	}

	arrPath := strings.Split(importSongs.Path, "/")
	if len(arrPath) < 5 {
		return msg, plterror.InvalidSongPath
	}

	genre := arrPath[4]
	//pegar lista dos Diretorios de Generos do Nas
	nasGenreURL, err := GetTwonkyGenre(genre)
	if err != nil {
		return msg, err
	}

	artist := ""
	if len(arrPath)-1 > 4 {
		artist = arrPath[5]
	}
	nasArtistURL, err := GetTwonkyArtistFolder(nasGenreURL, artist, &msg)
	if err != nil {
		return msg, err
	}

	var albumLinkList []string
	album := ""
	if len(arrPath)-1 > 5 {
		artist = arrPath[6]
	}
	albumLinkList, err = GetArtistAlbumList(nasArtistURL, album, &msg)
	if err != nil {
		return msg, err
	}

	songExtraTable, err := GetAlbumSongList(albumLinkList, album, &msg, importSongs.SongExtension)
	if err != nil {
		return msg, err
	}

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
			song, err = svc.processMp3(songPath, songExtraTable)
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
				return msg, nil
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

func (svc *PLTService) processMp3(songPath string, songExtraTable []dto.SongExtraTable) (model.Song, error) {
	var songMp3 model.Song
	f, err := os.Open(songPath)
	if err != nil {
		log.Fatal("Error while opening mp3 file: ", err)
	}
	defer f.Close()

	genreTagAux := ""
	albumAux := "Unknonwn"
	artistAux := "Unknonwn"
	titleAux := "Unknonwn"
	trackAux := 0
	var yearAux uint = 0
	var genreID uint = 1
	pathSplit := strings.Split(songPath, "/")
	pIndex := len(pathSplit) - 1
	fname := pathSplit[pIndex]

	mp3Tag, err := tag.ReadFrom(f)
	if err != nil {
		fmt.Printf("error reading file: %v\n", err)
		fmt.Println("file:", songPath)
		// return songMp3, err
	} else {
		genreTagAux = mp3Tag.Genre()
		albumAux = mp3Tag.Album()
		yearAux = uint(mp3Tag.Year())
		artistAux = mp3Tag.Artist()
		titleAux = mp3Tag.Title()
		trackAux, _ = mp3Tag.Track()
	}

	if genreTagAux == "" {
		genreTagAux = pathSplit[4]
	}
	genreID, err = svc.DB.SearchGenre(genreTagAux)
	if err != nil {
		return songMp3, plterror.ErrServerError
	}
	if genreID == 0 {
		genreID = 1
	}

	// mp3Sec := utils.GetMp3Time(f)  -- it takes forever to retrieve the duration
	// songMp3.Lenght = utils.SecondsToMinutes(int(mp3Sec))
	// songMp3.LenghtSec = mp3Sec

	songMp3.Album = albumAux
	songMp3.AlbumDate = yearAux
	songMp3.Artist = artistAux
	songMp3.GenreID = genreID
	songMp3.GenreTag = genreTagAux
	songMp3.Title = titleAux
	songMp3.UpdatedAt = time.Now()
	songMp3.FilePath = songPath
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
	songFlac.LenghtSec = float64(flacLenght)
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

func GetAlbumSongList(albumUrlList []string, album string, msg *dto.ImportedSongs, ext string) ([]dto.SongExtraTable, error) {
	songList := make([]dto.SongExtraTable, 0)
	takeAll := false
	if len(album) == 0 {
		takeAll = true
	}

	for _, album := range albumUrlList {

		req, err := http.NewRequest(http.MethodGet, album, nil)
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
			if takeAll || album == item.Title {
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
	}

	return songList, nil
}

func (svc *PLTService) GetSongsByTitle(title string) ([]dto.Songs, error) {
	var songs []dto.Songs
	songsDB, err := svc.DB.GetSongsByTitle(title)
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
