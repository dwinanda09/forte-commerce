package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dwinanda09/forte-commerce/internal/domain"
	"github.com/dwinanda09/forte-commerce/internal/mocks"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func newProductEcho(method, path string, body []byte) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	var req *http.Request
	if body != nil {
		req = httptest.NewRequest(method, path, bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}

func sampleProduct() *domain.Product {
	return &domain.Product{
		ID:           uuid.New(),
		SKU:          "TEST-SKU",
		Name:         "Test Product",
		Price:        99.99,
		InventoryQty: 100,
		ReservedQty:  10,
	}
}

// ── List ──────────────────────────────────────────────────────────────────────

func TestProductHandler_List_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	products := []domain.Product{*sampleProduct(), *sampleProduct()}
	repo := mocks.NewMockProductRepository(ctrl)
	repo.EXPECT().FindAll(gomock.Any()).Return(products, nil)

	h := NewProductHandler(repo)
	c, rec := newProductEcho(http.MethodGet, "/products", nil)

	require.NoError(t, h.List(c))
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.True(t, resp["success"].(bool))
}

func TestProductHandler_List_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockProductRepository(ctrl)
	repo.EXPECT().FindAll(gomock.Any()).Return(nil, errors.New("db error"))

	h := NewProductHandler(repo)
	c, rec := newProductEcho(http.MethodGet, "/products", nil)

	require.NoError(t, h.List(c))
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

// ── Create ────────────────────────────────────────────────────────────────────

func TestProductHandler_Create_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockProductRepository(ctrl)
	repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	h := NewProductHandler(repo)
	body, _ := json.Marshal(CreateProductRequest{
		SKU:          "NEW-SKU",
		Name:         "New Product",
		Price:        49.99,
		InventoryQty: 50,
	})
	c, rec := newProductEcho(http.MethodPost, "/products", body)

	require.NoError(t, h.Create(c))
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.True(t, resp["success"].(bool))
}

func TestProductHandler_Create_MissingSKU(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockProductRepository(ctrl)
	h := NewProductHandler(repo)
	body, _ := json.Marshal(CreateProductRequest{
		Name:  "No SKU",
		Price: 49.99,
	})
	c, rec := newProductEcho(http.MethodPost, "/products", body)

	require.NoError(t, h.Create(c))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestProductHandler_Create_MissingName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockProductRepository(ctrl)
	h := NewProductHandler(repo)
	body, _ := json.Marshal(CreateProductRequest{
		SKU:   "SKU-001",
		Price: 49.99,
	})
	c, rec := newProductEcho(http.MethodPost, "/products", body)

	require.NoError(t, h.Create(c))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestProductHandler_Create_MissingPrice(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockProductRepository(ctrl)
	h := NewProductHandler(repo)
	body, _ := json.Marshal(CreateProductRequest{
		SKU:  "SKU-001",
		Name: "Product",
	})
	c, rec := newProductEcho(http.MethodPost, "/products", body)

	require.NoError(t, h.Create(c))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestProductHandler_Create_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockProductRepository(ctrl)
	repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("db error"))

	h := NewProductHandler(repo)
	body, _ := json.Marshal(CreateProductRequest{
		SKU:   "SKU-001",
		Name:  "Product",
		Price: 49.99,
	})
	c, rec := newProductEcho(http.MethodPost, "/products", body)

	require.NoError(t, h.Create(c))
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

// ── Update ────────────────────────────────────────────────────────────────────

func TestProductHandler_Update_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prod := sampleProduct()
	repo := mocks.NewMockProductRepository(ctrl)
	repo.EXPECT().FindByID(gomock.Any(), prod.ID).Return(prod, nil)
	repo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	h := NewProductHandler(repo)
	body, _ := json.Marshal(UpdateProductRequest{Name: "Updated Product"})
	c, rec := newProductEcho(http.MethodPut, "/products/"+prod.ID.String(), body)
	c.SetParamNames("id")
	c.SetParamValues(prod.ID.String())

	require.NoError(t, h.Update(c))
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.True(t, resp["success"].(bool))
}

func TestProductHandler_Update_InvalidUUID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockProductRepository(ctrl)
	h := NewProductHandler(repo)
	body, _ := json.Marshal(UpdateProductRequest{Name: "Updated"})
	c, rec := newProductEcho(http.MethodPut, "/products/not-a-uuid", body)
	c.SetParamNames("id")
	c.SetParamValues("not-a-uuid")

	require.NoError(t, h.Update(c))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestProductHandler_Update_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := uuid.New()
	repo := mocks.NewMockProductRepository(ctrl)
	repo.EXPECT().FindByID(gomock.Any(), id).Return(nil, errNotFound("product not found"))

	h := NewProductHandler(repo)
	body, _ := json.Marshal(UpdateProductRequest{Name: "Updated"})
	c, rec := newProductEcho(http.MethodPut, "/products/"+id.String(), body)
	c.SetParamNames("id")
	c.SetParamValues(id.String())

	require.NoError(t, h.Update(c))
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestProductHandler_Update_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prod := sampleProduct()
	repo := mocks.NewMockProductRepository(ctrl)
	repo.EXPECT().FindByID(gomock.Any(), prod.ID).Return(prod, nil)
	repo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errors.New("db error"))

	h := NewProductHandler(repo)
	body, _ := json.Marshal(UpdateProductRequest{Name: "Updated"})
	c, rec := newProductEcho(http.MethodPut, "/products/"+prod.ID.String(), body)
	c.SetParamNames("id")
	c.SetParamValues(prod.ID.String())

	require.NoError(t, h.Update(c))
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

// ── Delete ────────────────────────────────────────────────────────────────────

func TestProductHandler_Delete_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := uuid.New()
	repo := mocks.NewMockProductRepository(ctrl)
	repo.EXPECT().Delete(gomock.Any(), id).Return(nil)

	h := NewProductHandler(repo)
	c, rec := newProductEcho(http.MethodDelete, "/products/"+id.String(), nil)
	c.SetParamNames("id")
	c.SetParamValues(id.String())

	require.NoError(t, h.Delete(c))
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "deleted", data["status"])
}

func TestProductHandler_Delete_InvalidUUID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockProductRepository(ctrl)
	h := NewProductHandler(repo)
	c, rec := newProductEcho(http.MethodDelete, "/products/not-a-uuid", nil)
	c.SetParamNames("id")
	c.SetParamValues("not-a-uuid")

	require.NoError(t, h.Delete(c))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestProductHandler_Delete_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := uuid.New()
	repo := mocks.NewMockProductRepository(ctrl)
	repo.EXPECT().Delete(gomock.Any(), id).Return(errNotFound("not found"))

	h := NewProductHandler(repo)
	c, rec := newProductEcho(http.MethodDelete, "/products/"+id.String(), nil)
	c.SetParamNames("id")
	c.SetParamValues(id.String())

	require.NoError(t, h.Delete(c))
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestProductHandler_Delete_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := uuid.New()
	repo := mocks.NewMockProductRepository(ctrl)
	repo.EXPECT().Delete(gomock.Any(), id).Return(errors.New("db error"))

	h := NewProductHandler(repo)
	c, rec := newProductEcho(http.MethodDelete, "/products/"+id.String(), nil)
	c.SetParamNames("id")
	c.SetParamValues(id.String())

	require.NoError(t, h.Delete(c))
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
