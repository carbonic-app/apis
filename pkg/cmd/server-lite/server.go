package main

import (
	"flag"
	"fmt"
	"log"
	"net"

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
	// gRPC server start parameters section
	// gRPC is TCP port to listen by gRPC server
	GRPCPort int

	// DB Datastore parameters section
	// DatastoreDBFile is the sqlite db file
	DatastoreDBFile string
}

// RunServer runs gRPC server and HTTP gateway
func RunServer() error {
	// get configuration
	var cfg Config
	flag.IntVar(&cfg.GRPCPort, "grpc-port", 0, "gRPC port to bind")
	flag.StringVar(&cfg.DatastoreDBFile, "db-file", "", "Database file path")
	flag.Parse()

	if cfg.GRPCPort <= 0 || cfg.GRPCPort >= 65565 {
		return fmt.Errorf("invalid TCP port for gRPC server: '%d'", cfg.GRPCPort)
	}
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", cfg.GRPCPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	db, err := gorm.Open("sqlite3", cfg.DatastoreDBFile)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	hasher := password.NewPlaintextHasher()
	session := auth.NewInMemorySession()
	v0API := account.NewAccountServiceServer(db, hasher, session)

	v0.RegisterAccountServiceServer(s, v0API)
	return s.Serve(lis)
}
