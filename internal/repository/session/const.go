package session

const (
	queryCreateActiveSession = `
	INSERT INTO sessions (token, user_id, user_agent, ip, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id
	`

	queryFindActiveSessionByTokenID = `
	SELECT id, token, user_id, user_agent, ip, created_at, updated_at
	FROM sessions
	WHERE token = $1
	`

	queryDeleteActiveSessionByTokenID = `
	DELETE FROM sessions
	WHERE token = $1
	`

	queryUpdateRefreshDateByTokenID = `
	UPDATE sessions
	SET updated_at = $2
	WHERE token = $1
	`
)
