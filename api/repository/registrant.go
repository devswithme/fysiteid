package repository

import (
	"context"
	"fybe/model/domain"
	"fybe/model/dto"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RegistrantRepository interface {
	GetByTicketID(ctx context.Context, ticketID string, userID uint, page int, limit int, search string) ([]dto.RegistrantGetByTicketID, int64, error)
	GetByUserID(ctx context.Context, userID uint) ([]dto.RegistrantGetByUserID, error)
	Create(ctx context.Context, registrant *domain.Registrant) error
	Update(ctx context.Context, id string, ticketID string, userID uint) error
}

type registrantRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewRegistrantRepository(db *gorm.DB, logger *zap.Logger) RegistrantRepository {
	return &registrantRepository{
		db:     db,
		logger: logger,
	}
}

func (r *registrantRepository) GetByTicketID(ctx context.Context, ticketID string, userID uint, page int, limit int, search string) ([]dto.RegistrantGetByTicketID, int64, error) {
	var registrants []dto.RegistrantGetByTicketID
	var total int64

	offset := (page - 1) * limit

	query := r.db.WithContext(ctx).Table("registrants").
		Joins("JOIN users ON users.id = registrants.user_id").
		Where("registrants.ticket_id = ? AND registrants.ticket_owner_id = ?", ticketID, userID)

	if search != "" {
		query = query.Where("users.username ILIKE ?", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		r.logger.Error("failed to get total registrants",
			zap.String("ticketID", ticketID),
			zap.Uint("userID", userID),
			zap.String("search", search),
			zap.Error(err))
		return nil, 0, err
	}

	if err := query.Select("users.username, users.avatar_url, registrants.is_verified, registrants.created_at").
		Offset(offset).
		Limit(limit).
		Scan(&registrants).Error; err != nil {
		r.logger.Error("failed to get registrants by ticket id",
			zap.String("ticketID", ticketID),
			zap.Uint("userID", userID),
			zap.Int("page", page),
			zap.Int("limit", limit),
			zap.String("search", search),
			zap.Error(err))
		return nil, 0, err
	}
	return registrants, total, nil
}

func (r *registrantRepository) GetByUserID(ctx context.Context, userID uint) ([]dto.RegistrantGetByUserID, error) {
	var registrants []dto.RegistrantGetByUserID
	if err := r.db.WithContext(ctx).Table("registrants").
		Select("tickets.title, tickets.id, tickets.image_url, registrants.id, registrants.ticket_id, registrants.is_verified, registrants.created_at").
		Joins("JOIN tickets ON tickets.id = registrants.ticket_id").Where("registrants.user_id = ?", userID).Scan(&registrants).Error; err != nil {
		r.logger.Error("failed to get registrants by user id", zap.Uint("userID", userID), zap.Error(err))
		return nil, err
	}

	return registrants, nil
}

func (r *registrantRepository) Create(ctx context.Context, registrant *domain.Registrant) error {
	if err := r.db.WithContext(ctx).Create(&registrant).Error; err != nil {
		r.logger.Error("failed to create registrant", zap.String("regID", registrant.ID), zap.Error(err))
		return err
	}

	return nil
}

func (r *registrantRepository) Update(ctx context.Context, id string, ticketID string, userID uint) error {
	if err := r.db.WithContext(ctx).Table("registrants").Where("id = ? AND ticket_id = ? AND ticket_owner_id = ?", id, ticketID, userID).Update("is_verified", true).Error; err != nil {
		r.logger.Error("failed to update registrant", zap.String("id", id), zap.Uint("userID", userID), zap.Error(err))
		return err
	}

	return nil
}
