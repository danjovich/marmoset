package compiler

import (
	"fmt"
	"marmoset/ast"
	"marmoset/code"
	"marmoset/object"
)

type Compiler struct {
	Constants   []object.Object     // constant pool
	SymbolTable *SymbolTable        // symbol table
	scopes      []*CompilationScope // scope stack
	AllScopes   []*CompilationScope // all scopes, to allow future inspection
	scopeIndex  int                 // scope stack "pointer"
}

type EmittedInstruction struct {
	Opcode   code.Opcode
	Position int
}

// Scope for allowing function compilation to result in another scope than the main code
type CompilationScope struct {
	Instructions        code.Instructions  // generated bytecode
	lastInstruction     EmittedInstruction // last instruction
	previousInstruction EmittedInstruction // instruction before the last
	Name                string             // the scope's function name
	IsMain              bool               // whether it is the main scope or not
	Args                []string           // list of arguments passed to the scope
	SymbolTable         *SymbolTable       // symbol table
}

func New() *Compiler {
	symbolTable := NewSymbolTable()

	mainScope := CompilationScope{
		Instructions:        code.Instructions{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
		Name: `.global _start
_start`,
		IsMain:      true,
		Args:        []string{},
		SymbolTable: symbolTable,
	}

	for i, v := range object.Builtins {
		symbolTable.DefineBuiltin(i, v.Name)
	}

	return &Compiler{
		Constants:   []object.Object{},
		SymbolTable: symbolTable,
		scopes:      []*CompilationScope{&mainScope},
		AllScopes:   []*CompilationScope{&mainScope},
		scopeIndex:  0,
	}
}

// For preserving state between REPL runs
func NewWithState(s *SymbolTable, constants []object.Object) *Compiler {
	compiler := New()
	compiler.SymbolTable = s
	compiler.AllScopes[0].SymbolTable = s
	compiler.Constants = constants
	return compiler
}

func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}

	case *ast.ExpressionStatement:
		// hanging ; should do nothing
		if node.Expression == nil {
			return nil
		}
		err := c.Compile(node.Expression)
		if err != nil {
			return err
		}
		// expressions should pop the computed value from the stack
		c.emit(code.OpPop)

	case *ast.InfixExpression:
		if node.Operator == "<" {
			// reorder operators and emit OpGreaterThan
			err := c.Compile(node.Right)
			if err != nil {
				return err
			}

			err = c.Compile(node.Left)
			if err != nil {
				return err
			}

			c.emit(code.OpGreaterThan)
			return nil
		}
		err := c.Compile(node.Left)
		if err != nil {
			return err
		}
		err = c.Compile(node.Right)
		if err != nil {
			return err
		}
		switch node.Operator {
		case "+":
			c.emit(code.OpAdd)
		case "-":
			c.emit(code.OpSub)
		case "*":
			c.emit(code.OpMul)
		case "/":
			c.emit(code.OpDiv)
		case ">":
			c.emit(code.OpGreaterThan)
		case "==":
			c.emit(code.OpEqual)
		case "!=":
			c.emit(code.OpNotEqual)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}

	case *ast.PrefixExpression:
		err := c.Compile(node.Right)
		if err != nil {
			return err
		}
		switch node.Operator {
		case "!":
			c.emit(code.OpBang)
		case "-":
			c.emit(code.OpMinus)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}

	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(integer))

	case *ast.Boolean:
		if node.Value {
			c.emit(code.OpTrue)
		} else {
			c.emit(code.OpFalse)
		}

	case *ast.IfExpression:
		err := c.Compile(node.Condition)
		if err != nil {
			return err
		}
		// Emit an `OpJumpNotTruthy` with a bogus value
		jumpNotTruthyPos := c.emit(code.OpJumpNotTruthy, 9999)
		err = c.Compile(node.Consequence)
		if err != nil {
			return err
		}
		// remove last pop (conditionals should store result in the stack)
		if c.lastInstructionIs(code.OpPop) {
			c.removeLastInstruction()
		}
		// Emit an `OpJump` with a bogus value
		jumpPos := c.emit(code.OpJump, 9999)
		afterConsequencePos := len(c.currentInstructions())
		// changes the first bogus 9999 to the actual jump position
		c.changeOperand(jumpNotTruthyPos, afterConsequencePos)
		// if no else
		if node.Alternative == nil {
			// alternative is just pushing null in the case of no else statement
			c.emit(code.OpNull)
		} else {
			err := c.Compile(node.Alternative)
			if err != nil {
				return err
			}
			if c.lastInstructionIs(code.OpPop) {
				c.removeLastInstruction()
			}
		}
		afterAlternativePos := len(c.currentInstructions())
		// changes the second bogus 9999 to the actual jump position
		c.changeOperand(jumpPos, afterAlternativePos)

	case *ast.BlockStatement:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}

	case *ast.LetStatement:
		err := c.Compile(node.Value)
		if err != nil {
			return err
		}
		symbol := c.SymbolTable.Define(node.Name.Value)
		if symbol.Scope == GlobalScope {
			c.emit(code.OpSetGlobal, symbol.Index)
		} else {
			c.emit(code.OpSetLocal, symbol.Index)
		}

	case *ast.Identifier:
		symbol, ok := c.SymbolTable.Resolve(node.Value)
		if !ok {
			return fmt.Errorf("undefined variable %s", node.Value)
		}
		c.loadSymbol(symbol)

	case *ast.StringLiteral:
		str := &object.String{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(str))

	case *ast.ArrayLiteral:
		for _, el := range node.Elements {
			// compile each element...
			err := c.Compile(el)
			if err != nil {
				return err
			}
		}
		// ...and then compile the array
		c.emit(code.OpArray, len(node.Elements))

	case *ast.IndexExpression:
		err := c.Compile(node.Left)
		if err != nil {
			return err
		}
		err = c.Compile(node.Index)
		if err != nil {
			return err
		}
		c.emit(code.OpIndex)

	case *ast.FunctionStatement:
		// defines symbol before compiling function to allow recursive functions to know about their own existence
		symbol := c.SymbolTable.Define(node.Name.Value)
		// enters new function scope
		c.enterScope(node.Name.Value, node.Parameters)
		// defines function arguments in the current scope symbol table
		for _, p := range node.Parameters {
			c.SymbolTable.Define(p.Value)
		}
		err := c.Compile(node.Body)
		if err != nil {
			return err
		}
		// if last instruction is pop, remove it
		if c.lastInstructionIs(code.OpPop) {
			c.replaceLastInstructionWithReturn()
		}
		// we should always have a return at the end (void functions return null)
		if !c.lastInstructionIs(code.OpReturnValue) {
			c.emit(code.OpReturn)
		}
		// gets number of local variables in the function
		numLocals := c.SymbolTable.numDefinitions
		// gets instructions from the function scope
		instructions := c.leaveScope()
		// creates a CompiledFunction object with the instructions
		compiledFn := &object.CompiledFunction{
			Instructions:  instructions,
			NumLocals:     numLocals,
			NumParameters: len(node.Parameters),
			Name:          node.Name.Value,
		}
		// emits an instruction to push the CompiledFunction as a constant
		c.emit(code.OpConstant, c.addConstant(compiledFn))
		// emits definition of identifier
		if symbol.Scope == GlobalScope {
			c.emit(code.OpSetGlobal, symbol.Index)
		} else {
			c.emit(code.OpSetLocal, symbol.Index)
		}

	case *ast.ReturnStatement:
		err := c.Compile(node.ReturnValue)
		if err != nil {
			return err
		}
		c.emit(code.OpReturnValue)

	case *ast.CallExpression:
		err := c.Compile(node.Function)
		if err != nil {
			return err
		}
		for _, a := range node.Arguments {
			err := c.Compile(a)
			if err != nil {
				return err
			}
		}
		c.emit(code.OpCall, len(node.Arguments))
	}

	return nil
}

