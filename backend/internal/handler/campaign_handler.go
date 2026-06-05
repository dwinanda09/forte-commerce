package handler

import (
	"encoding/json"

	"github.com/dwinanda09/forte-commerce/internal/domain"
	"github.com/dwinanda09/forte-commerce/util"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CampaignHandler struct {
	campaignRepo domain.CampaignRepository
}

func NewCampaignHandler(campaignRepo domain.CampaignRepository) *CampaignHandler {
	return &CampaignHandler{campaignRepo: campaignRepo}
}

type CampaignRequest struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	IsActive    bool            `json:"is_active"`
	Priority    int             `json:"priority"`
	Conditions  json.RawMessage `json:"conditions"`
	Actions     json.RawMessage `json:"actions"`
}

func (h *CampaignHandler) List(c echo.Context) error {
	campaigns, err := h.campaignRepo.FindAll(c.Request().Context())
	if err != nil {
		return util.Fail(c, 500, "ERR-HDL-CAM-001", "Failed to fetch campaigns")
	}
	return util.OK(c, campaigns)
}

func (h *CampaignHandler) GetByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return util.Fail(c, 400, "ERR-HDL-CAM-002", "Invalid campaign ID")
	}
	campaign, err := h.campaignRepo.FindByID(c.Request().Context(), id)
	if err != nil {
		if util.IsNotFound(err) {
			return util.Fail(c, 404, "ERR-HDL-CAM-003", "Campaign not found")
		}
		return util.Fail(c, 500, "ERR-HDL-CAM-004", "Failed to fetch campaign")
	}
	return util.OK(c, campaign)
}

func (h *CampaignHandler) Create(c echo.Context) error {
	var req CampaignRequest
	if err := c.Bind(&req); err != nil {
		return util.Fail(c, 400, "ERR-HDL-CAM-005", "Invalid request body")
	}
	if req.Name == "" {
		return util.Fail(c, 400, "ERR-HDL-CAM-006", "Campaign name is required")
	}
	if len(req.Conditions) == 0 {
		req.Conditions = json.RawMessage("[]")
	}
	if len(req.Actions) == 0 {
		req.Actions = json.RawMessage("[]")
	}
	campaign := &domain.Campaign{
		Name:        req.Name,
		Description: req.Description,
		IsActive:    req.IsActive,
		Priority:    req.Priority,
		Conditions:  req.Conditions,
		Actions:     req.Actions,
	}
	if err := h.campaignRepo.Create(c.Request().Context(), campaign); err != nil {
		return util.Fail(c, 500, "ERR-HDL-CAM-007", "Failed to create campaign")
	}
	return util.OK(c, campaign)
}

func (h *CampaignHandler) Update(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return util.Fail(c, 400, "ERR-HDL-CAM-008", "Invalid campaign ID")
	}
	existing, err := h.campaignRepo.FindByID(c.Request().Context(), id)
	if err != nil {
		if util.IsNotFound(err) {
			return util.Fail(c, 404, "ERR-HDL-CAM-009", "Campaign not found")
		}
		return util.Fail(c, 500, "ERR-HDL-CAM-010", "Failed to fetch campaign")
	}
	var req CampaignRequest
	if err := c.Bind(&req); err != nil {
		return util.Fail(c, 400, "ERR-HDL-CAM-011", "Invalid request body")
	}
	if req.Name != "" {
		existing.Name = req.Name
	}
	existing.Description = req.Description
	existing.IsActive = req.IsActive
	existing.Priority = req.Priority
	if len(req.Conditions) > 0 {
		existing.Conditions = req.Conditions
	}
	if len(req.Actions) > 0 {
		existing.Actions = req.Actions
	}
	if err := h.campaignRepo.Update(c.Request().Context(), existing); err != nil {
		return util.Fail(c, 500, "ERR-HDL-CAM-012", "Failed to update campaign")
	}
	return util.OK(c, existing)
}

func (h *CampaignHandler) Delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return util.Fail(c, 400, "ERR-HDL-CAM-013", "Invalid campaign ID")
	}
	if err := h.campaignRepo.Delete(c.Request().Context(), id); err != nil {
		if util.IsNotFound(err) {
			return util.Fail(c, 404, "ERR-HDL-CAM-014", "Campaign not found")
		}
		return util.Fail(c, 500, "ERR-HDL-CAM-015", "Failed to delete campaign")
	}
	return util.OK(c, map[string]string{"status": "deleted"})
}

func (h *CampaignHandler) Toggle(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return util.Fail(c, 400, "ERR-HDL-CAM-016", "Invalid campaign ID")
	}
	var body struct {
		Active bool `json:"active"`
	}
	if err := c.Bind(&body); err != nil {
		return util.Fail(c, 400, "ERR-HDL-CAM-017", "Invalid request body")
	}
	if err := h.campaignRepo.Toggle(c.Request().Context(), id, body.Active); err != nil {
		if util.IsNotFound(err) {
			return util.Fail(c, 404, "ERR-HDL-CAM-018", "Campaign not found")
		}
		return util.Fail(c, 500, "ERR-HDL-CAM-019", "Failed to toggle campaign")
	}
	return util.OK(c, map[string]interface{}{"id": id, "active": body.Active})
}
