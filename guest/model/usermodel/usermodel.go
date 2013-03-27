package usermodel

import (
	"database/sql"
	"guest/model/sqlmodel"
)

func IsDuplicateUser(username string) bool {
	query := "SELECT id FROM user WHERE username=?"
	return sqlmodel.Query(query, username).Next()
}
func AddGuestUser(username, passwd string, permiss uint32) int64 {
	query := "INSERT user set username=?,passwd=?,permiss=?"
	return sqlmodel.NonQuery(query, username, passwd, permiss)
}
func CheckLoginUser(username, passwd string) *sql.Rows {
	query := "SELECT id, username, permiss FROM user WHERE username=? AND passwd=?"
	rows := sqlmodel.Query(query, username, passwd)
	return rows
}