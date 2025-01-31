package database

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	domain "shorten_url/pkg/domains"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	instance  *gorm.DB
	once      sync.Once
	DB_SCHEMA string = "surl"
)

func getDBConnection() string {
	ENV := os.Getenv("ENV")
	if ENV == "local" {
		currentDir, _ := os.Getwd()
		log.Printf("Current working directory: %s", currentDir)
		dbenvPath := filepath.Join(currentDir, "../.dbenv")
		log.Printf("Current env file Path: %s", dbenvPath)

		_ = godotenv.Load(dbenvPath)
	}

	DB_DIALECT := os.Getenv("DB_DIALECT")
	DB_USER := os.Getenv("DB_USER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_HOST := os.Getenv("DB_HOST")
	DB_PORT := os.Getenv("DB_PORT")
	DB_NAME := os.Getenv("DB_NAME")
	DB_SCHEMA = os.Getenv("DB_SCHEMA")
	dsn := DB_DIALECT + "://" + DB_USER + ":" + DB_PASSWORD + "@" + DB_HOST + ":" + DB_PORT + "/" + DB_NAME + "?sslmode=disable&search_path=" + DB_SCHEMA

	return dsn
}

func GetDB() *gorm.DB {
	once.Do(func() {
		db_dsn := getDBConnection()
		db, err := gorm.Open(postgres.Open(db_dsn), &gorm.Config{PrepareStmt: true})
		if err != nil {
			log.Fatal("Error connecting to database: ", err)
		}
		// schema_sql := fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", DB_SCHEMA)
		db.Exec("CREATE SCHEMA IF NOT EXISTS " + DB_SCHEMA)
		// search_path_sql := fmt.Sprintf("SET search_path TO %s", DB_SCHEMA)
		db.Exec("SET search_path TO " + DB_SCHEMA)
		db.AutoMigrate(&domain.Url{})
		instance = db
	})
	return instance
}

type UnitOfWork struct {
	DB *gorm.DB
}

func NewUnitOfWork(db *gorm.DB) *UnitOfWork {
	// dB: pointer to the type gorm.DB
	return &UnitOfWork{DB: db} // return a pointer of a new UnitOfWork struct
}

func (uow *UnitOfWork) Begin() *gorm.DB {
	return uow.DB.Begin()
}

func (u *UnitOfWork) Commit(tx *gorm.DB) error {
	return tx.Commit().Error
}

func (u *UnitOfWork) Rollback(tx *gorm.DB) error {
	return tx.Rollback().Error
}

func WithTransaction(fn func(tx *gorm.DB) error) error {
	uow := NewUnitOfWork(GetDB())
	tx := uow.Begin()
	var err error

	defer func() {
		if err != nil {
			uow.Rollback(tx)
		} else {
			err = uow.Commit(tx)
		}
	}()

	err = fn(tx)
	return err
}
