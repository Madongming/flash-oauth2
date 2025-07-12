// Package database provides PostgreSQL database connection and migration functionality.
// It handles database initialization, connection pooling, and automatic schema migration
// for the OAuth2 server.
package database

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

// Init creates and configures a new PostgreSQL database connection.
// It sets up connection pooling parameters and tests the connection.
//
// Parameters:
//   - databaseURL: PostgreSQL connection string (e.g., "postgres://user:pass@host:5432/dbname?sslmode=disable")
//
// Returns:
//   - *sql.DB: Database connection pool
//   - error: Connection or configuration error
//
// The connection pool is configured with:
//   - Max open connections: 25
//   - Max idle connections: 25
//   - Connection max lifetime: 5 minutes
func Init(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	// 设置连接池
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(time.Minute * 5)

	// 测试连接
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// Migrate creates all necessary database tables and inserts default data.
// This function is idempotent and can be safely run multiple times.
//
// Created tables:
//   - users: User accounts with phone numbers
//   - oauth_clients: Registered OAuth2 client applications
//   - auth_codes: Short-lived authorization codes
//   - access_tokens: Access token records (for audit)
//   - refresh_tokens: Long-lived refresh tokens
//
// It also inserts a default OAuth2 client with ID "default-client" for development.
//
// Parameters:
//   - db: Database connection
//
// Returns:
//   - error: Migration error if any table creation fails
func Migrate(db *sql.DB) error {
	// 用户表
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		phone VARCHAR(20) UNIQUE NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	// OAuth2客户端表
	createOAuthClientsTable := `
	CREATE TABLE IF NOT EXISTS oauth_clients (
		id VARCHAR(255) PRIMARY KEY,
		secret VARCHAR(255) NOT NULL,
		name VARCHAR(255) NOT NULL,
		redirect_uris TEXT[] NOT NULL,
		grant_types TEXT[] NOT NULL,
		response_types TEXT[] NOT NULL,
		scope VARCHAR(255) DEFAULT 'openid profile',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	// 授权码表
	createAuthCodesTable := `
	CREATE TABLE IF NOT EXISTS auth_codes (
		code VARCHAR(255) PRIMARY KEY,
		client_id VARCHAR(255) NOT NULL,
		user_id INTEGER NOT NULL,
		redirect_uri VARCHAR(512) NOT NULL,
		scope VARCHAR(255),
		expires_at TIMESTAMP NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (client_id) REFERENCES oauth_clients(id),
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`

	// 访问令牌表
	createAccessTokensTable := `
	CREATE TABLE IF NOT EXISTS access_tokens (
		token VARCHAR(512) PRIMARY KEY,
		client_id VARCHAR(255) NOT NULL,
		user_id INTEGER NOT NULL,
		scope VARCHAR(255),
		expires_at TIMESTAMP NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (client_id) REFERENCES oauth_clients(id),
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`

	// 刷新令牌表
	createRefreshTokensTable := `
	CREATE TABLE IF NOT EXISTS refresh_tokens (
		token VARCHAR(512) PRIMARY KEY,
		client_id VARCHAR(255) NOT NULL,
		user_id INTEGER NOT NULL,
		scope VARCHAR(255),
		expires_at TIMESTAMP NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (client_id) REFERENCES oauth_clients(id),
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`

	// 开发者表
	createDevelopersTable := `
	CREATE TABLE IF NOT EXISTS developers (
		id VARCHAR(255) PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		phone VARCHAR(20),
		status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'suspended')),
		api_quota INTEGER DEFAULT 10000,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	// 外部应用表
	createExternalAppsTable := `
	CREATE TABLE IF NOT EXISTS external_apps (
		id VARCHAR(255) PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		description TEXT,
		developer_id VARCHAR(255) NOT NULL,
		status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'suspended', 'revoked')),
		callback_url VARCHAR(512) NOT NULL,
		scopes VARCHAR(512) DEFAULT 'openid profile',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		revoked_at TIMESTAMP,
		FOREIGN KEY (developer_id) REFERENCES developers(id)
	);`

	// 应用密钥对表
	createAppKeyPairsTable := `
	CREATE TABLE IF NOT EXISTS app_key_pairs (
		id VARCHAR(255) PRIMARY KEY,
		app_id VARCHAR(255) NOT NULL,
		key_id VARCHAR(255) UNIQUE NOT NULL,
		private_key TEXT NOT NULL,
		public_key TEXT NOT NULL,
		algorithm VARCHAR(10) DEFAULT 'RS256' CHECK (algorithm IN ('RS256', 'RS384', 'RS512')),
		status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'expired', 'revoked')),
		expires_at TIMESTAMP,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		revoked_at TIMESTAMP,
		last_used_at TIMESTAMP,
		FOREIGN KEY (app_id) REFERENCES external_apps(id)
	);`

	// 执行所有表创建语句
	tables := []string{
		createUsersTable,
		createOAuthClientsTable,
		createAuthCodesTable,
		createAccessTokensTable,
		createRefreshTokensTable,
		createDevelopersTable,
		createExternalAppsTable,
		createAppKeyPairsTable,
	}

	for _, table := range tables {
		if _, err := db.Exec(table); err != nil {
			return err
		}
	}

	// 插入默认客户端
	insertDefaultClient := `
	INSERT INTO oauth_clients (id, secret, name, redirect_uris, grant_types, response_types, scope) 
	VALUES ('default-client', 'default-secret', 'Default Client', 
			ARRAY['http://localhost:3000/callback'], 
			ARRAY['authorization_code', 'refresh_token'], 
			ARRAY['code'], 
			'openid profile email phone')
	ON CONFLICT (id) DO NOTHING;`

	_, err := db.Exec(insertDefaultClient)
	return err
}
