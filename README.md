# 🐒 Monkey compiler

This is the interpreter from the book [Writing A Compiler In Go](https://compilerbook.com/), by Thorsten Ball.

## Progress

- [x] Introduction
  - [x] Evolving Monkey
    - [x] The Past and Present
    - [x] The Future
  - [x] Use This Book
- [x] 1. Compilers & Virtual Machines
  - [x] Compilers
  - [x] Virtual and Real Machines
    - [x] Real Machines
    - [x] What Is a Virtual Machine?
    - [x] Why Build One?
    - [x] Bytecode
  - [x] What We’re Going to Do, or: the Duality of VM and Compiler
- [x] 2. Hello Bytecode!
  - [x] First Instructions
    - [x] Starting With Bytes
    - [x] The Smallest Compiler
    - [x] Bytecode, Disassemble!
    - [x] Back to the Task at Hand
    - [x] Powering On the Machine
  - [x] Adding on the Stack
  - [x] Hooking up the REPL
- [x] 3. Compiling Expressions
  - [x] Cleaning Up the Stack
  - [x] Infix Expressions
  - [x] Booleans
  - [x] Comparison Operators
  - [x] Prefix Expressions
- [x] 4. Conditionals
  - [x] Jumps
  - [x] Compiling Conditionals
  - [x] Executing Jumps
  - [x] Welcome Back, Null!
- [x] 5. Keeping Track of Names
  - [x] The Plan
  - [x] Compiling Bindings
    - [x] Introducing: the Symbol Table
    - [x] Using Symbols in the Compiler
  - [x] Adding Globals to the VM
- [x] 6. String, Array and Hash
  - [x] String
  - [x] Array
  - [x] Hash
  - [x] Adding the index operator
- [ ] 7. Functions
  - [ ] Dipping Our Toes: a Simple Function
    - [x] Representing Functions
    - [x] Opcodes to Execute Functions
    - [x] Compiling Function Literals
    - [x] Compiling Function Calls
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