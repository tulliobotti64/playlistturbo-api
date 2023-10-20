package service

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/dhowden/tag"
	"playlistturbo.com/dto"
	"playlistturbo.com/model"
	"playlistturbo.com/plterror"
	"playlistturbo.com/utils"
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

	//ler todos as extensoes
	fileList := make([]string, 0)
	for _, ext := range importSongs.SongExtension {
		fList := make([]string, 0)
		fmt.Println("ext:", ext)
		extension := "*." + ext
		var err error
		fList, err = svc.WalkMatch(importSongs.Path, extension)
		if err != nil {
			return nil, err
		}
		fileList2 = copy(fileList, fList)
	}

	songsList := []model.Song{}
	// Ler os tags de todos os arquivos

	for _, song := range fileList {
		f, err := os.Open(song)
		if err != nil {
			log.Fatal("Error while opening mp3 file: ", err)
		}
		defer f.Close()

		m, err := tag.ReadFrom(f)
		if err != nil {
			fmt.Printf("error reading file: %v\n", err)
			return songsList, nil
		}
		utils.PrintMetadata(m)
		sec := utils.GetMp3Time(f)
		genreID, err := svc.DB.SearchGenre(m.Genre())
		if err != nil {
			return songsList, plterror.ErrServerError
		}
		if genreID == 0 {
			genreID = 1
		}
		song := model.Song{
			GenreID:   genreID,
			Artist:    m.Artist(),
			Album:     m.Album(),
			Title:     m.Title(),
			AlbumYear: uint(m.Year()),
			FilePath:  song,
			Favorite:  false,
			ListenQty: 0,
			UpdatedAt: time.Now(),
			LenghtSec: sec,
			Lenght:    utils.SecondsToMinutes(int(sec)),
		}
		songDB, err := svc.DB.AddSong(song)
		if err != nil {
			fmt.Printf("error adding song: %s", m.Title())
			fmt.Printf("error : %v\n", err)
			return songsList, nil
		}
		songsList = append(songsList, songDB)
	}

	return songsList, nil
}

func (svc *PLTService) WalkMatch(root, pattern string) ([]string, error) {
	matches := make([]string, 0)
	var qt int
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			qt++
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}
