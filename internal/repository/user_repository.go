package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/srinivasarynh/age_calculator/internal/models"
	"go.uber.org/zap"
)

type UserRepository interface {
	Create(ctx context.Context, name string, dob time.Time) (*models.User, error)
	GetById(ctx context.Context, id int32) (*models.User, error)
	List(ctx context.Context, limit, offset int32) ([]models.User, error)
	Update(ctx context.Context, id int32, name string, dob time.Time) (*models.User, error)
	Delete(ctx context.Context, id int32) error
	Count(ctx context.Context) (int64, error)
}

type userRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewUserRepository(db *sql.DB, logger *zap.Logger) UserRepository {
	return &userRepository{
		db:     db,
		logger: logger,
	}
}

func (r *userRepository) Create(ctx context.Context, name string, dob time.Time) (*models.User, error) {
	query := `INSERT INTO users (name, dob) VALUES ($1, $2) RETURNING id, name, dob, created_at, updated_at`

	var user models.User
	err := r.db.QueryRowContext(ctx, query, name, dob).Scan(
		&user.ID,
		&user.Name,
		&user.DOB,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		r.logger.Error("Failed to create user", zap.Error(err))
		return nil, err
	}

	r.logger.Info("User created", zap.Int32("id", user.ID))
	return &user, nil
}

func (r *userRepository) GetById(ctx context.Context, id int32) (*models.User, error) {
	query := `SELECT id, name, dob, created_at, updated_at FROM users WHERE id = $1`

	var user models.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.DOB,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		r.logger.Error("Failed to get user", zap.Error(err), zap.Int32("id", id))
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) List(ctx context.Context, limit, offset int32) ([]models.User, error) {
	query := `SELECT id, name, dob, created_at, updated_at FROM users ORDER BY id LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("Failed to list users", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	users := make([]models.User, 0)
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.DOB, &user.CreatedAt, &user.UpdatedAt); err != nil {
			r.logger.Error("Failed to scan user", zap.Error(err))
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepository) Update(ctx context.Context, id int32, name string, dob time.Time) (*models.User, error) {
	query := `UPDATE users SET name = $1, dob = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3 RETURNING id, name, dob, created_at, updated_at`

	var user models.User
	err := r.db.QueryRowContext(ctx, query, name, dob, id).Scan(
		&user.ID,
		&user.Name,
		&user.DOB,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		r.logger.Error("Failed to update user", zap.Error(err), zap.Int32("id", id))
		return nil, err
	}

	r.logger.Info("User updated", zap.Int32("id", user.ID))
	return &user, nil
}

func (r *userRepository) Delete(ctx context.Context, id int32) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.Error("Failed to delete user", zap.Error(err), zap.Int32("id", id))
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	r.logger.Info("User deleted", zap.Int32("id", id))
	return nil
}

func (r *userRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM users`

	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		r.logger.Error("Failed to count users", zap.Error(err))
		return 0, err
	}

	return count, nil
}
