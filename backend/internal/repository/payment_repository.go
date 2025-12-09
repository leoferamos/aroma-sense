package repository

import (
	"github.com/leoferamos/aroma-sense/internal/model"
	"gorm.io/gorm"
)

// PaymentRepository manages payment records for reconciliation.
type PaymentRepository interface {
	Create(payment *model.Payment) error
	FindByIntentID(intentID string) (*model.Payment, error)
	UpdateStatusByIntentID(intentID string, status model.PaymentStatus, errorCode, errorMessage string) error
	AttachOrderPublicID(intentID string, orderPublicID string) error
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

// Create inserts a new payment record into the database.
func (r *paymentRepository) Create(payment *model.Payment) error {
	return r.db.Create(payment).Error
}

// FindByIntentID retrieves a payment by its intent ID.
func (r *paymentRepository) FindByIntentID(intentID string) (*model.Payment, error) {
	var p model.Payment
	if err := r.db.Where("intent_id = ?", intentID).First(&p).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

// UpdateStatusByIntentID updates the status and error details of a payment by its intent ID.
func (r *paymentRepository) UpdateStatusByIntentID(intentID string, status model.PaymentStatus, errorCode, errorMessage string) error {
	updates := map[string]interface{}{
		"status":        status,
		"error_code":    errorCode,
		"error_message": errorMessage,
	}
	return r.db.Model(&model.Payment{}).Where("intent_id = ?", intentID).Updates(updates).Error
}

// AttachOrderPublicID links a payment to an order by updating the order_public_id field.
func (r *paymentRepository) AttachOrderPublicID(intentID string, orderPublicID string) error {
	return r.db.Model(&model.Payment{}).Where("intent_id = ?", intentID).Update("order_public_id", orderPublicID).Error
}
