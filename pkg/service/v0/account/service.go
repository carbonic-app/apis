package account

import (
	"context"
	"fmt"

	v0 "github.com/carbonic-app/apis/pkg/api/v0"
	"github.com/carbonic-app/apis/pkg/service/v0/common"
	"github.com/jinzhu/gorm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	// apiVersion is the version of API provided by this server
	apiVersion = "v0"
)

type accountServiceServer struct {
	db   *gorm.DB
	app  internal
	auth common.Auth
}

type internal interface {
	HashPassword(string) string
}

// NewAccountServiceServer creates an Account Service
func NewAccountServiceServer(db *gorm.DB, app internal, auth common.Auth) *accountServiceServer {
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

// Create builds an account and returns it
func (s *accountServiceServer) Create(ctx context.Context, req *v0.CreateRequest) (*v0.TokenResponse, error) {
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// TODO: Add a Password hash function and
	user := User{Username: req.Username, PasswordHash: req.Password}
	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}

	// TODO: Add actual JWT library to serialize tokens
	t := v0.Token{Data: fmt.Sprint(user.ID)}
	return &v0.TokenResponse{Api: apiVersion, Token: &t}, nil
}
