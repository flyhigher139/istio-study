// Copyright 2023 igevin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Handler struct{}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.serveHttpV1(w, r)
	//h.serveHttpV2(w, r)
}

func (h *Handler) serveHttpV1(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, "/validate-header") {
		msg := "Url Not found"
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(msg))
	}

	myHeader := r.Header.Get("X-My-Header")
	if !validateHeader(myHeader) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Invalid header value"))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Valid header"))
	return
}

func (h *Handler) serveHttpV2(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, "/validate-header") {
		msg := "Url Not found"
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(msg))
	}

	myHeader := r.Header.Get("X-My-Header")

	valid := false
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		valid = validateHeader(myHeader)
		wg.Done()
	}()
	wg.Wait()
	if !valid {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Invalid header value"))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Valid header"))
	return
}

func validateHeader(header string) bool {
	t := rand.Intn(100)
	time.Sleep(time.Duration(t) * time.Millisecond)
	return header == "expected-value"
}

func main() {
	h := &Handler{}
	err := http.ListenAndServe(":8080", h)
	if err != nil {
		log.Fatal(err)
	}
}
