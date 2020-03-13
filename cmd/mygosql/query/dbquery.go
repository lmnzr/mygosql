package database

import (
	// "fmt"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/lmnzr/simpleshop/cmd/simpleshop/helper/env"
	logutil "github.com/lmnzr/simpleshop/cmd/simpleshop/helper/log"
)

//DBQuery : Necessary Input for DB Query
type DBQuery struct {
	QueryString string
	QueryParam  []interface{}
	SelectField []string
	SelectType  []interface{}
}

func queryRows(dbQuery DBQuery, qm QueryModel) (results [][]byte, err error) {
	var rows *sql.Rows

	rows, err = qm.DBConn.Query(dbQuery.QueryString, dbQuery.QueryParam...)

	resultmap := make(map[string]interface{})
	for i := 0; i < len(dbQuery.SelectField); i++ {
		resultmap[dbQuery.SelectField[i]] = dbQuery.SelectType[i]
	}

	defer rows.Close()

	if err == nil {

		for rows.Next() {
			err = rows.Scan(dbQuery.SelectType...)
			if err == nil {

				var res []byte
				res, err = json.Marshal(resultmap)

				if err == nil {
					results = append(results, res)
				}
			}
		}
	}

	return results, err
}

func queryRow(dbQuery DBQuery, qm QueryModel) (result []byte, err error) {
	var row *sql.Row

	row = qm.DBConn.QueryRow(dbQuery.QueryString, dbQuery.QueryParam...)

	resultmap := make(map[string]interface{})
	for i := 0; i < len(dbQuery.SelectField); i++ {
		resultmap[dbQuery.SelectField[i]] = dbQuery.SelectType[i]
	}

	err = row.Scan(dbQuery.SelectType...)

	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
	} else {
		result, err = json.Marshal(resultmap)
	}

	return result, err
}

func queryExec(dbQuery DBQuery, qm QueryModel) (result sql.Result, err error) {
	var stmt *sql.Stmt

	stmt, err = qm.DBConn.Prepare(dbQuery.QueryString)

	result, err = stmt.Exec(dbQuery.QueryParam...)

	defer stmt.Close()

	return result, err
}

func (dbq *DBQuery) toString() string {
	str := ""

	kvmap := make(map[int]interface{})
	for i := 0; i < len(dbq.QueryParam); i++ {
		kvmap[i] = dbq.QueryParam[i]
	}
	result, err := json.Marshal(kvmap)

	str += fmt.Sprintf("query=\"%s\"", dbq.QueryString)
	str += " , "
	if err == nil {
		str += fmt.Sprintf("param=%s", string(result))
	}

	return str
}

func command(cmdtype string, multipleval bool, qm QueryModel) (result interface{}, err error) {
	dbQuery, builderr := queryBuild(cmdtype, qm)

	debug(dbQuery.toString())

	if builderr == nil {
		result, err = execute(cmdtype, multipleval, qm, dbQuery)

	} else {
		err = builderr
	}

	return result, err
}

func commandCustom(cmdtype string, multipleval bool, qm QueryModel, dbQuery DBQuery) (result interface{}, err error) {
	debug(dbQuery.toString())

	return execute(cmdtype, multipleval, qm, dbQuery)
}

func execute(cmdtype string, multipleval bool, qm QueryModel, dbQuery DBQuery) (result interface{}, err error) {
	switch cmdtype {
	case "SELECT":
		if multipleval {
			result, err = queryRows(dbQuery, qm)
		} else {
			result, err = queryRow(dbQuery, qm)
		}

	case "INSERT":
		var sqlres sql.Result
		sqlres, err = queryExec(dbQuery, qm)
		if err == nil {
			result, err = sqlres.LastInsertId()
		}

	case "UPDATE", "DELETE":
		var sqlres sql.Result
		sqlres, err = queryExec(dbQuery, qm)
		if err == nil {
			result, err = sqlres.RowsAffected()
		}

	default:
		err = errors.New("Undefined Query Type")
	}

	return result, err
}

func debug(debug string) {
	environment := env.Getenv("ENVIRONMENT", "development")

	if environment == "development" {
		logutil.LoggerDB().Debug(debug)
	}
}
