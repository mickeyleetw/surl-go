package database

import (
	"fmt"
	"log"
	"net/url"
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
	dbSchema  string = "surl"
	dbDialect string = "postgres"
	dbName    string = "surl"
)

func getDBConnection() string {
	env := os.Getenv("ENV")
	if env == "local" {
		currentDir, _ := os.Getwd()
		log.Printf("Current working directory: %s", currentDir)
		dbenvPath := filepath.Join(currentDir, "../.dbenv")
		log.Printf("Current env file Path: %s", dbenvPath)

		_ = godotenv.Load(dbenvPath)
	}

	dbUser := url.QueryEscape(os.Getenv("DB_USER"))
	dbPassword := url.QueryEscape(os.Getenv("DB_PASSWORD"))
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dsn := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=require&search_path=%s",
		dbDialect, dbUser, dbPassword, dbHost, dbPort, dbName, dbSchema)

	return dsn
}

// GetDB returns a singleton instance of the database
func GetDB() *gorm.DB {
	once.Do(func() {
		dbDsn := getDBConnection()
		db, err := gorm.Open(postgres.Open(dbDsn), &gorm.Config{PrepareStmt: true})
		if err != nil {
			log.Fatal("Error connecting to database: ", err)
		}
		db.Exec("DROP SCHEMA IF EXISTS " + dbSchema + " CASCADE")
		// schema_sql := fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", DB_SCHEMA)
		db.Exec("CREATE SCHEMA IF NOT EXISTS " + dbSchema)
		// search_path_sql := fmt.Sprintf("SET search_path TO %s", DB_SCHEMA)
		db.Exec("SET search_path TO " + dbSchema)
		err = db.AutoMigrate(&domain.URL{})
		if err != nil {
			log.Printf("Failed to auto migrate database: %v", err)
			return
		}
		instance = db
	})
	return instance
}

// UnitOfWork is a struct that implements the unit of work pattern
type UnitOfWork struct {
	DB *gorm.DB
}

// NewUnitOfWork is a function that returns a pointer to a new UnitOfWork struct
func NewUnitOfWork(db *gorm.DB) *UnitOfWork {
	// dB: pointer to the type gorm.DB
	return &UnitOfWork{DB: db} // return a pointer of a new UnitOfWork struct
}

// Begin is a method that returns a pointer to a new gorm.DB transaction
func (uow *UnitOfWork) Begin() *gorm.DB {
	return uow.DB.Begin()
}

// Commit is a method that commits the transaction
func (uow *UnitOfWork) Commit(tx *gorm.DB) error {
	return tx.Commit().Error
}

// Rollback is a method that rolls back the transaction
func (uow *UnitOfWork) Rollback(tx *gorm.DB) error {
	return tx.Rollback().Error
}

// WithTransaction is a function that executes a function within a transaction
func WithTransaction(fn func(tx *gorm.DB) error) error {
	uow := NewUnitOfWork(GetDB())
	tx := uow.Begin()
	var err error

	defer func() {
		if err != nil {
			if rbErr := uow.Rollback(tx); rbErr != nil {
				err = fmt.Errorf("original error: %v, rollback error: %v", err, rbErr)
			}
		} else {
			err = uow.Commit(tx)
		}
	}()

	err = fn(tx)
	return err
}
