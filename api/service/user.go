package service

import (
	"context"
	"fybe/helper"
	"fybe/model/domain"
	"fybe/model/dto"
	"fybe/repository"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type UserService interface {
	Create(ctx context.Context, user dto.UserCreate) (*uint, error)
	GetByUsername(ctx context.Context, username string) (*dto.UserGet, error)
	GetPublicByTicketID(ctx context.Context, ticketID string) (*dto.UserGet, error)
	GetById(ctx context.Context, userID uint) (*dto.UserGet, error)
	Update(ctx context.Context, c *fiber.Ctx, user dto.UserUpdate, userID uint) error
}

type userService struct {
	userRepository repository.UserRepository
	upload         helper.UploadHelper
	logger         *zap.Logger
}

func NewUserService(userRepository repository.UserRepository, upload helper.UploadHelper, logger *zap.Logger) UserService {
	return &userService{
		userRepository: userRepository,
		upload:         upload,
		logger:         logger,
	}
}

func (s *userService) Create(ctx context.Context, user dto.UserCreate) (*uint, error) {
	u := &domain.User{
		Name:      user.Name,
		GoogleID:  user.GoogleID,
		Email:     user.Email,
		AvatarURL: &user.AvatarURL,
		Username:  helper.ExtractUsername(user.Email),
	}

	if err := s.userRepository.Create(ctx, u); err != nil {
		s.logger.Error("failed to create user", zap.Error(err))
		return nil, err
	}

	s.logger.Info("user created successfully", zap.Uint("userID", u.ID))
	return &u.ID, nil
}

func (s *userService) GetByUsername(ctx context.Context, username string) (*dto.UserGet, error) {
	user, err := s.userRepository.GetByUsername(ctx, username)

	if err != nil {
		s.logger.Error("failed to get user by username", zap.Error(err))
		return nil, err
	}

	return &dto.UserGet{
		Name:      user.Name,
		Username:  user.Username,
		AvatarURL: helper.FormatURL(user.AvatarURL),
	}, nil
}

func (s *userService) GetPublicByTicketID(ctx context.Context, ticketID string) (*dto.UserGet, error) {
	user, err := s.userRepository.GetByTicketID(ctx, ticketID)

	if err != nil {
		s.logger.Error("failed to get user by ticket id", zap.Error(err))
		return nil, err
	}

	return &dto.UserGet{
		Name:      user.Name,
		Username:  user.Username,
		AvatarURL: helper.FormatURL(user.AvatarURL),
	}, nil
}

func (s *userService) GetById(ctx context.Context, userID uint) (*dto.UserGet, error) {
	user, err := s.userRepository.GetById(ctx, userID)

	if err != nil {
		s.logger.Error("failed to get user by id", zap.Error(err))
		return nil, err
	}

	return &dto.UserGet{
		Name:      user.Name,
		Username:  user.Username,
		AvatarURL: helper.FormatURL(user.AvatarURL),
	}, nil
}

func (s *userService) Update(ctx context.Context, c *fiber.Ctx, user dto.UserUpdate, userID uint) error {
	u := &domain.User{
		Name:     user.Name,
		Username: user.Username,
	}

	path, err := s.upload.Save(c, "avatar", strconv.Itoa(int(userID)))

	if err != nil {
		s.logger.Error("failed saving avatar file", zap.Error(err))
		return err
	}

	if path != "" {
		u.AvatarURL = &path
	}

	if err := s.userRepository.Update(ctx, u, userID); err != nil {
		s.logger.Error("failed to update user", zap.Error(err))
		return err
	}

	return nil
}
