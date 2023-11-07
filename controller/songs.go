package controller

import (
	"net/http"

	"playlistturbo.com/dto"
	"playlistturbo.com/model"
	"playlistturbo.com/plterror"
)

type SongsController interface {
	AddSong(w http.ResponseWriter, r *http.Request)
	GetMainList(w http.ResponseWriter, r *http.Request)
	ImportSongs(w http.ResponseWriter, r *http.Request)
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
	body, ok = ctrl.GetBody(r).(dto.ImportSongs)
	if !ok {
		ctrl.EncodeEmptyResponse(r, w, plterror.ErrBadSyntax)
		return
	}

	songs, err := ctrl.Svc.ImportSongs(body)
	if err != nil {
		ctrl.EncodeEmptyResponse(r, w, err)
		return
	}

	if len(body.Genre) == 0 {
		ctrl.EncodeEmptyResponse(r, w, plterror.InvalidGenre)
		return
	}
	// if len(body.Artist) == 0 {
	// 	ctrl.EncodeEmptyResponse(r, w, plterror.InvalidArtist)
	// 	return
	// }

	ctrl.EncodeDataResponse(r, w, songs, nil)
}
