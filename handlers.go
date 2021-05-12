package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
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
