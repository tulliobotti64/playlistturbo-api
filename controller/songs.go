package controller

import (
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

func (ctrl *HTTPController) GetSongsByTitle(w http.ResponseWriter, r *http.Request) {
	title := ctrl.GetParam(r, "title")
	// ctrl.Utils.GetParam()
	enableCors(&w)
	if title == "" {
		ctrl.EncodeEmptyResponse(r, w, plterror.ErrBadSyntax)
		return
	}
	resp, err := ctrl.Svc.GetSongsByTitle(title)
	if err != nil {
		ctrl.EncodeEmptyResponse(r, w, err)
		return
	}
	ctrl.EncodeDataResponse(r, w, resp, nil)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
