package repository

import (
	"time"

	"github.com/leoferamos/aroma-sense/internal/model"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(order *model.Order) error
	FindByID(id uint) (*model.Order, error)
	FindByUserID(userID string) ([]model.Order, error)
	FindByPublicIDWithItems(publicID string) (*model.Order, error)
	ListOrders(status *string, startDate *time.Time, endDate *time.Time, page int, perPage int) ([]model.Order, int64, float64, error)
	HasUserDeliveredOrderWithProduct(userID string, productID uint) (bool, error)
	UpdateStatusByPublicID(publicID string, status model.OrderStatus) error
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(order *model.Order) error {
	return r.db.Create(order).Error
}

func (r *orderRepository) FindByID(id uint) (*model.Order, error) {
	var order model.Order
	err := r.db.Preload("Items").Preload("Items.Product").First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) FindByUserID(userID string) ([]model.Order, error) {
	var orders []model.Order
	err := r.db.Preload("Items").Preload("Items.Product").Where("user_id = ?", userID).Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// FindByPublicIDWithItems retrieves an order by public_id including items.
func (r *orderRepository) FindByPublicIDWithItems(publicID string) (*model.Order, error) {
	var order model.Order
	err := r.db.Preload("Items").Where("public_id = ?", publicID).First(&order).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}

// ListOrders implements OrderRepository.ListOrders
func (r *orderRepository) ListOrders(status *string, startDate *time.Time, endDate *time.Time, page int, perPage int) ([]model.Order, int64, float64, error) {
	var orders []model.Order

	// Base query for count and list
	base := r.db.Model(&model.Order{})
	if status != nil && *status != "" {
		base = base.Where("status = ?", *status)
	}
	if startDate != nil {
		base = base.Where("created_at >= ?", *startDate)
	}
	if endDate != nil {
		base = base.Where("created_at <= ?", *endDate)
	}

	// Total count
	var totalCount int64
	if err := base.Count(&totalCount).Error; err != nil {
		return nil, 0, 0, err
	}

	// Total revenue
	revQuery := r.db.Model(&model.Order{})
	if status != nil && *status != "" {
		revQuery = revQuery.Where("status = ?", *status)
	}
	if startDate != nil {
		revQuery = revQuery.Where("created_at >= ?", *startDate)
	}
	if endDate != nil {
		revQuery = revQuery.Where("created_at <= ?", *endDate)
	}
	var totalRevenue float64
	if err := revQuery.Select("COALESCE(SUM(total_amount),0)").Scan(&totalRevenue).Error; err != nil {
		return nil, 0, 0, err
	}

	// Pagination defaults
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 25
	}
	offset := (page - 1) * perPage

	// Fetch records with preloads
	listQuery := r.db.Preload("Items").Preload("Items.Product").Model(&model.Order{})
	if status != nil && *status != "" {
		listQuery = listQuery.Where("status = ?", *status)
	}
	if startDate != nil {
		listQuery = listQuery.Where("created_at >= ?", *startDate)
	}
	if endDate != nil {
		listQuery = listQuery.Where("created_at <= ?", *endDate)
	}
	if err := listQuery.Order("created_at desc").Offset(offset).Limit(perPage).Find(&orders).Error; err != nil {
		return nil, 0, 0, err
	}

	return orders, totalCount, totalRevenue, nil
}

// HasUserDeliveredOrderWithProduct returns true if the user has at least one delivered order containing the given product.
func (r *orderRepository) HasUserDeliveredOrderWithProduct(userID string, productID uint) (bool, error) {
	var exists bool
	raw := `
		SELECT EXISTS(
			SELECT 1
			FROM orders o
			JOIN order_items oi ON oi.order_id = o.id
			WHERE o.user_id = ?
			  AND o.status = 'delivered'
			  AND oi.product_id = ?
		)`
	if err := r.db.Raw(raw, userID, productID).Scan(&exists).Error; err != nil {
		return false, err
	}
	return exists, nil
}

// UpdateStatusByPublicID updates the status of an order identified by its public_id.
func (r *orderRepository) UpdateStatusByPublicID(publicID string, status model.OrderStatus) error {
	return r.db.Model(&model.Order{}).
		Where("public_id = ?", publicID).
		Update("status", status).Error
}
