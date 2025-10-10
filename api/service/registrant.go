package service

import (
	"context"
	"fmt"
	"fybe/helper"
	"fybe/model/domain"
	"fybe/model/dto"
	"fybe/repository"
	"math"
	"os"
	"strconv"
	"time"

	"go.uber.org/zap"
)

type RegistrantService interface {
	GetByUserID(ctx context.Context, userID uint) ([]dto.RegistrantGetByUserID, error)
	GetByTicketID(ctx context.Context, ticketID string, userID uint, page int, limit int, search string) (*dto.PaginatedRegistrants, error)
	Create(ctx context.Context, ticketID string, state string, userID uint) error
	Verify(ctx context.Context, id string, ticketID string, userID uint) error
	Generate(ctx context.Context, ticketID string, userID uint) (string, error)
}

type registrantService struct {
	registrantRepository repository.RegistrantRepository
	ticketRepository     repository.TicketRepository
	redis                helper.RedisHelper
	logger               *zap.Logger
}

func NewRegistrantService(registrantRepository repository.RegistrantRepository, ticketRepository repository.TicketRepository, redis helper.RedisHelper, logger *zap.Logger) RegistrantService {
	return &registrantService{
		registrantRepository: registrantRepository,
		ticketRepository:     ticketRepository,
		redis:                redis,
		logger:               logger,
	}
}

func (s *registrantService) GetByUserID(ctx context.Context, userID uint) ([]dto.RegistrantGetByUserID, error) {
	registrants, err := s.registrantRepository.GetByUserID(ctx, userID)

	if err != nil {
		s.logger.Error("failed to get registrants by user id", zap.Error(err))
		return nil, err
	}

	var result []dto.RegistrantGetByUserID

	for _, registrant := range registrants {
		result = append(result, dto.RegistrantGetByUserID{
			ID:         registrant.ID,
			URL:        os.Getenv("BE_DOMAIN") + "/api/v1/registrant/verify/" + registrant.ID + "/" + registrant.TicketID,
			Title:      registrant.Title,
			ImageURL:   helper.FormatURL(registrant.ImageURL),
			IsVerified: registrant.IsVerified,
			CreatedAt:  registrant.CreatedAt,
		})
	}

	return result, nil
}

func (s *registrantService) GetByTicketID(ctx context.Context, ticketID string, userID uint, page int, limit int, search string) (*dto.PaginatedRegistrants, error) {
	data, total, err := s.registrantRepository.GetByTicketID(ctx, ticketID, userID, page, limit, search)

	if err != nil {
		s.logger.Error("failed to get registrants by ticket id", zap.Error(err))
		return nil, err
	}

	mappedData := make([]dto.RegistrantGetByTicketID, len(data))
	for i, item := range data {
		mappedData[i] = dto.RegistrantGetByTicketID{
			Username:   item.Username,
			AvatarURL:  helper.FormatURL(item.AvatarURL),
			IsVerified: item.IsVerified,
			CreatedAt:  item.CreatedAt,
		}
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &dto.PaginatedRegistrants{
		Data:       mappedData,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *registrantService) Create(ctx context.Context, ticketID string, state string, userID uint) error {
	ticket, err := s.ticketRepository.GetPublicByID(ctx, ticketID)
	if err != nil {
		s.logger.Error("failed to get public ticket by id", zap.Error(err))
		return fmt.Errorf("failed to get ticket: %w", err)
	}

	if *ticket.Mode {
		if state == "" {
			return fmt.Errorf("state is required for private mode tickets")
		}

		cache, err := s.redis.Get(ticketID + strconv.Itoa(int(userID)))
		if err != nil {
			s.logger.Error("failed to get state from cache", zap.Error(err))
			return fmt.Errorf("failed to verify state: %w", err)
		}

		if cache != state {
			return fmt.Errorf("invalid state provided")
		}
	}

	exist, err := s.redis.Exist(ticketID)
	if err != nil {
		s.logger.Error("failed to check quota existence", zap.Error(err))
		return fmt.Errorf("failed to check quota: %w", err)
	}

	if !exist {
		remaining := ticket.Quota - ticket.RegisteredCount
		if remaining <= 0 {
			return fmt.Errorf("ticket quota has been exhausted")
		}

		if err := s.redis.Set(ticketID, remaining, 12*time.Hour); err != nil {
			s.logger.Error("failed to initialize quota in redis", zap.Error(err))
			return fmt.Errorf("failed to set initial quota: %w", err)
		}
	}

	remaining, err := s.redis.Script(ticketID)
	if err != nil {
		s.logger.Error("failed to update quota", zap.Error(err))
		return fmt.Errorf("failed to update quota: %w", err)
	}
	if remaining <= 0 {
		return fmt.Errorf("ticket quota has been exhausted")
	}

	r := &domain.Registrant{
		ID:            helper.GenerateState(6),
		UserID:        userID,
		TicketID:      ticketID,
		TicketOwnerID: ticket.UserID,
	}

	if err := s.registrantRepository.Create(ctx, r); err != nil {
		s.redis.Incr(ticketID)
		s.logger.Error("failed to create registrant", zap.Error(err))
		return fmt.Errorf("failed to create registration: %w", err)
	}

	if err := s.ticketRepository.UpdateRegisteredCount(ctx, ticketID); err != nil {
		s.redis.Incr(ticketID)
		s.logger.Error("failed to update registered count", zap.Error(err))
		return fmt.Errorf("failed to update ticket count: %w", err)
	}

	if *ticket.Mode {
		if err := s.redis.Delete(ticketID + strconv.Itoa(int(userID))); err != nil {
			s.logger.Warn("failed to cleanup state", zap.Error(err))
		}
	}

	return nil
}

func (s *registrantService) Verify(ctx context.Context, id string, ticketID string, userID uint) error {
	if err := s.registrantRepository.Update(ctx, id, ticketID, userID); err != nil {
		s.logger.Error("failed updating registrant", zap.Error(err))
		return err
	}
	return nil
}

func (s *registrantService) Generate(ctx context.Context, ticketID string, userID uint) (string, error) {
	mode, err := s.ticketRepository.GetModeByID(ctx, ticketID)

	if err != nil {
		return "", err
	}

	if mode {
		state := helper.GenerateState(6)
		if err := s.redis.Set(ticketID+strconv.Itoa(int(userID)), state, time.Minute); err != nil {
			s.logger.Error("failed generating registrant", zap.Error(err))
			return "", err
		}
		return state, nil
	}

	return "", nil
}
