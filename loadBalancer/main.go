package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"sync/atomic"
)

type RoundRobin struct {
	addrs      []string
	reqCounter uint64
}

func (rr *RoundRobin) NextServer() string {
	counter := atomic.AddUint64(&rr.reqCounter, 1)
	return rr.addrs[counter%uint64(len(rr.addrs))]
}

type ProxyHandler struct {
	loadBalancer *RoundRobin
	client       *http.Client
}

func (ph *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	targetAddr := ph.loadBalancer.NextServer()
	targetURL := &url.URL{
		Scheme: "http",
		Host:   targetAddr,
		Path:   r.URL.Path,
	}
	r.URL = targetURL

	resp, err := ph.client.Do(r)
	if err != nil {
		log.Printf("Error proxying request: %v", err)
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for name, values := range resp.Header {
		w.Header()[name] = values
	}
	w.WriteHeader(resp.StatusCode)

	io.Copy(w, resp.Body)
}

func handleTCPConnection(listener net.Listener, loadBalancer *RoundRobin) {
	for {
		downstream, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		go func(downstream net.Conn) {
			defer downstream.Close()

			upstreamAddr := loadBalancer.NextServer()
			upstream, err := net.Dial("tcp", upstreamAddr)
			if err != nil {
				log.Println("Error connecting to upstream:", err)
				return
			}
			defer upstream.Close()

			go io.Copy(upstream, downstream)
			io.Copy(downstream, upstream)
		}(downstream)
	}
}

func main() {
	addrs := []string{"api01:3000", "api02:3000"}
	roundRobin := &RoundRobin{addrs: addrs}
	client := &http.Client{}

	http.Handle("/", &ProxyHandler{loadBalancer: roundRobin, client: client})

	port := ":9998"
	fmt.Println("Proxy server listening on port", port)

	tcpListener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("Error starting TCP proxy server:", err)
	}
	defer tcpListener.Close()

	go handleTCPConnection(tcpListener, roundRobin)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("Error starting proxy server:", err)
	}
}
