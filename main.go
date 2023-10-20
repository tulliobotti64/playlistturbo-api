package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"playlistturbo.com/config"
	"playlistturbo.com/controller"
	"playlistturbo.com/jobs"
	"playlistturbo.com/router"
	"playlistturbo.com/service"
)

const Version = "0.0.1"

func main() {
	config.SetupConfig()

	ds := service.Init()
	dgsComp := ds.ExportComponents()

	defer dgsComp.DB.Close()

	// initializes schedules
	sched := &jobs.Scheduler{
		DB: dgsComp.DB,
	}

	if config.Config.StartupJobs {
		sched.StartupJobs()
	}

	// Initialize the controller
	ctrl := controller.NewController(ds)

	// Initialize the router
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Mount("/api", router.Get(ctrl))
	// VERSION
	// plterror.Logger.Info("Version: ", Version)

	// Start server
	// plterror.Logger.Info("starting plt-backend on port:", config.Config.Server.Port)
	log.Println("starting plt-backend on port:", config.Config.Server.Port)
	if err := http.ListenAndServe(":"+config.Config.Server.Port, r); err != nil {
		log.Println("error", err)
		// plterror.Logger.Error("http server error:", err)
	}
}