// Loads a symbol
func (c *Compiler) loadSymbol(s Symbol) {
	switch s.Scope {
	case GlobalScope:
		c.emit(code.OpGetGlobal, s.Index)
	case LocalScope:
		c.emit(code.OpGetLocal, s.Index)
	case BuiltinScope:
		c.emit(code.OpGetBuiltin, s.Index)
	}
}

// Enters a new scope
func (c *Compiler) enterScope(name string, args []*ast.Identifier) {
	argsList := []string{}
	for _, arg := range args {
		argsList = append(argsList, arg.Value)
	}
	newSymbolTable := NewEnclosedSymbolTable(c.SymbolTable)

	scope := CompilationScope{
		Instructions:        code.Instructions{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
		Name:                name,
		IsMain:              false,
		Args:                argsList,
		SymbolTable:         newSymbolTable,
	}

	c.scopes = append(c.scopes, &scope)
	c.AllScopes = append(c.AllScopes, &scope)
	c.scopeIndex++
	c.SymbolTable = newSymbolTable
}

// Leaves current scope, returning its instructions
func (c *Compiler) leaveScope() code.Instructions {
	instructions := c.currentInstructions()

	c.scopes = c.scopes[:len(c.scopes)-1]
	c.scopeIndex--
	c.SymbolTable = c.SymbolTable.Outer

	return instructions
}

// Adds to the constant pool
func (c *Compiler) addConstant(obj object.Object) int {
	c.Constants = append(c.Constants, obj)
	// return identifier (index) for the constant
	return len(c.Constants) - 1
}

// Emits bytecode
func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	ins := code.Make(op, operands...)
	pos := c.addInstruction(ins)

	c.setLastInstruction(op, pos)

	// returns the starting position of the emitted instruction
	return pos
}

