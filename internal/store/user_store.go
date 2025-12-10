package store

import (
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type password struct {
	plaintext *string
	Hash      []byte
}

func (p *password) Set(plaintext string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintext), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	p.plaintext = &plaintext
	p.Hash = hash
	return nil
}
func (p *password) Matches(plaintext string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.Hash, []byte(plaintext))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

func (p *password) Clear() {
	p.plaintext = nil
	p.Hash = nil
}

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	Password     string    `json:"-"`
	PasswordHash password  `json:"-"`
	Bio          string    `json:"bio"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{db: db}
}

type UserStore interface {
	CreateUser(*User) (*User, error)
	GetUserByUserName(username string) (*User, error)
	UpdateUser(*User) error
	GetUserByID(id int64) (*User, error)
	DeleteUser(id int64) error
}

func (pg *PostgresUserStore) CreateUser(user *User) (*User, error) {
	query := `
		INSERT INTO users (username, email, password_hash, bio, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, created_at, updated_at
		`
	err := pg.db.QueryRow(query, user.Username, user.Email, user.PasswordHash.Hash, user.Bio).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (pg *PostgresUserStore) GetUserByUserName(username string) (*User, error) {
	user := &User{
		PasswordHash: password{},
	}
	query := `SELECT id, username, email, password_hash, bio, created_at, updated_at
			  FROM users
			  WHERE username = $1`
	err := pg.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash.Hash,
		&user.Bio,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (pg *PostgresUserStore) GetUserByID(id int64) (*User, error) {
	user := &User{
		PasswordHash: password{},
	}
	query := `SELECT id, username, email, password_hash, bio, created_at, updated_at
			  FROM users
			  WHERE id = $1`
	err := pg.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash.Hash,
		&user.Bio,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (pg *PostgresUserStore) UpdateUser(user *User) error {
	query := `
		UPDATE users
		SET username = $1, email = $2, bio = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
		Returning updated_at,id
	`
	result, err := pg.db.Exec(query, user.Username, user.Email, user.Bio, user.ID)
	if err != nil {
		return err
	}
	rowAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (pg *PostgresUserStore) DeleteUser(id int64) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := pg.db.Exec(query, id)
	return err
}
