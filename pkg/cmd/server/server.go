package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	v0 "github.com/carbonic-app/apis/pkg/api/v0"
	"github.com/carbonic-app/apis/pkg/service/v0/account"
	"github.com/carbonic-app/apis/pkg/service/v0/account/password"
	"github.com/jinzhu/gorm"
	"google.golang.org/grpc"
)

// Config is configuration for Server
type Config struct {
	// gRPC server start parameters section
	// gRPC is TCP port to listen by gRPC server
	GRPCPort int

	// DB Datastore parameters section
	// DatastoreDBHost is host of database
	DatastoreDBHost string
	// DatastoreDBUser is username to connect to database
	DatastoreDBUser string
	// DatastoreDBPassword password to connect to database
	DatastoreDBPassword string
	// DatastoreDBSchema is schema of database
	DatastoreDBSchema string
}

// RunServer runs gRPC server and HTTP gateway
func RunServer() error {
	// get configuration
	var cfg Config
	flag.IntVar(&cfg.GRPCPort, "grpc-port", 0, "gRPC port to bind")
	flag.StringVar(&cfg.DatastoreDBHost, "db-host", "", "Database host")
	flag.StringVar(&cfg.DatastoreDBUser, "db-user", "", "Database user")
	flag.StringVar(&cfg.DatastoreDBPassword, "db-password", "", "Database password")
	flag.StringVar(&cfg.DatastoreDBSchema, "db-schema", "", "Database schema")
	flag.Parse()

	if cfg.GRPCPort <= 0 || cfg.GRPCPort >= 65565 {
		return fmt.Errorf("invalid TCP port for gRPC server: '%d'", cfg.GRPCPort)
	}
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", cfg.GRPCPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	// add MySQL driver specific parameter to parse date/time
	// Drop it for another database
	param := "parseTime=true"

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s",
		cfg.DatastoreDBUser,
		cfg.DatastoreDBPassword,
		cfg.DatastoreDBHost,
		cfg.DatastoreDBSchema,
		param)
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	hasher := password.NewPlaintextHasher()
	v0API := account.NewAccountServiceServer(db, hasher, nil)

	v0.RegisterAccountServiceServer(s, v0API)
	return s.Serve(lis)
}
