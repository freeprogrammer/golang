package contentmodel

import (
	"database/sql"
	"guest/model/sqlmodel"
)

func AddGuest(content string, userID uint32) int64 {
	query := "INSERT content set content=?,userID=?"
	return sqlmodel.NonQuery(query, content, userID)
}
func GetGuest(offset, perpage int) *sql.Rows {
	query := "SELECT content.id,content,content.userID, username,time FROM content,user WHERE content.userID=user.id"
	query += " ORDER BY time DESC LIMIT ?, ?"
	return sqlmodel.Query(query, offset, perpage)
}
func GetTotalCount() *sql.Rows {
	query := "SELECT count(id) AS count FROM content"
	return sqlmodel.Query(query)
}
func DelGuest(contentID int) int64 {
	query := "DELETE FROM content WHERE id=?"
	return sqlmodel.NonQuery(query, contentID)
}
func CheckUserID(contentID int, userID uint32) *sql.Rows {
	query := "SELECT id FROM content WHERE id=? AND userID=?"
	return sqlmodel.Query(query, contentID, userID)
}
