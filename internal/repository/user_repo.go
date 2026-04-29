package repository

import (
	"context"
	"database/sql"

	"WalletApps/internal/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create 插入用户数据
func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	query := `
        INSERT INTO users (id, name, password_hash, salt, created_at) 
        VALUES ($1, $2, $3, $4, NOW())
    `

	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Name,
		user.PasswordHash,
		user.Salt,
	)
	return err
}

// GetByUsername 根据用户名查询用户（用于注册查重或登录验证）
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	query := `SELECT id, name, password_hash, salt FROM users WHERE name = $1 LIMIT 1`

	var u model.User
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&u.ID, &u.Name, &u.PasswordHash, &u.Salt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // 用户不存在，返回 nil
	}
	if err != nil {
		return nil, err
	}

	return &u, nil
}
