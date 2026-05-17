package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/superplanehq/superplane/pkg/database"
	"gorm.io/gorm"
)

type App struct {
	ID                    uuid.UUID
	OrganizationID        uuid.UUID
	DisplayName           string
	Slug                  string
	Description           string
	CanvasID              *uuid.UUID
	CodeStorageRepoID     string
	CodeStorageRemoteURL  string
	DefaultBranch         string
	LiveCommitSHA         string
	EditSessionBranch     *string
	SyncStatus            string
	SyncError             *string
	CreatedBy             *uuid.UUID
	CreatedAt             *time.Time
	UpdatedAt             *time.Time
	DeletedAt             gorm.DeletedAt `gorm:"index"`
}

func FindAppByID(id uuid.UUID) (*App, error) {
	return FindAppByIDInTransaction(database.Conn(), id)
}

func FindAppByIDInTransaction(tx *gorm.DB, id uuid.UUID) (*App, error) {
	var app App
	err := tx.
		Where("id = ?", id).
		First(&app).
		Error

	if err != nil {
		return nil, err
	}

	return &app, nil
}

func FindAppsByOrganizationID(orgID uuid.UUID) ([]App, error) {
	return FindAppsByOrganizationIDInTransaction(database.Conn(), orgID)
}

func FindAppsByOrganizationIDInTransaction(tx *gorm.DB, orgID uuid.UUID) ([]App, error) {
	var apps []App
	err := tx.
		Where("organization_id = ?", orgID).
		Order("display_name ASC").
		Find(&apps).
		Error

	if err != nil {
		return nil, err
	}

	return apps, nil
}
