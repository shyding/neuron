package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/emqx/neuronex-go/internal/auth"
	"github.com/emqx/neuronex-go/internal/config"
	"github.com/emqx/neuronex-go/internal/server"
	"github.com/emqx/neuronex-go/internal/supervisor"
)

var (
	configFile string
	aloneMode  bool
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
	fs.BoolVar(&aloneMode, "alone", false, "run alone mode (no child process supervisor)")

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
		fmt.Println("NeuronEX Enterprise Gateway v4.0.2-self-impl")
	case "help", "-h", "--help":
		printUsage()
	default:
		printUsage()
	}
}

func printUsage() {
	fmt.Println(`NeuronEX - IIoT Edge Connectivity & Stream Processing Platform

Usage:
  neuronex [command] [flags]

Available Commands:
  run             Run NeuronEX in foreground console (starts Neuron + eKuiper + Web Gateway)
  alone           Run NeuronEX Gateway alone (without managing child processes)
  start           Start NeuronEX
  reset-password  Reset admin password to default (0000)
  version         Show version information
  help            Show this help message

Flags:
  -c, --config    Path to neuronex.yaml configuration file (default "etc/neuronex.yaml")`)
}

func runServer(standalone bool) {
	fmt.Println("==================================================================")
	fmt.Println("   NeuronEX IIoT Edge Platform (Neuron + eKuiper + Web Dashboard)  ")
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
		// try fallback
		webDir = filepath.Join(baseDir, "neuronex", "web", "common")
	}

	// Start supervisor if not in alone mode
	sup := supervisor.NewSupervisor()
	if !standalone {
		// 1. Launch Neuron core
		neuronBin := filepath.Join(baseDir, "bin", "neuron")
		if _, err := os.Stat(neuronBin); err != nil {
			neuronBin = filepath.Join(baseDir, "neuron", "bin", "neuron")
		}

		runDir := filepath.Join(baseDir, "run", "neuron")
		_ = os.MkdirAll(runDir, 0755)

		_ = sup.StartProcess(supervisor.ProcessSpec{
			Name:       "neuron",
			Executable: neuronBin,
			Args:       []string{"--log", "--disable_auth"},
			Dir:        runDir,
			Env:        []string{"LD_LIBRARY_PATH=" + filepath.Join(baseDir, "lib", "neuron")},
		})

		// 2. Launch eKuiper if enabled
		if !cfg.Server.DisableKuiper {
			kuiperBin := filepath.Join(baseDir, "bin", "kuiperd")
			_ = sup.StartProcess(supervisor.ProcessSpec{
				Name:       "ekuiper",
				Executable: kuiperBin,
				Args: []string{
					"-loadFileType", "absolute",
					"-etc", filepath.Join(baseDir, "share", "ekuiper", "etc"),
					"-data", filepath.Join(baseDir, "data", "ekuiper"),
					"-log", filepath.Join(baseDir, "log", "ekuiper"),
				},
				Dir: baseDir,
				Env: []string{
					"KUIPER__BASIC__RESTIP=127.0.0.1",
					"KUIPER__BASIC__PROMETHEUSPORT=9081",
				},
			})
		}
	}

	// Start HTTP Server
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
		fmt.Printf("[NeuronEX] Web Dashboard available at: http://localhost:%d/web/common/\n", cfg.Server.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("[NeuronEX] Server listen error: %v\n", err)
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
