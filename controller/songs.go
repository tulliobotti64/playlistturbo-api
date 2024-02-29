package controller

import (
	"log"
	"net/http"

	"playlistturbo.com/dto"
	"playlistturbo.com/model"
	"playlistturbo.com/plterror"
	"playlistturbo.com/utils"
)

type SongsController interface {
	AddSong(w http.ResponseWriter, r *http.Request)
	GetMainList(w http.ResponseWriter, r *http.Request)
	ImportSongs(w http.ResponseWriter, r *http.Request)
	GetSongsByTitle(w http.ResponseWriter, r *http.Request)
	MoveOneSong(w http.ResponseWriter, r *http.Request)
	RemoveSongs(w http.ResponseWriter, r *http.Request)
	SetFavoriteSong(w http.ResponseWriter, r *http.Request)
	GetGenres(w http.ResponseWriter, r *http.Request)
	GetArtistByGenre(w http.ResponseWriter, r *http.Request)
	GetAlbumByArtist(w http.ResponseWriter, r *http.Request)
	GetSongsByAlbum(w http.ResponseWriter, r *http.Request)
	GetFavorites(w http.ResponseWriter, r *http.Request)
}

func (ctrl *HTTPController) AddSong(w http.ResponseWriter, r *http.Request) {
	body, ok := ctrl.GetBody(r).(model.Song)
	if !ok {
		panic(plterror.ErrBadSyntax)
	}
	err := ctrl.Svc.AddSong(body)
	ctrl.EncodeDataResponse(r, w, nil, err)
}

func (ctrl *HTTPController) GetMainList(w http.ResponseWriter, r *http.Request) {
	songList, err := ctrl.Svc.GetMainList()
	ctrl.EncodeDataResponse(r, w, songList, err)
}

func (ctrl *HTTPController) ImportSongs(w http.ResponseWriter, r *http.Request) {
	var body dto.ImportSongs
	var ok bool
	log.Println("entrou no controller")
	body, ok = ctrl.GetBody(r).(dto.ImportSongs)
	if !ok {
		ctrl.EncodeEmptyResponse(r, w, plterror.ErrBadSyntax)
		return
	}

	err := utils.ValidSongExtension(body.SongExtension)
	if err != nil {
		ctrl.EncodeEmptyResponse(r, w, err)
		return
	}

	songs, err := ctrl.Svc.ImportSongs(body)
	if err != nil {
		ctrl.EncodeEmptyResponse(r, w, err)
		return
	}

	ctrl.EncodeDataResponse(r, w, songs, nil)
}

func (ctrl *HTTPController) MoveOneSong(w http.ResponseWriter, r *http.Request) {
	var body dto.MoveSongs
	var ok bool
	body, ok = ctrl.GetBody(r).(dto.MoveSongs)
	if !ok {
		ctrl.EncodeEmptyResponse(r, w, plterror.ErrBadSyntax)
		return
	}

	err := utils.ValidSongExtension(body.SongExtension)
	if err != nil {
		ctrl.EncodeEmptyResponse(r, w, err)
		return
	}

	err = ctrl.Svc.MoveOneSong(body)
	if err != nil {
		ctrl.EncodeEmptyResponse(r, w, err)
		return
	}

	ctrl.EncodeDataResponse(r, w, nil, nil)
}

func (ctrl *HTTPController) GetSongsByTitle(w http.ResponseWriter, r *http.Request) {
	title := ctrl.GetParam(r, "title")
	if title == "" {
		ctrl.EncodeEmptyResponse(r, w, plterror.ErrBadSyntax)
		return
	}
	limit := ctrl.GetParamInt(r, "limit")
	getHide := ctrl.GetParamBool(r, "gethide")

	resp, err := ctrl.Svc.GetSongsByTitle(title, limit, getHide)
	if err != nil {
		ctrl.EncodeEmptyResponse(r, w, err)
		return
	}
	ctrl.EncodeDataResponse(r, w, resp, nil)
}

func (ctrl *HTTPController) RemoveSongs(w http.ResponseWriter, r *http.Request) {
	var body dto.ImportSongs
	var ok bool
	body, ok = ctrl.GetBody(r).(dto.ImportSongs)
	if !ok {
		ctrl.EncodeEmptyResponse(r, w, plterror.ErrBadSyntax)
		return
	}

	err := utils.ValidSongExtension(body.SongExtension)
	if err != nil {
		ctrl.EncodeEmptyResponse(r, w, err)
		return
	}

	err = ctrl.Svc.RemoveSong(body)
	if err != nil {
		ctrl.EncodeEmptyResponse(r, w, err)
		return
	}

	ctrl.EncodeDataResponse(r, w, nil, nil)
}

func (ctrl *HTTPController) SetFavoriteSong(w http.ResponseWriter, r *http.Request) {
	id := ctrl.GetParamUUID(r, "id")
	err := ctrl.Svc.SetFavoriteSong(id)
	if err != nil {
		ctrl.EncodeEmptyResponse(r, w, err)
		return
	}

	ctrl.EncodeDataResponse(r, w, id, nil)
}

func (ctrl *HTTPController) GetGenres(w http.ResponseWriter, r *http.Request) {
	genres, err := ctrl.Svc.GetGenres()
	if err != nil {
		ctrl.EncodeEmptyResponse(r, w, err)
		return
	}

	ctrl.EncodeDataResponse(r, w, genres, nil)
}

func (ctrl *HTTPController) GetArtistByGenre(w http.ResponseWriter, r *http.Request) {
	genreID := ctrl.GetParamInt(r, "id")

	if genreID == 0 {
		ctrl.EncodeEmptyResponse(r, w, plterror.ErrBadSyntax)
		return
	}
	resp, err := ctrl.Svc.GetArtistByGenre(genreID)
	if err != nil {
		ctrl.EncodeEmptyResponse(r, w, err)
		return
	}
	ctrl.EncodeDataResponse(r, w, resp, nil)
}

func (ctrl *HTTPController) GetAlbumByArtist(w http.ResponseWriter, r *http.Request) {
	artist := ctrl.GetParam(r, "id")

	if artist == "" {
		ctrl.EncodeEmptyResponse(r, w, plterror.ErrBadSyntax)
		return
	}
	resp, err := ctrl.Svc.GetAlbumByArtist(artist)
	if err != nil {
		ctrl.EncodeEmptyResponse(r, w, err)
		return
	}
	ctrl.EncodeDataResponse(r, w, resp, nil)
}

func (ctrl *HTTPController) GetSongsByAlbum(w http.ResponseWriter, r *http.Request) {
	album := ctrl.GetParam(r, "id")
	if album == "" {
		ctrl.EncodeEmptyResponse(r, w, plterror.ErrBadSyntax)
		return
	}

	limit := ctrl.GetParamInt(r, "limit")

	resp, err := ctrl.Svc.GetSongsByAlbum(album, limit)
	if err != nil {
		ctrl.EncodeEmptyResponse(r, w, err)
		return
	}
	ctrl.EncodeDataResponse(r, w, resp, nil)
}

func (ctrl *HTTPController) GetFavorites(w http.ResponseWriter, r *http.Request) {
	var body dto.Favorites
	var ok bool
	body, ok = ctrl.GetBody(r).(dto.Favorites)
	if !ok {
		ctrl.EncodeEmptyResponse(r, w, plterror.ErrBadSyntax)
		return
	}

	resp, err := ctrl.Svc.GetFavorites(body.Genre, body.Artist)
	if err != nil {
		ctrl.EncodeEmptyResponse(r, w, err)
		return
	}
	ctrl.EncodeDataResponse(r, w, resp, nil)
}
