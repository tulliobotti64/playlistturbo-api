package database

import (
	"strings"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"playlistturbo.com/model"
)

type MusicDatabase interface {
	AddSong(model.Song) (model.Song, error)
	GetMainList() ([]model.Song, error)
	SearchGenre(genre string) (uint, error)
	SearchFilePath(filePath string) (bool, error)
	GetSongsByTitle(title string) ([]model.Song, error)
	GetSongsByPath(path string) ([]model.Song, error)
	UpdateFilePath(songID uuid.UUID, path string) error
}

func (p *PostgresDB) AddSong(Song model.Song) (model.Song, error) {
	err := p.Gorm.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&Song).Error; err != nil {
			return err
		}

		return tx.Preload(clause.Associations).Find(&Song).Error
	})
	return Song, handleError(err)
}

func (p *PostgresDB) GetMainList() ([]model.Song, error) {
	var songs []model.Song
	if err := p.Gorm.Model(&songs).Find(&songs).Error; err != nil {
		return nil, err
	}
	return songs, nil
}

func (p *PostgresDB) SearchGenre(gSearch string) (uint, error) {
	var genre model.Genre
	gs := strings.ToLower(gSearch)
	p.Gorm.Where("lower(name) = ?", gs).Find(&genre)
	return genre.ID, nil
}

func (p *PostgresDB) SearchFilePath(filePath string) (bool, error) {
	var exist bool
	err := p.Gorm.Model(model.Song{}).Select("count(*) > 0").Where("file_path = ?", filePath).Find(&exist).Error
	if err != nil {
		return false, err
	}
	return exist, nil
}
func (p *PostgresDB) GetSongsByTitle(title string) ([]model.Song, error) {
	var songs []model.Song
	titlex := "%" + title + "%"
	err := p.Gorm.Model(&songs).
		Where("title ilike ?", titlex).
		Find(&songs).
		Error
	if err != nil {
		return songs, err
	}
	return songs, nil
}

func (p *PostgresDB) GetSongsByPath(path string) ([]model.Song, error) {
	var songs []model.Song
	pathx := "%" + path + "%"
	err := p.Gorm.Model(&songs).
		Where("file_path like ?", pathx).
		Find(&songs).
		Error
	if err != nil {
		return songs, err
	}
	return songs, nil
}

func (p *PostgresDB) UpdateFilePath(songID uuid.UUID, path string) error {
	err := p.Gorm.Model(model.Song{}).
		Select("file_path", "twonky_link", "album_art_uri").
		Where("id = ?", songID).
		Updates(map[string]interface{}{"file_path": path, "twonky_link": "", "album_art_uri": ""}).
		Error
	if err != nil {
		return handleError(err)
	}
	return nil
}
