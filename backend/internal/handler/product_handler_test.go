package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/dto"
	"github.com/leoferamos/aroma-sense/internal/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ---- MOCK SERVICE ----
type MockProductService struct {
	mock.Mock
}

func (m *MockProductService) CreateProduct(ctx context.Context, input dto.ProductFormDTO, file dto.FileUpload) error {
	args := m.Called(ctx, input, file)
	return args.Error(0)
}

func (m *MockProductService) GetProductByID(ctx context.Context, id uint) (dto.ProductResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return dto.ProductResponse{}, args.Error(1)
	}
	return args.Get(0).(dto.ProductResponse), args.Error(1)
}

func (m *MockProductService) GetLatestProducts(ctx context.Context, page int, limit int) ([]dto.ProductResponse, int, error) {
	args := m.Called(ctx, page, limit)
	if len(args.Get(0).([]dto.ProductResponse)) == 0 && args.Error(2) != nil {
		return []dto.ProductResponse{}, 0, args.Error(2)
	}
	return args.Get(0).([]dto.ProductResponse), args.Int(1), args.Error(2)
}

func (m *MockProductService) SearchProducts(ctx context.Context, query string, page int, limit int, sort string) ([]dto.ProductResponse, int, error) {
	args := m.Called(ctx, query, page, limit, sort)

	var items []dto.ProductResponse
	if v := args.Get(0); v != nil {
		items = v.([]dto.ProductResponse)
	} else {
		items = []dto.ProductResponse{}
	}

	total := 0
	if len(args) > 1 {
		if v, ok := args.Get(1).(int); ok {
			total = v
		}
	}

	var err error
	if len(args) > 2 {
		err = args.Error(2)
	}

	return items, total, err
}

func (m *MockProductService) UpdateProduct(ctx context.Context, id uint, input dto.UpdateProductRequest) error {
	args := m.Called(ctx, id, input)
	return args.Error(0)
}

func (m *MockProductService) DeleteProduct(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ---- SETUP ROUTER ----
func setupProductRouter() (*gin.Engine, *MockProductService) {
	mockService := new(MockProductService)
	productHandler := handler.NewProductHandler(mockService, nil, nil)

	router := gin.Default()
	// Public routes
	router.GET("/products/:id", productHandler.GetProduct)
	router.GET("/products", productHandler.GetLatestProducts)

	// Admin routes
	adminGroup := router.Group("/admin")
	{
		adminGroup.POST("/products", productHandler.CreateProduct)
		adminGroup.PUT("/products/:id", productHandler.UpdateProduct)
		adminGroup.DELETE("/products/:id", productHandler.DeleteProduct)
	}

	return router, mockService
}

// ---- HELPERS ----
func performProductRequest(t *testing.T, router *gin.Engine, method, url string, payload interface{}) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	var bodyBytes []byte
	if payload != nil {
		var err error
		bodyBytes, err = json.Marshal(payload)
		if err != nil {
			t.Fatalf("failed to marshal payload: %v", err)
		}
	}

	req, _ := http.NewRequest(method, url, bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	return w
}

func performMultipartRequest(t *testing.T, router *gin.Engine, method, url string, form map[string]string, withFile bool) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for key, val := range form {
		_ = writer.WriteField(key, val)
	}

	if withFile {
		part, err := writer.CreateFormFile("image", "test.jpg")
		if err != nil {
			t.Fatalf("Failed to create form file: %v", err)
		}
		_, err = part.Write([]byte("fake-image-content"))
		if err != nil {
			t.Fatalf("Failed to write to form file: %v", err)
		}
	}

	writer.Close()

	req, _ := http.NewRequest(method, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(w, req)
	return w
}

// ---- TESTS ----
func TestProductHandler_CreateProduct(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		router, mockService := setupProductRouter()

		form := map[string]string{
			"name":           "Test Perfume",
			"brand":          "Test Brand",
			"weight":         "100",
			"price":          "120.50",
			"category":       "Woody",
			"notes":          "Sandalwood,Cedar",
			"stock_quantity": "50",
		}

		mockService.On(
			"CreateProduct",
			mock.Anything,
			mock.AnythingOfType("dto.ProductFormDTO"),
			mock.AnythingOfType("dto.FileUpload"),
		).Return(nil)

		w := performMultipartRequest(t, router, http.MethodPost, "/admin/products", form, true)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), "Product created successfully")
		mockService.AssertExpectations(t)
	})

	t.Run("Missing Required Field", func(t *testing.T) {
		router, _ := setupProductRouter()

		form := map[string]string{
			"brand": "Test Brand",
		}

		w := performMultipartRequest(t, router, http.MethodPost, "/admin/products", form, true)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Missing Image", func(t *testing.T) {
		router, _ := setupProductRouter()
		form := map[string]string{
			"name":           "Test Fragance",
			"brand":          "Test Brand",
			"weight":         "100",
			"price":          "120.50",
			"category":       "Woody",
			"notes":          "Sandalwood,Cedar",
			"stock_quantity": "50",
		}
		w := performMultipartRequest(t, router, http.MethodPost, "/admin/products", form, false)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "image is required")
	})

	t.Run("Service Error", func(t *testing.T) {
		router, mockService := setupProductRouter()
		form := map[string]string{
			"name":           "Test Fragance",
			"brand":          "Test Brand",
			"weight":         "100",
			"price":          "120.50",
			"category":       "Woody",
			"notes":          "Sandalwood,Cedar",
			"stock_quantity": "50",
		}

		mockService.On(
			"CreateProduct",
			mock.Anything,
			mock.AnythingOfType("dto.ProductFormDTO"),
			mock.AnythingOfType("dto.FileUpload"),
		).Return(fmt.Errorf("service error"))

		w := performMultipartRequest(t, router, http.MethodPost, "/admin/products", form, true)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("File Read Error", func(t *testing.T) {
		router, _ := setupProductRouter()
		form := map[string]string{
			"name":           "Test Perfume",
			"brand":          "Test Brand",
			"weight":         "100",
			"price":          "120.50",
			"category":       "Woody",
			"notes":          "Sandalwood,Cedar",
			"stock_quantity": "50",
		}

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		for key, val := range form {
			_ = writer.WriteField(key, val)
		}

		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition",
			fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "image", "test.jpg"))
		_, err := writer.CreatePart(h)
		if err != nil {
			t.Fatalf("Failed to create form part: %v", err)
		}

		writer.Close()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/admin/products", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "failed to read file")
	})
}

