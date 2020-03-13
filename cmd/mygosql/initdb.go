package database

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/lmnzr/simpleshop/cmd/simpleshop/helper/config"
	"github.com/lmnzr/simpleshop/cmd/simpleshop/helper/env"
	logutil "github.com/lmnzr/simpleshop/cmd/simpleshop/helper/log"
)

func setup() map[string]interface{} {
	configmap := make(map[string]interface{})

	config, conferr := config.GetConfig()

	if conferr != nil {
		configmap["maxconn"] = 25
		configmap["maxiddle"] = 25
		configmap["maxlifetime"] = 5
		configmap["timezone"] = "Asia%2FJakarta"
	} else {
		configmap["maxconn"] = config.GetInt("dbMaxConn")
		configmap["maxiddle"] = config.GetInt("dbMaxIddle")
		configmap["maxlifetime"] = config.GetInt("dbMaxLifeTime")
		configmap["timezone"] = strings.Replace(config.GetString("timeZone"), "/", "%2F", 1)
	}
	return configmap

}

//OpenDbConnection :
func OpenDbConnection() (*sql.DB, error) {
	configmap := setup()

	dbMaxConns := configmap["maxconn"].(int)
	dbMaxIdleConns := configmap["maxiddle"].(int)
	dbMaxLifeTime := configmap["maxlifetime"].(int)
	timezone := configmap["timezone"].(string)

	logutil.LoggerDB().Info("open db connection")

	dbUser := env.Getenv("DB_USER", "lmnzr")
	dbPass := env.Getenv("DB_PASS", "root")
	dbHost := env.Getenv("DB_HOST", "localhost")
	dbPort := env.Getenv("DB_PORT", "3306")
	dbSchema := env.Getenv("DB_SCHEMA", "simpleshop")
	dbConnString :=
		fmt.Sprintf("%[1]s:%[2]s@tcp(%[3]s:%[4]s)/%[5]s?parseTime=true&charset=utf8&loc=%[6]s",
			dbUser, dbPass, dbHost, dbPort, dbSchema, timezone)

	db, err := sql.Open("mysql", dbConnString)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(dbMaxConns)
	db.SetMaxIdleConns(dbMaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(dbMaxLifeTime) * time.Minute)

	return db, nil
}
