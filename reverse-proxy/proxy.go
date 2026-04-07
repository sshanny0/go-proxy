package reverse

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
	"time"
)

type ReverseProxy struct {
	backends []*url.URL
	current  uint64
}

func NewReverseProxy(backendsURLs []string) (*ReverseProxy, error) {
	urls := make([]*url.URL, len(backendsURLs))
	for i, u := range backendsURLs {
		parsedURL, err := url.Parse(u)
		if err != nil {
			return nil, err
		}
		urls[i] = parsedURL
	}
	return &ReverseProxy{
		backends: urls,
	}, nil
}

func (rp *ReverseProxy) serveHTTP(w http.ResponseWriter, r *http.Request) {
	index := atomic.AddUint64(&rp.current, 1) % uint64(len(rp.backends))
	proxy := httputil.NewSingleHostReverseProxy(rp.backends[index])
	proxy.ServeHTTP(w, r)
}

func (p *ReverseProxy) healthCheck() {
	for {
		for _, backend := range p.backends {
			resp, err := http.Get(backend.String())
			if err != nil {
				log.Printf("Error checking backend %s: %v", backend.String(), err)
			} else {
				resp.Body.Close()
			}
		}
		time.Sleep(10 * time.Second)
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
