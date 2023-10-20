package service

import "playlistturbo.com/database"

type Components struct {
	DB database.Database
}

type PLTService struct {
	Components
}

type Service interface {
	SongsService

	ExportComponents() *Components
}

func Init() Service {
	var plts PLTService

	plts.DB = database.SetupDB()

	return &plts
}

func (svc *PLTService) ExportComponents() *Components {
	return &svc.Components
}
