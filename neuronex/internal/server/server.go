package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/emqx/neuronex-go/internal/auth"
	"github.com/emqx/neuronex-go/internal/config"
	"github.com/emqx/neuronex-go/internal/proxy"
)

type Server struct {
	cfg     *config.Config
	proxies []*proxy.ProxyHandler
	webDir  string
}

func NewServer(cfg *config.Config, webDir string) (*Server, error) {
	s := &Server{
		cfg:    cfg,
		webDir: webDir,
	}

	// Register Neuron proxies
	for _, rp := range cfg.Neuron.ReverseProxies {
		if rp.Location != "" && rp.ProxyPath != "" {
			h, err := proxy.NewProxyHandler(rp.Location, rp.ProxyPath)
			if err != nil {
				return nil, err
			}
			s.proxies = append(s.proxies, h)
			fmt.Printf("[NeuronEX] Registered proxy: %s -> %s\n", rp.Location, rp.ProxyPath)
		}
	}

	// Register Ekuiper proxies
	for _, rp := range cfg.Ekuiper.ReverseProxies {
		if rp.Location != "" && rp.ProxyPath != "" {
			h, err := proxy.NewProxyHandler(rp.Location, rp.ProxyPath)
			if err != nil {
				return nil, err
			}
			s.proxies = append(s.proxies, h)
			fmt.Printf("[NeuronEX] Registered proxy: %s -> %s\n", rp.Location, rp.ProxyPath)
		}
	}

	return s, nil
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// CORS Headers
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
	rw.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

	if req.Method == http.MethodOptions {
		rw.WriteHeader(http.StatusNoContent)
		return
	}

	p := req.URL.Path

	// Root redirect
	if p == "/" || p == "/web" || p == "/web/" {
		http.Redirect(rw, req, "/web/common/", http.StatusSeeOther)
		return
	}

	// Internal APIs
	switch {
	case p == "/api/login":
		s.handleLogin(rw, req)
		return
	case p == "/api/password":
		s.handlePassword(rw, req)
		return
	case p == "/api/license":
		s.handleLicense(rw, req)
		return
	case p == "/api/ping":
		rw.Header().Set("Content-Type", "application/json")
		rw.Write([]byte(`{"status":"ok"}`))
		return
	case strings.HasPrefix(p, "/api/users/role/"):
		rw.Header().Set("Content-Type", "application/json")
		rw.Write([]byte(`{"role":"Administrator"}`))
		return
	}

	// Check Auth for other /api/ routes
	if strings.HasPrefix(p, "/api/") && !s.cfg.Server.DisableAuth {
		authHeader := req.Header.Get("Authorization")
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" || authHeader == token {
			// check cookie
			if c, err := req.Cookie("token"); err == nil {
				token = c.Value
			}
		}

		if token == "" {
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte(`{"error":20203,"message":"missing token"}`))
			return
		}

		if _, err := auth.VerifyToken(token); err != nil {
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(rw, `{"error":20203,"message":"invalid token: %v"}`, err)
			return
		}
	}

	// Proxy matching
	for _, prx := range s.proxies {
		if strings.HasPrefix(p, prx.Location) {
			prx.ServeHTTP(rw, req)
			return
		}
	}

	// Serve Static Files (/web/common/...)
	if strings.HasPrefix(p, "/web/") {
		s.handleStatic(rw, req)
		return
	}

	http.NotFound(rw, req)
}

func (s *Server) handleLogin(rw http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(rw, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(rw, `{"error": 400, "message": "cannot read body"}`, http.StatusBadRequest)
		return
	}

	var creds struct {
		Name     string `json:"name"`
		Pass     string `json:"pass"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.Unmarshal(body, &creds); err != nil {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(`{"error":20100,"message":"Request body invalid"}`))
		return
	}

	u := creds.Name
	if u == "" {
		u = creds.Username
	}
	pw := creds.Pass
	if pw == "" {
		pw = creds.Password
	}

	if !auth.ValidateCredentials(u, pw) {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte(`{"error":20200,"message":"username or password error"}`))
		return
	}

	token, _ := auth.GenerateToken(u, "Administrator")
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(map[string]interface{}{
		"token": token,
		"role":  "Administrator",
	})
}

func (s *Server) handlePassword(rw http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(rw, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var pReq struct {
		OldPass string `json:"oldPass"`
		NewPass string `json:"newPass"`
	}
	if err := json.NewDecoder(req.Body).Decode(&pReq); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(`{"error": 400, "message": "invalid request"}`))
		return
	}

	if !auth.ValidateCredentials("admin", pReq.OldPass) {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(`{"error": 400, "message": "old password incorrect"}`))
		return
	}

	_ = auth.SetPassword("admin", pReq.NewPass)
	rw.Header().Set("Content-Type", "application/json")
	rw.Write([]byte(`{"status":"success"}`))
}

func (s *Server) handleLicense(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(map[string]interface{}{
		"valid_until":   "2999-01-01 00:00:00",
		"valid_since":   "2024-01-01 00:00:00",
		"license_type":  "official",
		"max_node_tags": 1000000,
		"used_tags":     0,
		"error":         0,
	})
}

func (s *Server) handleStatic(rw http.ResponseWriter, req *http.Request) {
	relPath := strings.TrimPrefix(req.URL.Path, "/web/common")
	relPath = strings.TrimPrefix(relPath, "/")
	if relPath == "" {
		relPath = "index.html"
	}

	filePath := filepath.Join(s.webDir, relPath)
	info, err := os.Stat(filePath)
	if err != nil || info.IsDir() {
		// SPA routing fallback: serve index.html if no file extension
		if !strings.Contains(filepath.Base(relPath), ".") {
			http.ServeFile(rw, req, filepath.Join(s.webDir, "index.html"))
			return
		}
		http.NotFound(rw, req)
		return
	}

	http.ServeFile(rw, req, filePath)
}
