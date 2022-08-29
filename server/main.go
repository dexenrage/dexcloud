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
	"context"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"server/api"
	"server/catcherr"
	"server/directory"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

func init() {
	directory.CreateCriticalDirectories()
}

func initHandlers(r *mux.Router) {
	// Pages
	r.HandleFunc(directory.IndexHTTP, indexHandler).Methods(http.MethodGet)
	r.HandleFunc(directory.RegisterHTTP, registerHandler).Methods(http.MethodGet)
	r.HandleFunc(directory.LoginHTTP, loginHandler).Methods(http.MethodGet)
	r.HandleFunc(directory.ProfileHTTP, profileHandler).Methods(http.MethodGet)

	// FileServers
	r.PathPrefix(directory.UploadsHTTP).HandlerFunc(uploadsFileServerHandler).Methods(http.MethodGet)
	r.PathPrefix(directory.StaticHTTP).HandlerFunc(staticFileServerHandler).Methods(http.MethodGet)

	// API
	api.HandleApi(r)
}

func main() {
	r := mux.NewRouter()
	initHandlers(r)

	srv := http.Server{
		Addr:         "localhost:80",
		Handler:      r,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Panicln(srv.ListenAndServe())
}

func getFileServerHandler(w http.ResponseWriter, r *http.Request, prefix, dir string) http.Handler {
	if strings.HasSuffix(r.URL.Path, "/") {
		err := errors.New(`The user is not allowed to enter this directory`)
		catcherr.HandleError(w, catcherr.Forbidden, err)
	}

	fs := http.FileServer(http.Dir(dir))
	return http.StripPrefix(prefix, fs)
}

func staticFileServerHandler(w http.ResponseWriter, r *http.Request) {
	defer catcherr.RecoverState(`main.staticFileServerHandler`)
	handler := getFileServerHandler(w, r, directory.StaticHTTP, directory.StaticFiles())
	handler.ServeHTTP(w, r)
}

func uploadsFileServerHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	defer catcherr.RecoverState(`main.uploadsFileServerHandler`)

	var (
		userID  = api.GetUserID(ctx, w, r)
		userDir = fmt.Sprint(directory.UploadsHTTP, userID, `/`)
	)

	if !strings.HasPrefix(r.URL.Path, userDir) {
		err := errors.New(`The user is not allowed to enter this directory`)
		catcherr.HandleError(w, catcherr.Forbidden, err)
	}

	handler := getFileServerHandler(w, r, directory.UploadsHTTP, directory.UserUploads())
	handler.ServeHTTP(w, r)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	defer catcherr.RecoverState(`main.indexHandler`)
	executeTemplate(w, directory.IndexPage())
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	defer catcherr.RecoverState(`main.profileHandler`)
	executeTemplate(w, directory.ProfilePage())
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	defer catcherr.RecoverState(`main.registerHandler`)
	executeTemplate(w, directory.RegisterPage())
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	defer catcherr.RecoverState(`main.loginHandler`)
	executeTemplate(w, directory.LoginPage())
}

func executeTemplate(w http.ResponseWriter, directory string) {
	if directory == `` {
		err := errors.New(`Template directory is empty`)
		catcherr.HandleError(w, catcherr.InternalServerError, err)
	}

	tmpl, err := template.ParseFiles(directory)
	catcherr.HandleError(w, catcherr.InternalServerError, err)

	err = tmpl.Execute(w, nil)
	catcherr.HandleError(w, catcherr.InternalServerError, err)
}
