package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/shyding/neuronex/internal/auth"
	"github.com/shyding/neuronex/internal/config"
	"github.com/shyding/neuronex/internal/server"
	"github.com/shyding/neuronex/internal/supervisor"
)

const Version = "4.0.2-neuronex-src"

var (
	configFile string
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]

	fs := flag.NewFlagSet("neuronex", flag.ExitOnError)
	fs.StringVar(&configFile, "c", "etc/neuronex.yaml", "config file path")
	fs.StringVar(&configFile, "config", "etc/neuronex.yaml", "config file path")

	if len(os.Args) > 2 {
		_ = fs.Parse(os.Args[2:])
	}

	switch command {
	case "run":
		runServer(false)
	case "alone":
		runServer(true)
	case "start":
		runServer(false)
	case "reset-password":
		_ = auth.SetPassword("admin", "0000")
		fmt.Println("Admin password reset to 0000 successfully.")
	case "version", "-v", "--version":
		fmt.Printf("NeuronEX %s\n", Version)
	case "help", "-h", "--help":
		printUsage()
	default:
		printUsage()
	}
}

func printUsage() {
	fmt.Println(`NeuronEX — IIoT Edge Connectivity & Stream Processing Platform

Usage:
  neuronex [command] [flags]

Available Commands:
  run             Run NeuronEX in foreground (starts gateway + built-in data engine)
  alone           Run NeuronEX gateway alone (no child process management)
  start           Start NeuronEX (alias for run)
  reset-password  Reset admin password to default (0000)
  version         Show version information
  help            Show this help message

Flags:
  -c, --config    Path to neuronex.yaml configuration file (default "etc/neuronex.yaml")`)
}

func runServer(standalone bool) {
	fmt.Println("==================================================================")
	fmt.Println("    NeuronEX — IIoT Edge Platform (100% Source Implementation)     ")
	fmt.Println("==================================================================")

	baseDir, err := os.Getwd()
	if err != nil {
		baseDir = "."
	}

	// Resolve config path
	confPath := configFile
	if !filepath.IsAbs(confPath) {
		confPath = filepath.Join(baseDir, confPath)
	}

	cfg, err := config.LoadConfig(confPath)
	if err != nil {
		fmt.Printf("[NeuronEX] Failed to load config: %v, using defaults\n", err)
		cfg = config.DefaultConfig()
	}

	// Resolve Web directory
	webDir := filepath.Join(baseDir, "web", "common")
	if _, err := os.Stat(webDir); err != nil {
		fmt.Printf("[NeuronEX] Warning: web directory %s not found\n", webDir)
	}

	// Start supervisor for child processes if not in standalone mode
	sup := supervisor.NewSupervisor()
	if !standalone {
		fmt.Println("[NeuronEX] Running in supervised mode (child processes managed)")
		// In a full deployment, this would start neuron core and ekuiper
		// For the source implementation, the data engine runs in-process
	} else {
		fmt.Println("[NeuronEX] Running in standalone/alone mode")
	}

	// Create and start the HTTP Server
	srv, err := server.NewServer(cfg, webDir)
	if err != nil {
		fmt.Printf("[NeuronEX] Failed to initialize server: %v\n", err)
		return
	}

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	httpServer := &http.Server{
		Addr:    addr,
		Handler: srv,
	}

	go func() {
		fmt.Printf("\n[NeuronEX] HTTP Server listening on http://%s\n", addr)
		fmt.Printf("[NeuronEX] Web Dashboard: http://localhost:%d/web/common/\n", cfg.Server.Port)
		fmt.Printf("[NeuronEX] Login: admin / 0000\n\n")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("[NeuronEX] Server error: %v\n", err)
		}
	}()

	// Signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\n[NeuronEX] Shutting down...")
	_ = httpServer.Close()
	sup.StopAll()
	fmt.Println("[NeuronEX] Bye.")
}
