package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/go-chi/chi/v5"
)

func main() {
	log.Printf("starting example")

	go func() {
		err := http.ListenAndServe(":8080", router())
		if err != nil {
			panic(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	signal.Notify(stop, syscall.SIGTERM)
	<-stop
}

func router() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Get("/fib", func(w http.ResponseWriter, r *http.Request) {
		nQp := r.URL.Query().Get("n")
		n, err := strconv.Atoi(nQp)
		if err != nil {
			log.Printf("error parsing query parameter for request: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		res := fibMemoization(n, make(map[int]int))
		bArr, err := json.Marshal(res)
		if err != nil {
			log.Printf("error marshalling request: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = w.Write(bArr)
		if err != nil {
			return
		}
		if err != nil {
			log.Printf("error writing response: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
	return r
}

func fibMemoization(n int, m map[int]int) int {
	if n < 0 {
		return 0
	} else if n <= 1 {
		return n
	} else if val, ok := m[n]; ok {
		return val
	}

	val := fibMemoization(n-2, m) + fibMemoization(n-1, m)
	m[n] = val
	return val
}
