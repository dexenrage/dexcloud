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
	"strings"
	"time"

	"github.com/gorilla/mux"
)

func init() {
	err := os.Mkdir(directory.StaticFiles(), os.ModePerm)
	if os.IsExist(err) {
		err = nil
	}
	if err != nil {
		logger.Fatalln(err)
	}

	err = os.Mkdir(directory.UserUploads(), os.ModePerm)
	if os.IsExist(err) {
		err = nil
	}
	if err != nil {
		logger.Fatalln(err)
	}
}

func main() {
	defer logger.Sync()

	r := mux.NewRouter()

	r.PathPrefix(directory.UploadsHTTP).Handler(apperror.Middleware(uploadsFileServer)).Methods(http.MethodGet)
	r.PathPrefix(directory.StaticHTTP).Handler(apperror.Middleware(staticFileServer)).Methods(http.MethodGet)

	r.HandleFunc(directory.IndexHTTP, apperror.Middleware(indexHandler)).Methods(http.MethodGet)
	r.HandleFunc(directory.ProfileHTTP, apperror.Middleware(profileHandler)).Methods(http.MethodGet)
	api.HandleApi(r)

	srv := http.Server{
		Addr:         "localhost:80",
		Handler:      r,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	logger.Panicln(srv.ListenAndServe())
}

func fileServerHandler(prefix, dir string) http.Handler {
	fs := http.FileServer(http.Dir(dir))
	return http.StripPrefix(prefix, fs)
}

func staticFileServer(w http.ResponseWriter, r *http.Request) (err error) {
	handler := fileServerHandler(directory.StaticHTTP, directory.StaticFiles())
	if strings.HasSuffix(r.URL.Path, "/") {
		return apperror.ErrForbidden
	}
	handler.ServeHTTP(w, r)
	return err
}

func uploadsFileServer(w http.ResponseWriter, r *http.Request) (err error) {
	handler := fileServerHandler(directory.UploadsHTTP, directory.UserUploads())
	if strings.HasSuffix(r.URL.Path, "/") {
		return apperror.ErrForbidden
	}
	handler.ServeHTTP(w, r)
	return err
}

func indexHandler(w http.ResponseWriter, r *http.Request) (err error) {
	tmpl, err := template.ParseFiles(directory.IndexPage())
	if err != nil {
		logger.Panicln(err)
		return apperror.ErrInternalServerError
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		logger.Panicln(err)
		return apperror.ErrInternalServerError
	}
	return err
}

func profileHandler(w http.ResponseWriter, r *http.Request) (err error) {
	tmpl, err := template.ParseFiles(directory.ProfilePage())
	if err != nil {
		logger.Panicln(err)
		return apperror.ErrInternalServerError
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		logger.Panicln(err)
		return apperror.ErrInternalServerError
	}
	return err
}
