package database

//ViewQuery : Query For SQL Database View (Only SELECT is Possible)
type ViewQuery struct {
	QueryModel
}

//NewViewQuery : Create Database Access Object With Generated Query
func NewViewQuery(dbcon DBConn, tablename string, model interface{}) ViewQuery {
	return ViewQuery{
		NewQuery(dbcon, tablename, model),
	}
}

//NewViewQueryCustom : Create Database Access Object With Custom Query
func NewViewQueryCustom(dbcon DBConn) ViewQuery {
	return ViewQuery{
		NewQueryCustom(dbcon),
	}
}

//RetrieveAll :
func (tm *ViewQuery) RetrieveAll() (res [][]byte, err error) {
	var cmdres interface{}
	cmdres, err = command("SELECT", true, tm.QueryModel)

	if cmdres != nil {
		res = cmdres.([][]byte)
	}

	return res, err
}

//RetrieveAllCustom :
func (tm *ViewQuery) RetrieveAllCustom(dbQuery DBQuery) (res [][]byte, err error) {
	var cmdres interface{}
	cmdres, err = commandCustom("SELECT", true, tm.QueryModel, dbQuery)

	if cmdres != nil {
		res = cmdres.([][]byte)
	}

	return res, err
}

//Retrieve :
func (tm *ViewQuery) Retrieve() (res []byte, err error) {
	var cmdres interface{}
	cmdres, err = command("SELECT", false, tm.QueryModel)

	if cmdres != nil {
		res = cmdres.([]byte)
	}

	return res, err
}

//RetrieveCustom :
func (tm *ViewQuery) RetrieveCustom(dbQuery DBQuery) (res []byte, err error) {
	var cmdres interface{}
	cmdres, err = commandCustom("SELECT", false, tm.QueryModel, dbQuery)

	if cmdres != nil {
		res = cmdres.([]byte)
	}

	return cmdres.([]byte), err
}
