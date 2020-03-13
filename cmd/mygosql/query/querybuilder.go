package database

import (
	"errors"
	"fmt"
	"strings"

	"github.com/lmnzr/simpleshop/cmd/simpleshop/database/filter"
	reflectutil "github.com/lmnzr/simpleshop/cmd/simpleshop/helper/reflect"
)

type queryBuilder struct {
	QueryModel
	query      string
	field      []string
	param      []interface{}
	value      []interface{}
	softdelete bool
}

func (q *queryBuilder) commandSelect() *queryBuilder {
	if len(q.field) > 0 {
		fields := strings.Join(q.field, ",")
		if !(strings.Contains(strings.ToUpper(q.QueryModel.Table), "DROP") ||
			strings.Contains(strings.ToUpper(fields), "DROP")) {
			q.query = "SELECT " + fields + " FROM " + q.QueryModel.Table
		} else {
			q.Error = append(q.Error, "SQL INSERT : dangerous table or field name")
		}

	} else {
		q.Error = append(q.Error, "SQL SELECT : No Field Selected")
	}

	return q
}

func (q *queryBuilder) commandInsert() *queryBuilder {
	if !strings.Contains(strings.ToUpper(q.QueryModel.Table), "DROP") {
		q.query = "INSERT INTO " + q.QueryModel.Table
	} else {
		q.Error = append(q.Error, "SQL INSERT : dangerous table name")
	}

	return q
}

func (q *queryBuilder) commandUpdate() *queryBuilder {
	if !strings.Contains(strings.ToUpper(q.QueryModel.Table), "DROP") {
		q.query = "UPDATE " + q.QueryModel.Table
	} else {
		q.Error = append(q.Error, "SQL UPDATE : dangerous table name")
	}

	return q
}

func (q *queryBuilder) commandDelete() *queryBuilder {
	if !strings.Contains(strings.ToUpper(q.QueryModel.Table), "DROP") {
		q.query = "DELETE FROM " + q.QueryModel.Table
	} else {
		q.Error = append(q.Error, "SQL DELETE : dangerous table name")
	}

	return q
}

func (q *queryBuilder) where() *queryBuilder {
	if q.QueryModel.Filters != nil {
		firstcall := true

		q.query += " WHERE "

		var filter string
		for i := range q.QueryModel.Filters {
			k := q.QueryModel.Filters[i].Field
			o := q.QueryModel.Filters[i].GetOperator()
			v := q.QueryModel.Filters[i].Value
			t := q.QueryModel.Filters[i].Type

			if reflectutil.IsFieldExist(q.QueryModel.Model, k, "field") ||
				reflectutil.IsFieldExist(q.QueryModel.Model, k, "") {

				if firstcall {
					firstcall = false
				} else {
					q.query += fmt.Sprintf(" %s ", t)
				}

				filter = fmt.Sprintf("%[1]s %[2]s ?", k, o)
				q.param = append(q.param, v)
				q.query += filter
			} else {
				q.Error = append(q.Error, fmt.Sprintf("SQL WHERE : %s is not a valid database field", k))
			}
		}
	}

	return q
}

func (q *queryBuilder) groupBy() *queryBuilder {
	if q.QueryModel.Groups != nil {
		firstcall := true

		q.query += " GROUP BY "

		var group string
		for i := range q.QueryModel.Groups {
			k := q.QueryModel.Groups[i].Field

			if reflectutil.IsFieldExist(q.QueryModel.Model, k, "field") ||
				reflectutil.IsFieldExist(q.QueryModel.Model, k, "") {

				if firstcall {
					firstcall = false
				} else {
					q.query += " , "
				}

				group = fmt.Sprintf("%s", k)
				q.query += group
			} else {
				q.Error = append(q.Error, fmt.Sprintf("SQL GROUP BY : %s is not a valid database field", k))
			}
		}
	}

	return q
}

func (q *queryBuilder) orderBy() *queryBuilder {
	if q.QueryModel.Orders != nil {
		firstcall := true

		q.query += " ORDER BY "

		var order string
		for i := range q.QueryModel.Orders {
			k := q.QueryModel.Orders[i].Field
			t := q.QueryModel.Orders[i].Type

			if reflectutil.IsFieldExist(q.QueryModel.Model, k, "field") ||
				reflectutil.IsFieldExist(q.QueryModel.Model, k, "") {

				if firstcall {
					firstcall = false
				} else {
					q.query += " , "
				}

				order = fmt.Sprintf("%[1]s %[2]s", k, t)
				q.query += order
			} else {
				q.Error = append(q.Error, fmt.Sprintf("SQL ORDER BY : %s is not a valid database field", k))
			}
		}
	}

	return q
}

