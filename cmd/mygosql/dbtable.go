package database

//TableQuery : Query For SQL Database Table
type TableQuery struct {
	QueryModel
}

//NewTableQuery : Create Database Access Object With Generated Query
func NewTableQuery(dbcon DBConn, tablename string, model interface{}) TableQuery {
	return TableQuery{
		NewQuery(dbcon, tablename, model),
	}
}

//NewTableQueryCustom : Create Database Access Object With Custom Query
func NewTableQueryCustom(dbcon DBConn) TableQuery {
	return TableQuery{
		NewQueryCustom(dbcon),
	}
}

//RetrieveAll :
func (tm *TableQuery) RetrieveAll() (res [][]byte, err error) {
	var cmdres interface{}
	cmdres, err = command("SELECT", true, tm.QueryModel)

	if cmdres != nil {
		res = cmdres.([][]byte)
	}

	return res, err
}

//RetrieveAllCustom :
func (tm *TableQuery) RetrieveAllCustom(dbQuery DBQuery) (res [][]byte, err error) {
	var cmdres interface{}
	cmdres, err = commandCustom("SELECT", true, tm.QueryModel, dbQuery)

	if cmdres != nil {
		res = cmdres.([][]byte)
	}

	return res, err
}

//Retrieve :
func (tm *TableQuery) Retrieve() (res []byte, err error) {
	var cmdres interface{}
	cmdres, err = command("SELECT", false, tm.QueryModel)

	if cmdres != nil {
		res = cmdres.([]byte)
	}

	return res, err
}

//RetrieveCustom :
func (tm *TableQuery) RetrieveCustom(dbQuery DBQuery) (res []byte, err error) {
	var cmdres interface{}
	cmdres, err = commandCustom("SELECT", false, tm.QueryModel, dbQuery)

	if cmdres != nil {
		res = cmdres.([]byte)
	}

	return cmdres.([]byte), err
}

//Insert :
func (tm *TableQuery) Insert() (insertedID int64, err error) {
	var cmdres interface{}
	cmdres, err = command("INSERT", false, tm.QueryModel)

	if cmdres != nil {
		insertedID = cmdres.(int64)
	}

	return insertedID, err
}

//InsertCustom :
func (tm *TableQuery) InsertCustom(dbQuery DBQuery) (insertedID int64, err error) {
	var cmdres interface{}
	cmdres, err = commandCustom("INSERT", false, tm.QueryModel, dbQuery)

	if cmdres != nil {
		insertedID = cmdres.(int64)
	}

	return insertedID, err
}

//Update :
func (tm *TableQuery) Update() (updatedRows int64, err error) {
	var cmdres interface{}
	cmdres, err = command("UPDATE", false, tm.QueryModel)

	if cmdres != nil {
		updatedRows = cmdres.(int64)
	}

	return updatedRows, err
}

//UpdateCustom :
func (tm *TableQuery) UpdateCustom(dbQuery DBQuery) (updatedRows int64, err error) {
	var cmdres interface{}
	cmdres, err = commandCustom("UPDATE", false, tm.QueryModel, dbQuery)

	if cmdres != nil {
		updatedRows = cmdres.(int64)
	}

	return updatedRows, err
}

//Delete :
func (tm *TableQuery) Delete() (updatedRows int64, err error) {
	var cmdres interface{}
	cmdres, err = command("DELETE", false, tm.QueryModel)

	if cmdres != nil {
		updatedRows = cmdres.(int64)
	}

	return updatedRows, err
}

//DeleteCustom :
func (tm *TableQuery) DeleteCustom(dbQuery DBQuery) (updatedRows int64, err error) {
	var cmdres interface{}
	cmdres, err = commandCustom("DELETE", false, tm.QueryModel, dbQuery)

	if cmdres != nil {
		updatedRows = cmdres.(int64)
	}

	return updatedRows, err
}
