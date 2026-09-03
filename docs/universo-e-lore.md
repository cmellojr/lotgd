# The Legend of the Go Dragon — Universo e Lore

> **Nota de uso:** este arquivo é a fonte da verdade sobre o *universo narrativo* do jogo.
> Referencie-o a partir do `AGENTS.md` (ex: `@docs/universo-e-lore.md`) para que qualquer
> agente de IA — Copilot, OpenCode, Antigravity/Gemini, Jules — tenha o mesmo contexto de
> lore ao gerar diálogos, textos de evento, nomes de itens ou telas do jogo.
> O `AGENTS.md` continua sendo o lugar certo para comandos de build, convenções de código
> e regras técnicas; este arquivo cobre apenas **mundo, tom e nomenclatura narrativa**.

---

## 1. Prompt de contexto (cole isto no início de sessões de IA focadas em conteúdo)

```
Você está ajudando a desenvolver "The Legend of the Go Dragon", um RPG de texto para
terminal (TUI), escrito em Go, inspirado nos clássicos jogos de porta de BBS dos anos
1980-90 (Legend of the Red Dragon e derivados como Legend of the Green Dragon).

O jogo é multiusuário e persistente: vários jogadores compartilham o mesmo vilarejo,
o mesmo Dragão do dia e o mesmo ranking, conectando via SSH.

Contexto de idioma, IMPORTANTE:
- Toda a INTERFACE visível ao jogador (textos, menus, diálogos, mensagens de sistema)
  deve ser escrita em português do Brasil, com tom leve, nostálgico e levemente cômico
  — nunca sombrio ou gráfico demais. É um jogo de fantasia cordial, não grimdark.
- Todo o CÓDIGO e dados internos (nomes de tipos, campos de struct, chaves de mapas,
  identificadores de constantes/enums) devem permanecer em inglês, seguindo convenção
  idiomática de Go. Textos em português nunca devem aparecer como identificador de
  código — sempre como valor de string, carregado de uma camada de i18n.
- Use os nomes de personagens, monstros e locais definidos neste documento como
  referência canônica. Não invente nomes de licenças existentes (D&D, Warcraft,
  Diablo etc.) — o universo é original, apenas inspirado nos arquétipos clássicos
  de fantasia e no espírito nostálgico dos BBS door games.
```

---

## 2. Premissa

Um vilarejo isolado vive sob a sombra de um Dragão que desperta a cada amanhecer
("o Novo Dia"). Aventureiros saem para caçar na Floresta Sombria, treinam na Guilda,
bebem na taverna, se curam na capela e, ao fim do dia, um único herói pode desafiar
o Dragão — ou todos falham e o ciclo recomeça amanhã.

## 3. O vilarejo e seus NPCs

| Papel | Nome (PT-BR) | Função no jogo |
|---|---|---|
| Taverneira | **Dona Rosalinda** | Fofocas do dia, rumores sobre o Dragão, ponto de encontro social |
| Curandeiro | **Frei Anselmo** | Cura HP mediante oferenda; diálogos de sabedoria/humor |
| Ferreiro | **Mestre Torin** | Compra/venda e melhoria de armas e armaduras |
| Mercadora ambulante | **Yolanda, a Cigana** | Poções, itens raros, "sorte do dia" |
| Interesse romântico | **Cassandra** | Mecânica de cortejo — vários heróis competem por sua atenção |
| Rival recorrente | **Cavaleiro Vermelho** | Duelos PvP amistosos, zomba do jogador no ranking |
| Bibliotecário da Guilda | **Mestre Tobias** | Explica regras, tutorial, lore adicional sob demanda |

## 4. Bestiário (Floresta Sombria, por nível de dificuldade)

**Nível 1 — Iniciante**
Rato-do-Esgoto · Aranha-de-Musgo · Kobold Trapalhão · Slime Verde · Corvo Sinistro

**Nível 2 — Aprendiz**
Goblin Batedor · Orc Recruta · Lobo das Sombras · Bandido de Estrada · Esqueleto Enferrujado

**Nível 3 — Veterano**
Troll da Ponte · Ogro Malcriado · Necromante Iniciante · Harpia Cantora · Golem de Pedra Rachada

**Nível 4 — Lendário**
Cavaleiro Caído · Wyvern da Montanha · Lich Menor · Basilisco de Olhar Frio · Espectro do Pântano

> Sugestão de mecânica (como no LORD original): combine um **prefixo aleatório**
> ("Feroz", "Covarde", "Enfurecido", "Sortudo") com o nome-base pra gerar variações
> sem precisar de um bestiário infinito escrito à mão.

## 5. O Dragão

**Nome:** o Dragão não tem nome fixo revelado — é chamado apenas de **"o Dragão"**
pelos moradores, com reverência/medo. Isso preserva o mistério e evita choque de tom
com o easter egg opcional abaixo.

Cada "Novo Dia" ele renasce com atributos sorteados (HP, ataque, defesa), então
narrativamente ele é sempre o mesmo Dragão eterno, não um monstro genérico.

## 6. Antagonista secundário (arco de médio prazo, opcional)

**O Arquiteto das Sombras** — uma figura encapuzada que aparece em rumores da
taverna, sugerindo que o ciclo do Dragão não é natural, mas mantido por um pacto
antigo. Bom gancho pra conteúdo futuro (quests, final alternativo) sem comprometer
o loop principal agora.

## 7. Easter eggs para devs (opcional, use com moderação)

Camada extra só pra quem conhece o contexto do projeto — nunca deve substituir os
nomes "sérios" acima na experiência principal, apenas aparecer em itens raros,
conquistas ou mensagens secretas:

- Poção rara: **"Elixir do Garbage Collector"** (cura total, "libera memória perdida")
- Investida especial do Dragão: **"Fúria da Goroutine"** (ele "se multiplica" e ataca em dobro)
- Item de humor: **"Ponteiro Nulo Amaldiçoado"** (arma que às vezes "erra o golpe" — pane cômica)
- Conquista secreta: **"Zero Downtime"** — vencer o Dragão sem perder HP

## 8. Convenção técnica: identificadores em inglês, texto em português

Exemplo de padrão para manter código idiomático em Go e conteúdo em PT-BR:

```go
// internal/bestiary/monster.go — identificadores sempre em inglês
type MonsterID string

const (
    MonsterSewerRat   MonsterID = "sewer_rat"
    MonsterMossSpider MonsterID = "moss_spider"
    MonsterGoblinScout MonsterID = "goblin_scout"
)

// internal/i18n/pt_br.go — strings de exibição carregadas à parte
var MonsterNamesPTBR = map[MonsterID]string{
    MonsterSewerRat:    "Rato-do-Esgoto",
    MonsterMossSpider:  "Aranha-de-Musgo",
    MonsterGoblinScout: "Goblin Batedor",
}
```

Essa separação deixa a base pronta pra outros idiomas no futuro (mesmo que hoje só
exista `pt_br.go`) e evita que texto em português vaze pro meio da lógica de jogo.

## 9. Tom e diretrizes de escrita

- Leve, nostálgico, ocasionalmente cômico — nunca violência gráfica ou humor cruel.
- Frases curtas, no estilo de texto de terminal dos anos 90 (evitar parágrafos longos).
- Evitar qualquer referência direta a obras protegidas (D&D, Warcraft, Tolkien etc.);
  arquétipos genéricos de fantasia (goblin, troll, dragão) são livres de uso.
- Consistência de nomes: sempre os mesmos NPCs/nomes deste documento — se a IA
  precisar de um nome novo, deve seguir o mesmo padrão (nome próprio + título/ofício).