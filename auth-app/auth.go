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
	"net/http"
	"strings"
)

type Handler struct{}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.serveHttpV1(w, r)
	//h.serveHttpV2(w, r)
}

func (h *Handler) serveHttpV1(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/validate-header") {
		myHeader := r.Header.Get("X-My-Header")
		log.Printf("%s %s, X-My-Header: %s", r.Method, r.URL.Path, myHeader)
		if myHeader != "expected-value" {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("Invalid header value"))
			return
		}

		// 如果校验成功，则返回 200 OK
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Valid header"))
		return
	}
	msg := "Url Not found"
	log.Printf("%s %s, msg: %s",
		r.Method, r.URL.Path, msg)
	w.WriteHeader(http.StatusNotFound)
	_, _ = w.Write([]byte(msg))
}

func (h *Handler) serveHttpV2(w http.ResponseWriter, r *http.Request) {
	//log.Printf("%s %s", r.Method, r.URL.Path)
	myHeader := r.Header.Get("X-My-Header")
	log.Printf("%s %s, X-My-Header: %s", r.Method, r.URL.Path, myHeader)
	if myHeader != "expected-value" {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Invalid header value"))
		return
	}

	// 如果校验成功，则返回 200 OK
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Valid header"))
	return
}

func main() {
	h := &Handler{}
	err := http.ListenAndServe(":8080", h)
	if err != nil {
		log.Fatal(err)
	}
}
