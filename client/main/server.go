package main

import (
	"log"
	"github.com/gorilla/mux"
	"blockchain_at_insurtech/client/handle"
	"net/http"
	"time"
	"fmt"
	"html"
	"flag"
)

const 	addr = ":8001"

func main() {
	log.Printf("starting REST server")
	router := mux.NewRouter().StrictSlash(true)
	// http://localhost:81/api/creditors
	// handle root
	router.HandleFunc("/", Index)
	// handle UI
	// This will serve files under http://localhost:8000/ui/<filename>
	var dir string
	flag.StringVar(&dir, "dir", "./ui/", "the directory to serve files from. Defaults to the current dir")
	flag.Parse()
	router.PathPrefix("/ui/").Handler(
		http.StripPrefix("/ui/", http.FileServer(http.Dir(dir))))

	//
	// handle API
	//
	// insuranceProduct API
	router.HandleFunc("/api/insuranceProduct", handle.CreateInsuranceProductHandler).Methods("POST")
	router.HandleFunc("/api/insuranceProduct/{id}", handle.GetInsuranceProductHandler).Methods("GET")
	router.HandleFunc("/api/insuranceProduct", handle.GetInsuranceProductsHandler).Methods("GET")
	// entry API
	router.HandleFunc("/api/insuranceEntry", handle.CreateInsuranceEntryHandler).Methods("POST")
	router.HandleFunc("/api/insuranceEntry/{id}", handle.GetInsuranceEntryHandler).Methods("GET")
	router.HandleFunc("/api/insuranceEntry", handle.GetInsuranceEntriesHandler).Methods("GET")

	//log.Fatal(http.ListenAndServe(":81", router))
	httpServer := &http.Server{
		Addr:           addr,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	err := httpServer.ListenAndServe()
	log.Fatal(err)
}

func Index(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	fmt.Fprintf(w, "Hello from %q", html.EscapeString(r.Host))
}
