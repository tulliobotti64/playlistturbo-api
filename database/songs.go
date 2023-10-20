package database

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"playlistturbo.com/model"
)

type MusicDatabase interface {
	AddSong(model.Song) (model.Song, error)
	GetMainList() ([]model.Song, error)
	SearchGenre(genre string) (uint, error)
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
	p.Gorm.Where("name = ?", gSearch).Find(&genre)
	return genre.ID, nil
}
