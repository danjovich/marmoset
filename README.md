# 🐒 Monkey compiler

This is the interpreter from the book [Writing A Compiler In Go](https://compilerbook.com/), by Thorsten Ball.

## Progress

- [x] Introduction
  - [x] Evolving Monkey
    - [x] The Past and Present
    - [x] The Future
  - [x] Use This Book
- [ ] 1. Compilers & Virtual Machines
  - [ ] Compilers
  - [ ] Virtual and Real Machines
    - [ ] Real Machines
    - [ ] What Is a Virtual Machine?
    - [ ] Why Build One?
    - [ ] Bytecode
  - [ ] What We’re Going to Do, or: the Duality of VM and Compiler
- [ ] 2. Hello Bytecode!
  - [ ] First Instructions
    - [ ] Starting With Bytes
    - [ ] The Smallest Compiler
    - [ ] Bytecode, Disassemble!
    - [ ] Back to the Task at Hand
    - [ ] Powering On the Machine
  - [ ] Adding on the Stack
  - [ ] Hooking up the REPL
- [ ] 3. Compiling Expressions
  - [ ] Cleaning Up the Stack
  - [ ] Infix Expressions
  - [ ] Booleans
  - [ ] Comparison Operators
  - [ ] Prefix Expressions
- [ ] 4. Conditionals
  - [ ] Jumps
  - [ ] Compiling Conditionals
  - [ ] Executing Jumps
  - [ ] Welcome Back, Null!
- [ ] 5. Keeping Track of Names
  - [ ] The Plan
  - [ ] Compiling Bindings
    - [ ] Introducing: the Symbol Table
    - [ ] Using Symbols in the Compiler
  - [ ] Adding Globals to the VM
- [ ] 6. String, Array and Hash
  - [ ] String
  - [ ] Array
  - [ ] Hash
  - [ ] Adding the index operator
- [ ] 7. Functions
  - [ ] Dipping Our Toes: a Simple Function
    - [ ] Representing Functions
    - [ ] Opcodes to Execute Functions
    - [ ] Compiling Function Literals
    - [ ] Compiling Function Calls
    - [ ] Functions in the VM
    - [ ] A Little Bonus
  - [ ] Local Bindings
    - [ ] Opcodes for Local Bindings
    - [ ] Compiling Locals
    - [ ] Implementing Local Bindings in the VM
  - [ ] Arguments
    - [ ] Compiling Calls With Arguments
    - [ ] Resolving References to Arguments
    - [ ] Arguments in the VM
- [ ] 8. Built-in Functions
  - [ ] Making the Change Easy
  - [ ] Making the Change: the Plan
  - [ ] A New Scope for Built-in Functions
  - [ ] Executing built-in functions
- [ ] 9. Closures
  - [ ] The Problem
  - [ ] The Plan
  - [ ] Everything’s a closure
  - [ ] Compiling and resolving free variables
  - [ ] Creating real closures at run time
  - [ ] Recursive Closures
- [ ] 10. Taking Time