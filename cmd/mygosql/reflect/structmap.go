package database

import (
	// "fmt"
	"fmt"
	"reflect"

	nulltype "github.com/lmnzr/simpleshop/cmd/simpleshop/types"
)

//StructMap : Map struct to arrays
type StructMap struct {
	Fields  []string
	Values  []interface{}
	Indexes []int
	Filter  map[string]string
}

type tag struct {
	field  string
	value  string
	remove bool
	unique bool
	filter bool
}

func maptype(types []interface{}, t reflect.Type, i int) []interface{} {
	fieldtype := "string"

	if val, found := t.Field(i).Tag.Lookup("type"); found {
		fieldtype = val
	}

	switch fieldtype {
	case "int":
		types = append(types, new(nulltype.NullInt))
	case "float":
		types = append(types, new(nulltype.NullFloat))
	case "boolean":
		types = append(types, new(nulltype.NullBool))
	case "datetime":
		types = append(types, new(nulltype.NullTime))
	default:
		types = append(types, new(nulltype.NullString))
	}

	return types
}

func mapvalue(fields []string, values []interface{}, indexes []int, tagname string, v reflect.Value, t reflect.Type, i int) ([]string, []interface{}, []int) {
	fieldtype := "string"

	if val, found := t.Field(i).Tag.Lookup("type"); found {
		fieldtype = val
	}

	val := v.Field(i).Interface()

	switch fieldtype {
	case "int":
		if val.(nulltype.NullInt).Valid {
			fields = append(fields, mapfield(t.Field(i), tagname))
			values = append(values, val)
			indexes = append(indexes,i)
		}
	case "float":
		if val.(nulltype.NullString).Valid {
			fields = append(fields, mapfield(t.Field(i), tagname))
			values = append(values, val)
			indexes = append(indexes,i)
		}
	case "boolean":
		if val.(nulltype.NullBool).Valid {
			fields = append(fields, mapfield(t.Field(i), tagname))
			values = append(values, val)
			indexes = append(indexes,i)
		}
	case "datetime":
		if val.(nulltype.NullTime).Valid {
			fields = append(fields, mapfield(t.Field(i), tagname))
			values = append(values, val)
			indexes = append(indexes,i)
		}
	case "string":
		if val.(nulltype.NullString).Valid {
			fields = append(fields, mapfield(t.Field(i), tagname))
			values = append(values, val)
			indexes = append(indexes,i)
		}
	}

	return fields, values, indexes
}

func mapfield(tag reflect.StructField, tagname string) (field string) {
	if val, found := tag.Tag.Lookup(tagname); found {
		field = val
	} else {
		field = tag.Name
	}

	return field
}

func mapping(model interface{}, tagname string, flag tag, hidden bool, cmdtype string) StructMap {
	var fields []string
	var values []interface{}
	var indexes []int
	filter := make(map[string]string)

	t := reflect.TypeOf(model)
	v := reflect.ValueOf(model)

	for i := 0; i < t.NumField(); i++ {
		swtch, ok := t.Field(i).Tag.Lookup(flag.field)
		match := ok && swtch == flag.value

		var cond bool

		if flag.remove {
			cond = !match || !hidden
		} else if flag.unique {
			cond = match
		}

		if flag.filter && match {
			var field, fieldtype string
			var value string

			field = mapfield(t.Field(i),tagname)

			if val, found := t.Field(i).Tag.Lookup("type"); found {
				fieldtype = val
			}
			
			val := v.Field(i).Interface()

			switch fieldtype {
			case "int":
				if val.(nulltype.NullInt).Valid {
					value = fmt.Sprintf("%d", val.(nulltype.NullInt).Int64)
				}
			case "string":
				if val.(nulltype.NullString).Valid {
					value = val.(nulltype.NullString).String.String
				}
			}

			filter[field] = value

		} else if cond {
			switch cmdtype {
			case "SELECT":
				indexes = append(indexes,i)
				fields = append(fields,mapfield(t.Field(i),tagname))
				values = maptype(values, t, i)

			case "INSERT", "UPDATE", "DELETE":
				fields, values, indexes = mapvalue(fields, values, indexes, tagname, v, t, i)

			}
		}
	}

	return StructMap{
		Fields:  fields,
		Values:  values,
		Indexes: indexes,
		Filter:  filter,
	}
}

//MapModel : Create Mapping of a struct
func MapModel(model interface{}, cmdtype string) (structmap StructMap) {
	var tagselector tag

	switch cmdtype {
	case "SELECT":
		tagselector = tag{
			field:  "hidden",
			value:  "true",
			remove: true,
			unique: false,
			filter: false,
		}
	case "INSERT":
		tagselector = tag{
			field:  "increment",
			value:  "auto",
			remove: true,
			unique: false,
			filter: false,
		}
	case "UPDATE":
		tagselector = tag{
			field:  "pkey",
			value:  "true",
			remove: true,
			unique: false,
			filter: true,
		}
	case "DELETE":
		tagselector = tag{
			field:  "pkey",
			value:  "true",
			remove: false,
			unique: true,
			filter: true,
		}
	}

	return mapping(model, "field", tagselector, true, cmdtype)
}
