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
	"github.com/dwinanda09/forte-commerce/util"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func newCampaignEcho(method, path string, body []byte) (echo.Context, *httptest.ResponseRecorder) {
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

func sampleCampaign() *domain.Campaign {
	return &domain.Campaign{
		ID:          uuid.New(),
		Name:        "Test Campaign",
		Description: "desc",
		IsActive:    true,
		Priority:    1,
		Conditions:  json.RawMessage(`[]`),
		Actions:     json.RawMessage(`[]`),
	}
}

// ── List ──────────────────────────────────────────────────────────────────────

func TestCampaignHandler_List_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	campaigns := []domain.Campaign{*sampleCampaign()}
	repo := mocks.NewMockCampaignRepository(ctrl)
	repo.EXPECT().FindAll(gomock.Any()).Return(campaigns, nil)

	h := NewCampaignHandler(repo)
	c, rec := newCampaignEcho(http.MethodGet, "/campaigns", nil)

	require.NoError(t, h.List(c))
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.True(t, resp["success"].(bool))
}

func TestCampaignHandler_List_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockCampaignRepository(ctrl)
	repo.EXPECT().FindAll(gomock.Any()).Return(nil, errors.New("db error"))

	h := NewCampaignHandler(repo)
	c, rec := newCampaignEcho(http.MethodGet, "/campaigns", nil)

	require.NoError(t, h.List(c))
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

// ── GetByID ───────────────────────────────────────────────────────────────────

func TestCampaignHandler_GetByID_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cam := sampleCampaign()
	repo := mocks.NewMockCampaignRepository(ctrl)
	repo.EXPECT().FindByID(gomock.Any(), cam.ID).Return(cam, nil)

	h := NewCampaignHandler(repo)
	c, rec := newCampaignEcho(http.MethodGet, "/campaigns/"+cam.ID.String(), nil)
	c.SetParamNames("id")
	c.SetParamValues(cam.ID.String())

	require.NoError(t, h.GetByID(c))
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestCampaignHandler_GetByID_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockCampaignRepository(ctrl)
	h := NewCampaignHandler(repo)
	c, rec := newCampaignEcho(http.MethodGet, "/campaigns/not-a-uuid", nil)
	c.SetParamNames("id")
	c.SetParamValues("not-a-uuid")

	require.NoError(t, h.GetByID(c))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCampaignHandler_GetByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := uuid.New()
	repo := mocks.NewMockCampaignRepository(ctrl)
	// wrap error so IsNotFound returns true
	repo.EXPECT().FindByID(gomock.Any(), id).Return(nil, errNotFound("campaign not found"))

	h := NewCampaignHandler(repo)
	c, rec := newCampaignEcho(http.MethodGet, "/campaigns/"+id.String(), nil)
	c.SetParamNames("id")
	c.SetParamValues(id.String())

	require.NoError(t, h.GetByID(c))
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// ── Create ────────────────────────────────────────────────────────────────────

func TestCampaignHandler_Create_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockCampaignRepository(ctrl)
	repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	h := NewCampaignHandler(repo)
	body, _ := json.Marshal(CampaignRequest{
		Name:        "New Campaign",
		Description: "desc",
		IsActive:    true,
		Conditions:  json.RawMessage(`[{"type":"cart_has_sku","sku":"ABC"}]`),
		Actions:     json.RawMessage(`[{"type":"fixed_discount","amount":10}]`),
	})
	c, rec := newCampaignEcho(http.MethodPost, "/campaigns", body)

	require.NoError(t, h.Create(c))
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestCampaignHandler_Create_MissingName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockCampaignRepository(ctrl)
	h := NewCampaignHandler(repo)
	body, _ := json.Marshal(CampaignRequest{Description: "no name"})
	c, rec := newCampaignEcho(http.MethodPost, "/campaigns", body)

	require.NoError(t, h.Create(c))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCampaignHandler_Create_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockCampaignRepository(ctrl)
	repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("db error"))

	h := NewCampaignHandler(repo)
	body, _ := json.Marshal(CampaignRequest{Name: "Camp"})
	c, rec := newCampaignEcho(http.MethodPost, "/campaigns", body)

	require.NoError(t, h.Create(c))
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

