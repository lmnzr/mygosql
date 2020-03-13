package group

//Group : SQL GROUP BY
type Group struct {
	Field string
}

//NewGroup : Create New Entry for SQL GROUP BY
func NewGroup(field string) Group {
	return Group{
		Field: field,
	}
}
