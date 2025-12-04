package repository

import (
	"github.com/leoferamos/aroma-sense/internal/model"
	"gorm.io/gorm"
)

type UserContestationRepository interface {
	Create(contestation *model.UserContestation) error
	FindByID(id uint) (*model.UserContestation, error)
	FindByUserID(userID uint) ([]model.UserContestation, error)
	ListPending(limit, offset int) ([]model.UserContestation, int64, error)
	Update(contestation *model.UserContestation) error
}

type userContestationRepository struct {
	db *gorm.DB
}

func NewUserContestationRepository(db *gorm.DB) UserContestationRepository {
	return &userContestationRepository{db: db}
}

// Create inserts a new user contestation into the database
func (r *userContestationRepository) Create(contestation *model.UserContestation) error {
	return r.db.Create(contestation).Error
}

// FindByID retrieves a user contestation by its ID
func (r *userContestationRepository) FindByID(id uint) (*model.UserContestation, error) {
	var c model.UserContestation
	err := r.db.First(&c, id).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// FindByUserID retrieves all contestations for a specific user
func (r *userContestationRepository) FindByUserID(userID uint) ([]model.UserContestation, error) {
	var cs []model.UserContestation
	err := r.db.Where("user_id = ?", userID).Order("requested_at desc").Find(&cs).Error
	return cs, err
}

// ListPending retrieves pending user contestations with pagination
func (r *userContestationRepository) ListPending(limit, offset int) ([]model.UserContestation, int64, error) {
	var cs []model.UserContestation
	var count int64

	db := r.db.Model(&model.UserContestation{}).Where("status = ?", "pending")
	db.Count(&count)
	err := db.Order("requested_at asc").Limit(limit).Offset(offset).Find(&cs).Error
	return cs, count, err
}

// Update modifies an existing user contestation record
func (r *userContestationRepository) Update(contestation *model.UserContestation) error {
	return r.db.Save(contestation).Error
}
