package proxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"
)

type NeuronTokenManager struct {
	mu       sync.RWMutex
	token    string
	expireAt time.Time
	loginUrl string
}

func (m *NeuronTokenManager) GetToken() string {
	m.mu.RLock()
	if m.token != "" && time.Now().Before(m.expireAt) {
		tok := m.token
		m.mu.RUnlock()
		return tok
	}
	m.mu.RUnlock()

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.token != "" && time.Now().Before(m.expireAt) {
		return m.token
	}

	body := []byte(`{"name":"admin","pass":"0000"}`)
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Post(m.loginUrl, "application/json", bytes.NewReader(body))
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	var res struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil || res.Token == "" {
		return ""
	}

	m.token = res.Token
	m.expireAt = time.Now().Add(12 * time.Hour)
	return m.token
}

type ProxyHandler struct {
	Location     string
	Target       *url.URL
	Proxy        *httputil.ReverseProxy
	tokenManager *NeuronTokenManager
}

func NewProxyHandler(location, targetUrl string) (*ProxyHandler, error) {
	target, err := url.Parse(targetUrl)
	if err != nil {
		return nil, fmt.Errorf("invalid target url %s: %w", targetUrl, err)
	}

	var tokMgr *NeuronTokenManager
	if strings.HasPrefix(location, "/api/neuron") {
		tokMgr = &NeuronTokenManager{
			loginUrl: strings.TrimSuffix(targetUrl, "/") + "/login",
		}
	}

	p := httputil.NewSingleHostReverseProxy(target)

	p.Director = func(req *http.Request) {
		origPath := req.URL.Path
		origRawQuery := req.URL.RawQuery

		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.Host = target.Host

		// Strip location prefix from the original path
		relPath := strings.TrimPrefix(origPath, location)
		if !strings.HasPrefix(relPath, "/") {
			relPath = "/" + relPath
		}

		targetBasePath := strings.TrimSuffix(target.Path, "/")
		req.URL.Path = targetBasePath + relPath

		// Pass query parameters
		if target.RawQuery == "" {
			req.URL.RawQuery = origRawQuery
		} else if origRawQuery == "" {
			req.URL.RawQuery = target.RawQuery
		} else {
			req.URL.RawQuery = target.RawQuery + "&" + origRawQuery
		}

		req.Header.Set("X-Real-IP", "127.0.0.1")
		req.Header.Set("X-Forwarded-For", "127.0.0.1")

		// If proxying to Neuron Core, inject valid Neuron Core token
		if tokMgr != nil {
			if upstreamTok := tokMgr.GetToken(); upstreamTok != "" {
				req.Header.Set("Authorization", "Bearer "+upstreamTok)
			}
		}
	}

	p.ModifyResponse = func(resp *http.Response) error {
		// Ensure CORS headers on all proxied responses
		resp.Header.Set("Access-Control-Allow-Origin", "*")
		resp.Header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		resp.Header.Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		return nil
	}

	p.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(rw, `{"error": 502, "message": "NeuronEX upstream connection failed: %v"}`, err)
	}

	return &ProxyHandler{
		Location:     location,
		Target:       target,
		Proxy:        p,
		tokenManager: tokMgr,
	}, nil
}

func (ph *ProxyHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodOptions {
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		rw.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		rw.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		rw.WriteHeader(http.StatusNoContent)
		return
	}
	ph.Proxy.ServeHTTP(rw, req)
}
