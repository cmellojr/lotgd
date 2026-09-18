# Relatório de Avaliação Criteriosa de Pendências Reportadas

**Data**: 2026-09-18
**Repositório**: `lotgd` (`develop` branch)
**Avaliado por**: Jules (Agente de IA)

Este documento registra a avaliação criteriosa das 3 pendências reportadas anteriormente em relação ao estado atual do repositório na branch `develop`.

---

## 1. Status do ADR-0007 (`design/adr/0007-equipment-sale-and-trade-in-credit.md`)

- **Pendência Reportada**: O ADR-0007 ainda estaria com `Status: Proposed` e precisaria ter seu status atualizado para `Approved`.
- **Resultado da Avaliação**: **Já resolvido na branch `develop`**.
- **Detalhamento**: Verificação no arquivo `design/adr/0007-equipment-sale-and-trade-in-credit.md` confirma que o cabeçalho já está definido como `- Status: Approved` na linha 3 do documento.

---

## 2. Opção de Recusa da Oferta de Troca (Opção D do ADR-0007)

- **Pendência Reportada**: A opção de recusar a oferta de troca não teria sido implementada.
- **Resultado da Avaliação**: **Parcialmente implementado / Esclarecimento de fluxo**.
- **Detalhamento**:
  - **Venda Avulsa (`[V]`)**: O fluxo de negociação com recusa de oferta **está totalmente implementado** em `internal/tui/screens/smith.go`. Ao pressionar a tecla `[V]`, Mestre Torin gera uma cotação flutuante via `engine.CalculateTradeInQuote` e entra no estado de confirmação (`confirmingSale`). Se o jogador pressionar `[N]`, `[R]` ou `[ESC]`, a função `handleConfirmSale(false)` é acionada, exibindo a mensagem `MsgSmithSellOfferRejected` e cancelando a transação. Uma nova tentativa em `[V]` gera uma nova cotação.
  - **Crédito de Troca na Compra Direta**: Ao comprar um item novo via catálogo (`ENTER`), o abatimento do equipamento equipado é calculado e aplicado automaticamente no custo líquido (`netCost`). Não é exibido um modal intermediário perguntando "Aceita X moedas pelo seu item atual?" durante a compra direta.
- **Recomendação**:
  - A mecânica atual atende os requisitos de reciclagem e venda de equipamentos. Caso seja desejado um modal de recusa/aceitação explícita de trade-in credit no ato da compra de novos itens, essa funcionalidade pode ser mapeada como uma melhoria UX em issue de acompanhamento para a milestone v0.0.3.

---

## 3. Tabela "Achados que permanecem abertos" em `docs/plano-fixes-auditoria.md`

- **Pendência Reportada**: A tabela estaria desatualizada e precisaria de uma nota indicativa para evitar que seja consultada como fonte primária da verdade.
- **Resultado da Avaliação**: **Corrigido**.
- **Detalhamento**: Foi adicionada a nota de aviso `> ⚠️ **Nota**: Tabela desatualizada, ver CHANGELOG.md para o status real das correções.` diretamente acima da tabela no arquivo `docs/plano-fixes-auditoria.md`.

---
