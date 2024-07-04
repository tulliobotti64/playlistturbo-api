package jobs

import (
	"fmt"
	"log"
	"reflect"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"playlistturbo.com/database"
	"playlistturbo.com/model"
	"playlistturbo.com/plterror"
)

type AutomigrateJob struct {
	DB database.Database
}

func (j *AutomigrateJob) Run() {
	log.Println("Started automigrate:")
	db := j.DB.GormDB()

	createUUIDExtension(db)
	automigrateTables(db)
	migrateForeignKeyConstraints(db)
	createGenre(db)

	log.Printf("Completed automigrate")
}

func createUUIDExtension(db *gorm.DB) {
	log.Println("Add extensions...")
	err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Error
	if err != nil {
		err = fmt.Errorf("failed to create uuid extension: %w", err)
		plterror.LogFatalError(err.Error())
	}
	err = db.Exec(`CREATE EXTENSION IF NOT EXISTS "unaccent"`).Error
	if err != nil {
		err = fmt.Errorf("failed to create unaccent extension: %w", err)
		plterror.LogFatalError(err.Error())
	}
}

func automigrateTables(db *gorm.DB) {
	log.Println("Creating/updating tables...")

	models := []interface{}{
		&model.Song{},
		&model.Genre{},
	}

	for _, m := range models {
		fmt.Printf("\tautomigrating %v\n", reflect.TypeOf(m))
		done := migrateIgnoreCachedPlanError(db, m)
		if !done {
			err := db.AutoMigrate(m)
			if err != nil {
				err = fmt.Errorf("error creating automigrating tables: %w", err)
				plterror.LogFatalError(err.Error())
			}
		}
	}
}

func migrateForeignKeyConstraints(db *gorm.DB) {
	log.Println("Creating/updating foreign key constraints...")

	type foreignKeyConstraint struct {
		modelStruct      interface{}
		relationshipName string
	}

	fkConstraints := []foreignKeyConstraint{
		{&model.Song{}, "Genre"},
	}

	// delete fk constraint if already present and then recreates it (so if the cascade option is changed, it is updated)
	for _, c := range fkConstraints {
		if db.Migrator().HasConstraint(c.modelStruct, c.relationshipName) {
			err := db.Migrator().DropConstraint(c.modelStruct, c.relationshipName)
			if err != nil {
				err = fmt.Errorf("error deleting old constraint (struct %v , relation %s): %w", c.modelStruct, c.relationshipName, err)
				plterror.LogFatalError(err.Error())
			}
		}

		err := db.Migrator().CreateConstraint(c.modelStruct, c.relationshipName)
		if err != nil {
			err = fmt.Errorf("error creating FKs (struct %v , relation %s): %w", c.modelStruct, c.relationshipName, err)
			plterror.LogFatalError(err.Error())
		}
	}
}

func migrateIgnoreCachedPlanError(db *gorm.DB, m interface{}) bool {
	defer func() {
		_ = recover()
	}()

	if err := db.AutoMigrate(m); err != nil {
		return false
	}

	return true
}

func createGenre(db *gorm.DB) {
	log.Println("Creating Genres ...")

	var genresInt = [...]string{
		"Not Found", "Blues", "Classic rock", "Country", "Dance", "Disco", "Funk", "Grunge",
		"Hip-hop", "Jazz", "Metal", "New age", "Oldies", "Other", "Pop", "Rhythm and blues",
		"Rap", "Reggae", "Rock", "Techno", "Industrial", "Alternative", "Ska", "Death metal",
		"Pranks", "Soundtrack", "Euro-techno", "Ambient", "Trip-hop", "Vocal", "Jazz & funk",
		"Fusion", "Trance", "Classical", "Instrumental", "Acid", "House", "Game", "Sound clip",
		"Gospel", "Noise", "Alternative rock", "Bass", "Soul", "Punk", "Space", "Meditative",
		"Instrumental pop", "Instrumental rock", "Ethnic", "Gothic", "Darkwave",
		"Techno-industrial", "Electronic", "Pop-folk", "Eurodance", "Dream", "Southern rock",
		"Comedy", "Cult", "Gangsta", "Top 40", "Christian rap", "Pop/funk", "Jungle music",
		"Native US", "Cabaret", "New wave", "Psychedelic", "Rave", "Showtunes", "Trailer", "Lo-fi",
		"Tribal", "Acid punk", "Acid jazz", "Polka", "Retro", "Musical", "Rock 'n' roll",
		"Hard rock", "Folk", "Folk rock", "National folk", "Swing", "Fast fusion", "Bebop",
		"Latin", "Revival", "Celtic", "Bluegrass", "Avantgarde", "Gothic rock", "Progressive rock",
		"Psychedelic rock", "Symphonic rock", "Slow rock", "Big band", "Chorus", "Easy listening",
		"Acoustic", "Humour", "Speech", "Chanson", "Opera", "Chamber music", "Sonata", "Symphony",
		"Booty bass", "Primus", "Porn groove", "Satire", "Slow jam", "Club", "Tango", "Samba", "Folklore",
		"Ballad", "Power ballad", "Rhythmic Soul", "Freestyle", "Duet", "Punk rock", "Drum solo",
		"A cappella", "Euro-house", "Dance hall", "Goa music", "Drum & bass", "Club-house",
		"Hardcore techno", "Terror", "Indie", "Britpop", "Negerpunk", "Polsk punk", "Beat",
		"Christian gangsta rap", "Heavy metal", "Black metal", "Crossover", "Contemporary Christian",
		"Christian rock", "Merengue", "Salsa", "Thrash metal", "Anime", "Jpop", "Synthpop", "Christmas",
		"Art rock", "Baroque", "Bhangra", "Big beat", "Breakbeat", "Chillout", "Downtempo", "Dub", "EBM",
		"Eclectic", "Electro", "Electroclash", "Emo", "Experimental", "Garage", "Global", "IDM", "Illbient",
		"Industro-Goth", "Jam Band", "Krautrock", "Leftfield", "Lounge", "Math rock", "New romantic",
		"Nu-breakz", "Post-punk", "Post-rock", "Psytrance", "Shoegaze", "Space rock", "Trop rock",
		"World music", "Neoclassical", "Audiobook", "Audio theatre", "Neue Deutsche Welle", "Podcast",
		"Indie-rock", "G-Funk", "Dubstep", "Garage rock", "Psybient",
	}

	var err error
	genre := []model.Genre{}

	err = db.Table("genres").Find(&genre).Error
	if err != nil {
		err = fmt.Errorf("error opening genre: %w", err)
		plterror.LogFatalError(err.Error())
	}

	if len(genre) == 0 {
		genre = []model.Genre{}
		for _, genreInt := range genresInt {
			g := model.Genre{Name: genreInt}
			genre = append(genre, g)
		}

		err = db.Table("genres").Clauses(clause.OnConflict{DoNothing: true}).Create(&genre).Error
		if err != nil {
			err = fmt.Errorf("error creating Genres: %w", err)
			plterror.LogFatalError(err.Error())
		}
	}
}
