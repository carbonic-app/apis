package main

import (
	"fmt"
	"log"
	"net"
	"os"

	v0 "github.com/carbonic-app/apis/pkg/api/v0"
	"github.com/carbonic-app/apis/pkg/service/v0/account"
	"github.com/carbonic-app/apis/pkg/service/v0/account/password"
	"github.com/carbonic-app/apis/pkg/service/v0/common/auth"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"google.golang.org/grpc"
)

// Config is configuration for Server
type Config struct {
	// PasswordSecret is the password hashing secret
	PasswordSecret string

	// DB Datastore parameters section
	// DBFilePath is the sqlite db file path
	DBFilePath string
	// DBMigrate is if the db should be migrated on start
	DBMigrate bool
}

// RunServer runs gRPC server and HTTP gateway
func RunServer(l net.Listener) error {
	var cfg Config

	// Default: false, anything else: true
	cfg.DBMigrate = os.Getenv("DB_MIGRATE") != ""
	mustMapEnv(&cfg.PasswordSecret, "PASSWORD_SECRET")
	mustMapEnv(&cfg.DBFilePath, "DB_FILE_PATH")

	db, err := initDB(cfg)
	if err != nil {
		log.Fatal("Failed to init DB:", err)
	}
	defer db.Close()

	s := grpc.NewServer()

	hasher := password.NewCryptHasher(cfg.PasswordSecret)
	session := auth.NewInMemorySession()
	v0API := account.NewAccountServiceServer(db, hasher, session)

	v0.RegisterAccountServiceServer(s, v0API)
	return s.Serve(l)
}

func mustMapEnv(target *string, envKey string) {
	v := os.Getenv(envKey)
	if v == "" {
		log.Panic("Environment variable not set:", envKey)
	}
	*target = v
}

func initDB(cfg Config) (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", cfg.DBFilePath)
	if err != nil {
		return nil, fmt.Errorf("Failed to open database: %v", err)
	}
	if cfg.DBMigrate {
		if err := autoMigrate(db); err != nil {
			return nil, fmt.Errorf("Failed AutoMigration: %v", err)
		}
	}
	return db, nil
}

func autoMigrate(db *gorm.DB) error {
	var err error
	if !db.HasTable(&account.User{}) {
		err = db.AutoMigrate(&account.User{}).Error
	}
	return err
}
