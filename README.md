# Marmoset: Uma linguagem de programação compilada para ARM

Essa é uma linguagem de programação simples inspirada na implementada no livro [_Writing A Compiler In Go_](https://compilerbook.com/), escrito por Thorsten Ball.

Já que a linguagem implementada no livro se chama "monkey", o nome da linguagem vem do fato de que _marmoset_ (sagui) é uma espécie de macaco e a linguagem é compilada para ARM (mARMoset).

## Progresso de desenvolvimento

- [x] Implementar funções recursivas;
- [ ] Tipagem estática;
  - [ ] Strings são apenas arrays de inteiros tratados como `char`'s;
    - [ ] Sem suporte para Unicode.
  - [ ] Estilo C: "tudo" é um `int`.
- [ ] Mudar declaração de funções para o nome vir depois de `fn`;
  - [ ] Mudar compilação de `let` para evitar que referência a si mesmo cause "null pointer dereference";
  - [ ] Para facilitar compilação para ARM, fazer funções não serem mais "cidadãs de primeira classe" (impedir uso como argumento de outras funções).
- [ ] Talvez remover retornos implícitos;
- [ ] Remover implementação de hashes;
- [ ] Criar built-in `gets`;
- [ ] Se sobrar tempo, implementar `while`.