// ── Update ────────────────────────────────────────────────────────────────────

func TestCampaignHandler_Update_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cam := sampleCampaign()
	repo := mocks.NewMockCampaignRepository(ctrl)
	repo.EXPECT().FindByID(gomock.Any(), cam.ID).Return(cam, nil)
	repo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	h := NewCampaignHandler(repo)
	body, _ := json.Marshal(CampaignRequest{Name: "Updated"})
	c, rec := newCampaignEcho(http.MethodPut, "/campaigns/"+cam.ID.String(), body)
	c.SetParamNames("id")
	c.SetParamValues(cam.ID.String())

	require.NoError(t, h.Update(c))
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestCampaignHandler_Update_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := uuid.New()
	repo := mocks.NewMockCampaignRepository(ctrl)
	repo.EXPECT().FindByID(gomock.Any(), id).Return(nil, errNotFound("not found"))

	h := NewCampaignHandler(repo)
	body, _ := json.Marshal(CampaignRequest{Name: "X"})
	c, rec := newCampaignEcho(http.MethodPut, "/campaigns/"+id.String(), body)
	c.SetParamNames("id")
	c.SetParamValues(id.String())

	require.NoError(t, h.Update(c))
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// ── Delete ────────────────────────────────────────────────────────────────────

func TestCampaignHandler_Delete_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := uuid.New()
	repo := mocks.NewMockCampaignRepository(ctrl)
	repo.EXPECT().Delete(gomock.Any(), id).Return(nil)

	h := NewCampaignHandler(repo)
	c, rec := newCampaignEcho(http.MethodDelete, "/campaigns/"+id.String(), nil)
	c.SetParamNames("id")
	c.SetParamValues(id.String())

	require.NoError(t, h.Delete(c))
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestCampaignHandler_Delete_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := uuid.New()
	repo := mocks.NewMockCampaignRepository(ctrl)
	repo.EXPECT().Delete(gomock.Any(), id).Return(errNotFound("not found"))

	h := NewCampaignHandler(repo)
	c, rec := newCampaignEcho(http.MethodDelete, "/campaigns/"+id.String(), nil)
	c.SetParamNames("id")
	c.SetParamValues(id.String())

	require.NoError(t, h.Delete(c))
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// ── Toggle ────────────────────────────────────────────────────────────────────

func TestCampaignHandler_Toggle_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := uuid.New()
	repo := mocks.NewMockCampaignRepository(ctrl)
	repo.EXPECT().Toggle(gomock.Any(), id, false).Return(nil)

	h := NewCampaignHandler(repo)
	body, _ := json.Marshal(map[string]bool{"active": false})
	c, rec := newCampaignEcho(http.MethodPatch, "/campaigns/"+id.String()+"/toggle", body)
	c.SetParamNames("id")
	c.SetParamValues(id.String())

	require.NoError(t, h.Toggle(c))
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, false, data["active"])
}

func TestCampaignHandler_Toggle_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := uuid.New()
	repo := mocks.NewMockCampaignRepository(ctrl)
	repo.EXPECT().Toggle(gomock.Any(), id, true).Return(errNotFound("not found"))

	h := NewCampaignHandler(repo)
	body, _ := json.Marshal(map[string]bool{"active": true})
	c, rec := newCampaignEcho(http.MethodPatch, "/campaigns/"+id.String()+"/toggle", body)
	c.SetParamNames("id")
	c.SetParamValues(id.String())

	require.NoError(t, h.Toggle(c))
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// ── helpers ───────────────────────────────────────────────────────────────────

// errNotFound returns an AppError whose code ends in "-404", satisfying util.IsNotFound.
func errNotFound(msg string) error {
	return util.Wrap("ERR-TEST-404", msg, errors.New(msg))
}