// Gets current scope instructions
func (c *Compiler) currentInstructions() code.Instructions {
	return c.scopes[c.scopeIndex].Instructions
}

// Adds instruction to the instructions list
func (c *Compiler) addInstruction(ins []byte) int {
	posNewInstruction := len(c.currentInstructions())
	updatedInstructions := append(c.currentInstructions(), ins...)

	c.scopes[c.scopeIndex].Instructions = updatedInstructions

	return posNewInstruction
}

// Sets last and previous instructions
func (c *Compiler) setLastInstruction(op code.Opcode, pos int) {
	previous := c.scopes[c.scopeIndex].lastInstruction
	last := EmittedInstruction{Opcode: op, Position: pos}

	c.scopes[c.scopeIndex].previousInstruction = previous
	c.scopes[c.scopeIndex].lastInstruction = last
}

// Check if last emitted instruction (in the current scope) is op
func (c *Compiler) lastInstructionIs(op code.Opcode) bool {
	if len(c.currentInstructions()) == 0 {
		return false
	}

	return c.scopes[c.scopeIndex].lastInstruction.Opcode == op
}

// Removes last emitted instruction
func (c *Compiler) removeLastInstruction() {
	last := c.scopes[c.scopeIndex].lastInstruction
	previous := c.scopes[c.scopeIndex].previousInstruction

	old := c.currentInstructions()
	new := old[:last.Position]

	c.scopes[c.scopeIndex].Instructions = new
	c.scopes[c.scopeIndex].lastInstruction = previous
}

// Replaces an already emitted instruction to another of the same type
func (c *Compiler) replaceInstruction(pos int, newInstruction []byte) {
	ins := c.currentInstructions()

	for i := 0; i < len(newInstruction); i++ {
		ins[pos+i] = newInstruction[i]
	}
}

// Replaces last instruction with a return instruction
func (c *Compiler) replaceLastInstructionWithReturn() {
	lastPos := c.scopes[c.scopeIndex].lastInstruction.Position
	c.replaceInstruction(lastPos, code.Make(code.OpReturnValue))
	c.scopes[c.scopeIndex].lastInstruction.Opcode = code.OpReturnValue
}

// Changes the operand of an already emitted instruction
func (c *Compiler) changeOperand(opPos int, operand int) {
	op := code.Opcode(c.currentInstructions()[opPos])
	newInstruction := code.Make(op, operand)

	c.replaceInstruction(opPos, newInstruction)
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.currentInstructions(),
		Constants:    c.Constants,
	}
}

type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}
