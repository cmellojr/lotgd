# ADR-0005: Centralize Portuguese (PT-BR) Localization in Dedicated i18n Layer

- Status: Approved
- Date: 2026-09-05
- Author(s): Equipe de Desenvolvimento
- Deciders: Mantenedores do The Legend of the Go Dragon

## 1. Context
O *The Legend of the Go Dragon* adota o português do Brasil (PT-BR) como idioma de interface oficial, trazendo tom acolhedor, nostálgico e levemente bem-humorado, conforme estabelecido em `docs/universo-e-lore.md`.

No entanto, o código de engenharia deve seguir convenções universais em inglês para tipos, funções e interfaces. Se as strings em português forem espalhadas diretamente dentro da lógica do motor de combate (`internal/engine`) ou em repositórios de dados (`internal/storage`), ocorrem os seguintes problemas:
- Mistura de responsabilidades (lógica e apresentação textual).
- Dificuldade em manter consistência terminológica (nomes de monstros, itens e NPCs).
- Impossibilidade de internacionalizar o jogo no futuro para novos idiomas (como inglês ou espanhol).

## 2. Decision
Adotamos uma camada dedicada de localização centralizada em `internal/i18n`:

1. **Constantes e IDs em Inglês no Domínio**: As entidades do jogo possuem identificadores canônicos (ex: `MonsterSewerRat = "sewer_rat"`, `WeaponWoodenSword = "wooden_sword"`).
2. **Mapeamentos Dicionarizados em `internal/i18n`**: A camada de localização associa os identificadores a suas respectivas strings e templates formatados em PT-BR (ex: `MonsterNamesPTBR`, `ItemDescriptionsPTBR`).
3. **Proibição de Strings Hardcoded de Exibição**: Nenhuma lógica de domínio ou camada de dados deve formatar ou armazenar texto final de interface; a TUI consome as mensagens exclusivamente a partir do pacote `i18n`.

Alternativas consideradas:
- **Strings diretamente nas structs de modelo (`engine.Monster.Name = "Rato do Esgoto"`)**: Rejeitado por acoplar o idioma da UI à estrutura de dados do jogo e inviabilizar troca dinâmica de idioma.
- **Arquivos externos JSON/YAML/PO com runtime de templates**: Rejeitado por adicionar parsing em tempo de execução e risco de arquivos ausentes na distribuição de binário estático. Dicionários estáticos em Go compilados diretamente no binário oferecem type safety, zero dependência de filesystem externo e performance imediata.

## 3. Consequences
### Positivas
- Todo o conteúdo textual do jogo fica centralizado e fácil de revisar editorialmente por redatores e designers de narrativa.
- A base fica pronta para adicionar novos idiomas (ex: `en_us.go`) sem tocar em uma única linha de regras de combate.
- Compilação estática com verificação em tempo de compilação de chaves de tradução.

### Negativas / Restrições
- Requer disciplina contínua dos desenvolvedores e agentes para não inserir strings diretas em novas telas ou mensagens.
- Adicionar um novo monstro ou item requer um passo a mais (definir o ID no domínio e o texto correspondente no dicionário).

## 4. Compliance and Verification
- Regra de ouro documentada no `AGENTS.md`.
- Code review e varreduras com linters procurando por strings literais de mensagens em arquivos sob `internal/engine/`.
- Testes unitários em `internal/i18n` validando que todos os IDs registrados no bestiário e catálogo de itens possuem mapeamento correspondente.
