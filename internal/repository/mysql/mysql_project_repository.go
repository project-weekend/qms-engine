package mysql

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/project-weekend/qms-engine/internal/entity"
)

type ProjectRepository struct {
	Logger *slog.Logger
}

func NewProjectRepository(logger *slog.Logger) *ProjectRepository {
	return &ProjectRepository{
		Logger: logger,
	}
}

// Save creates a new project in the database
func (p *ProjectRepository) Save(tx *sqlx.Tx, project *entity.Project) (*entity.Project, error) {
	query := `
		INSERT INTO projects (name, description, created_at, updated_at)
		VALUES (?, ?, ?, ?)
	`

	now := time.Now()
	result, err := tx.Exec(query,
		project.Name,
		project.Description,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert project: %w", err)
	}

	// Get the last inserted ID
	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	project.ID = int(id)
	project.CreatedAt = now
	project.UpdatedAt = now

	return project, nil
}

// GetByName retrieves a project by its name
func (p *ProjectRepository) GetByName(tx *sqlx.Tx, name string) (*entity.Project, error) {
	query := `
		SELECT id, name, description, created_at, updated_at, deleted_at
		FROM projects
		WHERE name = ? AND deleted_at IS NULL
	`

	var project entity.Project
	err := tx.Get(&project, query, name)
	if err != nil {
		return nil, err
	}

	return &project, nil
}
