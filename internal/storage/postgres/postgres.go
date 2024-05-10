package postgres

import (
	"auth_service/internal/config"
	"auth_service/internal/domain/models"
	"auth_service/internal/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	_ "github.com/lib/pq" // PostgreSQL driver
)

type Storage struct {
	db *sql.DB
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// Конструктор Storage для PostgreSQL
func New(cfg config.DBConfig) (*Storage, error) {
	const op = "storage.postgres.New"
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &Storage{db: db}, nil
}

func (s *Storage) SaveUser(ctx context.Context, socialId string, email string, passHash []byte) (int64, error) {
	const op = "storage.postgres.SaveUser"

	// Запрос на добавление пользователя для таблица Users
	query := "INSERT INTO users(social_club_id,email, password_hash) VALUES($1, $2, $3)"
	var id int64

	err := s.db.QueryRowContext(ctx, query, socialId, email, passHash).Scan(&id)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" { // 23505 - код ошибки unique_violation в PostgreSQL
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

// User возвращает пользователя по email.
func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "storage.postgres.User"

	query := "SELECT id, email, pass_hash FROM users WHERE email = $1"

	row := s.db.QueryRowContext(ctx, query, email)

	var user models.User
	err := row.Scan(&user.ID, &user.Email, &user.PassHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}

		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

// App возвращает приложение по его ID.
func (s *Storage) App(ctx context.Context) (models.App, error) {
	const op = "storage.postgres.App"

	query := "SELECT name, secret FROM apps WHERE id = $1"

	row := s.db.QueryRowContext(ctx, query)

	var app models.App
	err := row.Scan(&app.ID, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}

		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}
