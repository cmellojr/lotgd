package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"lotgd/internal/storage"
	"lotgd/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	wishbubbletea "github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
)

func main() {
	// Flags de configuração do servidor SSH BBS.
	// Permite customizar host, porta de escuta, caminho da chave de host SSH e arquivo SQLite.
	host := flag.String("host", "0.0.0.0", "Endereço de escuta do servidor SSH")
	port := flag.Int("port", 2222, "Porta de escuta do servidor SSH")
	hostKeyPath := flag.String("host-key", ".wish_host_key", "Caminho para o arquivo de chave privada do host SSH")
	dbPath := flag.String("db", "lotgd.db", "Caminho para o arquivo de banco de dados SQLite")
	flag.Parse()

	// Inicializamos o logger estruturado com informações contextuais.
	logger := log.NewWithOptions(os.Stderr, log.Options{
		ReportTimestamp: true,
		Prefix:          "LOTGD-BBS",
	})

	// Inicializamos o banco de dados compartilhado em modo WAL.
	// O SQLite com modo WAL e busy_timeout suporta múltiplas leituras e escritas concorrentes
	// entre todas as sessões SSH conectadas simultaneamente sem corromper os dados.
	db, err := storage.OpenDB(*dbPath)
	if err != nil {
		logger.Fatal("Falha ao abrir banco de dados SQLite", "db", *dbPath, "err", err)
	}

	// Garantimos o encerramento gracioso do pool de conexões SQLite ao finalizar o servidor.
	defer func() {
		if err := db.Close(); err != nil {
			logger.Warn("Aviso ao fechar banco de dados", "err", err)
		}
	}()

	// Handler do Bubble Tea para o middleware do Wish:
	// Para cada conexão SSH estabelecida por um cliente, essa função é executada para
	// instanciar uma máquina de estados TUI isolada (MainModel) vinculada à sessão SSH.
	teaHandler := func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
		// PTY (Pseudo-Terminal): O Wish aloca automaticamente um terminal virtual para a sessão.
		_, _, active := s.Pty()
		if !active {
			wish.Fatalln(s, "Terminal interativo (PTY) é obrigatório para jogar LOTGD.")
			return nil, nil
		}

		// Instanciação de um novo modelo para a sessão do jogador conectado.
		model := tui.NewMainModel(db)

		// Opções específicas da sessão Bubble Tea sobre o túnel SSH:
		// - WithAltScreen(): Usa a tela secundária para limpar a interface ao desconectar.
		// O middleware do Wish (wishbubbletea.Middleware) já conecta automaticamente os streams
		// de entrada e saída (I/O) da sessão SSH ao programa Bubble Tea.
		return model, []tea.ProgramOption{
			tea.WithAltScreen(),
		}
	}

	// Construímos o servidor SSH com o framework Wish e middlewares essenciais:
	// 1. wishbubbletea.Middleware: Transforma conexões SSH em sessões do Bubble Tea.
	// 2. activeterm.Middleware: Garante que o cliente SSH forneça um terminal ativo (PTY).
	// 3. logging.Middleware: Registra conexões, desconexões e tempo de sessão.
	serverAddr := net.JoinHostPort(*host, fmt.Sprintf("%d", *port))
	server, err := wish.NewServer(
		wish.WithAddress(serverAddr),
		wish.WithHostKeyPath(*hostKeyPath),
		wish.WithMiddleware(
			wishbubbletea.Middleware(teaHandler),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)
	if err != nil {
		logger.Fatal("Falha ao configurar servidor Wish SSH", "err", err)
	}

	// Gerenciamento de sinais do Sistema Operacional (SIGINT, SIGTERM) para encerramento gracioso (Graceful Shutdown).
	// Usamos context.WithCancel e uma goroutine dedicada para escutar interrupções sem interromper transações abruptamente.
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info("Iniciando servidor SSH The Legend of the Go Dragon (BBS)...", "addr", serverAddr)
		logger.Info("Para conectar, use o comando:", "cmd", fmt.Sprintf("ssh localhost -p %d", *port))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			logger.Error("Erro na execução do servidor SSH", "err", err)
			done <- nil
		}
	}()

	// Bloqueia a execução da main goroutine até que um sinal de parada seja recebido no canal 'done'.
	<-done
	logger.Info("Sinal de encerramento recebido. Desligando servidor SSH graciosamente...")

	// Definimos um timeout de 10 segundos para desconectar as sessões ativas e liberar recursos.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		logger.Error("Erro ao desligar o servidor SSH", "err", err)
	}
	logger.Info("Servidor LOTGD BBS finalizado com sucesso.")
}
