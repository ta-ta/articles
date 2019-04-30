package db

import (
	"fmt"
	"log"

	"../env"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Database struct {
	database *sqlx.DB
}

var (
	DB Database
)

func Init() {
	// DBの設定
	BDConfig := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		env.DBUser,
		env.DBPassword,
		env.DBHost,
		"3306",
		env.DBDatabase)

	// DB準備
	database, err := sqlx.Open(env.DefaultDriver, BDConfig)
	if err != nil {
		log.Fatalf("DB Connection Error: %v", err)
		return
	}
	DB.database = database
}
