package linkfy

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"github.com/yantology/linkfy/pkg/customerror"
)

type linkfyPostgres struct {
	db *sql.DB
}

func NewLinkfyPostgres(db *sql.DB) LinkfyDBInterface {
	return &linkfyPostgres{db: db}
}

// CreateLinkfy creates a new linkfy profile
func (l *linkfyPostgres) CreateLinkfy(linkfy *LinkfyCreated) *customerror.CustomError {
	query := `INSERT INTO linkfy (user_id, username, avatar_url, name, bio) 
			  VALUES ($1, $2, $3, $4, $5)`

	_, err := l.db.Exec(
		query,
		linkfy.UserID,
		linkfy.Username,
		linkfy.AvatarURL,
		linkfy.Name,
		linkfy.Bio,
	)

	if err != nil {
		// Check for specific PostgreSQL errors
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return customerror.NewCustomError(err, "Username already exists", 409)
		}
		return customerror.NewPostgresError(err)
	}

	return nil
}

// GetLinkfyByID retrieves a linkfy profile by its ID
func (l *linkfyPostgres) GetLinkfyByID(id uuid.UUID) (*Linkfy, *customerror.CustomError) {
	query := `SELECT id, user_id, username, avatar_url, name, bio, created_at, updated_at 
			  FROM linkfy 
			  WHERE id = $1`

	linkfy := &Linkfy{}
	err := l.db.QueryRow(query, id).Scan(
		&linkfy.ID,
		&linkfy.UserID,
		&linkfy.Username,
		&linkfy.AvatarURL,
		&linkfy.Name,
		&linkfy.Bio,
		&linkfy.CreatedAt,
		&linkfy.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, customerror.NewCustomError(err, "Linkfy profile not found", 404)
		}
		return nil, customerror.NewPostgresError(err)
	}

	return linkfy, nil
}

// GetLinkfyByUsername retrieves a linkfy profile by username
func (l *linkfyPostgres) GetLinkfyByUsername(username string) (*Linkfy, *customerror.CustomError) {
	query := `SELECT id, user_id, username, avatar_url, name, bio, created_at, updated_at 
			  FROM linkfy 
			  WHERE username = $1`

	linkfy := &Linkfy{}
	err := l.db.QueryRow(query, username).Scan(
		&linkfy.ID,
		&linkfy.UserID,
		&linkfy.Username,
		&linkfy.AvatarURL,
		&linkfy.Name,
		&linkfy.Bio,
		&linkfy.CreatedAt,
		&linkfy.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, customerror.NewCustomError(err, "Linkfy profile not found", 404)
		}
		return nil, customerror.NewPostgresError(err)
	}

	return linkfy, nil
}

// GetAllLinkfyByUserID retrieves all linkfy profiles for a specific user
func (l *linkfyPostgres) GetAllLinkfyByUserID(userID string) ([]*Linkfy, *customerror.CustomError) {
	query := `SELECT id, user_id, username, avatar_url, name, bio, created_at, updated_at 
			  FROM linkfy 
			  WHERE user_id = $1
			  ORDER BY created_at DESC`

	rows, err := l.db.Query(query, userID)
	if err != nil {
		return nil, customerror.NewPostgresError(err)
	}
	defer rows.Close()

	var linkfyProfiles []*Linkfy
	for rows.Next() {
		linkfy := &Linkfy{}
		err := rows.Scan(
			&linkfy.ID,
			&linkfy.UserID,
			&linkfy.Username,
			&linkfy.AvatarURL,
			&linkfy.Name,
			&linkfy.Bio,
			&linkfy.CreatedAt,
			&linkfy.UpdatedAt,
		)
		if err != nil {
			return nil, customerror.NewPostgresError(err)
		}
		linkfyProfiles = append(linkfyProfiles, linkfy)
	}

	if err = rows.Err(); err != nil {
		return nil, customerror.NewPostgresError(err)
	}

	return linkfyProfiles, nil
}

// UpdateLinkfy updates an existing linkfy profile
func (l *linkfyPostgres) UpdateLinkfy(linkfy *LinkfyUpdated) *customerror.CustomError {
	query := `UPDATE linkfy 
			  SET username = $1, avatar_url = $2, name = $3, bio = $4
			  WHERE id = $5 `

	result, err := l.db.Exec(
		query,
		linkfy.Username,
		linkfy.AvatarURL,
		linkfy.Name,
		linkfy.Bio,
		linkfy.ID,
	)

	if err != nil {
		// Check for specific PostgreSQL errors
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return customerror.NewCustomError(err, "Username already exists", 409)
		}
		return customerror.NewPostgresError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return customerror.NewPostgresError(err)
	}

	if rowsAffected == 0 {
		return customerror.NewCustomError(fmt.Errorf("no rows affected"), "Linkfy profile not found", 404)
	}

	return nil
}

// CheckUsernameExists checks if a username already exists (for debounce)
func (l *linkfyPostgres) CheckUsernameExists(username string) *customerror.CustomError {
	query := `SELECT EXISTS(SELECT 1 FROM linkfy WHERE username = $1)`

	var exists bool
	err := l.db.QueryRow(query, username).Scan(&exists)
	if err != nil {
		return customerror.NewPostgresError(err)
	}

	return nil
}

func (l *linkfyPostgres) CheckUsernameNotExists(username string) *customerror.CustomError {
	query := `SELECT COUNT(*) FROM linkfy WHERE username = $1`

	var count int
	err := l.db.QueryRow(query, username).Scan(&count)
	if err != nil {
		return customerror.NewPostgresError(err)
	}

	if count > 0 {
		return customerror.NewCustomError(fmt.Errorf("username already exists"), "Username already exists", 409)
	}

	return nil
}

// DeleteLinkfy deletes a linkfy profile by its ID and user ID
func (l *linkfyPostgres) DeleteLinkfy(id uuid.UUID, userID string) *customerror.CustomError {
	query := `DELETE FROM linkfy WHERE id = $1 AND user_id = $2`

	result, err := l.db.Exec(query, id, userID)
	if err != nil {
		return customerror.NewPostgresError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return customerror.NewPostgresError(err)
	}

	if rowsAffected == 0 {
		return customerror.NewCustomError(
			fmt.Errorf("no rows affected"),
			"Linkfy profile not found or unauthorized",
			404,
		)
	}

	return nil
}
