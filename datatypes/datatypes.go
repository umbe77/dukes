// Copyright (c) 2023 Robeto Ughi
// 
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package datatypes

type DataType uint8

const (
	String DataType = iota
	Int
	Decimal
	Bool
	Date
)
