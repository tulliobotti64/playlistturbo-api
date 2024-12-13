package database

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"playlistturbo.com/config"
)

type Database interface {
	MusicDatabase
	GormDB() *gorm.DB
	Close()
}

type PostgresDB struct {
	Gorm *gorm.DB
}

func (p *PostgresDB) GormDB() *gorm.DB {
	return p.Gorm
}

func SetupDB() *PostgresDB {
	var db PostgresDB
	var err error

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,  // Slow SQL threshold
			LogLevel:                  logger.Error, // Log level logger.Info or Error
			IgnoreRecordNotFoundError: true,         // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,         // Disable color
		},
	)

	gormConfig := gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   newLogger,
	}

	if db.Gorm, err = gorm.Open(postgres.New(postgres.Config{
		DSN: "host=" + config.Config.Database.DBAddr +
			" port=" + config.Config.Database.DBPort +
			" dbname=" + config.Config.Database.DBName +
			" user=" + config.Config.Database.DBUser +
			" password=" + config.Config.Database.DBPass +
			" sslmode=" + config.Config.Database.DBType,
	}), &gormConfig); err != nil { // disable logger = logger.silence
		log.Fatal("Error connecting with the DB", err)
	}

	sqlDB, _ := db.Gorm.DB()
	sqlDB.SetMaxOpenConns(config.Config.Database.MaxOpenConns)

	log.Println("Successfully connected with database")

	return &db
}

func (p *PostgresDB) Close() {
	sqlDB, _ := p.Gorm.DB()
	_ = sqlDB.Close()
}

var ErrNotFound = fmt.Errorf("record not found")

func handleError(err error) error {
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = ErrNotFound
		}
	}
	return err
}
