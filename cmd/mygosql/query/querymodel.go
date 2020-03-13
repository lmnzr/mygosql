package database

import (
	"github.com/lmnzr/simpleshop/cmd/simpleshop/database/filter"
	"github.com/lmnzr/simpleshop/cmd/simpleshop/database/group"
	"github.com/lmnzr/simpleshop/cmd/simpleshop/database/order"
)

//QueryModel : Database Access Object
type QueryModel struct {
	DBConn  DBConn
	Model   interface{}
	Table   string
	Limit   int
	Offset  int
	Filters []filter.Filter
	Orders  []order.Order
	Groups  []group.Group
	Error   []string
}

//NewQuery : Create Database Access Object With Generated Query
func NewQuery(dbcon DBConn, tablename string, model interface{}) QueryModel {
	return QueryModel{
		DBConn: dbcon,
		Table:  tablename,
		Model:  model,
	}
}

//NewQueryCustom : Create Database Access Object With Custom Query
func NewQueryCustom(dbcon DBConn) QueryModel {
	return QueryModel{
		DBConn: dbcon,
	}
}

//SetLimit : Set Query LIMIT only for SELECT
func (qm *QueryModel) SetLimit(limit int) *QueryModel {
	qm.Limit = limit
	return qm
}

//SetOffset : Set Query OFSSET if LIMIT defined only for SELECT
func (qm *QueryModel) SetOffset(offset int) *QueryModel {
	qm.Offset = offset
	return qm
}

//SetFilters : Set WHERE filter in Query
func (qm *QueryModel) SetFilters(filters []filter.Filter) *QueryModel {
	qm.Filters = filters
	return qm
}

//SetGroups : Set GROUP BY only for SELECT
func (qm *QueryModel) SetGroups(groups []group.Group) *QueryModel {
	qm.Groups = groups
	return qm
}

//SetOrders : Set ORDER BY only for SELECT
func (qm *QueryModel) SetOrders(orders []order.Order) *QueryModel {
	qm.Orders = orders
	return qm
}
