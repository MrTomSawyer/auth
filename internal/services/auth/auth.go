package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/MrTomSawyer/sso/internal/domain/models"
	"github.com/MrTomSawyer/sso/internal/storage"
	"github.com/MrTomSawyer/sso/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

type Auth struct {
	log          *slog.Logger
	UserSaver    UserSaver
	UserProvider UserProvider
	tokenTTL     time.Duration
	secret       string
}

type UserSaver interface {
	SaveUser(ctx context.Context, email string, passwordHash []byte) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

func New(
	log *slog.Logger,
	UserSaver UserSaver,
	UserProvider UserProvider,
	tokenTTL time.Duration) *Auth {
	return &Auth{
		log:          log,
		UserSaver:    UserSaver,
		UserProvider: UserProvider,
		tokenTTL:     tokenTTL,
	}
}

func (a *Auth) Login(ctx context.Context, email string, password string) (string, error) {
	const op = "Auth.Login"

	log := a.log.With(slog.String("op", op), slog.String("username", email))
	log.Info("attempting to login user")

	user, err := a.UserProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found", err.Error())
		}

		a.log.Error("failed to get user", err.Error())
		return "", fmt.Errorf("%s: %w", op, storage.ErrInvalidCredentials)
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)); err != nil {
		a.log.Info("invalid credentials", err.Error())
		return "", fmt.Errorf("%s: %w", op, storage.ErrInvalidCredentials)
	}

	a.log.Info("user has successfully logged")

	token, err := jwt.NewToken(user, a.tokenTTL, a.secret)
	if err != nil {
		a.log.Error("failed to generate token", err.Error())
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (a *Auth) Register(ctx context.Context, email string, password string) (int64, error) {
	const op = "auth.RegisterNewUser"

	log := a.log.With(slog.String("op", op), slog.String("email", email))
	log.Info("registering user")

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to hash password", err.Error())
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.UserSaver.SaveUser(ctx, email, hash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn("user already exists", storage.ErrUserExists)
			return 0, fmt.Errorf("%s: %w", op, err)
		}
		log.Error("failed to save user", err.Error())
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "Auth.IsAdmin"
	log := a.log.With(
		slog.String("op", op),
		slog.Int64("user_id", userID),
	)

	log.Info("checking if user is an admin")

	isAdmin, err := a.UserProvider.IsAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found", err.Error())
			return false, fmt.Errorf("%s: %w", op, err)
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("checked if user is admin", slog.Bool("is_admin", isAdmin))

	return isAdmin, nil
}

func (a *Auth) User(ctx context.Context, email string) (bool, error) {
	const op = "Auth.User"

	log := a.log.With(slog.String("op", op), slog.String("email", email))
	log.Info("registering user")

	_, err := a.UserProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found", err.Error())
			return false, fmt.Errorf("%s: %w", op, err)
		}
		log.Warn("error when finding user", err.Error())
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return true, nil
}
