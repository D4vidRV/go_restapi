package db

import (
	"log"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dns  = "postgres://postgresuser:postgrespass@localhost:5433/localpostgres?sslmode=disable"
	db   *gorm.DB
	once sync.Once
)

func NewDB() {
	once.Do(func() {
		var err error
		db, err = gorm.Open(postgres.Open(dns), &gorm.Config{})

		if err != nil {
			log.Fatalf("Error al intentar conectar con la DB: %v", err)
		}

		log.Println("DB connected!")
	})
}

// return a unique instance of DB
func DB() *gorm.DB {
	return db
}
