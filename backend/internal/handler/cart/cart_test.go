package cart_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/apperror"
	"github.com/leoferamos/aroma-sense/internal/dto"
	"github.com/leoferamos/aroma-sense/internal/handler/cart"
	"github.com/leoferamos/aroma-sense/internal/model"
	cartservice "github.com/leoferamos/aroma-sense/internal/service/cart"
	"github.com/stretchr/testify/assert"
)

type mockCartService struct {
	createCartForUserResult        error
	getCartByUserIDResult          *model.Cart
	getCartByUserIDErr             error
	getCartResponseResult          *dto.CartResponse
	getCartResponseErr             error
	addItemToCartResult            *dto.CartResponse
	addItemToCartErr               error
	updateItemQuantityResult       *dto.CartResponse
	updateItemQuantityErr          error
	updateItemQuantityBySlugResult *dto.CartResponse
	updateItemQuantityBySlugErr    error
	removeItemResult               *dto.CartResponse
	removeItemErr                  error
	removeItemBySlugResult         *dto.CartResponse
	removeItemBySlugErr            error
	clearCartResult                *dto.CartResponse
	clearCartErr                   error
}

func (m *mockCartService) CreateCartForUser(userID string) error {
	return m.createCartForUserResult
}

func (m *mockCartService) GetCartByUserID(userID string) (*model.Cart, error) {
	return m.getCartByUserIDResult, m.getCartByUserIDErr
}

func (m *mockCartService) GetCartResponse(userID string) (*dto.CartResponse, error) {
	return m.getCartResponseResult, m.getCartResponseErr
}

func (m *mockCartService) AddItemToCart(userID string, productSlug string, quantity int) (*dto.CartResponse, error) {
	return m.addItemToCartResult, m.addItemToCartErr
}

func (m *mockCartService) UpdateItemQuantity(userID string, itemID uint, quantity int) (*dto.CartResponse, error) {
	return m.updateItemQuantityResult, m.updateItemQuantityErr
}

func (m *mockCartService) UpdateItemQuantityBySlug(userID string, productSlug string, quantity int) (*dto.CartResponse, error) {
	return m.updateItemQuantityBySlugResult, m.updateItemQuantityBySlugErr
}

func (m *mockCartService) RemoveItem(userID string, itemID uint) (*dto.CartResponse, error) {
	return m.removeItemResult, m.removeItemErr
}

func (m *mockCartService) RemoveItemBySlug(userID string, productSlug string) (*dto.CartResponse, error) {
	return m.removeItemBySlugResult, m.removeItemBySlugErr
}

func (m *mockCartService) ClearCart(userID string) (*dto.CartResponse, error) {
	return m.clearCartResult, m.clearCartErr
}

func setupCartRouter(svc cartservice.CartService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Add middleware to simulate authentication
	r.Use(func(c *gin.Context) {
		c.Set("userID", "user-123")
		c.Next()
	})

	handler := cart.NewCartHandler(svc)
	r.GET("/cart", handler.GetCart)
	r.POST("/cart", handler.AddItem)
	r.PATCH("/cart/items/:productSlug", handler.UpdateItemQuantity)
	r.DELETE("/cart/items/:productSlug", handler.RemoveItem)
	r.DELETE("/cart", handler.ClearCart)
	return r
}

func createTestCartResponse() *dto.CartResponse {
	return &dto.CartResponse{
		Items: []dto.CartItemResponse{
			{
				Product: &dto.ProductResponse{
					Slug:  "test-product",
					Name:  "Test Product",
					Price: 29.99,
				},
				Quantity: 2,
				Price:    29.99,
				Total:    59.98,
			},
		},
		Total:     59.98,
		ItemCount: 1,
	}
}

// --- Tests ---
func TestCartHandler_GetCart(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		cartResponse := createTestCartResponse()
		svc := &mockCartService{
			getCartResponseResult: cartResponse,
		}
		r := setupCartRouter(svc)

		req, _ := http.NewRequest("GET", "/cart", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dto.CartResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Items, 1)
		assert.Equal(t, 59.98, response.Total)
		assert.Equal(t, 1, response.ItemCount)
	})

	t.Run("cart not found", func(t *testing.T) {
		svc := &mockCartService{
			getCartResponseErr: errors.New("cart not found"),
		}
		r := setupCartRouter(svc)

		req, _ := http.NewRequest("GET", "/cart", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "cart not found", response.Error)
	})

	t.Run("unauthenticated", func(t *testing.T) {
		svc := &mockCartService{}
		gin.SetMode(gin.TestMode)
		r := gin.New()
		handler := cart.NewCartHandler(svc)
		r.GET("/cart", handler.GetCart)

		req, _ := http.NewRequest("GET", "/cart", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "unauthenticated", response.Error)
	})
}

