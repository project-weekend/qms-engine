package repository

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/project-weekend/qms-engine/server/config"
)

// Entity interface that all entities must implement
type Entity interface {
	GetTableName() string
}

// Repository provides generic CRUD operations for entities
type Repository[T Entity] struct {
	DB *config.DBConnections
}

// NewRepository creates a new repository instance
func NewRepository[T Entity](db *config.DBConnections) *Repository[T] {
	return &Repository[T]{DB: db}
}

// Save inserts a single entity into the database
func (r *Repository[T]) Save(entity *T) error {
	tableName := (*entity).GetTableName()

	query, args, err := r.buildInsertQuery(tableName, entity)
	if err != nil {
		return fmt.Errorf("failed to build insert query: %w", err)
	}

	result, err := r.DB.Master.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to insert entity into %s: %w", tableName, err)
	}

	// Get the last inserted ID and set it to the entity if it has an ID field
	lastID, err := result.LastInsertId()
	if err == nil && lastID > 0 {
		r.setIDField(entity, lastID)
	}

	return nil
}

// SaveBatch inserts multiple entities in a single transaction
func (r *Repository[T]) SaveBatch(entities []*T) error {
	if len(entities) == 0 {
		return nil
	}

	tableName := (*entities[0]).GetTableName()

	tx, err := r.DB.Master.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	for _, entity := range entities {
		query, args, err := r.buildInsertQuery(tableName, entity)
		if err != nil {
			return fmt.Errorf("failed to build insert query: %w", err)
		}

		_, err = tx.Exec(query, args...)
		if err != nil {
			return fmt.Errorf("failed to insert entity into %s: %w", tableName, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Update updates an entity by ID
func (r *Repository[T]) Update(id int64, entity *T) error {
	tableName := (*entity).GetTableName()

	query, args, err := r.buildUpdateQuery(tableName, entity, id)
	if err != nil {
		return fmt.Errorf("failed to build update query: %w", err)
	}

	result, err := r.DB.Master.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update entity in %s: %w", tableName, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected, entity with id %d not found", id)
	}

	return nil
}

// UpdateBatch updates multiple entities in a single transaction
func (r *Repository[T]) UpdateBatch(entities []*T) error {
	if len(entities) == 0 {
		return nil
	}

	tableName := (*entities[0]).GetTableName()

	tx, err := r.DB.Master.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	for _, entity := range entities {
		id := r.getIDField(entity)
		if id == 0 {
			return fmt.Errorf("entity must have a valid ID for batch update")
		}

		query, args, err := r.buildUpdateQuery(tableName, entity, id)
		if err != nil {
			return fmt.Errorf("failed to build update query: %w", err)
		}

		_, err = tx.Exec(query, args...)
		if err != nil {
			return fmt.Errorf("failed to update entity in %s: %w", tableName, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// UpdateEntity updates specific fields of an entity
func (r *Repository[T]) UpdateEntity(id int64, updates map[string]interface{}) error {
	var entity T
	tableName := entity.GetTableName()

	if len(updates) == 0 {
		return fmt.Errorf("no fields to update")
	}

	// Add updated_at timestamp if not present
	if _, exists := updates["updated_at"]; !exists {
		updates["updated_at"] = time.Now()
	}

	setClauses := make([]string, 0, len(updates))
	args := make([]interface{}, 0, len(updates)+1)

	for field, value := range updates {
		setClauses = append(setClauses, fmt.Sprintf("%s = ?", field))
		args = append(args, value)
	}
	args = append(args, id)

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = ?", tableName, strings.Join(setClauses, ", "))

	result, err := r.DB.Master.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update entity in %s: %w", tableName, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected, entity with id %d not found", id)
	}

	return nil
}

// Upsert inserts or updates an entity based on ID or unique constraint
func (r *Repository[T]) Upsert(entity *T) error {
	id := r.getIDField(entity)
	if id > 0 {
		// Check if entity exists
		exists, err := r.Exists(id)
		if err != nil {
			return fmt.Errorf("failed to check if entity exists: %w", err)
		}

		if exists {
			return r.Update(id, entity)
		}
	}

	return r.Save(entity)
}

// Find retrieves entities based on conditions
func (r *Repository[T]) Find(conditions map[string]interface{}, limit, offset int) ([]*T, error) {
	var entity T
	tableName := entity.GetTableName()

	query := fmt.Sprintf("SELECT * FROM %s", tableName)
	args := make([]interface{}, 0)

	if len(conditions) > 0 {
		whereClauses := make([]string, 0, len(conditions))
		for field, value := range conditions {
			whereClauses = append(whereClauses, fmt.Sprintf("%s = ?", field))
			args = append(args, value)
		}
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}

	if offset > 0 {
		query += fmt.Sprintf(" OFFSET %d", offset)
	}

	var results []*T
	err := r.DB.Slave.Select(&results, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to find entities in %s: %w", tableName, err)
	}

	return results, nil
}

// FindByID retrieves a single entity by ID
func (r *Repository[T]) FindByID(id int64) (*T, error) {
	var entity T
	tableName := entity.GetTableName()

	query := fmt.Sprintf("SELECT * FROM %s WHERE id = ?", tableName)

	err := r.DB.Slave.Get(&entity, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("entity with id %d not found in %s", id, tableName)
		}
		return nil, fmt.Errorf("failed to find entity in %s: %w", tableName, err)
	}

	return &entity, nil
}

// FindOne retrieves a single entity based on conditions
func (r *Repository[T]) FindOne(conditions map[string]interface{}) (*T, error) {
	var entity T
	tableName := entity.GetTableName()

	query := fmt.Sprintf("SELECT * FROM %s", tableName)
	args := make([]interface{}, 0)

	if len(conditions) > 0 {
		whereClauses := make([]string, 0, len(conditions))
		for field, value := range conditions {
			whereClauses = append(whereClauses, fmt.Sprintf("%s = ?", field))
			args = append(args, value)
		}
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	query += " LIMIT 1"

	err := r.DB.Slave.Get(&entity, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("entity not found in %s", tableName)
		}
		return nil, fmt.Errorf("failed to find entity in %s: %w", tableName, err)
	}

	return &entity, nil
}

// Delete performs a soft delete by setting deleted_at timestamp
func (r *Repository[T]) Delete(id int64) error {
	var entity T
	tableName := entity.GetTableName()

	query := fmt.Sprintf("UPDATE %s SET deleted_at = ? WHERE id = ?", tableName)

	result, err := r.DB.Master.Exec(query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to delete entity in %s: %w", tableName, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected, entity with id %d not found", id)
	}

	return nil
}

// HardDelete permanently removes an entity from the database
func (r *Repository[T]) HardDelete(id int64) error {
	var entity T
	tableName := entity.GetTableName()

	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", tableName)

	result, err := r.DB.Master.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to hard delete entity in %s: %w", tableName, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected, entity with id %d not found", id)
	}

	return nil
}

