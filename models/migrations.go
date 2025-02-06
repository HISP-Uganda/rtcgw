package models

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	log "github.com/sirupsen/logrus"
	"rtcgw/config"
	"runtime"
)

func init() {
	var migrationsDir string
	currentOS := runtime.GOOS
	switch currentOS {
	case "windows":
		migrationsDir = "file:///C:\\ProgramData\\Rtcgw"
	case "darwin", "linux":
		migrationsDir = config.RTCGwConf.Server.MigrationsDirectory
	default:
		migrationsDir = "file://db/migrations"
	}
	m, err := migrate.New(
		migrationsDir,
		config.RTCGwConf.Database.URI)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Error running migration:", err)
	}

	if err != nil {
		log.Fatalln(err)
	}
}
