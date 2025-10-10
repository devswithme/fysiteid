package service

import (
	"context"
	"encoding/json"
	"errors"
	"fybe/config"
	"fybe/helper"
	"fybe/model/dto"
	"fybe/repository"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type AuthService interface {
	SaveState(redirect string) (string, error)
	Store(ctx context.Context, c *fiber.Ctx, user dto.UserCreate) error
	Callback(c *fiber.Ctx) (string, error)
	VerifyState(state string) (string, error)
	FetchProfile(ctx context.Context, code string) (*dto.UserCreate, error)
	Refresh(c *fiber.Ctx, userID uint) error
}

type authService struct {
	userRepository repository.UserRepository
	userService    UserService
	redis          helper.RedisHelper
	logger         *zap.Logger
}

func NewAuthService(
	userRepository repository.UserRepository,
	userService UserService,
	redis helper.RedisHelper,
	logger *zap.Logger,
) AuthService {
	return &authService{
		userRepository: userRepository,
		userService:    userService,
		redis:          redis,
		logger:         logger,
	}
}

func (s *authService) SaveState(redirect string) (string, error) {
	state := helper.GenerateState(6)

	if err := s.redis.Set(state, redirect, time.Minute); err != nil {
		s.logger.Error("redis failed storing state", zap.Error(err))
		return "", err
	}

	s.logger.Info("save state to redis",
		zap.String("key", state),
		zap.String("value", redirect),
		zap.Duration("exp", time.Minute),
	)

	return state, nil
}

func (s *authService) VerifyState(state string) (string, error) {
	redirect, err := s.redis.Get(state)

	if err != nil {
		if errors.Is(err, redis.Nil) {
			s.logger.Error("state not found in redis", zap.String("state", state))
		}
		return "", err
	}

	if err := s.redis.Delete(state); err != nil {
		s.logger.Error("err deleting state in redis", zap.String("state", state))
		return "", err
	}

	return redirect, nil
}

func (s *authService) FetchProfile(ctx context.Context, code string) (*dto.UserCreate, error) {
	token, err := config.GoogleOAuthConfig.Exchange(ctx, code)

	if err != nil {
		return nil, err
	}

	client := config.GoogleOAuthConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	user := &dto.UserCreate{}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Store(ctx context.Context, c *fiber.Ctx, user dto.UserCreate) error {
	userID, err := s.userService.Create(ctx, user)

	if err != nil {
		s.logger.Error("failed to create user", zap.Error(err))
		return err
	}

	aToken, err := helper.NewTokenHelper(os.Getenv("A_SECRET")).Generate(*userID, time.Hour*24)
	if err != nil {
		s.logger.Error("failed to generate access token", zap.Error(err))
		return err
	}

	rToken, err := helper.NewTokenHelper(os.Getenv("R_SECRET")).Generate(*userID, time.Hour*24)
	if err != nil {
		s.logger.Error("failed to generate refresh token", zap.Error(err))
		return err
	}

	helper.NewCookieHelper(c, "access_token").Set(aToken, int(24*time.Hour/time.Second))
	helper.NewCookieHelper(c, "refresh_token").Set(rToken, int(7*24*time.Hour/time.Second))

	return nil
}

func (s *authService) Callback(c *fiber.Ctx) (string, error) {
	state := c.Query("state")
	code := c.Query("code")

	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	redirect, err := s.VerifyState(state)

	if err != nil {
		s.logger.Error("failed to verify redis", zap.Error(err))
		return "", err
	}

	user, err := s.FetchProfile(ctx, code)
	if err != nil {
		s.logger.Error("failed to fetch profile", zap.Error(err))
		return "", err
	}

	if err := s.Store(ctx, c, *user); err != nil {
		s.logger.Error("failed to fetch profile", zap.Error(err))
		return "", err
	}

	return redirect, nil
}

func (s *authService) Refresh(c *fiber.Ctx, userID uint) error {
	aToken, err := helper.NewTokenHelper(os.Getenv("A_SECRET")).Generate(userID, time.Hour*24)

	if err != nil {
		s.logger.Error("failed to refresh token", zap.Error(err))
		return err
	}

	helper.NewCookieHelper(c, "access_token").Set(aToken, int(24*time.Hour/time.Second))

	return nil
}
