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
	"net/http"
	"os"
	"path/filepath"
	"server/api"
	"server/apperror"
	"server/logger"
	"time"

	"github.com/gorilla/mux"
)

func init() {
	err := os.Mkdir(`./uploads`, os.ModePerm)
	if os.IsExist(err) {
		err = nil
	}
	if err != nil {
		logger.Panicln(err)
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

/*
 * UNSAFE UNSAFE UNSAFE UNSAFE UNSAFE
 */
func prepareUserFileServer(r *mux.Router) {
	const (
		files = `./uploads`
		link  = `/uploads`
	)
	var (
		dir    = http.Dir(files)
		fs     = http.FileServer(dir)
		prefix = http.StripPrefix(link, fs)
		path   = r.PathPrefix(link + `/`)
	)
	path.Handler(prefix).Methods("GET")
}

func main() {
	defer logger.Sync()

	r := mux.NewRouter()
	r.HandleFunc("/", apperror.Middleware(indexHandler)).Methods("GET")
	r.HandleFunc("/profile", apperror.Middleware(profileHandler)).Methods("GET")

	prepareFileServer(r)
	prepareUserFileServer(r)
	api.BuildApi(r)

	srv := http.Server{
		Addr:         "localhost:80",
		Handler:      r,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	logger.Panicln(srv.ListenAndServe())
}

func indexHandler(w http.ResponseWriter, r *http.Request) (err error) {
	path := filepath.Join("web", "index.html")
	tmpl, err := template.ParseFiles(path)

	if err != nil {
		logger.Panicln(err)
		return err
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		logger.Panicln(err)
		return err
	}
	return err
}

func profileHandler(w http.ResponseWriter, r *http.Request) (err error) {
	path := filepath.Join("web", "profile.html")
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		logger.Panicln(err)
		return err
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		logger.Panicln(err)
		return err
	}
	return err
}
