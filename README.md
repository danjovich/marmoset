# Marmoset: Uma linguagem de programação compilada para ARM

Essa é uma linguagem de programação simples inspirada na implementada no livro [_Writing A Compiler In Go_](https://compilerbook.com/), escrito por Thorsten Ball.

Já que a linguagem implementada no livro se chama "monkey", o nome da linguagem vem do fato de que _marmoset_ (sagui) é uma espécie de macaco e a linguagem é compilada para ARM (mARMoset).

## Como usar

Para compilar o compilador, basta rodar `make`. É necessário ter o copilador de go, com versão >=1.20 instalado.

Até o momento, o compilador apenas imprime na tela o assembly gerado a partir do código fonte. Para gerar um executável ELF, então, é necessário usar o `make`.

Antes de gerar ps binários, é preciso ter instalados os binutils necessários. Se você está numa máquina com CPU amd64 e com o APT como package manager, é preciso rodar:

```bash
sudo apt install binutils-arm-linux-gnueabihf binutils-arm-linux-gnueabihf-dbg
```

Além disso, é necessário, em uma máquina amd64, instalar o qemu para AArch32 em modo usuário, isso pode ser feito com:

```bash
sudo apt install qemu-user qemu-user-static
```

Em seguida, para compilar um programa em Marmoset, basta colocar o arquivo .marm na pasta examples e rodar:

```bash
make examples/bin/<nome_do_programa>.out
```

Para rodar o programa, então, basta rodar:

```bash
./examples/bin/<nome_do_programa>.out
```

Por exemplo, para o programa `fibonacci.marm` fornecido, que imprime na tela o sétimo, depois o primeiro, depois o 15º número da sequência de Fibonacci calculados através de uma função recursiva, a saída deve ser:

```
13
1
610
```

## A linguagem

### Declaração e atribuição de variáveis

A declaração de variáveis deve ser feita com a palavra reservada `let`. Toda variável deve ser inicializada com um valor, não é possível simplesmente declará-las (por exemplo `let a;` não é uma expressão válida em Marmoset). A atribuição também ocorre com o `let`, e é permitido o shadowing, expressões idênticas podem atuar como declaração junto a atribuição ou apenas atribuição. Note que o ponto e vírgula (`;`) ao fim de cada expressão é opcional.

```
let a = 1; 
let b = 2;
let a = b;
let c = true
```

### Comentários

Comentários podem ser feitos no estilo de C++:

```
// comentário
```

### Funções

A declaração de funções começa com a palavra reservada `fn` e a chamada ocorre como em C. Não é necessário declarar tipos de retorno ou dos parâmetros, tudo é tratado como inteiro ou booleano. Funções que não retornam nenhum valor explicitamente na verdade retornam `null` (que equivale a 0). Os retornos de função não precisam ser explícitos: se a última expressão a ser executada não for guardada em um identificador ou usada de alguma forma, o valor dela será retornado. Returns vazios (`return;`) não são permitidos.

```
fn identity(x) { x; }; 
identity(7); // Retorna 7

fn add(a, b) { a + b }; 
add(1, 2); // retorna 3

fn sub(a, b) { a - b }; 
sub(1, 2); // retorna -1
```

### Condicionais

Os condicionais em Marmoset são expressões, então eles retornam valores. Se o código dentro de um `if` (ou `else`) terminar com uma expressão sem atribuição, o valor dessa expressão será retornado, se não, o condicional retornará `null` (0).

```
if (x < y) { 
    z 
} else { 
    w 
}
	
let a = if (true) {5;} // a recebe 5

let b = if (true) {
	if (false) {
		10;
	} else {
		20;
	}
} // b recebe 20
```

### Precedência de Operadores

A precedência de operadores é a mesma que a encontrada em outras linguagens, em ordem crescente:
  * Comparação de igual (`==`) ou diferente (`!=`);
  * Comparações de menor (`<`) ou maior (`>`);
  * Somas e subtrações;
  * Multiplicações, divisões e resto (`%`);
  * Negação (`!x`) ou inversão de sinal (`-x`);
  * Chamadas de funções.

Nesses exemplos, é possível ver um pouco de como o compilador trata as precedências (o que é mostrado nos comentários é realmente o resultado que o parser irá gerar quando analisar a precedência):

```
3 + 4 * 5 == 3 * 1 + 4 * 5 
// ((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))

add(a + b + c * d / f + g % h * i)
// add((((a + b) + ((c * d) / f)) + ((g % h) * i)))
```

### Builtins

Foram implementadas outras funções:

* `put`: Usa a chamada de sistema `write` do Linux para escrever um caractere na tela (isto é, no stdout). Como não há chars em Marmoset, ela deve receber um inteiro e portanto irá imprimir o caractere ASCII equivalente a esse inteiro.
* `get`: Usa a chamada de sistema `read` do Linux para ler um caractere da tela (isto é, do stdin). Como não há chars em Marmoset, retorna um inteiro equivalente ao caractere ASCII digitado.
* `putint`: Para facilitar a impressão de inteiros na tela, já que passar um inteiro para a função `put` irá imprimir o caractere ASCII equivalente, não o inteiro, essa função é fornecida. Ela foi implementada em Marmoset e compilada, e então quando for usada por um programa a rotina gerada em assembly é incluída no assembly do novo programa. A implementação dela em Marmoset é a seguinte (note que 45 é o valor do caractere ASCII “-” e 48 do “0”):

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

* `putintln`: Como a função `putint` não insere um `\n` depois de imprimir o inteiro, e é frequente que se queira fazer isso, foi implementada também essa função builtin, da seguinte forma (10 é `\n` em ASCII):

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