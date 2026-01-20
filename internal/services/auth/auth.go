package auth

import (
	"context"
	"errors"
	"fmt"
	"go-gRPC-sso/internal/domain/models"
	"go-gRPC-sso/internal/lib/jwt"
	"go-gRPC-sso/internal/lib/logger/sl"
	"go-gRPC-sso/internal/storage"
	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrorInvalidCredentials = errors.New("invalid credentials")
	ErrorInvalidAppID       = errors.New("invalid app id")
	ErrorUserExists         = errors.New("user already exists")
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
	IsAdmin(context context.Context, userID int64) (bool, error)
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
	const op = "auth.Login"

	logger := a.logger.With(
		slog.String("op", op),
		slog.String("email", email),
	)
	logger.Info("attempting to login user")

	user, err := a.userProvider.User(context, email)
	if err != nil {
		if errors.Is(err, storage.ErrorUserExists) {
			a.logger.Warn("user not found")
			return "", fmt.Errorf("%s: %w", op, ErrorInvalidCredentials)
		}
		a.logger.Error("failed to get user", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.logger.Info("invalid credentials", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, ErrorInvalidCredentials)
	}

	app, err := a.appProvider.App(context, appID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	logger.Info("successfully logged in")

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		a.logger.Error("failed to generate token", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return token, nil
}

func (a *Auth) RegisterNewUser(context context.Context, email, password string) (int64, error) {
	const op = "auth.RegisterNewUser"

	logger := a.logger.With(
		slog.String("op", op),
		slog.String("email", email),
	)
	logger.Info("register new user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("failed to generate password hash", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.userSaver.SaveUser(context, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrorUserExists) {
			logger.Warn("user already exists", sl.Err(err))
			return 0, fmt.Errorf("%s: %w", op, ErrorUserExists)
		}
		logger.Error("failed to save user", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (a *Auth) IsAdmin(context context.Context, userID int64) (bool, error) {
	const op = "auth.IsAdmin"

	logger := a.logger.With(
		slog.String("op", op),
		slog.Int64("user_id", userID),
	)
	logger.Info("checking if user is admin")

	isAdmin, err := a.userProvider.IsAdmin(context, userID)
	if err != nil {
		if errors.Is(err, storage.ErrorAppNotFound) {
			logger.Warn("user not found", sl.Err(err))
			return false, fmt.Errorf("%s: %w", op, ErrorInvalidAppID)
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}
	logger.Info("checked if user is admin", slog.Bool("is_admin", isAdmin))

	return isAdmin, nil
}
