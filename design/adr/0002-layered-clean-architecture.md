# ADR-0002: Adopt Layered Clean Architecture and English Identifiers

- Status: Approved
- Date: 2026-09-05
- Author(s): Equipe de Desenvolvimento
- Deciders: Mantenedores do The Legend of the Go Dragon

## 1. Context
O *The Legend of the Go Dragon* combina múltiplos domínios funcionais: combate de RPG, ciclo de dias, persistência relacional, localização para múltiplos idiomas e renderização em terminal com protocolo SSH. Sem uma demarcação clara de responsabilidades, o código rapidamente acumula acoplamento entre regras de negócio, queries de banco e códigos de formatação ANSI / escape sequences do terminal. Além disso, misturar termos de domínio em português nos identificadores de Go enfraquece o uso de tooling padrão do ecossistema Go e convenções globais da comunidade.

## 2. Decision
Adotamos uma arquitetura modular em camadas estritas organizada na estrutura canônica do Go:
- `cmd/lotgd` e `cmd/server`: Pontos de entrada para execução CLI local e servidor SSH BBS.
- `internal/engine`: Camada de domínio puro (entidades, atributos, `CombatEngine`, inventário e lógica de turnos) isolada de efeitos colaterais e banco de dados.
- `internal/bestiary`: Gerador estocástico de monstros e afixos, desacoplado da apresentação.
- `internal/storage`: Camada de persistência relacional e repositórios.
- `internal/tui`: Apresentação baseada em The Elm Architecture (Bubble Tea e Lip Gloss).
- `internal/i18n`: Dicionários e mensagens para os jogadores.

**Regra de Linguagem**:
1. Todo o código em Go (nomes de structs, interfaces, métodos, variáveis, testes, constantes e pacotes) é estritamente em **inglês idiomático**.
2. Nenhuma string voltada ao usuário pode ficar hardcoded na lógica de negócio ou modelos de dados.

Alternativas consideradas:
- **Monólito único em package `main`**: Descartado devido à impossibilidade de testar o domínio sem disparar a interface de terminal e por inviabilizar o reuso entre CLI e servidor SSH.
- **Identificadores em Português**: Descartado pois contraria o Google Go Style Guide, dificulta linters e gera inconsistência léxica com a biblioteca padrão (`context.Context`, `net.Listener`, `fmt.Sprintf`).

## 3. Consequences
### Positivas
- A lógica de combate (`CombatEngine`) é determinística e testável de forma unitária sem mock de terminal ou banco.
- Facilidade de reutilizar 100% da lógica e telas tanto no cliente local quanto no daemon SSH.
- Suporte a linters e ferramentas padrão de Go sem atrito.

### Negativas / Restrições
- Necessidade de conversões/mapeamentos explícitos entre modelos de storage e modelos de domínio (`storage.Player` vs `engine.Player`).
- Exige rigor contínuo para evitar que chamadas de I/O vazem para a camada de domínio.

## 4. Compliance and Verification
- Verificação automatizada de testes com `go test ./internal/engine/...` garantindo zero dependências de `internal/tui` ou `internal/storage`.
- Linters de análise estática e CI gates impedem import cycles (`go vet ./...`).
