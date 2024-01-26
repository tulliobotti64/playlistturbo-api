package router

import (
	"net/http"

	"playlistturbo.com/controller"
	"playlistturbo.com/dto"
	"playlistturbo.com/model"
	"playlistturbo.com/router/middlewares"
)

func SongsRoutes(ctrl controller.Controller) []Route {
	return []Route{
		{
			Path:    "/mainlist",
			Method:  http.MethodPost,
			Handler: ctrl.AddSong,
			Body:    model.Song{},
		},
		{
			Path:    "/mainlist",
			Method:  http.MethodGet,
			Handler: ctrl.GetMainList,
			Body:    nil,
		},
		{
			Path:    "/import",
			Method:  http.MethodPost,
			Handler: ctrl.ImportSongs,
			Body:    dto.ImportSongs{},
		},
		{
			Path:    "/import",
			Method:  http.MethodPut,
			Handler: ctrl.MoveSongs,
			Body:    dto.MoveSongs{},
		},
		{
			Path:    "/import",
			Method:  http.MethodDelete,
			Handler: ctrl.RemoveSongs,
			Body:    dto.ImportSongs{},
		},
		{
			Path:    "/{title}",
			Method:  http.MethodGet,
			Handler: ctrl.GetSongsByTitle,
			Body:    nil,
			Params: middlewares.Params{
				{
					Name:       "title",
					Validation: "required",
				},
				{
					Name:       "limit",
					Validation: "required",
				},
			},
		},
		{
			Path:    "/{id}",
			Method:  http.MethodPut,
			Handler: ctrl.SetFavoriteSong,
			Body:    nil,
			Params: middlewares.Params{
				{
					Name:       "id",
					Validation: "required",
				},
			},
		},
		{
			Path:    "/genres",
			Method:  http.MethodGet,
			Handler: ctrl.GetGenres,
			Body:    nil,
		},
		{
			Path:    "/artistbygenre/{id}",
			Method:  http.MethodGet,
			Handler: ctrl.GetArtistByGenre,
			Body:    nil,
			Params: middlewares.Params{
				{
					Name:       "id",
					Validation: "required",
				},
			},
		},
		{
			Path:    "/albumbyartist/{id}",
			Method:  http.MethodGet,
			Handler: ctrl.GetAlbumByArtist,
			Body:    nil,
			Params: middlewares.Params{
				{
					Name:       "id",
					Validation: "required",
				},
			},
		},
		{
			Path:    "/songsbyalbum/{id}",
			Method:  http.MethodGet,
			Handler: ctrl.GetSongsByAlbum,
			Body:    nil,
			Params: middlewares.Params{
				{
					Name:       "id",
					Validation: "required",
				},
				{
					Name:       "limit",
					Validation: "required",
				},
			},
		},
	}
}
