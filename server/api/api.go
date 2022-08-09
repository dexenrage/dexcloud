/*
Copyright 2022 dexenrage

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Account struct {
	Login    string `json:"email"`
	Password string `json:"password"`
}

func BuildApi(r *mux.Router) {
	r.HandleFunc(`/api/register`, registerHandler).Methods("POST")
	r.HandleFunc(`/api/upload`, uploadHandler).Methods("POST")
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	bodyBuffer, _ := io.ReadAll(r.Body)

	var acc Account

	err := json.Unmarshal(bodyBuffer, &acc)
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, `{ "message": "OK" }`)

	//database.RegisterUser(login, pass)
	//http.Redirect(w, r, "/profile", http.StatusFound)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	const (
		createFileDIR = `./uploads/%s`
		redirectPath  = `/profile`
	)

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, redirectPath, http.StatusFound)
		return
	}
	defer file.Close()

	var (
		filename = fileHeader.Filename
		filepath = fmt.Sprintf(createFileDIR, filename)
	)

	dst, err := os.Create(filepath)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, redirectPath, http.StatusFound)
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, redirectPath, http.StatusFound)
	}
	http.Redirect(w, r, redirectPath, http.StatusFound)
}