func (q *queryBuilder) limit() *queryBuilder {
	if q.QueryModel.Limit > 0 {
		limit := fmt.Sprintf(" LIMIT %d ", q.QueryModel.Limit)
		q.query += limit

		if q.QueryModel.Offset > 0 {
			offset := fmt.Sprintf(" OFFSET %d ", q.QueryModel.Offset)
			q.query += offset
		}
	}

	return q
}

func (q *queryBuilder) into() *queryBuilder {
	if q.field != nil {
		fields := strings.Join(q.field, " , ")
		fields = fmt.Sprintf(" (%s) ", fields)
		q.query += fields
	} else {
		q.Error = append(q.Error, "SQL INSERT :no inserted field")
	}
	return q
}

func (q *queryBuilder) values() *queryBuilder {
	var params []string

	q.query += " VALUES "

	if q.value != nil {
		for i := range q.value {
			k := q.field[i]

			if reflectutil.IsFieldExist(q.QueryModel.Model, k, "field") ||
				reflectutil.IsFieldExist(q.QueryModel.Model, k, "") {
				params = append(params, "?")
			} else {
				q.Error = append(q.Error, fmt.Sprintf("SQL INSERT : %s is not a valid database field", k))
			}
		}
	} else {
		q.Error = append(q.Error, "SQL INSERT :no inserted field")
	}

	q.query += fmt.Sprintf("(%s)", strings.Join(params, " , "))
	return q
}

func (q *queryBuilder) set() *queryBuilder {
	var params []string

	q.query += " SET "

	if q.value != nil {
		for i := range q.value {
			k := q.field[i]

			if reflectutil.IsFieldExist(q.QueryModel.Model, k, "field") ||
				reflectutil.IsFieldExist(q.QueryModel.Model, k, "") {
				params = append(params, fmt.Sprintf("%s=?", k))
			} else {
				q.Error = append(q.Error, fmt.Sprintf("SQL INSERT : %s is not a valid database field", k))
			}
		}
	} else {
		q.Error = append(q.Error, "SQL INSERT :no inserted field")
	}

	q.query += strings.Join(params, " , ")
	return q
}

func (q *queryBuilder) filterOutDeleted() *queryBuilder {
	deletedfield := reflectutil.IsFieldExist(q.QueryModel.Model, "is_deleted", "field")

	if deletedfield {
		if q.QueryModel.Filters == nil {
			var filters []filter.Filter
			q.QueryModel.Filters = filters
		}
		q.QueryModel.Filters = append(q.QueryModel.Filters, filter.NewAndFilter("is_deleted", "0"))
	}

	return q
}

func (q *queryBuilder) isSoftDelete() bool {
	return reflectutil.IsFieldExist(q.QueryModel.Model, "is_deleted", "field")
}

func (q *queryBuilder) setDeleteValue() *queryBuilder {
	if q.param == nil {
		var param []interface{}
		var value []interface{}
		q.param = param
		q.value = value
	}
	q.field = append(q.field, "is_deleted")
	q.value = append(q.value, true)

	return q
}

func (q *queryBuilder) checkFiltered() *queryBuilder {

	if q.QueryModel.Filters == nil {
		q.Error = append(q.Error, "SQL DELETE : dangerous delete / update without filter")
	}

	return q
}

func queryBuild(method string, qm QueryModel) (dbQuery DBQuery, err error) {
	qs := queryBuilder{
		QueryModel: qm,
		query:      "",
	}

	structmap := MapModel(qm.Model, method)
	qs.field = structmap.Fields
	qs.value = structmap.Values

	switch method {
	case "SELECT":
		qs.filterOutDeleted()
		qs.commandSelect().where().groupBy().orderBy().limit()

	case "INSERT":
		qs.param = structmap.Values
		qs.commandInsert().into().values()

	case "UPDATE":
		qs.param = structmap.Values

		qs.SetFilters(mapFilter(structmap.Filter))
		qs.checkFiltered()
		qs.commandUpdate().set().where()

	case "DELETE":
		if qs.isSoftDelete() {
			qs.setDeleteValue()
			qs.SetFilters(mapFilter(structmap.Filter))
			qs.param = qs.value
			qs.commandUpdate().set().where()
		} else {
			qs.SetFilters(mapFilter(structmap.Filter))
			qs.commandDelete().where()
		}

	default:
		err = errors.New("QUERY BUILDER : undefined query method")
	}

	if qs.Error != nil {
		err = errors.New(strings.Join(qs.Error, "\n"))
	}

	dbq := DBQuery{
		QueryString: qs.query,
		QueryParam:  qs.param,
		SelectField: structmap.Fields,
		SelectType:  structmap.Values,
	}
	return dbq, err
}

func mapFilter(filtermap map[string]string) []filter.Filter {
	var filters []filter.Filter
	for k, v := range filtermap {
		filters = append(filters, filter.NewAndFilter(k, v))
	}

	return filters
}
