package account

import (
	"context"
	"database/sql"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	// apiVersion is the version of API provided by this server
	apiVersion = "v0"
)

type accountServiceServer struct {
	db *sql.DB
}

// NewAccountServiceServer creates an Account Service
func NewAccountServiceServer(db *sql.DB) *accountServiceServer {
	return &accountServiceServer{db: db}
}

func (s *accountServiceServer) checkAPI(api string) error {
	// API version is "" means use current version of the service
	if len(api) > 0 {
		if apiVersion != api {
			return status.Errorf(
				codes.Unimplemented,
				"unsupported API version: service implements API version '%s', but asked for '%s'",
				apiVersion,
				api,
			)
		}
	}
	return nil
}

// connect returns SQL database connection from the pool
func (s *accountServiceServer) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := s.db.Conn(ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, "Failed to connect to to database-> "+err.Error())
	}
	return c, nil
}

// Create builds an account and returns it
