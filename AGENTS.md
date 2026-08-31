# Guia de Agentes de IA — The Legend of the Go Dragon

Este documento define as diretrizes de desenvolvimento e o contexto de atuação para **todos os agentes de Inteligência Artificial** que colaboram neste projeto.

---

## 🧭 Referências Canônicas Obrigatórias

Antes de gerar código, diálogos ou telas, consulte sempre a documentação especializada:

1. **Universo Narrativo e Lore**: [@docs/universo-e-prompt.md](file:///c:/GitHub/go/lotgd/docs/universo-e-prompt.md)
   - *Fonte da verdade sobre tom de voz, NPCs, bestiário e regras de nomenclatura.*
2. **Game Design Document (GDD)**: [@docs/GDD.md](file:///c:/GitHub/go/lotgd/docs/GDD.md)
   - *Core loop, balanceamento, fórmulas de combate, economia de turnos e mapeamento de telas TUI.*
3. **Arquitetura & Engenharia**: [@docs/architecture.md](file:///c:/GitHub/go/lotgd/docs/architecture.md)
   - *Diagrama de camadas, SQLite CGO-free, The Elm Architecture (Bubble Tea) e concorrência SSH (Wish).*
4. **Diretrizes Globais do Repositório**: [AGENTS.md Raiz](file:///c:/GitHub/go/AGENTS.md)
   - *Padrões do Google Go Style Guide e comentários didáticos explicativos.*

---

## 🛠️ Regra de Ouro da Linguagem e Código

- **Código em Inglês**: Todas as structs, métodos, funções, enums, constantes e nomes de pacotes devem ser 100% em **inglês idiomático**.
- **Interface e Conteúdo em Português do Brasil**: Nenhuma string de exibição para o jogador deve ficar hardcoded no código de domínio ou da TUI; utilize a camada de localização em `internal/i18n`.
- **Código em Produção Limpo e Idiomático**: Ao contrário do diretório raiz do repositório, o código dentro de `lotgd/` não necessita de comentários didáticos explicativos internos. Siga o estilo padrão e limpo da comunidade Go ([Google Go Style Guide](https://google.github.io/styleguide/go/)), mantendo apenas a documentação padrão de tipos e funções públicas exportadas quando necessário.
