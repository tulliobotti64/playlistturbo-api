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

	mp3Tag, err := tag.ReadFrom(f)
	if err != nil {
		fmt.Printf("error reading file: %v\n", err)
		return songMp3, nil
	}

	// mp3Sec := utils.GetMp3Time(f)  -- it takes forever to retrieve the duration
	genreTag := mp3Tag.Genre()
	genreID, err := svc.DB.SearchGenre(genreTag)
	if err != nil {
		return songMp3, plterror.ErrServerError
	}
	if genreID == 0 {
		genreID = 1
	}

	songMp3.Album = mp3Tag.Album()
	songMp3.AlbumDate = uint(mp3Tag.Year())
	songMp3.Artist = mp3Tag.Artist()
	songMp3.GenreID = genreID
	songMp3.GenreTag = genreTag
	// songMp3.Lenght = utils.SecondsToMinutes(int(mp3Sec))
	// songMp3.LenghtSec = mp3Sec
	songMp3.Title = mp3Tag.Title()
	songMp3.UpdatedAt = time.Now()
	songMp3.FilePath = songPath

	pathSplit := strings.Split(songPath, "/")
	pIndex := len(pathSplit) - 1
	fname := pathSplit[pIndex]

	track, _ := mp3Tag.Track()
	songMp3.TrackNumber = uint(track)
	songMp3.Format = string(mp3Tag.Format())

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

	URLGenreList := "http://192.168.32.5:9000/nmc/rss/server/RBuuid:55076f6e-6b79-1d65-a4cf-0000c03c4dd4,0/IBuuid:55076f6e-6b79-1d65-a4cf-0000c03c4dd4,_MCQxJDEz,,0,0,_Um9vdA==,0,,1,0,_TXVzaWM=,_MCQx,?start=0&count=40&fmt=json"
	url := ""

	req, err := http.NewRequest(http.MethodGet, URLGenreList, nil)
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
