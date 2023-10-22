package migrations

import (
	database2 "FioapiKafka/internal/repository/database"
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

func Migrate(cfg *database2.Config) (*sql.DB, error) {
	log.Infoln("Migrate database...")
	db, err := database2.InitdDB(cfg)
	if err != nil {
		log.Errorf("Failed init database:%s", err.Error())
	}
	mig, err := migrate.New(
		"file://internal/repository/database/pgsql",
		fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=%s",
			cfg.DriverName, cfg.Uname, cfg.Pass, cfg.Host, cfg.Port, cfg.NameDB, cfg.SSL))
	if err != nil {
		log.Errorf("Failed database migration:%s", err.Error())
		return nil, err
	}
	err = mig.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Errorf("Error Up migrate database:%s", err.Error())
		return nil, err
	}
	return db, nil
}
