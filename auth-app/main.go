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
)

func headerValidationHandler(w http.ResponseWriter, r *http.Request) {
	// 提取请求头
	myHeader := r.Header.Get("X-My-Header")

	// 执行自定义校验逻辑
	if myHeader != "expected-value" {
		// 如果校验失败，则返回 401 Unauthorized
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid header value"))
		return
	}

	// 如果校验成功，则返回 200 OK
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Valid header"))
}

func main() {
	http.HandleFunc("/validate-header", headerValidationHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
