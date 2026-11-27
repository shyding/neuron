package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/emqx/neuronex-go/internal/auth"
	"github.com/emqx/neuronex-go/internal/cctl"
	"github.com/emqx/neuronex-go/internal/config"
	"github.com/emqx/neuronex-go/internal/license"
	"github.com/emqx/neuronex-go/internal/proxy"
)

type UserInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Role        string `json:"role"`
}

type Server struct {
	cfg        *config.Config
	proxies    []*proxy.ProxyHandler
	webDir     string
	startTime  time.Time
	licenseMgr *license.Manager
	cctlConn   *cctl.Connector

	mu    sync.RWMutex
	users []UserInfo
}

func NewServer(cfg *config.Config, webDir string) (*Server, error) {
	cctlHost := cfg.Cctl.Host
	if cctlHost == "" {
		cctlHost = "127.0.0.1"
	}
	cctlPort := cfg.Cctl.Port
	if cctlPort <= 0 {
		cctlPort = 18884
	}
	cctlDev := cfg.Cctl.DeviceID
	if cctlDev == "" {
		cctlDev = "pump-03"
	}
	cctlKa := cfg.Cctl.Keepalive
	if cctlKa <= 0 {
		cctlKa = 30
	}

	s := &Server{
		cfg:        cfg,
		webDir:     webDir,
		startTime:  time.Now(),
		licenseMgr: license.NewManager("etc/license.json"),
		cctlConn: cctl.NewConnector(cctl.Config{
			Host:      cctlHost,
			Port:      cctlPort,
			DeviceID:  cctlDev,
			Keepalive: cctlKa,
		}),
		users: []UserInfo{
			{
				Name:        "admin",
				Description: "",
				Role:        "Administrator",
			},
		},
	}
	if cfg.Cctl.Enabled {
		_ = s.cctlConn.Start()
		fmt.Printf("[NeuronEX] CCTL Connector started for device: %s on %s:%d\n", cctlDev, cctlHost, cctlPort)
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

	// Register Monitor proxies
	for _, rp := range cfg.Monitor.ReverseProxies {
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

	// Transparently rewrite /web/common/api/* -> /api/*
	if strings.HasPrefix(p, "/web/common/api/") {
		p = strings.TrimPrefix(p, "/web/common")
		req.URL.Path = p
	}

	// Root and /web redirect to /web/common/
	if p == "/" || p == "/web" || p == "/web/" || p == "/web/common" {
		http.Redirect(rw, req, "/web/common/", http.StatusMovedPermanently)
		return
	}

	// Public APIs (No auth required)
	switch {
	case p == "/api/auth":
		rw.Header().Set("Content-Type", "application/json")
		rw.Write([]byte(`{"disableAuth":false}`))
		return
	case p == "/api/login":
		s.handleLogin(rw, req)
		return
	case p == "/api/sso" || strings.HasPrefix(p, "/api/sso/"):
		rw.Header().Set("Content-Type", "application/json")
		rw.Write([]byte(`[]`))
		return
	case p == "/api/system" || p == "/api/system/":
		s.handleSystem(rw, req)
		return
	case p == "/api/version":
		rw.Header().Set("Content-Type", "application/json")
		rw.Write([]byte(`{"version":"4.0.2"}`))
		return
	case p == "/api/ping":
		rw.Header().Set("Content-Type", "application/json")
		rw.Write([]byte(`{"status":"ok"}`))
		return
	case p == "/api/process/status":
		rw.Header().Set("Content-Type", "application/json")
		rw.Write([]byte(`{"ekuiper":true,"neuron":true}`))
		return
	case p == "/api/cctl/status":
		rw.Header().Set("Content-Type", "application/json")
		connected := s.cctlConn.IsConnected()
		rw.Write([]byte(fmt.Sprintf(`{"connected":%t,"host":"127.0.0.1","port":18884,"deviceId":"pump-03"}`, connected)))
		return
	case p == "/api/cctl/send":
		s.handleCCTLSend(rw, req)
		return
	}

	// Check if this route is handled by upstream proxies (Neuron, Ekuiper, Monitor)
	isProxied := false
	for _, prx := range s.proxies {
		if strings.HasPrefix(p, prx.Location) {
			isProxied = true
			break
		}
	}

	if isProxied {
		for _, prx := range s.proxies {
			if strings.HasPrefix(p, prx.Location) {
				prx.ServeHTTP(rw, req)
				return
			}
		}
	}

	// Check Auth for other internal /api/ routes
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

	// Protected Internal APIs
	switch {
	case p == "/api/rbac":
		s.handleRBAC(rw, req)
		return
	case p == "/api/password":
		s.handlePassword(rw, req)
		return
	case p == "/api/license" || p == "/api/license/":
		s.handleLicense(rw, req)
		return
	case p == "/api/license/activation_code":
		s.handleActivationCode(rw, req)
		return
	case p == "/api/license/trial" || p == "/api/license/reset":
		s.handleResetLicense(rw, req)
		return
	case p == "/api/system/tracing":
		s.handleSystemTracing(rw, req)
		return
	case p == "/api/process/status":
		rw.Header().Set("Content-Type", "application/json")
		rw.Write([]byte(`{"ekuiper":true,"neuron":true}`))
		return
	case p == "/api/log" || p == "/api/log/":
		rw.Header().Set("Content-Type", "application/json")
		rw.Write([]byte(`{"level":"error"}`))
		return
	case p == "/api/users" || p == "/api/users/":
		s.handleUsers(rw, req)
		return
	case strings.HasPrefix(p, "/api/users/role/"):
		s.handleUserRole(rw, req)
		return
	case strings.HasPrefix(p, "/api/users/password/"):
		s.handleUserPassword(rw, req)
		return
	case strings.HasPrefix(p, "/api/users/"):
		s.handleUserItem(rw, req)
		return
	}

	// Proxy matching
	for _, prx := range s.proxies {
		if strings.HasPrefix(p, prx.Location) {
			prx.ServeHTTP(rw, req)
			return
		}
	}

	// Serve Static Files (/web/common/..., or directly /css/, /js/, /favicon.ico, etc.)
	if strings.HasPrefix(p, "/web/") || strings.HasPrefix(p, "/css/") || strings.HasPrefix(p, "/js/") ||
		strings.HasPrefix(p, "/img/") || strings.HasPrefix(p, "/fonts/") || strings.HasPrefix(p, "/template/") ||
		p == "/favicon.ico" || strings.HasSuffix(p, ".js") || strings.HasSuffix(p, ".css") || strings.HasSuffix(p, ".png") || strings.HasSuffix(p, ".svg") {
		s.handleStatic(rw, req)
		return
	}

	http.NotFound(rw, req)
}

func (s *Server) handleRBAC(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	rw.Write([]byte(`{
	"admin_configuration": {
		"backup": true,
		"bacnet_scan": true,
		"edit": true,
		"ekuiper": true,
		"modbus_simulator": true,
		"restore": true,
		"tracing": true,
		"view": true
	},
	"admin_license": {
		"activation": true,
		"reset": true,
		"upload": true,
		"view": true
	},
	"admin_logs": {
		"download": true,
		"level": true
	},
	"admin_sso": {
		"create": true,
		"delete": true,
		"edit": true,
		"view": true
	},
	"admin_users": {
		"create": true,
		"delete": true,
		"edit": true,
		"password": true,
		"view": true
	},
	"application": {
		"ctl": true
	},
	"configuration_connector": {
		"create": true,
		"delete": true,
		"edit": true,
		"view": true
	},
	"configuration_files": {
		"create": true,
		"delete": true,
		"view": true
	},
	"configuration_resources": {
		"create": true,
		"delete": true,
		"edit": true,
		"view": true
	},
	"configuration_schema": {
		"create": true,
		"delete": true,
		"edit": true,
		"view": true
	},
	"extension_functions": {
		"create": true,
		"delete": true,
		"edit": true,
		"view": true
	},
	"extension_portable": {
		"create": true,
		"delete": true,
		"edit": true,
		"view": true
	},
	"extension_service": {
		"create": true,
		"delete": true,
		"edit": true,
		"view": true
	},
	"monitoring": {
		"view": true,
		"write": true
	},
	"north_driver": {
		"config": true,
		"create": true,
		"ctl": true,
		"debug": true,
		"delete": true,
		"edit": true,
		"import": true,
		"log": true,
		"view": true
	},
	"north_driver_group": {
		"create": true,
		"delete": true,
		"edit": true,
		"view": true
	},
	"plugin": {
		"create": true,
		"delete": true,
		"edit": true,
		"view": true
	},
	"plugin_cnc": {
		"create": true,
		"view": true
	},
	"rules": {
		"copy": true,
		"create": true,
		"ctl": true,
		"delete": true,
		"edit": true,
		"export": true,
		"import": true,
		"tag": true,
		"tracing": true,
		"view": true
	},
	"source_stream": {
		"create": true,
		"delete": true,
		"edit": true,
		"view": true
	},
	"source_table": {
		"create": true,
		"delete": true,
		"edit": true,
		"view": true
	},
	"south_driver": {
		"config": true,
		"copy": true,
		"create": true,
		"ctl": true,
		"debug": true,
		"delete": true,
		"edit": true,
		"export": true,
		"import": true,
		"log": true,
		"tag": true,
		"view": true
	},
	"south_driver_group": {
		"create": true,
		"delete": true,
		"edit": true,
		"export": true,
		"import": true,
		"view": true
	},
	"south_driver_group_tags": {
		"create": true,
		"delete": true,
		"edit": true,
		"view": true
	}
}`))
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

func (s *Server) handleUsers(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	switch req.Method {
	case http.MethodGet:
		s.mu.RLock()
		defer s.mu.RUnlock()
		json.NewEncoder(rw).Encode(s.users)
	case http.MethodPost:
		var u struct {
			Name        string `json:"name"`
			Password    string `json:"password"`
			Description string `json:"description"`
			Role        string `json:"role"`
		}
		if err := json.NewDecoder(req.Body).Decode(&u); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(`{"error": 20100, "message": "Request body invalid"}`))
			return
		}
		s.mu.Lock()
		defer s.mu.Unlock()
		s.users = append(s.users, UserInfo{
			Name:        u.Name,
			Description: u.Description,
			Role:        u.Role,
		})
		_ = auth.SetPassword(u.Name, u.Password)
		rw.Write([]byte(`{"error":0}`))
	default:
		http.Error(rw, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleUserItem(rw http.ResponseWriter, req *http.Request) {
	userName := strings.TrimPrefix(req.URL.Path, "/api/users/")
	rw.Header().Set("Content-Type", "application/json")

	if req.Method == http.MethodDelete {
		s.mu.Lock()
		defer s.mu.Unlock()
		newUsers := make([]UserInfo, 0, len(s.users))
		for _, u := range s.users {
			if u.Name != userName {
				newUsers = append(newUsers, u)
			}
		}
		s.users = newUsers
		rw.Write([]byte(`{"error":0}`))
		return
	}

	// GET single user
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, u := range s.users {
		if u.Name == userName {
			json.NewEncoder(rw).Encode(u)
			return
		}
	}
	http.NotFound(rw, req)
}

func (s *Server) handleUserRole(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	if req.Method == http.MethodPut {
		var rReq struct {
			Name string `json:"name"`
			Role string `json:"role"`
		}
		_ = json.NewDecoder(req.Body).Decode(&rReq)
		userName := strings.TrimPrefix(req.URL.Path, "/api/users/role/")
		if userName == "" {
			userName = rReq.Name
		}
		s.mu.Lock()
		defer s.mu.Unlock()
		for i := range s.users {
			if s.users[i].Name == userName {
				s.users[i].Role = rReq.Role
				break
			}
		}
		rw.Write([]byte(`{"error":0}`))
		return
	}
	rw.Write([]byte(`{"role":"Administrator"}`))
}

func (s *Server) handleUserPassword(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	if req.Method == http.MethodPut {
		var pReq struct {
			Password string `json:"password"`
		}
		_ = json.NewDecoder(req.Body).Decode(&pReq)
		userName := strings.TrimPrefix(req.URL.Path, "/api/users/password/")
		_ = auth.SetPassword(userName, pReq.Password)
		rw.Write([]byte(`{"error":0}`))
		return
	}
	http.Error(rw, "Method Not Allowed", http.StatusMethodNotAllowed)
}

func (s *Server) handleSystem(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	runningSecs := fmt.Sprintf("%d", int(time.Since(s.startTime).Seconds()))
	json.NewEncoder(rw).Encode(map[string]interface{}{
		"version":           "4.0.2",
		"systemStatus":      "0",
		"neuronRunningTime": runningSecs,
		"neuronStatus":      "0",
		"neuronMemoryUsed":  "15794176",
		"neuronCpuUsed":     "0.0",
	})
}

func (s *Server) handleSystemTracing(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	if req.Method == http.MethodPost {
		rw.Write([]byte(`{"status":0}`))
		return
	}
	json.NewEncoder(rw).Encode(map[string]interface{}{
		"action":           "",
		"control":          false,
		"data":             false,
		"data_sample_rate": 1,
	})
}

func (s *Server) handleLicense(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	if req.Method == http.MethodPost {
		var body struct {
			License string `json:"license"`
		}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(`{"error": 400, "message": "invalid license payload"}`))
			return
		}
		_, err := s.licenseMgr.ActivateLicense(body.License)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(fmt.Sprintf(`{"error": 400, "message": "%v"}`, err)))
			return
		}
		rw.Write([]byte(`{"status": 0}`))
		return
	}

	// GET: Return current permanent license info
	lic := s.licenseMgr.GetLicense()
	json.NewEncoder(rw).Encode(lic)
}

