package repo

import (
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"just_for_fun/internal/storage/db/models"
	"strings"
)

var errUser = errors.New("user")

type UserRepo struct {
	log zap.Logger
	db  *gorm.DB
}

func NewUserRepo(log zap.Logger, db *gorm.DB) *UserRepo {
	return &UserRepo{
		log: log,
		db:  db,
	}
}

func (r *UserRepo) Set(user *models.User) (*models.User, error) {
	result := r.db.Save(&user)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), ErrDuplicateKeyStr) {
			return nil, errors.Join(errUser, ErrIsExists)
		}

		r.log.Error("Failed to add user", zap.Error(result.Error))

		return nil, result.Error
	}

	r.log.Info("User added", zap.Int64("telegram_id", user.TgID))

	return user, nil
}

func (r *UserRepo) Update(user *models.User) (*models.User, error) {
	var updatedUser = &models.User{}

	result := r.db.Model(&models.User{}).Where("tg_id = ?", user.TgID).Updates(user)
	if result.Error != nil {
		r.log.Error("Failed to update user", zap.Error(result.Error))
		return nil, result.Error
	}

	result.Scan(updatedUser)

	r.log.Info("User updated", zap.Any("updated", updatedUser))

	return updatedUser, nil
}
