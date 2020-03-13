package filter

//Condition : SQL WHERE Condition OR/AND
type Condition struct {
	Type string
}

//Filter : SQL WHERE Filter
type Filter struct {
	Condition
	Field    string
	Value    string
	operator string
}

func conditionAnd() Condition {
	return Condition{Type: "AND"}
}

func conditionOr() Condition {
	return Condition{Type: "OR"}
}

//NewAndFilter : Create New AND Filter for SQL Where
func NewAndFilter(field string, value string) Filter {
	return Filter{
		Condition: conditionAnd(),
		Field:     field,
		Value:     value,
		operator:  "=",
	}
}

//NewOrFilter : Create New OR Filter for SQL Where
func NewOrFilter(field string, value string) Filter {
	return Filter{
		Condition: conditionOr(),
		Field:     field,
		Value:     value,
		operator:  "=",
	}
}

//SetOperator : Set Operation For Filtering
func (f Filter) SetOperator(operator string) Filter {
	switch operator {
	case "=", "!=", "<>", ">=", "<=", ">", "<":
		f.operator = operator
	}

	return f
}

//GetOperator : Get Operation For Filtering
func (f Filter) GetOperator() string {
	return f.operator
}
