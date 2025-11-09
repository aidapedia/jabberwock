package user

const (
	queryFindByID = `
		SELECT id, name, phone, email, password, type, status, is_verified, avatar_url, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	queryFindByPhone = `
		SELECT id, name, phone, email, password, type, status, is_verified, avatar_url, created_at, updated_at
		FROM users
		WHERE phone = $1
	`

	queryFindByEmail = `
		SELECT id, name, phone, email, password, type, status, is_verified, avatar_url, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	queryUpdateStatus = `
		UPDATE users
		SET status = $1, updated_at = $2
		WHERE id = $3
	`

	queryCreateUser = `
		INSERT INTO users (name, phone, email, password, type, status, is_verified, avatar_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
)
