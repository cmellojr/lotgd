# i18n

## Visão Geral
Catálogo de internacionalização responsável pelas traduções e textos em Português do Brasil (pt-BR) para locais do vilarejo, falas de NPCs, descrições do bestiário, menus e narrativas de combate.

## Interface
- Identificadores fortemente tipados (`types.go`): `LocationID`, `NPCID`, `MonsterID`, `ItemID`, `AffixID`, `UIKey`, `MessageKey`.
- Funções de consulta e formatação (`pt_br.go`):
  - `func GetLocationName(id LocationID) string`
  - `func GetNPCName(id NPCID) string`
  - `func GetMonsterName(id MonsterID) string`
  - `func GetItemName(id ItemID) string`
  - `func GetItemDescription(id ItemID) string`
  - `func GetAffixName(id AffixID) string`
  - `func GetUIText(key UIKey) string`
  - `func GetMessage(key MessageKey, args ...any) string`
- Dicionários (`pt_br.go`): `LocationNamesPTBR`, `NPCNamesPTBR`, `MonsterNamesPTBR`, `DragonTitlesPTBR`, `AffixNamesPTBR`, `UITextsPTBR`, `MessageTemplatesPTBR`, `ItemNamesPTBR`, `ItemDescriptionsPTBR`.

## Invariantes
- Todo o texto visível ao usuário permanece isolado da lógica de domínio e das telas TUI.
- Chaves ausentes no dicionário retornam valores de fallback sem causar pânico em tempo de execução.

## Efeitos Colaterais
- Nenhum. Todas as funções operam como consultas puras em tabelas de tradução em memória.

## Verificação
- Executar: `go test ./internal/i18n/...`
