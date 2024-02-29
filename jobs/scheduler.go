package jobs

import "playlistturbo.com/database"

type Scheduler struct {
	DB database.Database
}

func (s *Scheduler) StartupJobs() {
	automigrateJob := AutomigrateJob{DB: s.DB}
	automigrateJob.Run()
}
