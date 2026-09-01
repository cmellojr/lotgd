package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"lotgd/internal/bestiary"
	"lotgd/internal/storage"
	"lotgd/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Definição de flags de linha de comando para configuração flexível do banco de dados.
	// Em Go, o pacote 'flag' permite analisar argumentos passados pelo usuário no terminal (ex: -db=meujogo.db).
	dbPath := flag.String("db", "lotgd.db", "Caminho para o arquivo de banco de dados SQLite")
	flag.Parse()

	// Inicializamos a conexão com o banco de dados SQLite.
	// O módulo 'storage' configura pragmas de alta performance (WAL mode, busy_timeout)
	// e executa automaticamente as migrações de schema necessárias (DDL).
	db, err := storage.OpenDB(*dbPath)
	if err != nil {
		log.Fatalf("Erro ao inicializar o banco de dados SQLite em %q: %v", *dbPath, err)
	}

	// O uso da declaração 'defer' garante que o fechamento do banco de dados (db.Close)
	// seja executado de maneira confiável no momento em que a função main() retornar,
	// liberando descritores de arquivos e salvando estados pendentes mesmo se ocorrerem erros.
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Aviso: erro ao encerrar conexões com o banco de dados: %v", err)
		}
	}()

	// Instanciamos o modelo raiz da arquitetura TEA (The Elm Architecture) do Bubble Tea.
	// Esse modelo gerencia a máquina de estados, roteando entre telas (Login, Cidade, Floresta, etc.).
	// O DragonGenerator é injetado via closure para que o storage não precise importar bestiary
	// (evitando ciclo de dependência: storage → bestiary → engine → storage).
	dragonGen := func(dayDate string) (int, int, int, int) {
		d := bestiary.GenerateDragonOfDay(dayDate)
		return d.MaxHealth, d.Attack, d.Defense, d.GoldReward
	}
	mainModel := tui.NewMainModel(db, dragonGen)

	// Criamos o programa Bubble Tea com opções de execução:
	// - tea.WithAltScreen(): Ativa o buffer de tela alternativo do terminal, preservando
	//   o histórico anterior do terminal do usuário ao fechar o jogo.
	// - tea.WithMouseCellMotion(): Habilita eventos de mouse se necessário no terminal.
	program := tea.NewProgram(
		mainModel,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	// Executamos o loop de eventos principal do Bubble Tea.
	// A função program.Run() bloqueia até que o usuário saia do jogo (ex: pressionando Ctrl+C ou optando por sair).
	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Erro durante a execução do jogo: %v\n", err)
		os.Exit(1)
	}
}
