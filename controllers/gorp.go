package controllers

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/go-gorp/gorp"
	_ "github.com/go-sql-driver/mysql"
)

var dbmap = initDb()

func initDb() *gorp.DbMap {
	db, err := sql.Open("mysql", os.Getenv("GO_MYSQL"))
	if err != nil {
		fmt.Println(err, "sql.Open failed")
	}
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{Engine: "InnoDB", Encoding: "UTF8"}}
	// err = dbmap.CreateTablesIfNotExists()
	// checkErr(err, "Create tables failed")
	return dbmap
}

// Cors (Cross-origin resource sharing) 교차 출처 리소스 공유
// func Cors() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
// 		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
// 		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
// 		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET")
// 		c.Next()
// 	}
// }
