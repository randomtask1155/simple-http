package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

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

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/get/data", dataInResponseHandler)
	http.HandleFunc("/post/data", readBodyHandler)
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