func TestProductHandler_GetProduct(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		router, mockService := setupProductRouter()

		productResponse := dto.ProductResponse{
			ID:        1,
			Name:      "Test Fragrance",
			Brand:     "Test Brand",
			Price:     99.99,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockService.On("GetProductByID", mock.Anything, uint(1)).Return(productResponse, nil)

		w := performProductRequest(t, router, http.MethodGet, "/products/1", nil)

		assert.Equal(t, http.StatusOK, w.Code)

		var returnedProduct dto.ProductResponse
		json.Unmarshal(w.Body.Bytes(), &returnedProduct)
		assert.Equal(t, productResponse.ID, returnedProduct.ID)
		assert.Equal(t, productResponse.Name, returnedProduct.Name)

		mockService.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		router, mockService := setupProductRouter()

		mockService.On("GetProductByID", mock.Anything, uint(1)).Return(dto.ProductResponse{}, fmt.Errorf("Product not found"))

		w := performProductRequest(t, router, http.MethodGet, "/products/1", nil)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		router, _ := setupProductRouter()

		w := performProductRequest(t, router, http.MethodGet, "/products/abc", nil)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestProductHandler_GetLatestProducts(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		router, mockService := setupProductRouter()

		productResponses := []dto.ProductResponse{
			{ID: 1, Name: "Test Fragance 1"},
			{ID: 2, Name: "Test Fragance 2"},
		}

		mockService.On("GetLatestProducts", mock.Anything, 10).Return(productResponses, nil)

		w := performProductRequest(t, router, http.MethodGet, "/products", nil)

		assert.Equal(t, http.StatusOK, w.Code)

		var returnedProducts []dto.ProductResponse
		json.Unmarshal(w.Body.Bytes(), &returnedProducts)
		assert.Len(t, returnedProducts, 2)

		mockService.AssertExpectations(t)
	})

	t.Run("Service Error", func(t *testing.T) {
		router, mockService := setupProductRouter()

		mockService.On("GetLatestProducts", mock.Anything, 10).Return([]dto.ProductResponse{}, fmt.Errorf("database error"))

		w := performProductRequest(t, router, http.MethodGet, "/products", nil)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid Limit", func(t *testing.T) {
		router, _ := setupProductRouter()

		w := performProductRequest(t, router, http.MethodGet, "/products?limit=abc", nil)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestProductHandler_UpdateProduct(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		router, mockService := setupProductRouter()
		name := "Updated Name"
		payload := dto.UpdateProductRequest{Name: &name}

		mockService.On("UpdateProduct", mock.Anything, uint(1), payload).Return(nil)

		w := performProductRequest(t, router, http.MethodPut, "/admin/products/1", payload)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Product updated successfully")
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		router, _ := setupProductRouter()
		w := performProductRequest(t, router, http.MethodPut, "/admin/products/abc", nil)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Invalid Payload", func(t *testing.T) {
		router, _ := setupProductRouter()
		w := performProductRequest(t, router, http.MethodPut, "/admin/products/1", "invalid-json")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		router, mockService := setupProductRouter()
		name := "Updated Name"
		payload := dto.UpdateProductRequest{Name: &name}

		mockService.On("UpdateProduct", mock.Anything, uint(1), payload).Return(fmt.Errorf("service error"))

		w := performProductRequest(t, router, http.MethodPut, "/admin/products/1", payload)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestProductHandler_DeleteProduct(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		router, mockService := setupProductRouter()

		mockService.On("DeleteProduct", mock.Anything, uint(1)).Return(nil)

		w := performProductRequest(t, router, http.MethodDelete, "/admin/products/1", nil)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Product deleted successfully")
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		router, _ := setupProductRouter()
		w := performProductRequest(t, router, http.MethodDelete, "/admin/products/abc", nil)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		router, mockService := setupProductRouter()

		mockService.On("DeleteProduct", mock.Anything, uint(1)).Return(fmt.Errorf("service error"))

		w := performProductRequest(t, router, http.MethodDelete, "/admin/products/1", nil)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}
