package repository

import (
	"context"
	"fybe/model/domain"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetById(ctx context.Context, userID uint) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	GetByTicketID(ctx context.Context, ticketID string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User, userID uint) error
}

type userRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewUserRepository(db *gorm.DB, logger *zap.Logger) UserRepository {
	return &userRepository{
		db:     db,
		logger: logger,
	}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	if err := r.db.WithContext(ctx).Where("email = ?", user.Email).FirstOrCreate(&user).Error; err != nil {
		r.logger.Error("failed to create user", zap.Uint("userID", user.ID), zap.Error(err))
		return err
	}

	return nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	user := &domain.User{}
	if err := r.db.WithContext(ctx).First(user, "username = ?", username).Error; err != nil {
		r.logger.Error("failed to get by user username", zap.String("username", username), zap.Error(err))
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetByTicketID(ctx context.Context, ticketID string) (*domain.User, error) {
	user := &domain.User{}
	if err := r.db.WithContext(ctx).Model(&domain.User{}).Joins("JOIN tickets ON tickets.user_id = users.id").Where("tickets.id = ?", ticketID).First(user).Error; err != nil {
		r.logger.Error("failed to get by user ticket id", zap.String("ticketID", ticketID), zap.Error(err))
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetById(ctx context.Context, userID uint) (*domain.User, error) {
	user := &domain.User{}
	if err := r.db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", userID).First(user).Error; err != nil {
		r.logger.Error("failed to get by user id", zap.Uint("userID", userID), zap.Error(err))
		return nil, err
	}

	return user, nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User, userID uint) error {
	if err := r.db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", userID).Updates(user).Error; err != nil {
		r.logger.Error("failed to update user", zap.Uint("userID", userID), zap.Error(err))
		return err
	}

	return nil
}
