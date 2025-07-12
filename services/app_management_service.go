// Package services provides application and developer management services.
package services

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"flash-oauth2/models"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// AppManagementService provides services for managing external applications and their keys
type AppManagementService struct {
	db *sql.DB
}

// NewAppManagementService creates a new instance of AppManagementService
func NewAppManagementService(db *sql.DB) *AppManagementService {
	return &AppManagementService{db: db}
}

// RegisterDeveloper registers a new developer on the platform
func (s *AppManagementService) RegisterDeveloper(name, email, phone string) (*models.Developer, error) {
	developer := &models.Developer{
		ID:        uuid.New().String(),
		Name:      name,
		Email:     email,
		Phone:     phone,
		Status:    "active",
		APIQuota:  10000, // Default quota
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := s.db.Exec(`
		INSERT INTO developers (id, name, email, phone, status, api_quota, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, developer.ID, developer.Name, developer.Email, developer.Phone,
		developer.Status, developer.APIQuota, developer.CreatedAt, developer.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to register developer: %w", err)
	}

	return developer, nil
}

// RegisterExternalApp registers a new external application
func (s *AppManagementService) RegisterExternalApp(developerID, name, description, callbackURL, scopes string) (*models.ExternalApp, error) {
	app := &models.ExternalApp{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		DeveloperID: developerID,
		Status:      "active",
		CallbackURL: callbackURL,
		Scopes:      scopes,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err := s.db.Exec(`
		INSERT INTO external_apps (id, name, description, developer_id, status, callback_url, scopes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, app.ID, app.Name, app.Description, app.DeveloperID, app.Status,
		app.CallbackURL, app.Scopes, app.CreatedAt, app.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to register external app: %w", err)
	}

	return app, nil
}

// GenerateKeyPair generates a new RSA key pair for an application
func (s *AppManagementService) GenerateKeyPair(appID string, algorithm string, expiresAt *time.Time) (*models.AppKeyPair, error) {
	// Generate RSA key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("failed to generate RSA key: %w", err)
	}

	// Convert to PEM format
	privateKeyPEM, err := privateKeyToPEM(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encode private key: %w", err)
	}

	publicKeyPEM, err := publicKeyToPEM(&privateKey.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encode public key: %w", err)
	}

	keyPair := &models.AppKeyPair{
		ID:         uuid.New().String(),
		AppID:      appID,
		KeyID:      fmt.Sprintf("key_%s_%d", appID, time.Now().Unix()),
		PrivateKey: privateKeyPEM,
		PublicKey:  publicKeyPEM,
		Algorithm:  algorithm,
		Status:     "active",
		ExpiresAt:  expiresAt,
		CreatedAt:  time.Now(),
	}

	_, err = s.db.Exec(`
		INSERT INTO app_key_pairs (id, app_id, key_id, private_key, public_key, algorithm, status, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, keyPair.ID, keyPair.AppID, keyPair.KeyID, keyPair.PrivateKey, keyPair.PublicKey,
		keyPair.Algorithm, keyPair.Status, keyPair.ExpiresAt, keyPair.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to save key pair: %w", err)
	}

	return keyPair, nil
}

// GetAppKeyPairs retrieves all key pairs for an application
func (s *AppManagementService) GetAppKeyPairs(appID string) ([]*models.AppKeyPair, error) {
	rows, err := s.db.Query(`
		SELECT id, app_id, key_id, private_key, public_key, algorithm, status, 
			   expires_at, created_at, revoked_at, last_used_at
		FROM app_key_pairs 
		WHERE app_id = $1 
		ORDER BY created_at DESC
	`, appID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keyPairs []*models.AppKeyPair
	for rows.Next() {
		kp := &models.AppKeyPair{}
		err := rows.Scan(&kp.ID, &kp.AppID, &kp.KeyID, &kp.PrivateKey, &kp.PublicKey,
			&kp.Algorithm, &kp.Status, &kp.ExpiresAt, &kp.CreatedAt, &kp.RevokedAt, &kp.LastUsedAt)
		if err != nil {
			return nil, err
		}
		keyPairs = append(keyPairs, kp)
	}

	return keyPairs, nil
}

// RevokeKeyPair revokes a key pair
func (s *AppManagementService) RevokeKeyPair(keyID string) error {
	now := time.Now()
	_, err := s.db.Exec(`
		UPDATE app_key_pairs 
		SET status = 'revoked', revoked_at = $1 
		WHERE key_id = $2
	`, now, keyID)

	if err != nil {
		return fmt.Errorf("failed to revoke key pair: %w", err)
	}

	return nil
}

// GetDeveloperApps retrieves all applications for a developer
func (s *AppManagementService) GetDeveloperApps(developerID string) ([]*models.ExternalApp, error) {
	rows, err := s.db.Query(`
		SELECT id, name, description, developer_id, status, callback_url, scopes, 
			   created_at, updated_at, revoked_at
		FROM external_apps 
		WHERE developer_id = $1 
		ORDER BY created_at DESC
	`, developerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apps []*models.ExternalApp
	for rows.Next() {
		app := &models.ExternalApp{}
		err := rows.Scan(&app.ID, &app.Name, &app.Description, &app.DeveloperID,
			&app.Status, &app.CallbackURL, &app.Scopes, &app.CreatedAt, &app.UpdatedAt, &app.RevokedAt)
		if err != nil {
			return nil, err
		}
		apps = append(apps, app)
	}

	return apps, nil
}

// GetAllApps retrieves all applications (admin function)
func (s *AppManagementService) GetAllApps() ([]*models.ExternalApp, error) {
	rows, err := s.db.Query(`
		SELECT ea.id, ea.name, ea.description, ea.developer_id, ea.status, ea.callback_url, 
			   ea.scopes, ea.created_at, ea.updated_at, ea.revoked_at,
			   d.name as developer_name
		FROM external_apps ea
		JOIN developers d ON ea.developer_id = d.id
		ORDER BY ea.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apps []*models.ExternalApp
	for rows.Next() {
		app := &models.ExternalApp{}
		var developerName string
		err := rows.Scan(&app.ID, &app.Name, &app.Description, &app.DeveloperID,
			&app.Status, &app.CallbackURL, &app.Scopes, &app.CreatedAt, &app.UpdatedAt,
			&app.RevokedAt, &developerName)
		if err != nil {
			return nil, err
		}
		// Store developer name in description for display purposes
		app.Description = fmt.Sprintf("Developer: %s | %s", developerName, app.Description)
		apps = append(apps, app)
	}

	return apps, nil
}

// UpdateKeyLastUsed updates the last used timestamp for a key
func (s *AppManagementService) UpdateKeyLastUsed(keyID string) error {
	_, err := s.db.Exec(`
		UPDATE app_key_pairs 
		SET last_used_at = CURRENT_TIMESTAMP 
		WHERE key_id = $1
	`, keyID)
	return err
}

// Helper functions for PEM encoding
func privateKeyToPEM(key *rsa.PrivateKey) (string, error) {
	keyBytes, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return "", err
	}

	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: keyBytes,
	})

	return string(keyPEM), nil
}

func publicKeyToPEM(key *rsa.PublicKey) (string, error) {
	keyBytes, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return "", err
	}

	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: keyBytes,
	})

	return string(keyPEM), nil
}

// GetKeyPairByKeyID retrieves a key pair by its key ID
func (s *AppManagementService) GetKeyPairByKeyID(keyID string) (*models.AppKeyPair, error) {
	kp := &models.AppKeyPair{}
	err := s.db.QueryRow(`
		SELECT id, app_id, key_id, private_key, public_key, algorithm, status, 
			   expires_at, created_at, revoked_at, last_used_at
		FROM app_key_pairs 
		WHERE key_id = $1
	`, keyID).Scan(&kp.ID, &kp.AppID, &kp.KeyID, &kp.PrivateKey, &kp.PublicKey,
		&kp.Algorithm, &kp.Status, &kp.ExpiresAt, &kp.CreatedAt, &kp.RevokedAt, &kp.LastUsedAt)

	if err != nil {
		return nil, err
	}

	return kp, nil
}
