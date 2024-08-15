# Marmoset: Uma linguagem de programação compilada para ARM

Essa é uma linguagem de programação simples inspirada na implementada no livro [_Writing A Compiler In Go_](https://compilerbook.com/), escrito por Thorsten Ball.

Já que a linguagem implementada no livro se chama "monkey", o nome da linguagem vem do fato de que _marmoset_ (sagui) é uma espécie de macaco e a linguagem é compilada para ARM (mARMoset).

## Progresso de desenvolvimento

- [x] Implementar funções recursivas;
- [x] Mudar declaração de funções para o nome vir depois de `fn`;
  - [x] Mudar compilação de `let` para evitar que referência a si mesmo cause "null pointer dereference";
  - [x] Para facilitar compilação para ARM, fazer funções não serem mais "cidadãs de primeira classe" (impedir uso como argumento de outras funções).
- [ ] Talvez remover retornos implícitos;
  - [ ] Nesse caso, remover opcode de "pop".
- [x] Remover implementação de hashes;
- [ ] Criar built-ins `getc` e `putc`;
- [x] Compilar arquivos;
- [ ] Se sobrar tempo, implementar `while`;
- [x] Retornar código de erro em erros de compilação;
- [ ] Remover VM/REPL;
- [ ] Escrever em arquivo ao invés de imprimir;
- [ ] Fazer teste do arm_compiler.go (fibonacci e função separada).