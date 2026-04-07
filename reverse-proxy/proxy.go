package reverse

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
	"time"
)

type Backend struct {
	URL   *url.URL
	Alive bool
}

type ReverseProxy struct {
	backends []*Backend
	current  uint64
}

func NewReverseProxy(backendsURLs []string) (*ReverseProxy, error) {
	backends := make([]*Backend, len(backendsURLs))

	for i, u := range backendsURLs {
		parsedURL, err := url.Parse(u)
		if err != nil {
			return nil, err
		}

		backends[i] = &Backend{
			URL:   parsedURL,
			Alive: true,
		}
	}

	return &ReverseProxy{
		backends: backends,
	}, nil
}

// BACKEND SLEECTOR TO SKIP UNHEALTHY BACKENDS
func (p *ReverseProxy) getNextBackend() *Backend {
	n := len(p.backends)

	for i := 0; i < n; i++ {
		index := atomic.AddUint64(&p.current, 1) % uint64(n)
		backend := p.backends[index]

		if backend.Alive {
			return backend
		}
	}

	return nil
}

func (rp *ReverseProxy) serveHTTP(w http.ResponseWriter, r *http.Request) {
	backend := rp.getNextBackend()
	if backend == nil {
		http.Error(w, "No available backend", http.StatusServiceUnavailable)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(backend.URL)
	proxy.ServeHTTP(w, r)
}

func (p *ReverseProxy) healthCheck() {
	for {
		for _, backend := range p.backends {
			resp, err := http.Get(backend.URL.String())
			if err != nil {
				log.Printf("Backend DOWN: %s", backend.URL)
				backend.Alive = false
				continue
			}

			resp.Body.Close()
			backend.Alive = true
			log.Printf("Backend UP: %s", backend.URL)
		}

		time.Sleep(5 * time.Second)
	}
}

func Start() {
	backends := []string{
		"http://localhost:8081",
		"http://localhost:8082",
	}
	proxy, err := NewReverseProxy(backends)
	if err != nil {
		panic(err)
	}
	go proxy.healthCheck()

	http.HandleFunc("/", proxy.serveHTTP)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
