package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/estetiks/sso/internal/domain/models"
	jwt_sso "github.com/estetiks/sso/internal/lib/jwt"
	"github.com/estetiks/sso/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAppID       = errors.New("invalid appID")
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
			log.Warn("user not foung")

			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		log.Info("failed to get user")

		return "", fmt.Errorf("%s: %w", op, err)

	}

	if err := bcrypt.CompareHashAndPassword([]byte(password), user.PassHash); err != nil {
		log.Info("failed to compare hash password")

		return "", fmt.Errorf("%s: %w", op, err)
	}

	app, err := a.appProvider.App(ctx, appID)

	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			log.Warn("application not found")

			return "", fmt.Errorf("%s: %w", op, ErrInvalidAppID)
		}

		log.Info("failed to get application")

		return "", fmt.Errorf("%s: %w", op, err)
	}

	token, err := jwt_sso.NewToken(user, app, a.tokenTTL)

	if err != nil {
		log.Info("failed to create jwt token")

		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user logged in successfuly")

	return token, nil

}

func (a *Auth) RegisterNewUser(
	ctx context.Context,
	email string,
	password string,
) (int64, error) {
	const op = "auth.RegisterNewUser"

	log := a.log.With(slog.String("op", op), slog.String("email", email))

	log.Info("attempting to register new user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		log.Error("failed to generate password hash")

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	uid, err := a.userSaver.SaveUser(ctx, email, passHash)

	if err != nil {
		log.Info("failed to save new user")

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return uid, nil
}
