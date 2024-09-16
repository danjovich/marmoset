# Marmoset: A compiled programming language for ARM

This is a simple programming language inspired by the one implemented in the book [_Writing A Compiler In Go_](https://compilerbook.com/), written by Thorsten Ball.

Since the language implemented in the book is called "monkey," the name of this language comes from the fact that _marmoset_ is a type of monkey, and the language is compiled for ARM (mARMoset).

## How to use it

To compile the compiler, simply run `make`. You need to have a Go compiler installed, version >=1.20. To install the Go compiler, just follow the instructions on the [official website](https://go.dev/doc/install).

Currently, the Marmoset compiler only prints the assembly generated from the source code. To generate an ELF executable, you will need to use some recipes in the `Makefile`.

Before generating the binaries, the necessary binutils must be installed. If you are on an amd64 machine with APT as the package manager, you need to run:

```bash
sudo apt install binutils-arm-linux-gnueabihf binutils-arm-linux-gnueabihf-dbg
```

Additionally, on an amd64 machine, you need to install qemu for AArch32 in user mode, which can be done with:

```bash
sudo apt install qemu-user qemu-user-static
```

Then, to compile a Marmoset program, simply place the `.marm` file in the examples folder and run:

```bash
make examples/bin/<program_name>.out
```

To run the program, simply execute:

```bash
./examples/bin/<program_name>.out
```

For example, for the provided `fibonacci.marm` program, which prints the seventh, then the first, and then the 15th number in the Fibonacci sequence, calculated through a recursive function, the output should be:

```
13
1
610
```

It is also possible to run a program (with a `.marm` extension and inside the examples folder) directly with a single command `make run-<program_name>`, but `make` will delete the generated assembly and binary afterward (if they were not previously generated). For example, to run `fibonacci.marm` this way, simply run:

```bash
make run-fibonacci
```

## The language

### Variables declaration and assignment

Variables declaration must be done using the reserved keyword `let`. Every variable must be initialized with a value; it is not possible to simply declare them (e.g., `let a;` is not a valid expression in Marmoset). Assignment also occurs with `let`, and shadowing is allowed — identical expressions can act as both declaration and assignment or just assignment. Note that the semicolon (`;`) at the end of each expression is optional.

```
let a = 1; 
let b = 2;
let a = b;
let c = true
```

### Comments

Comments can be written in C++ style:

```
// comment
```

### Functions

Function declaration starts with the reserved keyword `fn`, and function calls are similar to C. There is no need to declare return types or parameter types; everything is treated as either an integer or a boolean. Functions that do not explicitly return any value actually return `null` (equivalent to 0). Function returns do not need to be explicit: if the last expression to be executed is not stored in a variable or used in any way, its value will be returned. Empty returns (`return;`) are not allowed.

```
fn identity(x) { x; }; 
identity(7); // returns 7

fn add(a, b) { a + b }; 
add(1, 2); // returns 3

fn sub(a, b) { a - b }; 
sub(1, 2); // returns -1
```

### Conditionals

Conditionals in Marmoset are expressions, so they return values. If the code inside an `if` (or `else`) ends with an expression without assignment, the value of that expression will be returned; otherwise, the conditional will return `null` (0).

```
if (x < y) { 
    z 
} else { 
    w 
}
	
let a = if (true) {5;} // a receives 5

let b = if (true) {
  if (false) {
    10;
  } else {
    20;
  }
} // b receives 20
```

### Operator Precedence

Operator precedence is the same as in most languages, in increasing order:
  * Equality (`==`) or inequality (`!=`);
  * Less than (`<`) or greater than (`>`);
  * Addition and subtraction;
  * Multiplication, division, and modulus (`%`);
  * Negation (`!x`) or sign inversion (`-x`);
  * Function calls.

In these examples, you can see how the compiler handles precedences (the comments show the result that the parser will generate when analyzing precedence):

```
3 + 4 * 5 == 3 * 1 + 4 * 5 
// ((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))

add(a + b + c * d / f + g % h * i)
// add((((a + b) + ((c * d) / f)) + ((g % h) * i)))
```

### Builtins

The following functions have been implemented:

* `put`: Uses the Linux `write` system call to print a character to the screen (i.e., to stdout). Since there are no chars in Marmoset, it must receive an integer and will print the ASCII character equivalent to that integer.
* `get`: Uses the Linux `read` system call to read a character from the screen (i.e., from stdin). Since there are no chars in Marmoset, it returns an integer equivalent to the ASCII character entered.
* `putint`: To simplify printing integers on the screen (since passing an integer to the `put` function will print the equivalent ASCII character, not the integer itself), this function is provided. It was implemented in Marmoset and compiled, so when it is used by a program, the generated assembly routine is included in the new program's assembly. Its implementation in Marmoset is as follows (note that 45 is the value of the ASCII character "-" and 48 of "0"):

  ```
  fn putint(i) {
    if (i < 0) {
      put(45);
      let i = -i;
    }

    if (i/10) {
      putint(i/10);
    }

    put((i % 10) + 48);
  }
  ```

* `putintln`: Since the `putint` function does not insert a newline (`\n`) after printing the integer, and it is often desired to do so, this builtin function was also implemented as follows (10 is `\n` in ASCII):

  ```
  fn putintln(i) {
    putint(i);
    put(10);
  }
  ```


<!-- BNF draft:
<identifier>      ::= (([a-z] | [A-Z])+ ([a-z] | [A-Z] | [0-9])*)
<integer>         ::= [0-9]+
<boolean>         ::= "true" | "false"
<spaces>          ::= " "+
<optional_spaces> ::= " "*


<expression>      ::= <optional_spaces> (<integer> | <boolean> | <identifier> | <fn_call>) <optional_spaces>
<expression_list> ::= <expression>* | ((<expression> "," <optional_spaces>)* <expression>)*
<fn_call>         ::= <identifier> <optional_spaces> "(" <expression_list> ")"


<let_statement>   ::= "let" <spaces> <identifier> <optional_spaces> "=" <expression> ";"*


<parameter_list>  ::=  <identifier>* | ((<identifier> "," <optional_spaces>)* <identifier>)*
<fn_statement>    ::=  "fn" <spaces> <identifier> <optional_spaces> "(" <parameter_list> ")" <optional_spaces>  "{" <optional_spaces> "}"

<statement>       ::= (<let_statement> | <fn_statement>) ";"*

 -->