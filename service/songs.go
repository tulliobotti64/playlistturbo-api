package service

import (
	"fmt"
	"log"
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
	ImportSongs(importSongs dto.ImportSongs) ([]model.Song, error)
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
func (svc *PLTService) ImportSongs(importSongs dto.ImportSongs) ([]model.Song, error) {
	if len(importSongs.Path) == 0 {
		log.Println("variable error")
	}

	fileList := make([]string, 0)
	fList := make([]string, 0)
	extension := "*." + importSongs.SongExtension
	var err error
	fList, err = svc.WalkMatch(importSongs.Path, extension, importSongs.Recursive)
	if err != nil {
		return nil, err
	}
	fileList = append(fileList, fList...)

	songsList := []model.Song{}
	// Ler os tags de todos os arquivos

	for _, songPath := range fileList {
		song := model.Song{}
		var err error
		if strings.Contains(songPath, "mp3") {
			song, err = svc.processMp3(songPath)
			if err != nil {
				return nil, err
			}
		}

		if strings.Contains(songPath, "flac") {
			song, err = svc.processFlac(songPath)
			if err != nil {
				return nil, err
			}
		}

		exist, err := svc.DB.SearchFilePath(songPath)
		if err != nil {
			return nil, err
		}

		var songDB model.Song
		if !exist {
			songDB, err = svc.DB.AddSong(song)
			if err != nil {
				fmt.Printf("error adding song: %s", song.Title)
				fmt.Printf("error : %v\n", err)
				return songsList, nil
			}
			songsList = append(songsList, songDB)
		}
	}

	return songsList, nil
}

func (svc *PLTService) processMp3(songPath string) (model.Song, error) {
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
	// utils.PrintMetadata(m)
	mp3Sec := utils.GetMp3Time(f)
	genreID, err := svc.DB.SearchGenre(mp3Tag.Genre())
	if err != nil {
		return songMp3, plterror.ErrServerError
	}
	if genreID == 0 {
		genreID = 1
	}
	songMp3.Album = mp3Tag.Album()
	songMp3.AlbumYear = uint(mp3Tag.Year())
	songMp3.Artist = mp3Tag.Artist()
	songMp3.GenreID = genreID
	songMp3.Lenght = utils.SecondsToMinutes(int(mp3Sec))
	songMp3.LenghtSec = mp3Sec
	songMp3.Title = mp3Tag.Title()
	songMp3.UpdatedAt = time.Now()
	songMp3.FilePath = songPath

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
	var genrex, albumx, titlex, artistx string
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
			// track, _ := cmt.Get("TRACKNUMBER")
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
	songFlac.AlbumYear = uint(yearx)
	songFlac.Artist = artistx
	songFlac.GenreID = genreID
	songFlac.Lenght = utils.SecondsToMinutes(int(flacLenght))
	songFlac.LenghtSec = float64(flacLenght)
	songFlac.Title = titlex
	songFlac.UpdatedAt = time.Now()
	songFlac.FilePath = songPath
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
