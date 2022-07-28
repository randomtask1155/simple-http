package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
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
		if i >= (l / 2) {
			panic("die half way through")
		}
	}
	w.Write([]byte(fmt.Sprintf("\"}")))
}

func readBodyHandler(w http.ResponseWriter, r *http.Request) {
	//time.Sleep(1 * time.Second)
	//panic("die stupid request body")
	fmt.Printf("%v\n", r)
	fmt.Printf("%v\n", r.Body)
	_, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Printf("%v\n", r)
	fmt.Printf("%v\n", r.Body)
	r.Body.Close()
	fmt.Printf("%v\n", r)
	fmt.Printf("%v\n", r.Body)
	//defer r.Body.Close()
}

func nestedHandler(w http.ResponseWriter, r *http.Request) {
	length := r.FormValue("length")
	if length == "" {
		length = "100"
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("http://127.0.0.1:%s/get/data?lenght=%s", os.Getenv("PORT"), length), nil)
	if err != nil {
		fmt.Printf("failed to create request %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("RESPONSE_BODY=%v: %s\n", resp.Body, err)
		w.WriteHeader(http.StatusBadGateway) // should i still close body?
		return
	}
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()
}
