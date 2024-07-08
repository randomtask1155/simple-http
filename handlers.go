package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var sleepTime int64

func init() {
	sleepTime = 0
}

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func DieHorriblyHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ouch")
	os.Exit(123)
}

/*
https://domain/health?sleep=2
adding sleep param will set sleep for all future requests
*/
func healthHandler(w http.ResponseWriter, r *http.Request) {
	sleep := r.FormValue("sleep")
	if sleep != "" {
		s, err := strconv.Atoi(sleep)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		sleepTime = int64(s)
	}
	time.Sleep(time.Duration(sleepTime) * time.Second)
}

// this is not synchronized so don't send too many requests or it will get weird
func shutdownHTTPServer(w http.ResponseWriter, r *http.Request) {
	var s int
	var err error
	sleep := r.FormValue("sleep")
	if sleep != "" {
		s, err = strconv.Atoi(sleep)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		s = 1 // default sleep for 1 second
	}
	fmt.Println("shutdown the http server")
	serverChan <- "stop"
	fmt.Printf("sleeping for %d seconds\n", s)
	time.Sleep(time.Duration(int64(s)) * time.Second)
	fmt.Println("starting the http server")
	serverChan <- "start"
}

func csvHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("\"columna\",\"columnb\"\n\"yup\",\"ok\"\n")))
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {

	type JsonTable struct {
		Ticket      string `json:"ticket"`
		Description string `json:"description"`
	}
	jt := JsonTable{"123456789", "nothing to see here"}
	b, err := json.Marshal(jt)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func dataInResponseHandler(w http.ResponseWriter, r *http.Request) {
	length := r.FormValue("length")
	if length == "" {
		length = "100"
	}
	sleep := r.FormValue("sleep")
	if sleep == "" {
		sleep = "0"
	}

	l, err := strconv.Atoi(length)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	s, err := strconv.Atoi(sleep)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write([]byte(fmt.Sprintf("{\"data\": \"")))
	for i := 1; i < l; i++ {
		w.Write([]byte(fmt.Sprintf("%s", RandStringRunes(1))))
		time.Sleep(time.Duration(int64(s)) * time.Second)
	}
	w.Write([]byte(fmt.Sprintf("\"}")))
}

func readBodyHandler(w http.ResponseWriter, r *http.Request) {
	_, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
}
