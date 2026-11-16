# internal

## Visão Geral
Pacotes centrais da aplicação que implementam a lógica de domínio do RPG, persistência em banco de dados, interface de usuário em modo texto, geradores procedurais e dicionários de localização.

## Subdiretórios e Fronteiras
- [bestiary](./bestiary/README.md) — Gerador procedural de monstros e definições do Dragão do Dia.
- [engine](./engine/README.md) — Regras de domínio: cálculos de combate, gerenciador de turnos, progressão e economia.
- [i18n](./i18n/README.md) — Dicionários em Português do Brasil (pt-BR), formatação de mensagens e diálogos.
- [storage](./storage/README.md) — Persistência SQLite CGO-free, migrações e repositório de jogadores.
- [tui](./tui/README.md) — Coordenador da interface de terminal em Bubble Tea e roteamento de telas.
- [ui](./ui/README.md) — Componentes visuais compartilhados do terminal e barra de status HUD.
