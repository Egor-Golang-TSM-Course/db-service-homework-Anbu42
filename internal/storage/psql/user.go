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

// Register регистрирует нового пользователя
func (pgdb *PostgreSQL) Register(user *models.User) (string, error) {
	// Генерируем уникальный идентификатор пользователя
	userID := uuid.New()

	// Хэшируем пароль с использованием bcrypt
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return "", err
	}

	// Вставляем пользователя в базу данных
	query := `INSERT INTO users (id, name, login, password, created) VALUES ($1, $2, $3, $4, $5)`
	_, err = pgdb.pool.Exec(context.Background(), query, userID, user.Name, user.Login, hashedPassword, time.Now())
	if err != nil {
		return "", err
	}

	// Возвращаем токен
	token, err := generateToken(userID.String())
	if err != nil {
		return "", err
	}

	return token, nil
}

// Login выполняет вход пользователя
func (pgdb *PostgreSQL) Login(user *models.User) (string, error) {
	// Получаем данные пользователя из базы данных
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

	// Возвращаем токен
	token, err := generateToken(userID.String())
	if err != nil {
		return "", err
	}

	return token, nil
}

// generateToken генерирует JWT токен
func generateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Токен действителен 24 часа
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := []byte("your-secret-key") // Замените на свой секретный ключ

	tokenString, err := token.SignedString(secretKey)
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
