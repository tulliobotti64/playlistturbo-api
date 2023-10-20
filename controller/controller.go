package controller

import (
	"playlistturbo.com/service"
	"playlistturbo.com/utils"
)

// const TokenHeader string = "X-Auth-Token"

// Controller interface
type Controller interface {
	SongsController
}

type HTTPController struct {
	Svc service.Service
	utils.Utils
}

func NewController(s service.Service) Controller {
	return &HTTPController{
		Svc:   s,
		Utils: utils.NewUtils(),
	}
}
