package integration

import (
	"database/sql"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

var (
	DB          *sql.DB
	RedisClient *redis.Client
)

func SetupTestData(t *testing.T, db *sql.DB) {
	_, err := db.Exec(`
		DELETE FROM users;
		DELETE FROM purchase_types;
		DELETE FROM purchases;
		DELETE FROM transactions;
	`)
	require.NoError(t, err)
}

func CreateTestUser(t *testing.T, username string, coins uint) uint {
	var id uint
	err := DB.QueryRow(
		"INSERT INTO users (username, coins, password_hash) VALUES ($1, $2, 'hash') RETURNING id",
		username, coins,
	).Scan(&id)
	require.NoError(t, err)
	return id
}

func CreatePurchaseType(t *testing.T, db *sql.DB, name string, cost uint) {
	_, err := db.Exec(
		"INSERT INTO purchase_types (name, cost) VALUES ($1, $2) ON CONFLICT DO NOTHING",
		name, cost,
	)
	require.NoError(t, err)
}