func (s *Server) handleActivationCode(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	if req.Method != http.MethodPost {
		http.Error(rw, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var body struct {
		ActivationCode string `json:"activationCode"`
		Code           string `json:"code"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(`{"error": 400, "message": "invalid activation code payload"}`))
		return
	}
	code := body.ActivationCode
	if code == "" {
		code = body.Code
	}
	_, err := s.licenseMgr.ActivateCode(code)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(fmt.Sprintf(`{"error": 400, "message": "%v"}`, err)))
		return
	}
	rw.Write([]byte(`{"status": 0}`))
}

func (s *Server) handleResetLicense(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	_ = s.licenseMgr.Reset()
	rw.Write([]byte(`{"status": 0}`))
}

func (s *Server) handleStatic(rw http.ResponseWriter, req *http.Request) {
	relPath := strings.TrimPrefix(req.URL.Path, "/web/common")
	relPath = strings.TrimPrefix(relPath, "/web")
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

func (s *Server) handleCCTLSend(rw http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(rw, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var body struct {
		DeviceID   string                 `json:"deviceId"`
		Properties map[string]interface{} `json:"properties"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	if err := s.cctlConn.SendTelemetry(body.DeviceID, body.Properties); err != nil {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(fmt.Sprintf(`{"error":"%v"}`, err)))
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.Write([]byte(`{"status":"success"}`))
}
