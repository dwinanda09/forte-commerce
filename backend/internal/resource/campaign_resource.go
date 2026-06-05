package resource

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dwinanda09/forte-commerce/internal/domain"
	"github.com/dwinanda09/forte-commerce/util"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CampaignResource struct {
	db     *sqlx.DB
	logger *util.Logger
}

func NewCampaignResource(db *sqlx.DB, logger *util.Logger) *CampaignResource {
	return &CampaignResource{db: db, logger: logger}
}

const campaignColumns = `id, name, description, is_active, priority, conditions, actions, created_at, updated_at`

func scanCampaign(row interface {
	Scan(dest ...any) error
}) (*domain.Campaign, error) {
	var c domain.Campaign
	var condRaw, actRaw []byte
	if err := row.Scan(&c.ID, &c.Name, &c.Description, &c.IsActive, &c.Priority, &condRaw, &actRaw, &c.CreatedAt, &c.UpdatedAt); err != nil {
		return nil, err
	}
	c.Conditions = json.RawMessage(condRaw)
	c.Actions = json.RawMessage(actRaw)
	return &c, nil
}

func scanCampaigns(rows *sql.Rows) ([]domain.Campaign, error) {
	var campaigns []domain.Campaign
	for rows.Next() {
		c, err := scanCampaign(rows)
		if err != nil {
			return nil, err
		}
		campaigns = append(campaigns, *c)
	}
	return campaigns, rows.Err()
}

func (r *CampaignResource) FindAll(ctx context.Context) ([]domain.Campaign, error) {
	start := r.logger.Start(ctx, "CampaignResource.FindAll")
	defer func() { r.logger.Finish(ctx, "CampaignResource.FindAll", start, nil) }()

	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(`SELECT %s FROM campaigns ORDER BY priority, created_at`, campaignColumns))
	if err != nil {
		r.logger.Finish(ctx, "CampaignResource.FindAll", start, err)
		return nil, util.Wrap("ERR-RS-CAM-001", "Failed to fetch campaigns", err)
	}
	defer rows.Close()

	campaigns, err := scanCampaigns(rows)
	if err != nil {
		r.logger.Finish(ctx, "CampaignResource.FindAll", start, err)
		return nil, util.Wrap("ERR-RS-CAM-001", "Failed to scan campaigns", err)
	}
	return campaigns, nil
}

func (r *CampaignResource) FindActive(ctx context.Context) ([]domain.Campaign, error) {
	start := r.logger.Start(ctx, "CampaignResource.FindActive")
	defer func() { r.logger.Finish(ctx, "CampaignResource.FindActive", start, nil) }()

	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(`SELECT %s FROM campaigns WHERE is_active = true ORDER BY priority, created_at`, campaignColumns))
	if err != nil {
		r.logger.Finish(ctx, "CampaignResource.FindActive", start, err)
		return nil, util.Wrap("ERR-RS-CAM-002", "Failed to fetch active campaigns", err)
	}
	defer rows.Close()

	campaigns, err := scanCampaigns(rows)
	if err != nil {
		r.logger.Finish(ctx, "CampaignResource.FindActive", start, err)
		return nil, util.Wrap("ERR-RS-CAM-002", "Failed to scan active campaigns", err)
	}
	return campaigns, nil
}

func (r *CampaignResource) FindByID(ctx context.Context, id uuid.UUID) (*domain.Campaign, error) {
	start := r.logger.Start(ctx, "CampaignResource.FindByID")
	defer func() { r.logger.Finish(ctx, "CampaignResource.FindByID", start, nil) }()

	row := r.db.QueryRowContext(ctx, fmt.Sprintf(`SELECT %s FROM campaigns WHERE id = $1`, campaignColumns), id)
	c, err := scanCampaign(row)
	if err != nil {
		r.logger.Finish(ctx, "CampaignResource.FindByID", start, err)
		return nil, util.Wrap("ERR-RS-CAM-003", "Campaign not found", err)
	}
	return c, nil
}

func (r *CampaignResource) Create(ctx context.Context, c *domain.Campaign) error {
	start := r.logger.Start(ctx, "CampaignResource.Create")
	defer func() { r.logger.Finish(ctx, "CampaignResource.Create", start, nil) }()

	c.ID = uuid.New()
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO campaigns (id, name, description, is_active, priority, conditions, actions, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, c.ID, c.Name, c.Description, c.IsActive, c.Priority, []byte(c.Conditions), []byte(c.Actions), c.CreatedAt, c.UpdatedAt)
	if err != nil {
		r.logger.Finish(ctx, "CampaignResource.Create", start, err)
		return util.Wrap("ERR-RS-CAM-004", "Failed to create campaign", err)
	}
	return nil
}

func (r *CampaignResource) Update(ctx context.Context, c *domain.Campaign) error {
	start := r.logger.Start(ctx, "CampaignResource.Update")
	defer func() { r.logger.Finish(ctx, "CampaignResource.Update", start, nil) }()

	c.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, `
		UPDATE campaigns
		SET name=$1, description=$2, is_active=$3, priority=$4, conditions=$5, actions=$6, updated_at=$7
		WHERE id=$8
	`, c.Name, c.Description, c.IsActive, c.Priority, []byte(c.Conditions), []byte(c.Actions), c.UpdatedAt, c.ID)
	if err != nil {
		r.logger.Finish(ctx, "CampaignResource.Update", start, err)
		return util.Wrap("ERR-RS-CAM-005", "Failed to update campaign", err)
	}
	return nil
}

func (r *CampaignResource) Delete(ctx context.Context, id uuid.UUID) error {
	start := r.logger.Start(ctx, "CampaignResource.Delete")
	defer func() { r.logger.Finish(ctx, "CampaignResource.Delete", start, nil) }()

	_, err := r.db.ExecContext(ctx, `DELETE FROM campaigns WHERE id = $1`, id)
	if err != nil {
		r.logger.Finish(ctx, "CampaignResource.Delete", start, err)
		return util.Wrap("ERR-RS-CAM-006", "Failed to delete campaign", err)
	}
	return nil
}

func (r *CampaignResource) Toggle(ctx context.Context, id uuid.UUID, active bool) error {
	start := r.logger.Start(ctx, "CampaignResource.Toggle")
	defer func() { r.logger.Finish(ctx, "CampaignResource.Toggle", start, nil) }()

	_, err := r.db.ExecContext(ctx, `UPDATE campaigns SET is_active=$1, updated_at=$2 WHERE id=$3`, active, time.Now(), id)
	if err != nil {
		r.logger.Finish(ctx, "CampaignResource.Toggle", start, err)
		return util.Wrap("ERR-RS-CAM-007", "Failed to toggle campaign", err)
	}
	return nil
}
