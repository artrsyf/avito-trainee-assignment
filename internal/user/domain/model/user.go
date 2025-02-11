package model

type User struct {
	ID           uint   `db:"id"`
	Username     string `db:"username"`
	Coins        uint   `db:"coins"`
	PasswordHash string `db:"password_hash"`
}
