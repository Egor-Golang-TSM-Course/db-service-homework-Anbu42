package psql

import (
	"context"
	"errors"
	"time"

	"blog/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (pgdb *PostgreSQL) Register(JwtSecretKey string, user *models.User) (string, error) {
	userID := uuid.New()

	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return "", err
	}

	query := `INSERT INTO users (id, name, login, password, created) VALUES ($1, $2, $3, $4, $5)`
	_, err = pgdb.pool.Exec(context.Background(), query, userID, user.Name, user.Login, hashedPassword, time.Now())
	if err != nil {
		return "", err
	}

	token, err := generateToken(JwtSecretKey, userID.String())
	if err != nil {
		return "", err
	}

	return token, nil
}

func (pgdb *PostgreSQL) Login(JwtSecretKey string, user *models.User) (string, error) {
	var userID uuid.UUID
	var hashedPassword string

	query := `SELECT id, password FROM users WHERE login = $1`
	err := pgdb.pool.QueryRow(context.Background(), query, user.Login).Scan(&userID, &hashedPassword)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Проверяем пароль с использованием bcrypt
	if !checkPassword(user.Password, hashedPassword) {
		return "", errors.New("invalid credentials")
	}

	token, err := generateToken(JwtSecretKey, userID.String())
	if err != nil {
		return "", err
	}

	return token, nil
}

// generateToken генерирует JWT токен
func generateToken(JwtSecretKey string, userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Токен действителен 24 часа
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(JwtSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// hashPassword хэширует пароль с использованием bcrypt
func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// checkPassword проверяет пароль с использованием bcrypt
func checkPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
