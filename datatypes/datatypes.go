package datatypes

type DataType uint8

const (
	String DataType = iota
	Int
	Decimal
	Bool
	Date
)
