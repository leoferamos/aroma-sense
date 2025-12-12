package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/leoferamos/aroma-sense/internal/dto"
	"github.com/leoferamos/aroma-sense/internal/model"
	"github.com/stretchr/testify/assert"
)

type mockOrderRepo struct {
	createErr        error
	listOrders       []model.Order
	listCount        int64
	listRevenue      float64
	listErr          error
	findByUserErr    error
	findByUserOrders []model.Order
}

func (m *mockOrderRepo) Create(order *model.Order) error {
	return m.createErr
}

func (m *mockOrderRepo) FindByID(id uint) (*model.Order, error) {
	return nil, nil
}

func (m *mockOrderRepo) FindByUserID(userID string) ([]model.Order, error) {
	return m.findByUserOrders, m.findByUserErr
}

func (m *mockOrderRepo) FindByPublicIDWithItems(publicID string) (*model.Order, error) {
	return nil, nil
}

func (m *mockOrderRepo) ListOrders(status *string, startDate *time.Time, endDate *time.Time, page int, perPage int) ([]model.Order, int64, float64, error) {
	return m.listOrders, m.listCount, m.listRevenue, m.listErr
}

func (m *mockOrderRepo) HasUserDeliveredOrderWithProduct(userID string, productID uint) (bool, error) {
	return false, nil
}

func (m *mockOrderRepo) UpdateStatusByPublicID(publicID string, status model.OrderStatus) error {
	return nil
}

type mockCartRepo struct {
	findByUserCart *model.Cart
	findByUserErr  error
	clearErr       error
}

func (m *mockCartRepo) Create(cart *model.Cart) error {
	return nil
}

func (m *mockCartRepo) FindByUserID(userID string) (*model.Cart, error) {
	return m.findByUserCart, m.findByUserErr
}

func (m *mockCartRepo) Update(cart *model.Cart) error {
	return nil
}

func (m *mockCartRepo) Delete(id uint) error {
	return nil
}

func (m *mockCartRepo) CreateCartItem(item *model.CartItem) error {
	return nil
}

func (m *mockCartRepo) UpdateCartItem(item *model.CartItem) error {
	return nil
}

func (m *mockCartRepo) FindCartItemByID(itemID uint) (*model.CartItem, error) {
	return nil, nil
}

func (m *mockCartRepo) DeleteCartItem(itemID uint) error {
	return nil
}

func (m *mockCartRepo) ClearCartItems(cartID uint) error {
	return m.clearErr
}

type mockProductRepo struct {
	findByIDProduct model.Product
	findByIDErr     error
}

func (m *mockProductRepo) Create(input dto.ProductFormDTO, imageURL string, thumbnailURL string) (uint, error) {
	return 0, nil
}

func (m *mockProductRepo) FindAll(limit int) ([]model.Product, error) {
	return nil, nil
}

func (m *mockProductRepo) FindAllPaginated(limit int, offset int) ([]model.Product, int, error) {
	return nil, 0, nil
}

func (m *mockProductRepo) FindByID(id uint) (model.Product, error) {
	return m.findByIDProduct, m.findByIDErr
}

func (m *mockProductRepo) FindBySlug(slug string) (model.Product, error) {
	return model.Product{}, nil
}

func (m *mockProductRepo) SearchProducts(ctx context.Context, query string, limit int, offset int, sort string) ([]model.Product, int, error) {
	return nil, 0, nil
}

func (m *mockProductRepo) SearchProductsByGender(ctx context.Context, query string, limit int, offset int, sort string, gender string) ([]model.Product, int, error) {
	return nil, 0, nil
}

func (m *mockProductRepo) Update(product *model.Product) error {
	return nil
}

func (m *mockProductRepo) Delete(id uint) error {
	return nil
}

func (m *mockProductRepo) DecrementStock(productID uint, quantity int) error {
	return nil
}

func (m *mockProductRepo) EnsureUniqueSlug(base string) (string, error) {
	return "", nil
}

func (m *mockProductRepo) UpsertProductEmbedding(productID uint, embedding []float32) error {
	return nil
}

func (m *mockProductRepo) HasProductEmbedding(productID uint) (bool, error) {
	return false, nil
}

