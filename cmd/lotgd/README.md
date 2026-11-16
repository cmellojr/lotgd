# cmd/lotgd

## Visão Geral
Ponto de entrada executável para o cliente de terminal local (single-player) de *The Legend of the Go Dragon*. Inicializa a persistência local em SQLite, constrói o modelo de interface em Bubble Tea e executa o loop interativo do terminal.

## Interface
- `main()`: Ponto de entrada em `main.go`. Processa flags da linha de comando, abre a conexão com o banco SQLite local (`lotgd.db`), inicializa repositórios e telas, e executa a aplicação Bubble Tea.

## Invariantes
- Executa em modo de terminal local utilizando a entrada e saída padrão do Bubble Tea (`tea.NewProgram`).
- Garante o fechamento limpo da conexão com o banco SQLite ao encerrar a aplicação.

## Efeitos Colaterais
- Cria ou abre o arquivo de banco de dados SQLite local `lotgd.db` no diretório de trabalho.
- Captura e restaura o buffer de tela alternativo do terminal e o modo raw.

## Verificação
- Executar: `go test ./cmd/lotgd/...`
