package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type ReverseProxy struct {
	Location  string `json:"location"`
	ProxyPath string `json:"proxyPath"`
}

type ServerConfig struct {
	Port          int    `json:"port"`
	Host          string `json:"host"`
	DisableAuth   bool   `json:"disableAuth"`
	DisableKuiper bool   `json:"disableKuiper"`
}

type ServiceProxyConfig struct {
	ReverseProxies []ReverseProxy `json:"reverseProxies"`
}

type Config struct {
	Server  ServerConfig       `json:"server"`
	Neuron  ServiceProxyConfig `json:"neuron"`
	Ekuiper ServiceProxyConfig `json:"ekuiper"`
	Monitor ServiceProxyConfig `json:"monitor"`
}

func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:          8085,
			Host:          "0.0.0.0",
			DisableAuth:   false,
			DisableKuiper: false,
		},
		Neuron: ServiceProxyConfig{
			ReverseProxies: []ReverseProxy{
				{
					Location:  "/api/neuron",
					ProxyPath: "http://127.0.0.1:7000/api/v2",
				},
			},
		},
		Ekuiper: ServiceProxyConfig{
			ReverseProxies: []ReverseProxy{
				{
					Location:  "/api/ekuiper",
					ProxyPath: "http://127.0.0.1:9081",
				},
				{
					Location:  "/api/sse/ekuiper",
					ProxyPath: "http://127.0.0.1:10081",
				},
				{
					Location:  "/ws/ekuiper",
					ProxyPath: "ws://127.0.0.1:10081",
				},
			},
		},
	}
}

// LoadConfig parses a simple yaml config or uses defaults
func LoadConfig(path string) (*Config, error) {
	cfg := DefaultConfig()
	f, err := os.Open(path)
	if err != nil {
		return cfg, nil // use default if file missing
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
				case "disableAuth":
					cfg.Server.DisableAuth = (v == "true" || v == "1")
				case "disableKuiper":
					cfg.Server.DisableKuiper = (v == "true" || v == "1")
				}
			}
		} else if currentSection == "neuron" || currentSection == "ekuiper" {
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
						if currentSection == "neuron" {
							cfg.Neuron.ReverseProxies = append(cfg.Neuron.ReverseProxies, p)
							currentProxy = &cfg.Neuron.ReverseProxies[len(cfg.Neuron.ReverseProxies)-1]
						} else {
							cfg.Ekuiper.ReverseProxies = append(cfg.Ekuiper.ReverseProxies, p)
							currentProxy = &cfg.Ekuiper.ReverseProxies[len(cfg.Ekuiper.ReverseProxies)-1]
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

	fmt.Printf("[NeuronEX] Config loaded: Port=%d, DisableAuth=%v, DisableKuiper=%v\n",
		cfg.Server.Port, cfg.Server.DisableAuth, cfg.Server.DisableKuiper)
	return cfg, nil
}
