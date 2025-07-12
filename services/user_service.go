// Package services provides business logic services for the OAuth2 server.
// It contains services for user management, OAuth2 operations, and JWT handling.
package services

import (
	"context"
	"crypto/rand"
	"database/sql"
	"flash-oauth2/models"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// UserService handles user authentication and management operations.
// It provides methods for verification code handling and user account management.
type UserService struct {
	db    *sql.DB       // Database connection for persistent user data
	redis *redis.Client // Redis client for temporary data (verification codes)
}

// NewUserService creates a new UserService instance with database and Redis connections.
//
// Parameters:
//   - db: Database connection for user data persistence
//   - redis: Redis client for temporary data storage
//
// Returns:
//   - *UserService: Configured user service instance
func NewUserService(db *sql.DB, redis *redis.Client) *UserService {
	return &UserService{
		db:    db,
		redis: redis,
	}
}

// SendVerificationCode generates and sends a 6-digit verification code to the specified phone number.
// The code is stored in Redis with a 5-minute expiration time.
//
// In a production environment, this method should integrate with an SMS service
// to send the verification code via text message. For development purposes,
// the code is printed to the server console.
//
// Parameters:
//   - phone: The phone number to send the verification code to
//
// Returns:
//   - error: An error if Redis operations fail, nil otherwise
//
// Example:
//
//	err := userService.SendVerificationCode("13800138000")
func (s *UserService) SendVerificationCode(phone string) error {
	ctx := context.Background()

	// 生成6位验证码
	code := generateVerificationCode()

	// 存储到Redis，有效期5分钟
	key := fmt.Sprintf("verification_code:%s", phone)
	err := s.redis.Set(ctx, key, code, 5*time.Minute).Err()
	if err != nil {
		return err
	}

	// 在实际应用中，这里应该调用短信服务发送验证码
	// 为了演示，我们只是打印验证码
	fmt.Printf("Verification code for %s: %s\n", phone, code)

	return nil
}

// VerifyCode validates a verification code against the one stored in Redis for the given phone number.
// The method checks if the provided code matches the stored code and if it hasn't expired.
// If valid, it returns the user associated with the phone number or creates a new user.
//
// Parameters:
//   - phone: The phone number associated with the verification code
//   - code: The verification code to validate
//
// Returns:
//   - *models.User: The authenticated user if verification succeeds
//   - error: An error if verification fails or database operations fail
//
// Example:
//
//	user, err := userService.VerifyCode("13800138000", "123456")
func (s *UserService) VerifyCode(phone, code string) (*models.User, error) {
	ctx := context.Background()

	// 从Redis获取验证码
	key := fmt.Sprintf("verification_code:%s", phone)
	storedCode, err := s.redis.Get(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("verification code not found or expired")
	}

	if storedCode != code {
		return nil, fmt.Errorf("invalid verification code")
	}

	// 验证成功，删除验证码
	s.redis.Del(ctx, key)

	// 查找或创建用户（幂等操作）
	user, err := s.findOrCreateUser(phone)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// findOrCreateUser searches for a user by phone number and creates one if not found.
// This method implements the user registration flow for phone-based authentication.
//
// Parameters:
//   - phone: The phone number to search for or create a user with
//
// Returns:
//   - *models.User: The found or newly created user
//   - error: An error if database operations fail
func (s *UserService) findOrCreateUser(phone string) (*models.User, error) {
	// 首先尝试查找用户
	user := &models.User{}
	err := s.db.QueryRow("SELECT id, phone, created_at, updated_at FROM users WHERE phone = $1", phone).
		Scan(&user.ID, &user.Phone, &user.CreatedAt, &user.UpdatedAt)

	if err == nil {
		// 用户存在，更新最后登录时间
		_, err = s.db.Exec("UPDATE users SET updated_at = CURRENT_TIMESTAMP WHERE id = $1", user.ID)
		return user, err
	}

	if err != sql.ErrNoRows {
		return nil, err
	}

	// 用户不存在，创建新用户
	err = s.db.QueryRow(
		"INSERT INTO users (phone) VALUES ($1) RETURNING id, phone, created_at, updated_at",
		phone,
	).Scan(&user.ID, &user.Phone, &user.CreatedAt, &user.UpdatedAt)

	return user, err
}

// GetUserByID retrieves a user from the database by their unique ID.
//
// Parameters:
//   - userID: The unique identifier of the user to retrieve
//
// Returns:
//   - *models.User: The user if found
//   - error: An error if the user is not found or database operations fail
//
// Example:
//
//	user, err := userService.GetUserByID(123)
func (s *UserService) GetUserByID(userID int) (*models.User, error) {
	user := &models.User{}
	err := s.db.QueryRow("SELECT id, phone, created_at, updated_at FROM users WHERE id = $1", userID).
		Scan(&user.ID, &user.Phone, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// generateVerificationCode creates a random 6-digit numeric verification code.
// This helper function uses the math/rand package to generate codes for SMS verification.
//
// Returns:
//   - string: A 6-digit numeric verification code (e.g., "123456")
func generateVerificationCode() string {
	// 生成6位数字验证码
	b := make([]byte, 3)
	rand.Read(b)

	code := ""
	for _, v := range b {
		code += fmt.Sprintf("%02d", v%100)
	}

	return code[:6]
}
