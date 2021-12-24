package dao

import (
	"database/sql"

	"go-example/common/db"
	"go-example/model"
)

func GetUserList(page int, size int) (*sql.Rows, error) {
	return db.DB.Query(`SELECT id, name, password, status FROM users LIMIT ? OFFSET ?`, size, (page-1)*size)
}

func GetUser(id *string, name *string, password *string, status *int) error {
	query := `SELECT id, name, password, status FROM users WHERE id = ?`
	err := db.DB.QueryRow(query, id).Scan(id, name, password, status)
	return err
}

func CreateUser(user model.User) (sql.Result, error) {
	if user.ID != "" {
		return db.DB.Exec(`INSERT INTO users (id, name, password, status) VALUES (?, ?, ?, ?)`, user.ID, user.Name, user.Password, user.Status)
	} else {
		return db.DB.Exec(`INSERT INTO users (name, password, status) VALUES (?, ?, ?)`, user.Name, user.Password, user.Status)
	}
}

func UpdateUser(user model.User) (sql.Result, error) {
	return db.DB.Exec(`UPDATE users set name = ?, password =?, status=? where id=?`, user.Name, user.Password, user.Status, user.ID)
}

func DeleteUser(id string) (sql.Result, error) {
	return db.DB.Exec(`DELETE FROM users WHERE id = ?`, id)
}
