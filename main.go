package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var server *http.Server
var serverChan chan (string)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	log.Printf("%s\n", body)
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
	server = &http.Server{Addr: ":" + os.Getenv("PORT"), Handler: nil}
	serverChan = make(chan (string), 0)
	serverError := make(chan (error), 0)

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/get/data", dataInResponseHandler)
	http.HandleFunc("/post/data", readBodyHandler)
	http.HandleFunc("/die", DieHorriblyHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/shutdown", shutdownHTTPServer)
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
