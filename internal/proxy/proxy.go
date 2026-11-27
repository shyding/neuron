package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// ProxyHandler handles reverse proxy for a specific location prefix
type ProxyHandler struct {
	Location string
	Target   *url.URL
	Proxy    *httputil.ReverseProxy
}

// NewProxyHandler creates a reverse proxy for location -> targetUrl
func NewProxyHandler(location, targetUrl string) (*ProxyHandler, error) {
	target, err := url.Parse(targetUrl)
	if err != nil {
		return nil, fmt.Errorf("invalid target url %s: %w", targetUrl, err)
	}

	p := httputil.NewSingleHostReverseProxy(target)
	originalDirector := p.Director

	p.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = target.Host

		// Strip location prefix and append target path
		relPath := strings.TrimPrefix(req.URL.Path, location)
		if !strings.HasPrefix(relPath, "/") {
			relPath = "/" + relPath
		}

		targetBasePath := strings.TrimSuffix(target.Path, "/")
		req.URL.Path = targetBasePath + relPath

		if target.RawQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = target.RawQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = target.RawQuery + "&" + req.URL.RawQuery
		}

		req.Header.Set("X-Real-IP", "127.0.0.1")
		req.Header.Set("X-Forwarded-For", "127.0.0.1")
	}

	p.ModifyResponse = func(resp *http.Response) error {
		resp.Header.Set("Access-Control-Allow-Origin", "*")
		resp.Header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		resp.Header.Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		return nil
	}

	p.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(rw, `{"error": 502, "message": "upstream connection failed: %v"}`, err)
	}

	return &ProxyHandler{
		Location: location,
		Target:   target,
		Proxy:    p,
	}, nil
}

// ServeHTTP handles incoming requests for this proxy
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
