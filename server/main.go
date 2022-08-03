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

package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	err = os.Mkdir("./uploads", os.ModePerm)
	if os.IsExist(err) {
		err = nil
	}
	if err != nil {
		log.Println(err)
	}

	dst, err := os.Create(fmt.Sprintf("./uploads/%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename)))
	if err != nil {
		log.Println(err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		log.Println(err)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join("web", "index.html")
	tmpl, err := template.ParseFiles(path)

	if err != nil {
		log.Println(err)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Println(err)
		return
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/upload", uploadHandler)

	fileServer := http.FileServer(http.Dir("./web/css"))
	r.PathPrefix("/css/").Handler(http.StripPrefix("/css", fileServer))

	srv := http.Server{
		Addr:         "localhost:80",
		Handler:      r,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatalln(srv.ListenAndServe())
}
