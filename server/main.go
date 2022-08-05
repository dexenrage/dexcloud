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
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"server/api"
	"time"

	"github.com/gorilla/mux"
)

func init() {
	err := os.Mkdir(`./uploads`, os.ModePerm)
	if os.IsExist(err) {
		err = nil
	}
	if err != nil {
		log.Fatalln(err)
	}
}

func prepareFileServer(r *mux.Router) {
	const (
		files  = `./web/static`
		static = `/static`
	)

	var (
		dir = http.Dir(files)
		fs  = http.FileServer(dir)

		prefix = http.StripPrefix(static, fs)
		path   = r.PathPrefix(static + `/`)
	)

	path.Handler(prefix).Methods("GET")
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler).Methods("GET")
	r.HandleFunc("/profile", profileHandler).Methods("GET")

	prepareFileServer(r)
	api.BuildApi(r)

	srv := http.Server{
		Addr:         "localhost:80",
		Handler:      r,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatalln(srv.ListenAndServe())
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

func profileHandler(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join("web", "profile.html")
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
