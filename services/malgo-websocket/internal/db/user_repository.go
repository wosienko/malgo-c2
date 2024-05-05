package db

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	if db == nil {
		panic("nil db")
	}

	return &UserRepository{db: db}
}

func (u *UserRepository) GetUserIdIfLoggedInAndOperator(ctx context.Context, sessionId string) (string, error) {
	// one query to check if sessionId is in the sessions table and the date is not expired and the user is an operator
	// if the user is an operator, return true
	row, err := u.db.QueryContext(ctx,
		`
				SELECT s.user_id FROM sessions s
				JOIN users u ON s.user_id = u.id
				JOIN roles r ON r.name = 'Operator'
				JOIN user_roles ur ON ur.user_id = u.id
				AND ur.role_id = r.id
				WHERE s.id = $1 AND expires_at > now()
	    `,
		sessionId,
	)
	if err != nil {
		return "", err
	}

	defer row.Close()

	if !row.Next() {
		return "", nil
	}
	var userId string
	if err := row.Scan(&userId); err != nil {
		return "", err
	}

	return userId, nil
}

func (u *UserRepository) IsUserAllowedToAccessProject(ctx context.Context, userId string, projectId string) (bool, error) {
	// one query to check if the user is an operator or has access to the project
	// if the user is an operator, return true
	row, err := u.db.QueryContext(ctx,
		`
				SELECT p.id FROM projects p
				JOIN user_projects pu ON pu.project_id = p.id
				JOIN users u ON u.id = pu.user_id
				WHERE u.id = $1 AND p.id = $2
	    `,
		userId, projectId,
	)
	if err != nil {
		return false, err
	}

	defer row.Close()

	if !row.Next() {
		return false, nil
	}

	return true, nil
}
