package account

import (
	"context"

	v0 "github.com/carbonic-app/apis/pkg/api/v0"
	"github.com/carbonic-app/apis/pkg/service/v0/account/password"
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
	db     *gorm.DB
	hasher password.Hasher
	auth   common.Auth
}

type internal interface {
	HashPassword(string) string
}

// NewAccountServiceServer creates an Account Service
func NewAccountServiceServer(db *gorm.DB, h password.Hasher, a common.Auth) *accountServiceServer {
	return &accountServiceServer{db: db, hasher: h, auth: a}
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

	h := s.hasher.HashPassword(req.Password)
	user := User{Username: req.Username, PasswordHash: h}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}

	t := v0.Token{Data: s.auth.GenerateToken(user.ID)}
	return &v0.TokenResponse{Token: &t}, nil
}

// Login tests an account
func (s *accountServiceServer) Login(ctx context.Context, req *v0.LoginRequest) (*v0.TokenResponse, error) {
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	var user User
	var token string

	if err := s.db.Where(&User{Username: req.Username}).First(&user).Error; err != nil {
		return nil, err
	}

	h := s.hasher.HashPassword(req.Password)
	if h == user.PasswordHash {
		token = s.auth.GenerateToken(user.ID)
	} else {
		return nil, status.Error(codes.PermissionDenied, "Invalid Password")
	}
	return &v0.TokenResponse{Token: &v0.Token{Data: token}}, nil
}
