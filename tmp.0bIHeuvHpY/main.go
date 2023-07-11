package main

import (
	"fmt"
	"log"
	"net/http"
	"crypto/tls"
	"os"
	"context"
	"sync"
)

func startHealthCheckHttpServer(wg *sync.WaitGroup) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		if r == nil || len(r.RemoteAddr) == 0 {
			log.Panic("Unkown host request for healthcheck")
		} else {
			log.Printf("Received healthcheck request from: %s", r.RemoteAddr)
		}
		fmt.Fprint(w, "OK")
	})

	srv := &http.Server{
		Addr:         ":8888",
		Handler:      mux,
		TLSConfig:    getCfg(),
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}

	go func() {
		defer wg.Done() // let main know we are done cleaning up

		// always returns error. ErrServerClosed on graceful close
		if err := srv.ListenAndServeTLS("", ""); err == http.ErrServerClosed {
			log.Printf("Health check server gracefully shutdown: %v", err)
		} else {
			// unexpected error. port in use?
			log.Fatalf("Health check server got an unexpected error: %v", err)
		}
	}()

	// returning reference so caller can call Shutdown()
	return srv
}

//REF:
// https://github.com/denji/golang-tls
// https://stackoverflow.com/questions/47857573/passing-certificate-and-key-as-string-to-listenandservetls
func hi(healthCheckHttpServerSrv *http.Server) {
	mux := http.NewServeMux()
	mux.HandleFunc("/hi", func(w http.ResponseWriter, req *http.Request) {
		if req == nil || len(req.RemoteAddr) == 0 {
			log.Panic("Unkown host request")
		} else {
			log.Printf("Received request from: %s", req.RemoteAddr)
		}
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		w.Write([]byte("This is an example server.\n"))
	})
	srv := &http.Server{
		Addr:         ":4443",
		Handler:      mux,
		TLSConfig:    getCfg(),
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}
	mux.HandleFunc("/poisonous", func(w http.ResponseWriter, req *http.Request) {
		if shutdownErr := srv.Shutdown(context.TODO()); shutdownErr == nil {
			log.Printf("TLS server gracefully shutdown")
		} else {
			panic(shutdownErr) // failure/timeout shutting down the server gracefully
		}
	})
	// Health check should be running/listening by now
	log.Fatal(srv.ListenAndServeTLS("", ""))
	// Shutdown healthcheck
	if shutdownErr := healthCheckHttpServerSrv.Shutdown(context.TODO()); shutdownErr == nil {
		log.Printf("TLS server gracefully shutdown")
	} else {
		panic(shutdownErr) // failure/timeout shutting down the server gracefully
	}
}

func getCfg() *tls.Config {
	chain, err_chain := os.ReadFile("certs/tls-chain.cert.pem")
	if nil == err_chain {
		key, err_key := os.ReadFile("certs/tls-key.pem")
		if nil == err_key {
			cert, err := tls.X509KeyPair(chain, key)
			if err == nil {
				cfg := &tls.Config{
					Certificates:             []tls.Certificate{cert},
					MinVersion:               tls.VersionTLS12,
					CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
					PreferServerCipherSuites: true,
					CipherSuites: []uint16{
						tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
						tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
						tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
						tls.TLS_RSA_WITH_AES_256_CBC_SHA,
					},
				}
				return cfg
			} 
			log.Fatal(err)
		} else {
			log.Fatal(err_key)
		}
	} else {
		log.Fatal(err_chain)
	}
	return nil
}

func main() {
	httpServerExitDone := &sync.WaitGroup{}
	httpServerExitDone.Add(1)
	srv := startHealthCheckHttpServer(httpServerExitDone)
	hi(srv)
	httpServerExitDone.Wait()
	log.Printf("main: done. exiting")
}
