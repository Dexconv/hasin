package db

import (
	"time"

	"github.com/dexconv/hasin/store/common"
	"github.com/dexconv/hasin/store/config"
	"github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	log = common.Log
	DB  = newDB()
)

func newDB() (db *gorm.DB) {
	var err error
	for i := 0; i < 5; i++ {
		err = nil
		db, err = gorm.Open(postgres.Open(config.GLB.DbArgs), &gorm.Config{})
		if err == nil {
			break
		}
		log.Warn("error connecting to postgres", "error", err.Error())
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		log.Error("error connecting to postgres", "error", err.Error())
		panic(err)
	}
	log.Info("connection to db stablished")
	err = db.AutoMigrate(&FileDesc{})
	if err != nil {
		log.Error("error migrating postgres", "error", err.Error())
		panic(err)
	}
	log.Info("migration successful")
	return db
}

func ToPqTextArray(s []string) pq.StringArray {
	return pq.StringArray(s)
}
