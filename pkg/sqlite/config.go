package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func SelectSubTarget(dbFilePath string) (forward []string) {
	db, err := sql.Open("sqlite3", dbFilePath)
	if err != nil {
		fmt.Println("读取数据库发生错误:", err)
		os.Exit(-1)
	}
	defer db.Close()

	rows, err := db.Query("SELECT forward FROM subTable")
	if err != nil {
		fmt.Println("数据库查询失败:", err)
		os.Exit(-1)
	}
	defer rows.Close()
	var col1 string
	for rows.Next() {
		err := rows.Scan(&col1)
		if err != nil {
			log.Print("[-]", err)
		}
		row := []string{col1}
		forward = append(forward, row...)
	}
	return forward
}