func (m *mockProductRepo) FindSimilarProductsByEmbedding(ctx context.Context, embedding []float32, limit int) ([]model.Product, error) {
	return nil, nil
}

func (m *mockProductRepo) FindSimilarProductsByEmbeddingAndGender(ctx context.Context, embedding []float32, limit int, gender string) ([]model.Product, error) {
	return nil, nil
}

type mockShippingSvc struct {
	calculateOptions []dto.ShippingOption
	calculateErr     error
}

func (m *mockShippingSvc) CalculateOptions(ctx context.Context, userID string, cep string) ([]dto.ShippingOption, error) {
	return m.calculateOptions, m.calculateErr
}

// --- Test helpers ---
func createTestCart() *model.Cart {
	return &model.Cart{
		ID:     1,
		UserID: "user123",
		Items: []model.CartItem{
			{
				ID:        1,
				CartID:    1,
				ProductID: 1,
				Quantity:  2,
			},
		},
	}
}

func createTestProduct() model.Product {
	return model.Product{
		ID:            1,
		Name:          "Test Product",
		Slug:          "test-product",
		Price:         10.0,
		StockQuantity: 10,
		ImageURL:      "http://example.com/image.jpg",
	}
}

func createTestOrder() model.Order {
	now := time.Now()
	return model.Order{
		ID:                        1,
		PublicID:                  "order123",
		UserID:                    "user123",
		TotalAmount:               20.0,
		Status:                    model.OrderStatusPending,
		ShippingAddress:           "Test Address",
		PaymentMethod:             model.PaymentMethodCreditCard,
		ShippingPrice:             5.0,
		ShippingCarrier:           "Test Carrier",
		ShippingServiceCode:       "standard",
		ShippingEstimatedDelivery: &now,
		Items: []model.OrderItem{
			{
				ProductID:       1,
				ProductSlug:     "test-product",
				ProductName:     "Test Product",
				ProductImageURL: "http://example.com/image.jpg",
				Quantity:        2,
				PriceAtPurchase: 10.0,
				Subtotal:        20.0,
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// --- Tests ---
func TestCreateOrderFromCart(t *testing.T) {
	cart := createTestCart()
	product := createTestProduct()
	shippingOptions := []dto.ShippingOption{
		{
			Carrier:       "Test Carrier",
			ServiceCode:   "standard",
			Price:         5.0,
			EstimatedDays: 3,
		},
	}

	t.Run("success", func(t *testing.T) {
		svc := NewOrderService(
			&mockOrderRepo{},
			&mockCartRepo{findByUserCart: cart},
			&mockProductRepo{findByIDProduct: product},
			&mockShippingSvc{calculateOptions: shippingOptions},
		)

		req := &dto.CreateOrderFromCartRequest{
			ShippingAddress: "12345-000",
			PaymentMethod:   string(model.PaymentMethodCreditCard),
			ShippingSelection: &dto.ShippingSelection{
				Carrier:     "Test Carrier",
				ServiceCode: "standard",
			},
		}

		resp, err := svc.CreateOrderFromCart("user123", req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, 25.0, resp.TotalAmount) // 20 + 5 shipping
		assert.Equal(t, "Test Carrier", resp.ShippingCarrier)
	})

	t.Run("empty cart", func(t *testing.T) {
		svc := NewOrderService(
			&mockOrderRepo{},
			&mockCartRepo{findByUserErr: errors.New("not found")},
			&mockProductRepo{},
			&mockShippingSvc{},
		)

		req := &dto.CreateOrderFromCartRequest{
			ShippingAddress: "12345-000",
			PaymentMethod:   string(model.PaymentMethodCreditCard),
		}

		resp, err := svc.CreateOrderFromCart("user123", req)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("insufficient stock", func(t *testing.T) {
		productLowStock := createTestProduct()
		productLowStock.StockQuantity = 1 // Less than cart quantity of 2

		svc := NewOrderService(
			&mockOrderRepo{},
			&mockCartRepo{findByUserCart: cart},
			&mockProductRepo{findByIDProduct: productLowStock},
			&mockShippingSvc{},
		)

		req := &dto.CreateOrderFromCartRequest{
			ShippingAddress: "12345-000",
			PaymentMethod:   string(model.PaymentMethodCreditCard),
		}

		resp, err := svc.CreateOrderFromCart("user123", req)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("invalid shipping selection", func(t *testing.T) {
		svc := NewOrderService(
			&mockOrderRepo{},
			&mockCartRepo{findByUserCart: cart},
			&mockProductRepo{findByIDProduct: product},
			&mockShippingSvc{calculateOptions: shippingOptions},
		)

		req := &dto.CreateOrderFromCartRequest{
			ShippingAddress: "12345-000",
			PaymentMethod:   string(model.PaymentMethodCreditCard),
			ShippingSelection: &dto.ShippingSelection{
				Carrier:     "Invalid Carrier",
				ServiceCode: "invalid",
			},
		}

		resp, err := svc.CreateOrderFromCart("user123", req)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("create order error", func(t *testing.T) {
		svc := NewOrderService(
			&mockOrderRepo{createErr: errors.New("db error")},
			&mockCartRepo{findByUserCart: cart},
			&mockProductRepo{findByIDProduct: product},
			&mockShippingSvc{},
		)

		req := &dto.CreateOrderFromCartRequest{
			ShippingAddress: "12345-000",
			PaymentMethod:   string(model.PaymentMethodCreditCard),
		}

		resp, err := svc.CreateOrderFromCart("user123", req)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("clear cart error", func(t *testing.T) {
		svc := NewOrderService(
			&mockOrderRepo{},
			&mockCartRepo{findByUserCart: cart, clearErr: errors.New("clear error")},
			&mockProductRepo{findByIDProduct: product},
			&mockShippingSvc{},
		)

		req := &dto.CreateOrderFromCartRequest{
			ShippingAddress: "12345-000",
			PaymentMethod:   string(model.PaymentMethodCreditCard),
		}

		resp, err := svc.CreateOrderFromCart("user123", req)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestAdminListOrders(t *testing.T) {
	orders := []model.Order{createTestOrder()}

	t.Run("success", func(t *testing.T) {
		svc := NewOrderService(
			&mockOrderRepo{listOrders: orders, listCount: 1, listRevenue: 20.0},
			&mockCartRepo{},
			&mockProductRepo{},
			&mockShippingSvc{},
		)

		resp, err := svc.AdminListOrders(nil, nil, nil, 1, 10)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.Orders, 1)
		assert.Equal(t, 1, resp.Meta.Pagination.TotalCount)
		assert.Equal(t, 20.0, resp.Meta.Stats.TotalRevenue)
		assert.Equal(t, 20.0, resp.Meta.Stats.AverageOrderValue)
	})

	t.Run("error", func(t *testing.T) {
		svc := NewOrderService(
			&mockOrderRepo{listErr: errors.New("db error")},
			&mockCartRepo{},
			&mockProductRepo{},
			&mockShippingSvc{},
		)

		resp, err := svc.AdminListOrders(nil, nil, nil, 1, 10)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestGetOrdersByUser(t *testing.T) {
	orders := []model.Order{createTestOrder()}

	t.Run("success", func(t *testing.T) {
		svc := NewOrderService(
			&mockOrderRepo{findByUserOrders: orders},
			&mockCartRepo{},
			&mockProductRepo{},
			&mockShippingSvc{},
		)

		resp, err := svc.GetOrdersByUser("user123")
		assert.NoError(t, err)
		assert.Len(t, resp, 1)
		assert.Equal(t, "order123", resp[0].PublicID)
	})

	t.Run("error", func(t *testing.T) {
		svc := NewOrderService(
			&mockOrderRepo{findByUserErr: errors.New("db error")},
			&mockCartRepo{},
			&mockProductRepo{},
			&mockShippingSvc{},
		)

		resp, err := svc.GetOrdersByUser("user123")
		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("empty result", func(t *testing.T) {
		svc := NewOrderService(
			&mockOrderRepo{findByUserOrders: []model.Order{}},
			&mockCartRepo{},
			&mockProductRepo{},
			&mockShippingSvc{},
		)

		resp, err := svc.GetOrdersByUser("user123")
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp, 0)
	})
}
