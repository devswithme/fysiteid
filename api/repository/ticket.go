package repository

import (
	"context"
	"fybe/model/domain"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TicketRepository interface {
	Get(ctx context.Context, userID uint) ([]domain.Ticket, error)
	GetPublicByUsername(ctx context.Context, username string) ([]domain.Ticket, error)
	GetPublicByID(ctx context.Context, ticketID string) (*domain.Ticket, error)
	GetByID(ctx context.Context, ticketID string, userID uint) (*domain.Ticket, error)
	GetOwnerIDByID(ctx context.Context, ticketID string) (uint, error)
	GetPathByID(ctx context.Context, ticketID string) (string, error)
	GetModeByID(ctx context.Context, ticketID string) (bool, error)
	Create(ctx context.Context, ticket *domain.Ticket) error
	Update(ctx context.Context, ticket *domain.Ticket, ticketID string, userID uint) (uint, error)
	UpdateRegisteredCount(ctx context.Context, ticketID string) error
	Delete(ctx context.Context, ticketID string, userID uint) error
}

type ticketRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewTicketRepository(db *gorm.DB, logger *zap.Logger) TicketRepository {
	return &ticketRepository{
		db:     db,
		logger: logger,
	}
}

func (r *ticketRepository) Get(ctx context.Context, userID uint) ([]domain.Ticket, error) {
	var tickets []domain.Ticket
	if err := r.db.WithContext(ctx).Find(&tickets, "user_id = ?", userID).Error; err != nil {
		r.logger.Error("failed to get ticket", zap.Uint("userID", userID), zap.Error(err))
		return nil, err
	}

	return tickets, nil
}

func (r *ticketRepository) GetPublicByUsername(ctx context.Context, username string) ([]domain.Ticket, error) {
	var tickets []domain.Ticket

	if err := r.db.WithContext(ctx).Joins("JOIN users ON users.id = tickets.user_id").Where("users.username = ?", username).Find(&tickets).Error; err != nil {
		r.logger.Error("failed to get public by username", zap.String("username", username), zap.Error(err))
		return nil, err
	}

	return tickets, nil
}

func (r *ticketRepository) GetPublicByID(ctx context.Context, ticketID string) (*domain.Ticket, error) {
	ticket := &domain.Ticket{}
	if err := r.db.WithContext(ctx).First(ticket, "id = ?", ticketID).Error; err != nil {
		r.logger.Error("failed to get public by ticket id", zap.String("ticketID", ticket.ID), zap.Error(err))
		return nil, err
	}

	return ticket, nil
}

func (r *ticketRepository) GetByID(ctx context.Context, ticketID string, userID uint) (*domain.Ticket, error) {
	ticket := &domain.Ticket{}
	if err := r.db.WithContext(ctx).First(ticket, "id = ? AND user_id = ?", ticketID, userID).Error; err != nil {
		r.logger.Error("failed to get by ticket id", zap.String("ticketID", ticket.ID), zap.Uint("userID", userID), zap.Error(err))
		return nil, err
	}

	return ticket, nil
}

func (r *ticketRepository) GetOwnerIDByID(ctx context.Context, ticketID string) (uint, error) {
	ticket := &domain.Ticket{}
	if err := r.db.WithContext(ctx).First(ticket, "id = ?", ticketID).Select("user_id").Error; err != nil {
		r.logger.Error("failed to get user id by ticket id", zap.String("ticketID", ticket.ID), zap.Error(err))
		return 0, err
	}

	return ticket.UserID, nil
}

func (r *ticketRepository) GetPathByID(ctx context.Context, ticketID string) (string, error) {
	ticket := &domain.Ticket{}
	if err := r.db.WithContext(ctx).First(ticket, "id = ?", ticketID).Select("image_url").Error; err != nil {
		r.logger.Error("failed to get path by ticket id", zap.String("ticketID", ticket.ID), zap.Error(err))
		return "", err
	}

	return *ticket.ImageURL, nil
}

func (r *ticketRepository) GetModeByID(ctx context.Context, ticketID string) (bool, error) {
	ticket := &domain.Ticket{}
	if err := r.db.WithContext(ctx).First(ticket, "id = ?", ticketID).Select("mode").Error; err != nil {
		r.logger.Error("failed to get path by ticket id", zap.String("ticketID", ticket.ID), zap.Error(err))
		return false, err
	}

	return *ticket.Mode, nil
}

func (r *ticketRepository) Create(ctx context.Context, ticket *domain.Ticket) error {
	if err := r.db.WithContext(ctx).Create(&ticket).Error; err != nil {
		r.logger.Error("failed to create ticket", zap.String("ticketID", ticket.ID), zap.Error(err))
		return err
	}

	return nil
}

func (r *ticketRepository) UpdateRegisteredCount(ctx context.Context, ticketID string) error {
	if err := r.db.WithContext(ctx).Model(&domain.Ticket{}).Where("id = ?", ticketID).Update("registered_count", gorm.Expr("registered_count + 1")).Error; err != nil {
		r.logger.Error("failed to update registered count", zap.String("ticketID", ticketID), zap.Error(err))
		return err
	}

	return nil
}

func (r *ticketRepository) Update(ctx context.Context, ticket *domain.Ticket, ticketID string, userID uint) (uint, error) {
	updated := &domain.Ticket{}
	if err := r.db.WithContext(ctx).Clauses(clause.Returning{Columns: []clause.Column{{Name: "registered_count"}}}).Model(&domain.Ticket{}).Where("id = ? AND user_id = ?", ticketID, userID).Updates(ticket).Scan(updated).Error; err != nil {
		r.logger.Error("failed to update ticket", zap.String("ticketID", ticketID), zap.Uint("userID", userID), zap.Error(err))
		return 0, err
	}

	return updated.RegisteredCount, nil
}

func (r *ticketRepository) Delete(ctx context.Context, ticketID string, userID uint) error {
	if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", ticketID, userID).Delete(&domain.Ticket{}).Unscoped().Error; err != nil {
		r.logger.Error("failed to delete ticket", zap.String("ticketID", ticketID), zap.Uint("userID", userID), zap.Error(err))
		return err
	}

	return nil
}
