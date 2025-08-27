package service

import (
	"auth-service/internal/jwt"
	"auth-service/internal/logger/sl"
	"auth-service/internal/model"
	"auth-service/internal/repository"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserSaver interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (uid int64, err error)
}

type UserProvider interface {
	GetUser(ctx context.Context, email string) (model.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (model.App, error)
}

type Auth struct {
	log         *slog.Logger
	usrSaver    UserSaver
	usrProvider UserProvider
	appProvider AppProvider
	tokenTTL    time.Duration
	jwtSecret   string
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAppID       = errors.New("invalid app id")
	ErrUserExists         = errors.New("user already exisits")
)

// New returns a new instance of the Auth service.
func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration,
	jwtSecret string,

) *Auth {
	return &Auth{
		usrSaver:    userSaver,
		usrProvider: userProvider,
		log:         log,
		appProvider: appProvider,
		tokenTTL:    tokenTTL,
		jwtSecret:   jwtSecret,
	}
}

func (a *Auth) Login(ctx context.Context, email, password string, appID int) (string, error) {
	const op = "auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("username", email),
	)

	log.Info("attempting to login user")

	// получить пользователя
	user, err := a.usrProvider.GetUser(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			a.log.Warn("user not found", sl.Err(err))

			return "", fmt.Errorf("%s:%w", op, ErrInvalidCredentials)
		}

		a.log.Error("failed to get user", sl.Err(err))

		return "", fmt.Errorf("%s:%w", op, err)

	}

	// проверка пароля
	err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(password))
	if err != nil {
		log.Error("invalid credentials", sl.Err(err))

		return "", fmt.Errorf("%s:%w", op, ErrInvalidCredentials)
	}

	// получить приложение в которое пользователь хочет залогинится
	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		return "", fmt.Errorf("%s:%w", op, ErrInvalidCredentials)
	}

	log.Info("user logged in succesfully")

	// создаем токен
	token, err := jwt.NewToken(user, app, a.jwtSecret, a.tokenTTL)
	if err != nil {
		return "", fmt.Errorf("%s:%w", op, ErrInvalidCredentials)
	}

	return token, nil

}

func (a *Auth) Register(ctx context.Context, email, password string) (int64, error) {
	const op = "auth.Register"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("registering user")

	// хешируем пароль
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", sl.Err(err))

		return 0, fmt.Errorf("%s:%w", op, err)
	}

	// сохранение в бд
	id, err := a.usrSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, repository.ErrUserExists) {
			a.log.Warn("user not found", sl.Err(err))

			return 0, fmt.Errorf("%s:%w", op, ErrUserExists)
		}

		a.log.Error("failed to save user", sl.Err(err))

		return 0, fmt.Errorf("%s:%w", op, err)
	}

	log.Info("user registered")

	return id, nil
}

func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "auth.IsAdmin"

	log := a.log.With(
		slog.String("op", op),
		slog.Int64("email", userID),
	)

	log.Info("checking if user is admin")

	isAdmin, err := a.usrProvider.IsAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrAppNotFound) {
			a.log.Warn("user not found", sl.Err(err))

			return false, fmt.Errorf("%s:%w", op, ErrInvalidAppID)
		}

		a.log.Error("failed to get user", sl.Err(err))

		return false, fmt.Errorf("%s:%w", op, err)

	}

	log.Info("checked if user is admin", slog.Bool("is_admin", isAdmin))

	return isAdmin, nil

}
