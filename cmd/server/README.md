# cmd/server

## Visão Geral
Ponto de entrada executável para o servidor SSH multi-usuário (BBS) de *The Legend of the Go Dragon*. Utiliza o Charm Wish (`charmbracelet/wish`) para servir sessões interativas de Bubble Tea através de conexões SSH remotas.

## Interface
- `main()`: Ponto de entrada em `main.go`. Configura o listener SSH do Wish, inicializa a chave de host do servidor, vincula a porta designada (padrão `:2222`) e gerencia as sessões SSH recebidas.

## Invariantes
- Cada conexão SSH recebida inicia uma sessão isolada com sua própria instância de `tea.Program` do Bubble Tea.
- O acesso ao banco de dados SQLite é seguro para concorrência entre múltiplas sessões ativas.

## Efeitos Colaterais
- Abre e escuta no socket de rede TCP na porta `:2222`.
- Lê ou cria o arquivo de chave privada do servidor SSH (`.ssh/id_ed25519` ou equivalente).
- Acessa o banco de dados compartilhado para autenticação de jogadores e persistência de estado.

## Verificação
- Executar: `go test ./cmd/server/...`