// Count returns the total number of entities based on conditions
func (r *Repository[T]) Count(conditions map[string]interface{}) (int64, error) {
	var entity T
	tableName := entity.GetTableName()

	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)
	args := make([]interface{}, 0)

	if len(conditions) > 0 {
		whereClauses := make([]string, 0, len(conditions))
		for field, value := range conditions {
			whereClauses = append(whereClauses, fmt.Sprintf("%s = ?", field))
			args = append(args, value)
		}
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	var count int64
	err := r.DB.Slave.Get(&count, query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to count entities in %s: %w", tableName, err)
	}

	return count, nil
}

// Exists checks if an entity with the given ID exists
func (r *Repository[T]) Exists(id int64) (bool, error) {
	var entity T
	tableName := entity.GetTableName()

	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE id = ?", tableName)

	var count int64
	err := r.DB.Slave.Get(&count, query, id)
	if err != nil {
		return false, fmt.Errorf("failed to check entity existence in %s: %w", tableName, err)
	}

	return count > 0, nil
}

// ExecuteInTransaction executes a function within a database transaction
func (r *Repository[T]) ExecuteInTransaction(fn func(*sqlx.Tx) error) error {
	tx, err := r.DB.Master.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if err := fn(tx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Helper functions

func (r *Repository[T]) buildInsertQuery(tableName string, entity *T) (string, []interface{}, error) {
	v := reflect.ValueOf(*entity)
	t := v.Type()

	columns := make([]string, 0)
	placeholders := make([]string, 0)
	args := make([]interface{}, 0)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		dbTag := field.Tag.Get("db")

		if dbTag == "" {
			dbTag = strings.ToLower(field.Name)
		}

		// Skip ID, created_at, updated_at, deleted_at for insert
		if dbTag == "id" || dbTag == "-" {
			continue
		}

		fieldValue := v.Field(i)

		// Handle timestamps
		if dbTag == "created_at" || dbTag == "updated_at" {
			columns = append(columns, dbTag)
			placeholders = append(placeholders, "?")
			args = append(args, time.Now())
			continue
		}

		if dbTag == "deleted_at" {
			continue
		}

		columns = append(columns, dbTag)
		placeholders = append(placeholders, "?")
		args = append(args, fieldValue.Interface())
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		tableName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "))

	return query, args, nil
}

func (r *Repository[T]) buildUpdateQuery(tableName string, entity *T, id int64) (string, []interface{}, error) {
	v := reflect.ValueOf(*entity)
	t := v.Type()

	setClauses := make([]string, 0)
	args := make([]interface{}, 0)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		dbTag := field.Tag.Get("db")

		if dbTag == "" {
			dbTag = strings.ToLower(field.Name)
		}

		// Skip ID, created_at, deleted_at for update
		if dbTag == "id" || dbTag == "created_at" || dbTag == "deleted_at" || dbTag == "-" {
			continue
		}

		fieldValue := v.Field(i)

		// Handle updated_at
		if dbTag == "updated_at" {
			setClauses = append(setClauses, fmt.Sprintf("%s = ?", dbTag))
			args = append(args, time.Now())
			continue
		}

		setClauses = append(setClauses, fmt.Sprintf("%s = ?", dbTag))
		args = append(args, fieldValue.Interface())
	}

	args = append(args, id)

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = ?",
		tableName,
		strings.Join(setClauses, ", "))

	return query, args, nil
}

func (r *Repository[T]) getIDField(entity *T) int64 {
	v := reflect.ValueOf(*entity)

	idField := v.FieldByName("ID")
	if idField.IsValid() && idField.CanInt() {
		return idField.Int()
	}

	return 0
}

func (r *Repository[T]) setIDField(entity *T, id int64) {
	v := reflect.ValueOf(entity).Elem()

	idField := v.FieldByName("ID")
	if idField.IsValid() && idField.CanSet() && idField.CanInt() {
		idField.SetInt(id)
	}
}
