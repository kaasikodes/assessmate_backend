package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/kaasikodes/assessmate_backend/internal/core/domain/user"
)

func (r *MySqlRepo) CreateUser(ctx context.Context, u *user.User) (*user.User, error) {
	query := `
		INSERT INTO users (name, email, status, password, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	res, err := r.db.ExecContext(ctx, query, u.GetName().String(), u.GetEmail().String(), u.GetStatus().String(), u.PasswordHash(), u.GetCreatedAt(), u.GetUpdatedAt())
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	parsedId, err := user.NewId(int(id))
	if err != nil {
		return nil, err
	}
	u.SetId(parsedId)

	return u, nil
}
func (r *MySqlRepo) UpdateUserPassword(ctx context.Context, u *user.User) (*user.User, error) {
	query := `UPDATE users SET password = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, u.PasswordHash(), u.GetUpdatedAt(), u.GetId().Value())
	return u, err
}
func (r *MySqlRepo) VerifyUser(ctx context.Context, u *user.User) (*user.User, error) {
	query := `UPDATE users SET verified_at = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, u.GetVerifiedAt(), u.GetUpdatedAt(), u.GetId().Value())
	if err != nil {
		return nil, err
	}
	return u, nil
}
func (r *MySqlRepo) ChangeUserStatus(ctx context.Context, userId user.Id, status user.UserStatus) (*user.User, error) {
	query := `UPDATE users SET status = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, status.String(), time.Now(), userId.Value())
	if err != nil {
		return nil, err
	}
	return r.GetUserById(ctx, userId)
}

func (r *MySqlRepo) CreateToken(ctx context.Context, value user.TokenValue, tokenType user.TokenType, userId user.Id) (token *user.Token, err error) {
	query := `INSERT INTO tokens (value, type, user_id, created_at) VALUES (?, ?, ?)`
	now := time.Now()
	res, err := r.db.ExecContext(ctx, query, value.String(), tokenType.String(), userId, now)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()

	parsedId, err := user.NewId(int(id))
	if err != nil {
		return nil, err
	}
	token, err = user.NewToken(value, tokenType, userId)
	if err != nil {
		return nil, err
	}
	token.SetId(parsedId)

	return token, nil
}

func (r *MySqlRepo) DeleteToken(ctx context.Context, id user.Id) error {
	query := `DELETE FROM tokens WHERE user_id = ?`
	_, err := r.db.ExecContext(ctx, query, id.Value())
	return err
}
func (r *MySqlRepo) GetToken(ctx context.Context, userId user.Id, value user.TokenValue) (*user.Token, error) {
	query := `SELECT id, user_id, value, type, created_at, expires_at FROM tokens WHERE user_id = ? AND value = ?`
	row := r.db.QueryRowContext(ctx, query, userId.Value(), string(value))

	var t user.Token
	var createdAt, expiresAt sql.NullTime
	var uid, tid int
	var val, typ string

	err := row.Scan(&tid, &uid, &val, &typ, &createdAt, &expiresAt)
	if err != nil {
		return nil, err
	}
	tokenId, err := user.NewId(tid)
	if err != nil {
		return nil, err
	}

	t.SetId(tokenId)
	t.SetUserId(userId)
	t.SetValue(value)
	t.SetCreatedAt(createdAt.Time)
	t.SetExpiresAt(expiresAt.Time)
	t.SetType(typ)

	return &t, nil
}

func (r *MySqlRepo) GetUserById(ctx context.Context, userId user.Id) (*user.User, error) {
	query := `SELECT id, name, email, status, password, created_at, updated_at, verified_at, deleted_at FROM users WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, userId.Value())
	return r.scanUser(row)
}
func (r *MySqlRepo) GetUserByEmail(ctx context.Context, email user.Email) (*user.User, error) {
	query := `SELECT id, name, email, status, password, created_at, updated_at, verified_at, deleted_at FROM users WHERE email = ?`
	row := r.db.QueryRowContext(ctx, query, email.String())
	return r.scanUser(row)
}

func (r *MySqlRepo) GetUsers(ctx context.Context, filter *user.UserFilter) ([]user.User, int, error) {
	baseQuery := `FROM users`
	var whereClause string
	var args []interface{}

	if filter != nil && filter.Status != nil {
		whereClause = ` WHERE status = ?`
		args = append(args, filter.Status.String())
	}

	// Count query
	countQuery := `SELECT COUNT(*) ` + baseQuery + whereClause
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Data query
	selectQuery := `SELECT id, name, email, status, password, created_at, updated_at, verified_at, deleted_at ` + baseQuery + whereClause
	rows, err := r.db.QueryContext(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []user.User
	for rows.Next() {
		u, err := r.scanUser(rows)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, *u)
	}

	return users, total, nil
}

func (r *MySqlRepo) scanUser(scanner interface {
	Scan(dest ...interface{}) error
}) (*user.User, error) {
	var (
		id         int
		name       string
		email      string
		status     string
		password   string
		createdAt  time.Time
		updatedAt  time.Time
		verifiedAt sql.NullTime
		deletedAt  sql.NullTime
	)

	err := scanner.Scan(&id, &name, &email, &status, &password, &createdAt, &updatedAt, &verifiedAt, &deletedAt)
	if err != nil {
		return nil, err
	}
	userName, err := user.NewName(name)
	if err != nil {
		return nil, err
	}
	userEmail, err := user.NewEmail(email)
	if err != nil {
		return nil, err
	}
	u, err := user.NewUser(userName, userEmail)
	if err != nil {
		return nil, err
	}
	uid, err := user.NewId(id)
	if err != nil {
		return nil, err
	}
	u.SetId(uid)

	return u, nil
}
