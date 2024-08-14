package compiler

import "marmoset/utils"

type SymbolScope string

const (
	LocalScope   SymbolScope = "LOCAL"
	GlobalScope  SymbolScope = "GLOBAL"
	BuiltinScope SymbolScope = "BUILTIN"
)

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

type SymbolTable struct {
	Outer               *SymbolTable // "parent" symbol table
	store               map[string]Symbol
	reverseStore        []string // allows fetching a symbol name by its index
	reverseBuiltinStore []string // same as reverseStore, but only for builtins
	numDefinitions      int
}

const SymbolsSize = 65536

func NewSymbolTable() *SymbolTable {
	s := make(map[string]Symbol)
	rs := make([]string, SymbolsSize)
	rsb := make([]string, SymbolsSize)
	return &SymbolTable{store: s, reverseStore: rs, reverseBuiltinStore: rsb}
}

func NewEnclosedSymbolTable(outer *SymbolTable) *SymbolTable {
	s := NewSymbolTable()
	s.Outer = outer
	return s
}

func (s *SymbolTable) Define(name string) Symbol {
	symbol := Symbol{Name: name, Index: s.numDefinitions}
	if s.Outer == nil {
		symbol.Scope = GlobalScope
	} else {
		symbol.Scope = LocalScope
	}

	s.store[name] = symbol
	s.reverseStore[s.numDefinitions] = name
	s.numDefinitions++
	return symbol
}

// finds a symbol by its name
func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	obj, ok := s.store[name]
	if !ok && s.Outer != nil {
		// if there is a parent scope, call Resolve again recursively
		obj, ok = s.Outer.Resolve(name)
		return obj, ok
	}
	return obj, ok
}

// finds a symbol name by its index (must indicate symbol's scope)
func (s *SymbolTable) ResolveName(index int, scope SymbolScope) (string, bool) {
	switch scope {
	case BuiltinScope:
		return utils.At(s.reverseBuiltinStore, index)
	case GlobalScope:
		if s.Outer != nil {
			return s.Outer.ResolveName(index, scope)
		}
		return utils.At(s.reverseStore, index)
	case LocalScope:
		return utils.At(s.reverseStore, index)
	}

	return "", false
}

func (s *SymbolTable) DefineBuiltin(index int, name string) Symbol {
	symbol := Symbol{Name: name, Index: index, Scope: BuiltinScope}
	s.store[name] = symbol
	s.reverseBuiltinStore[index] = name
	return symbol
}
