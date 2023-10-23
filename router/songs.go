package router

import (
	"net/http"

	"playlistturbo.com/controller"
	"playlistturbo.com/dto"
	"playlistturbo.com/model"
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
			// Params:  []middlewares.Param{},
			Body: nil,
		},
		{
			Path:    "/import",
			Method:  http.MethodPost,
			Handler: ctrl.ImportSongs,
			Body:    dto.ImportSongs{},
		},
	}
}