func TestCartHandler_AddItem(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		cartResponse := createTestCartResponse()
		svc := &mockCartService{
			addItemToCartResult: cartResponse,
		}
		r := setupCartRouter(svc)

		reqBody := `{
			"product_slug": "test-product",
			"quantity": 2
		}`
		req, _ := http.NewRequest("POST", "/cart", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dto.CartResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Items, 1)
	})

	t.Run("invalid json", func(t *testing.T) {
		svc := &mockCartService{}
		r := setupCartRouter(svc)

		reqBody := `{"product_slug": "", "quantity": 0}`
		req, _ := http.NewRequest("POST", "/cart", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid_request", response.Error)
	})

	t.Run("product not found", func(t *testing.T) {
		svc := &mockCartService{
			addItemToCartErr: apperror.NewCodeMessage("product_not_found", "product not found"),
		}
		r := setupCartRouter(svc)

		reqBody := `{
			"product_slug": "nonexistent-product",
			"quantity": 1
		}`
		req, _ := http.NewRequest("POST", "/cart", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "product_not_found", response.Error)
	})

	t.Run("insufficient stock", func(t *testing.T) {
		svc := &mockCartService{
			addItemToCartErr: apperror.NewCodeMessage("insufficient_stock", "insufficient stock"),
		}
		r := setupCartRouter(svc)

		reqBody := `{
			"product_slug": "test-product",
			"quantity": 100
		}`
		req, _ := http.NewRequest("POST", "/cart", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "insufficient_stock", response.Error)
	})

	t.Run("service error", func(t *testing.T) {
		svc := &mockCartService{
			addItemToCartErr: errors.New("database error"),
		}
		r := setupCartRouter(svc)

		reqBody := `{
			"product_slug": "test-product",
			"quantity": 1
		}`
		req, _ := http.NewRequest("POST", "/cart", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "internal_error", response.Error)
	})
}

func TestCartHandler_UpdateItemQuantity(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		cartResponse := createTestCartResponse()
		cartResponse.Items[0].Quantity = 5
		svc := &mockCartService{
			updateItemQuantityBySlugResult: cartResponse,
		}
		r := setupCartRouter(svc)

		reqBody := `{"quantity": 5}`
		req, _ := http.NewRequest("PATCH", "/cart/items/test-product", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dto.CartResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 5, response.Items[0].Quantity)
	})

	t.Run("invalid product slug", func(t *testing.T) {
		svc := &mockCartService{}
		r := setupCartRouter(svc)

		reqBody := `{"quantity": 3}`
		req, _ := http.NewRequest("PATCH", "/cart/items/", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("invalid json", func(t *testing.T) {
		svc := &mockCartService{}
		r := setupCartRouter(svc)

		reqBody := `{"quantity": -1}`
		req, _ := http.NewRequest("PATCH", "/cart/items/test-product", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid_request", response.Error)
	})

	t.Run("cart item not found", func(t *testing.T) {
		svc := &mockCartService{
			updateItemQuantityBySlugErr: apperror.NewCodeMessage("cart_item_not_found", "cart item not found"),
		}
		r := setupCartRouter(svc)

		reqBody := `{"quantity": 3}`
		req, _ := http.NewRequest("PATCH", "/cart/items/nonexistent-product", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "cart_item_not_found", response.Error)
	})

	t.Run("insufficient stock", func(t *testing.T) {
		svc := &mockCartService{
			updateItemQuantityBySlugErr: apperror.NewCodeMessage("insufficient_stock", "insufficient stock"),
		}
		r := setupCartRouter(svc)

		reqBody := `{"quantity": 1000}`
		req, _ := http.NewRequest("PATCH", "/cart/items/test-product", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "insufficient_stock", response.Error)
	})

	t.Run("service error", func(t *testing.T) {
		svc := &mockCartService{
			updateItemQuantityBySlugErr: errors.New("database error"),
		}
		r := setupCartRouter(svc)

		reqBody := `{"quantity": 3}`
		req, _ := http.NewRequest("PATCH", "/cart/items/test-product", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "internal_error", response.Error)
	})
}

func TestCartHandler_RemoveItem(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		cartResponse := &dto.CartResponse{
			Items:     []dto.CartItemResponse{},
			Total:     0,
			ItemCount: 0,
		}
		svc := &mockCartService{
			removeItemBySlugResult: cartResponse,
		}
		r := setupCartRouter(svc)

		req, _ := http.NewRequest("DELETE", "/cart/items/test-product", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dto.CartResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Items, 0)
	})

	t.Run("invalid product slug", func(t *testing.T) {
		svc := &mockCartService{}
		r := setupCartRouter(svc)

		req, _ := http.NewRequest("DELETE", "/cart/items/", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("cart item not found", func(t *testing.T) {
		svc := &mockCartService{
			removeItemBySlugErr: apperror.NewCodeMessage("cart_item_not_found", "cart item not found"),
		}
		r := setupCartRouter(svc)

		req, _ := http.NewRequest("DELETE", "/cart/items/nonexistent-product", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "cart_item_not_found", response.Error)
	})

	t.Run("service error", func(t *testing.T) {
		svc := &mockCartService{
			removeItemBySlugErr: errors.New("database error"),
		}
		r := setupCartRouter(svc)

		req, _ := http.NewRequest("DELETE", "/cart/items/test-product", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "internal_error", response.Error)
	})
}

func TestCartHandler_ClearCart(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		cartResponse := &dto.CartResponse{
			Items:     []dto.CartItemResponse{},
			Total:     0,
			ItemCount: 0,
		}
		svc := &mockCartService{
			clearCartResult: cartResponse,
		}
		r := setupCartRouter(svc)

		req, _ := http.NewRequest("DELETE", "/cart", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dto.CartResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Items, 0)
		assert.Equal(t, float64(0), response.Total)
		assert.Equal(t, 0, response.ItemCount)
	})

	t.Run("cart not found", func(t *testing.T) {
		svc := &mockCartService{
			clearCartErr: apperror.NewCodeMessage("cart_not_found", "cart not found"),
		}
		r := setupCartRouter(svc)

		req, _ := http.NewRequest("DELETE", "/cart", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "cart_not_found", response.Error)
	})

	t.Run("service error", func(t *testing.T) {
		svc := &mockCartService{
			clearCartErr: errors.New("database error"),
		}
		r := setupCartRouter(svc)

		req, _ := http.NewRequest("DELETE", "/cart", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "internal_error", response.Error)
	})
}
