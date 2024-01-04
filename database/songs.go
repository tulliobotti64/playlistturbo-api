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
	GetEmptyTwonkyLinks() ([]model.Song, error)
	UpdateTwonkyLinks(songID uuid.UUID, twonkyLink, albumUri string) error
	RemoveSong(path string) error
	SetFavoriteSong(id uuid.UUID) error
	GetGenres() ([]model.Song, error)
	GetArtistByGenre(genreID int) ([]model.Song, error)
	GetAlbumByArtist(artist string) ([]model.Song, error)
	GetSongsByAlbum(album string) ([]model.Song, error)
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
		Order("title ASC").
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

func (p *PostgresDB) GetEmptyTwonkyLinks() ([]model.Song, error) {
	var songs []model.Song
	err := p.Gorm.Model(&songs).
		Where(`twonky_link = ''`).
		Find(&songs).
		Order("file_path ASC").
		Error
	if err != nil {
		return songs, err
	}
	return songs, nil
}

func (p *PostgresDB) UpdateTwonkyLinks(songID uuid.UUID, twonkyLink, albumUri string) error {
	err := p.Gorm.Model(model.Song{}).
		Select("twonky_link", "album_art_uri").
		Where("id = ?", songID).
		Updates(map[string]interface{}{"twonky_link": twonkyLink, "album_art_uri": albumUri}).
		Error
	if err != nil {
		return handleError(err)
	}
	return nil
}

func (p *PostgresDB) RemoveSong(path string) error {
	path1 := "%" + strings.TrimSpace(path) + "%"
	err := p.Gorm.Delete(&model.Song{}, "file_path like ?", path1).Error
	if err != nil {
		return handleError(err)
	}
	return nil
}

func (p *PostgresDB) SetFavoriteSong(songID uuid.UUID) error {
	var song model.Song
	err := p.Gorm.Model(&song).
		Where("id = ?", songID).
		Find(&song).
		Error
	if err != nil {
		return handleError(err)
	}
	fav := !song.Favorite

	err = p.Gorm.Model(&song).
		Select("favorite").
		Where("id = ?", songID).
		Updates(map[string]interface{}{"favorite": fav}).
		Error
	if err != nil {
		return handleError(err)
	}
	return nil
}

func (p *PostgresDB) GetGenres() ([]model.Song, error) {
	var songs []model.Song
	err := p.Gorm.Model(&songs).
		Select("genres.name, songs.genre").
		Preload("Genre").
		Joins("join genres on genres.id = songs.genre").
		Group("genres.name, songs.genre").
		Order("genres.name").
		Find(&songs).Error
	if err != nil {
		return songs, handleError(err)
	}
	return songs, nil
}

func (p *PostgresDB) GetArtistByGenre(genreID int) ([]model.Song, error) {
	var songs []model.Song
	err := p.Gorm.Model(&songs).
		Select("artist").
		Where("genre = ?", genreID).
		Group("artist").
		Order("artist").
		Find(&songs).Error
	if err != nil {
		return songs, handleError(err)
	}
	return songs, nil
}
func (p *PostgresDB) GetAlbumByArtist(artist string) ([]model.Song, error) {
	var songs []model.Song
	err := p.Gorm.Model(&songs).
		Select("album").
		Where("artist = ?", artist).
		Group("album").
		Order("album").
		Find(&songs).Error
	if err != nil {
		return songs, handleError(err)
	}
	return songs, nil
}

func (p *PostgresDB) GetSongsByAlbum(album string) ([]model.Song, error) {
	var songs []model.Song
	err := p.Gorm.Model(&songs).
		Where("album = ?", album).
		Find(&songs).Error
	if err != nil {
		return songs, handleError(err)
	}
	return songs, nil
}
