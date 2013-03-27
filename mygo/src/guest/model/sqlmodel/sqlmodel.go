package sqlmodel

import (
	"database/sql"
	_ "github.com/Go-SQL-Driver/MySQL"
	"log"
)
type config struct {
	mysqlconfig
}
type mysqlconfig struct {
	driver string
	user string
	passwd string 
	dbname string 
}
var sqlconfig = config{mysqlconfig{driver: "mysql", user:"testaa", passwd: "testaa", dbname: "guest"}}

func Query(query string, args ...interface{}) *sql.Rows {
	db, err := sql.Open(sqlconfig.mysqlconfig.driver, sqlconfig.mysqlconfig.user+":"+sqlconfig.mysqlconfig.passwd+"@/"+sqlconfig.mysqlconfig.dbname+"?charset=utf8") 
	checkErr(err)
	defer db.Close()	
	
	rows, err := db.Query(query, args...)
	checkErr(err)
	return rows
}
func NonQuery(query string, args ...interface{}) int64 {
	db, err := sql.Open(sqlconfig.mysqlconfig.driver, sqlconfig.mysqlconfig.user+":"+sqlconfig.mysqlconfig.passwd+"@/"+sqlconfig.mysqlconfig.dbname+"?charset=utf8") 
	checkErr(err)
	defer db.Close()
	
	stmt, err := db.Prepare(query)
	checkErr(err)
	res, err := stmt.Exec(args...)
	checkErr(err)
	affected, err := res.RowsAffected();
	checkErr(err)
	return affected	
}

func checkErr(err error) {
    if err != nil {
		log.Fatal(err)
    }
}
