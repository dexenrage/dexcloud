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
	"server/api"
	"server/apperror"
	"server/directory"
	"server/logger"
	"time"

	"github.com/gorilla/mux"
)

func init() {
	err := os.Mkdir(directory.UserUploads(), os.ModePerm)
	if os.IsExist(err) {
		err = nil
	}
	if err != nil {
		logger.Panicln(err)
	}
}

func prepareFileServer(r *mux.Router) {
	var (
		dir    = http.Dir(directory.StaticFiles())
		fs     = http.FileServer(dir)
		prefix = http.StripPrefix(directory.StaticHTTP, fs)
		path   = r.PathPrefix(directory.StaticHTTP + directory.Slash)
	)
	path.Handler(prefix).Methods(http.MethodGet)
}

/*
 * UNSAFE UNSAFE UNSAFE UNSAFE UNSAFE
 */
func prepareUserFileServer(r *mux.Router) {
	var (
		dir    = http.Dir(directory.UserUploads())
		fs     = http.FileServer(dir)
		prefix = http.StripPrefix(directory.UploadsHTTP, fs)
		path   = r.PathPrefix(directory.UploadsHTTP + directory.Slash)
	)
	path.Handler(prefix).Methods(http.MethodGet)
}

func main() {
	defer logger.Sync()

	r := mux.NewRouter()
	r.HandleFunc(directory.IndexHTTP, apperror.Middleware(indexHandler)).Methods(http.MethodGet)
	r.HandleFunc(directory.ProfileHTTP, apperror.Middleware(profileHandler)).Methods(http.MethodGet)

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
	tmpl, err := template.ParseFiles(directory.IndexPage())
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
	tmpl, err := template.ParseFiles(directory.ProfilePage())
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
