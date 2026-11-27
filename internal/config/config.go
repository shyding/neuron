package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ReverseProxy defines a location-to-target proxy mapping
type ReverseProxy struct {
	Location  string `json:"location"`
	ProxyPath string `json:"proxyPath"`
}

// ServerConfig holds the main server settings
type ServerConfig struct {
	Port        int    `json:"port"`
	Host        string `json:"host"`
	DisableAuth bool   `json:"disableAuth"`
}

// ServiceProxyConfig groups reverse proxy definitions for a service
type ServiceProxyConfig struct {
	ReverseProxies []ReverseProxy `json:"reverseProxies"`
}

// Config is the top-level configuration structure
type Config struct {
	Server  ServerConfig       `json:"server"`
	Neuron  ServiceProxyConfig `json:"neuron"`
	Ekuiper ServiceProxyConfig `json:"ekuiper"`
	Monitor ServiceProxyConfig `json:"monitor"`
}

// DefaultConfig returns a sensible default configuration
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:        9085,
			Host:        "0.0.0.0",
			DisableAuth: false,
		},
	}
}

// LoadConfig parses a simple YAML-like config file
func LoadConfig(path string) (*Config, error) {
	cfg := DefaultConfig()
	f, err := os.Open(path)
	if err != nil {
		return cfg, nil // use defaults if file missing
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var currentSection string
	var inReverseProxies bool
	var currentProxy *ReverseProxy

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		rawLine := scanner.Text()
		indent := len(rawLine) - len(strings.TrimLeft(rawLine, " "))

		if indent == 0 && strings.HasSuffix(line, ":") {
			currentSection = strings.TrimSuffix(line, ":")
			inReverseProxies = false
			currentProxy = nil
			continue
		}

		if currentSection == "server" {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				k := strings.TrimSpace(parts[0])
				v := strings.TrimSpace(parts[1])
				switch k {
				case "port":
					if p, err := strconv.Atoi(v); err == nil {
						cfg.Server.Port = p
					}
				case "host":
					cfg.Server.Host = v
				case "disableAuth":
					cfg.Server.DisableAuth = (v == "true" || v == "1")
				}
			}
		} else if currentSection == "neuron" || currentSection == "ekuiper" || currentSection == "monitor" {
			if strings.HasPrefix(line, "reverseProxies:") {
				inReverseProxies = true
				continue
			}
			if inReverseProxies {
				if strings.HasPrefix(line, "- ") {
					item := strings.TrimPrefix(line, "- ")
					parts := strings.SplitN(item, ":", 2)
					if len(parts) == 2 {
						k := strings.TrimSpace(parts[0])
						v := strings.TrimSpace(parts[1])
						p := ReverseProxy{}
						if k == "location" {
							p.Location = v
						} else if k == "proxyPath" {
							p.ProxyPath = v
						}
						switch currentSection {
						case "neuron":
							cfg.Neuron.ReverseProxies = append(cfg.Neuron.ReverseProxies, p)
							currentProxy = &cfg.Neuron.ReverseProxies[len(cfg.Neuron.ReverseProxies)-1]
						case "ekuiper":
							cfg.Ekuiper.ReverseProxies = append(cfg.Ekuiper.ReverseProxies, p)
							currentProxy = &cfg.Ekuiper.ReverseProxies[len(cfg.Ekuiper.ReverseProxies)-1]
						case "monitor":
							cfg.Monitor.ReverseProxies = append(cfg.Monitor.ReverseProxies, p)
							currentProxy = &cfg.Monitor.ReverseProxies[len(cfg.Monitor.ReverseProxies)-1]
						}
					}
				} else if currentProxy != nil {
					parts := strings.SplitN(line, ":", 2)
					if len(parts) == 2 {
						k := strings.TrimSpace(parts[0])
						v := strings.TrimSpace(parts[1])
						if k == "location" {
							currentProxy.Location = v
						} else if k == "proxyPath" {
							currentProxy.ProxyPath = v
						}
					}
				}
			}
		}
	}

	fmt.Printf("[NeuronEX] Config loaded: Port=%d, Host=%s, DisableAuth=%v\n",
		cfg.Server.Port, cfg.Server.Host, cfg.Server.DisableAuth)
	return cfg, nil
}
