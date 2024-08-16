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


## Progresso de desenvolvimento

- [x] Implementar funções recursivas;
- [x] Mudar declaração de funções para o nome vir depois de `fn`;
  - [x] Mudar compilação de `let` para evitar que referência a si mesmo cause "null pointer dereference";
  - [x] Para facilitar compilação para ARM, fazer funções não serem mais "cidadãs de primeira classe" (impedir uso como argumento de outras funções).
- [x] Remover implementação de hashes;
- [x] Criar built-ins `get` e `put`;
- [x] Compilar arquivos;
- [x] Retornar código de erro em erros de compilação;
- [ ] Remover VM/REPL;
- [ ] Escrever em arquivo ao invés de imprimir;
- [ ] Fazer teste do arm_compiler.go (fibonacci e função separada).