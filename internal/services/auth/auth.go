package auth

import (
	"context"
	"go-gRPC-sso/internal/domain/models"
	"log/slog"
	"time"
)

type Auth struct {
	logger       *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}

type UserSaver interface {
	SaveUser(context context.Context, email string, passHash []byte) (uid int64, err error)
}

type UserProvider interface {
	User(context context.Context, email string) (models.User, error)
	IsAdmin(context context.Context, email string) (bool, error)
}

type AppProvider interface {
	App(context context.Context, appID int) (models.App, error)
}

func New(
	logger *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		userSaver:    userSaver,
		userProvider: userProvider,
		appProvider:  appProvider,
		logger:       logger,
		tokenTTL:     tokenTTL,
	}
}

func (a *Auth) Login(context context.Context, email string, password string, appID int) (string, error) {
	panic("implement me")
}

func (a *Auth) RegisterNewUser(context context.Context, email, password string) (int64, error) {
	panic("implement me")
}

func (a *Auth) IsAdmin(context context.Context, userID int64) (bool, error) {
	panic("implement me")
}
