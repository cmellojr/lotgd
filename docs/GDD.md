# Game Design Document (GDD) — The Legend of the Go Dragon

> **Documento Vivo de Design de Jogo**  
> **Gênero:** RPG em Modo Texto / BBS Door Game / TUI Multi-usuário  
> **Plataforma:** Terminal CLI / Servidor SSH (Go + Bubble Tea + Lip Gloss + Wish + SQLite)  
> **Referência Narrativa:** [universo-e-lore.md](universo-e-lore.md)

---

## 1. Visão Geral & Core Loop

### 1.1 Premissa & Elevator Pitch
Um vilarejo isolado vive sob a sombra de um Dragão que desperta a cada amanhecer ("o Novo Dia"). Jogadores conectam via terminal/SSH, exploram a Floresta Sombria, aprimoram seus heróis na Guilda e na Ferraria, interagem na Taverna e disputam quem será o herói a derrotar o Dragão antes que o ciclo diário recomece.

### 1.2 Core Gameplay Loop
```mermaid
graph TD
    A[Início do Dia: 15 Turnos na Floresta] --> B[Explorar a Floresta Sombria]
    B --> C{Combate por Turnos}
    C -->|Vitória| D[Ganho de XP e Ouro]
    C -->|Derrota| E[Morte: Envia ao Necrotério, perde ouro solto]
    D --> F[Vilarejo: Ferraria / Capela / Guilda]
    F --> G{Tem Nível para o Dragão?}
    G -->|Sim| H[Desafiar o Dragão]
    G -->|Não| B
    H -->|Venceu| I[Hall da Fama + Reinício do Ciclo do Dragão]
    H -->|Perdeu| E
```

---

## 2. Mecânicas de Jogo & Balanceamento

### 2.1 Economia de Recursos
- **Turnos Diários (*Forest Fights*)**: 15 lutas/dia por jogador (renovadas no "Novo Dia").
- **Pontos de Vida (*HP*)**: Recuperados na Capela (Frei Anselmo) mediante oferenda ou descanso na Pousada.
- **Ouro (*Gold*)**: Carregado na bolsa (em risco se morrer) ou guardado no Banco.
- **Experiência (*XP*)**: Necessária para subir de nível e ser promovido pelo *Mestre Tobias* na Guilda.

### 2.2 Sistema de Combate por Turnos
- **Fórmula de Dano**:
  $$\text{Dano} = \max(1, (\text{ATK}_{\text{atacante}} + \text{Rnd}(1, 4)) - \text{DEF}_{\text{defensor}})$$
- **Acerto Crítico**:
  - Há **10% de chance** em cada ataque de rolar um **Acerto Crítico**. Quando ocorre, o valor do ataque efetivo ($\text{ATK}_{\text{atacante}} + \text{Rnd}(1, 4)$) é multiplicado por **1,5** (+50% de dano base) antes de subtrair a defesa do oponente.
- **Ações no Turno**:
  - `[A]tacar`: Desfere golpe corpo a corpo aplicando a fórmula de dano e rolagem de crítico.
  - `[F]ugir`: Tenta recuar estrategicamente do combate.
    - **Chance base de sucesso**: 50%.
    - **Monstros com afixo "Covarde"**: 80% de chance de sucesso.
    - **Chefe Dragão do Dia**: 20% de chance de sucesso.
    - **Penalidade por falha**: Se a tentativa de fuga falhar, o inimigo recebe um contra-ataque livre de oportunidade.
  - `[P]oção`: Usa consumível (*Poção de Vida*) restaurando até +30 HP sem gastar o turno livre.

---

## 3. Estrutura de Telas & Mapeamento TUI (Bubble Tea)

| Tela | Rota TUI | Comandos / Atalhos Principais |
|---|---|---|
| **Login / Criação** | `ScreenLogin` | Digitar nome/senha, `[Enter]` confirma, `[Tab]` alterna |
| **Praça do Vilarejo** | `ScreenTown` | `[F]` Floresta, `[T]` Taverna, `[C]` Capela, `[M]` Ferraria, `[G]` Guilda, `[D]` Dragão, `[S]` Salvar e Sair (Logout) |
| **Floresta Sombria** | `ScreenForest` | `[P]` Procurar monstro, `[A]` Atacar, `[F]` Fugir, `[V]` Voltar à vila |
| **Ferraria (Torin)** | `ScreenSmith` | `[1]` Armas, `[2]` Armaduras, `[3]` Poções, `[Tab]` Trocar aba |
| **Capela (Anselmo)** | `ScreenChapel` | `[C]` Curar ferimentos, `[D]` Doação (10ouro), `[M]` Meditar |
| **Taverna (Rosalinda)**| `ScreenTavern` | `[R]/[F]` Ouvir fofocas, `[C]` Flertar com Cassandra |
| **Covil do Dragão** | `ScreenDragon` | `[D]` Desafiar o Dragão do Dia |
| **Guilda dos Aventureiros** | `ScreenGuild` | `[A]/[T]` Avançar/consultar, `[R]/[C]` Consultar tabela, `[P]/[L]` Ler pergaminhos |

---

## 4. Arquitetura e Engenharia em Go

- **`internal/engine`**: Lógica pura de regras, structs de domínio em inglês (`Player`, `CombatEngine`, `Stats`).
- **`internal/bestiary`**: Tiers de monstros (1 a 4) com gerador procedural de prefixos (*"Feroz"*, *"Covarde"*, etc.).
- **`internal/i18n`**: Dicionários `pt_br.go` isolando textos da interface.
- **`internal/storage`**: SQLite puro (`modernc.org/sqlite`) para multi-usuário com persistência em disco.
- **`cmd/lotgd` & `cmd/server`**: Suporte unificado tanto para cliente local quanto servidor SSH via Wish.
