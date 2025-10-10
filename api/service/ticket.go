package service

import (
	"context"
	"fybe/helper"
	"fybe/model/domain"
	"fybe/model/dto"
	"fybe/repository"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type TicketService interface {
	Get(ctx context.Context, userID uint) ([]dto.TicketGet, error)
	GetPublicByUsername(ctx context.Context, username string) ([]dto.TicketPublicGet, error)
	GetPublicByID(ctx context.Context, ticketID string) (*dto.TicketGet, error)
	GetByID(ctx context.Context, ticketID string, userID uint) (*dto.TicketGet, error)
	Create(ctx context.Context, c *fiber.Ctx, ticket dto.TicketCreate, userID uint) error
	Update(ctx context.Context, c *fiber.Ctx, ticket dto.TicketCreate, ticketID string, userID uint) error
	Delete(ctx context.Context, ticketID string, userID uint) error
}

type ticketService struct {
	ticketRepository repository.TicketRepository
	upload           helper.UploadHelper
	redis            helper.RedisHelper
	logger           *zap.Logger
}

func NewTicketService(ticketRepository repository.TicketRepository, upload helper.UploadHelper, redis helper.RedisHelper, logger *zap.Logger) TicketService {
	return &ticketService{
		ticketRepository: ticketRepository,
		upload:           upload,
		redis:            redis,
		logger:           logger,
	}
}

func (s *ticketService) Get(ctx context.Context, userID uint) ([]dto.TicketGet, error) {
	tickets, err := s.ticketRepository.Get(ctx, userID)

	if err != nil {
		s.logger.Error("failed getting ticket", zap.Error(err))
		return nil, err
	}

	var result []dto.TicketGet

	for _, ticket := range tickets {
		result = append(result, dto.TicketGet{
			ID:              ticket.ID,
			Title:           ticket.Title,
			Description:     ticket.Description,
			Mode:            *ticket.Mode,
			RegisteredCount: ticket.RegisteredCount,
			Quota:           ticket.Quota,
			Image:           helper.FormatURL(ticket.ImageURL),
		})
	}

	return result, nil
}

func (s *ticketService) GetPublicByUsername(ctx context.Context, username string) ([]dto.TicketPublicGet, error) {
	tickets, err := s.ticketRepository.GetPublicByUsername(ctx, username)

	if err != nil {
		s.logger.Error("failed getting ticket by username", zap.Error(err))
		return nil, err
	}

	var result []dto.TicketPublicGet

	for _, ticket := range tickets {
		result = append(result, dto.TicketPublicGet{
			ID:    ticket.ID,
			Image: helper.FormatURL(ticket.ImageURL),
			Title: ticket.Title,
		})
	}

	return result, nil
}

func (s *ticketService) GetPublicByID(ctx context.Context, ticketID string) (*dto.TicketGet, error) {
	ticket, err := s.ticketRepository.GetPublicByID(ctx, ticketID)

	if err != nil {
		s.logger.Error("failed getting ticket by id", zap.Error(err))
		return nil, err
	}

	return &dto.TicketGet{
		ID:              ticket.ID,
		Title:           ticket.Title,
		Description:     ticket.Description,
		Mode:            *ticket.Mode,
		Quota:           ticket.Quota,
		RegisteredCount: ticket.RegisteredCount,
		Image:           helper.FormatURL(ticket.ImageURL),
	}, nil
}

func (s *ticketService) GetByID(ctx context.Context, ticketID string, userID uint) (*dto.TicketGet, error) {
	ticket, err := s.ticketRepository.GetByID(ctx, ticketID, userID)

	if err != nil {
		s.logger.Error("failed getting ticket by id", zap.Error(err))
		return nil, err
	}

	return &dto.TicketGet{
		ID:              ticket.ID,
		Title:           ticket.Title,
		Description:     ticket.Description,
		Mode:            *ticket.Mode,
		Quota:           ticket.Quota,
		RegisteredCount: ticket.RegisteredCount,
		Image:           helper.FormatURL(ticket.ImageURL),
	}, nil
}

func (s *ticketService) Create(ctx context.Context, c *fiber.Ctx, ticket dto.TicketCreate, userID uint) error {
	ticketID := helper.GenerateState(6)
	path, err := s.upload.Save(c, "image", ticketID)

	if err != nil {
		s.logger.Error("failed saving ticket file", zap.Error(err))
		return err
	}

	t := &domain.Ticket{
		ID:          ticketID,
		Title:       ticket.Title,
		Description: ticket.Description,
		Mode:        &ticket.Mode,
		Quota:       ticket.Quota,
		ImageURL:    &path,
		UserID:      userID,
	}

	if err := s.ticketRepository.Create(ctx, t); err != nil {
		s.logger.Error("failed creating ticket", zap.Error(err))
		return err
	}

	if err := s.redis.Set(ticketID, ticket.Quota, 12*time.Hour); err != nil {
		s.logger.Error("failed creating ticket in redis", zap.Error(err))
		return err
	}

	return nil
}

func (s *ticketService) Update(ctx context.Context, c *fiber.Ctx, ticket dto.TicketCreate, ticketID string, userID uint) error {
	t := &domain.Ticket{
		Title:       ticket.Title,
		Description: ticket.Description,
		Mode:        &ticket.Mode,
		Quota:       ticket.Quota,
	}

	registeredCount, err := s.ticketRepository.Update(ctx, t, ticketID, userID)

	if err != nil {
		s.logger.Error("failed updating ticket", zap.Error(err))
		return err
	}

	path, err := s.upload.Save(c, "ticket", ticketID)

	if err != nil {
		s.logger.Error("failed saving ticket file", zap.Error(err))
		return err
	}

	if path != "" {
		u := &domain.Ticket{
			ImageURL: &path,
		}

		if _, err := s.ticketRepository.Update(ctx, u, ticketID, userID); err != nil {
			s.logger.Error("failed updating ticket file", zap.Error(err))
			return err
		}
	}

	if err := s.redis.Set(ticketID, t.Quota-registeredCount, 12*time.Hour); err != nil {
		s.logger.Error("failed updating ticket in redis", zap.Error(err))
		return err
	}

	return nil
}

func (s *ticketService) Delete(ctx context.Context, ticketID string, userID uint) error {
	path, err := s.ticketRepository.GetPathByID(ctx, ticketID)

	if err != nil {
		return nil
	}

	if err := s.ticketRepository.Delete(ctx, ticketID, userID); err != nil {
		s.logger.Error("failed deleting ticket file", zap.Error(err))
		return err
	}

	if err := s.redis.Delete(ticketID); err != nil {
		s.logger.Error("failed deleting ticket in redis", zap.Error(err))
		return err
	}

	if err := s.upload.Delete(path); err != nil {
		return nil
	}

	return nil
}
