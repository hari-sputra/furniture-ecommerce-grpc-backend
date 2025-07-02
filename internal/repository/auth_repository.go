package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/hari-sputra/furniture-ecommerce-grpc-be/internal/entity"
)

type IAuthRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	InsertUSer(ctx context.Context, user *entity.User) error
	UpdateUserPassword(ctx context.Context, updatePassword *entity.UpdateUserPassword) error
}

type authRepository struct {
	db *sql.DB
}

func (ar *authRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	row := ar.db.QueryRowContext(ctx, "SELECT id, email, password, full_name, role_code, created_at FROM \"user\" WHERE email = $1", email)

	if row.Err() != nil {
		return nil, row.Err()
	}

	var user entity.User
	err := row.Scan(
		&user.Id,
		&user.Email,
		&user.Password,
		&user.FullName,
		&user.RoleCode,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (ar *authRepository) InsertUSer(ctx context.Context, user *entity.User) error {
	_, err := ar.db.ExecContext(ctx, "INSERT INTO \"user\" (full_name, email, password, role_code, created_by) VALUES($1, $2, $3, $4, $5)", user.FullName, user.Email, user.Password, user.RoleCode, user.CreatedBy)

	if err != nil {
		return err
	}

	return nil
}

// UpdateUserPassword implements IAuthRepository.
func (ar *authRepository) UpdateUserPassword(ctx context.Context, updatePassword *entity.UpdateUserPassword) error {
	_, err := ar.db.ExecContext(ctx,
		"UPDATE \"user\" SET password = $1, updated_by = $2 WHERE id = $3",
		updatePassword.HashedNewPassword, updatePassword.UpdatedBy, updatePassword.UserId,
	)

	if err != nil {
		return err
	}

	return nil
}

func NewAuthRepository(db *sql.DB) IAuthRepository {
	return &authRepository{db: db}
}
