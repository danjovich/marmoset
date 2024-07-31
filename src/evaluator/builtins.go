package evaluator

import (
	"monkey/object"
)

var builtins = map[string]*object.Builtin{
	// length of array or string
	"len": object.GetBuiltinByName("len"),
	// get first element of an array
	"first": object.GetBuiltinByName("first"),
	// get last element of an array
	"last": object.GetBuiltinByName("last"),
	// get new array from existing one without the first element
	"rest": object.GetBuiltinByName("rest"),
	// pushes element copied array
	"push": object.GetBuiltinByName("push"),
	// print value(s) to STDOUT
	"puts": object.GetBuiltinByName("puts"),
}
