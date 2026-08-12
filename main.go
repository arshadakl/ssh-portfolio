package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
)

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, _ := s.Pty()
	renderer := bubbletea.MakeRenderer(s)

	clientIP := ""
	if addr := s.RemoteAddr(); addr != nil {
		if host, _, err := net.SplitHostPort(addr.String()); err == nil {
			clientIP = host
		} else {
			clientIP = addr.String()
		}
	}

	m := NewModel(renderer, clientIP)
	m.width = pty.Window.Width
	m.height = pty.Window.Height
	return m, []tea.ProgramOption{
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "22"
	}
	addr := "0.0.0.0:" + port

	keyPath := os.Getenv("HOST_KEY_PATH")
	if keyPath == "" {
		keyPath = "/app/keys/id_ed25519"
	}

	s, err := wish.NewServer(
		wish.WithAddress(addr),
		wish.WithHostKeyPath(keyPath),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			logging.Middleware(),
		),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not start server: %v\n", err)
		os.Exit(1)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	fmt.Printf("SSH portfolio v2 running on %s\n", addr)

	go func() {
		if err := s.ListenAndServe(); err != nil {
			fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		}
	}()

	<-done

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	s.Shutdown(ctx)
	fmt.Println("Server stopped")
}
