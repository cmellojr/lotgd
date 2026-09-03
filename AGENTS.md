# Guia de Agentes de IA — The Legend of the Go Dragon

Diretrizes de desenvolvimento e contexto de atuação para agentes de IA no projeto.

---

## 🧭 Referências Canônicas

Antes de gerar código ou conteúdo, consulte os documentos especializados:

- **Universo e Lore**: [`docs/universo-e-lore.md`](docs/universo-e-lore.md) — Tom de voz, NPCs, bestiário e nomenclatura.
- **Game Design Document**: [`docs/GDD.md`](docs/GDD.md) — Core loop, combate, economia de turnos e fluxo TUI.
- **Arquitetura & Engenharia**: [`docs/architecture.md`](docs/architecture.md) — Camadas, SQLite (CGO-free), Bubble Tea e Wish (SSH).

---

## 🛠️ Diretrizes Principais

1. **Código em Inglês**: Structs, métodos, funções, enums, variáveis e pacotes devem ser 100% em **inglês idiomático**.
2. **Interface em Português (PT-BR)**: Textos e mensagens para o jogador devem ser gerenciados pela camada de i18n (`internal/i18n`), sem strings de exibição hardcoded na lógica de domínio ou TUI.
3. **Estilo Go Idiomático**: Siga as recomendações do [Google Go Style Guide](https://google.github.io/styleguide/go/). Mantenha a documentação Go padrão (`godoc`) para tipos e funções exportadas.
