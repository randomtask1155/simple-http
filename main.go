package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

var server *http.Server
var serverChan chan (string)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	/*body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	log.Printf("%s\n", body)*/
	w.Write([]byte("<html>Hello!</html>"))
}

func monitorServerError(serverError chan (error)) {
	for {
		select {
		case err := <-serverError:
			fmt.Println(err)
			fmt.Println("server start failed.. attempting to restart in 30 seocnds")
			time.Sleep(30 * time.Second)
			go listenAndServeHTTP(serverError)
		}
	}
}

func listenAndServeHTTP(serverError chan (error)) {
	err := server.ListenAndServe()
	if err != nil {
		serverError <- err
	}
}

func main() {
	h2s := &http2.Server{
		// ...
	}
	//handler := http.HandlerFunc(rootHandler)
	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/get/data", dataInResponseHandler)
	mux.HandleFunc("/post/data", readBodyHandler)
	mux.HandleFunc("/die", DieHorriblyHandler)
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/shutdown", shutdownHTTPServer)
	mux.HandleFunc("/csv", csvHandler)
	mux.HandleFunc("/json", jsonHandler)
	mux.HandleFunc("/502", return502Handler)

	server = &http.Server{Addr: ":" + os.Getenv("PORT"), Handler: h2c.NewHandler(mux, h2s)}
	serverChan = make(chan (string), 0)
	serverError := make(chan (error), 0)

	go monitorServerError(serverError)

	//http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	go listenAndServeHTTP(serverError)
	for {
		time.Sleep(30 * time.Second)
		select {
		case msg := <-serverChan:
			if msg == "stop" {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				err := server.Shutdown(ctx)
				if err != nil {
					fmt.Println(err)
				}
				cancel()
			}
			if msg == "start" {
				server = &http.Server{Addr: ":" + os.Getenv("PORT"), Handler: nil}
				go listenAndServeHTTP(serverError)
			}
		}
	}
}
