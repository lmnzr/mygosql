package order

//Direction : SQL ORDER BY sort Direction ASC/DESC
type Direction struct{
	Type string
}

//Order : SQL ORDER BY
type Order struct {
	Direction
	Field string
}

func directionAscending() Direction{
	return Direction{Type: "ASC"}
}

func directionDescending() Direction{
	return Direction{Type: "DESC"}
}

//NewOrderAscending : New Entry SQL ORDER BY ASC
func NewOrderAscending(field string) Order{
	return Order{
		Direction : directionAscending(),
		Field: field,
	}
}

//NewOrderDescending : New Entry SQL ORDER BY DESC
func NewOrderDescending(field string) Order{
	return Order{
		Direction : directionDescending(),
		Field: field,
	}
}
