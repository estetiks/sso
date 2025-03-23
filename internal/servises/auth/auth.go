package auth

import (
	"context"
	"crypto"
	"errors"
	"fmt"
	"hash"
	"log/slog"
	"time"

	"github.com/estetiks/sso/internal/domain/models"
	"github.com/estetiks/sso/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Auth struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (uid int64, err error)
}

type UserProvider interface {
	User(
		ctx context.Context,
		email string,
	) (models.User, error)
}

type AppProvider interface {
	App(
		ctx context.Context,
		appID int,
	) (models.App, error)
}

// New returns a new instance of the Auth service
func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		log:          log,
		userSaver:    userSaver,
		userProvider: userProvider,
		appProvider:  appProvider,
		tokenTTL:     tokenTTL,
	}
}

func (a *Auth) Login(
	ctx context.Context,
	email string,
	password string,
	appID int,
) (string, error) {

	const op = "Auth.Login"

	log := a.log.With(slog.String("op", op), slog.String("email", email))

	log.Info("attempting to login user")

	user, err := a.userProvider.User(ctx, email)

	if err != nil {
		if errors.Is(err, storage.ErrUserAlreadyExist) {
			a.log.Warn("user not foung")

			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

	}

	if err := bcrypt.CompareHashAndPassword([]byte(password), user.PassHash); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword)
	}
}

func (a *Auth) RegisterNewUser(
	ctx context.Context,
	email string,
	password string,
	appID int,
) (string, error) {
	panic("not implimented")
}
