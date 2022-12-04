package repository

import (
	"database/sql"
	"fmt"
	"sync"

	. "github.com/apachejuice/eelchat/internal/config/keys"
	"github.com/apachejuice/eelchat/internal/logs"
	_ "github.com/go-sql-driver/mysql"
)

var (
	dbOnce sync.Once
	db     *sql.DB

	dbLog, repoLog logs.Logger
)

func connectDB() *sql.DB {
	dbOnce.Do(func() {
		dbLog = logs.NewLogger("database")
		repoLog = logs.NewLogger("repository")

		user, dbname := ConfigKeyDbUser.Get(), ConfigKeyDbName.Get()
		conn, err := sql.Open("mysql", fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?parseTime=true",
			user,
			ConfigKeyDbPass.Get(),
			ConfigKeyDbHost.Get(),
			ConfigKeyDbPort.Get(),
			dbname,
		))

		if err != nil {
			dbLog.Fatal("Failed to connect to database", "error", err.Error(), "user", user, "dbName", dbname)
		}

		err = conn.Ping()
		if err != nil {
			dbLog.Fatal("Failed to connect to database", "error", err.Error(), "user", user, "dbName", dbname)
		}

		db = conn
	})

	return db
}
