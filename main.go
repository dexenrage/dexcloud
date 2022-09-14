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
	"errors"
	"fmt"
	"net/http"
	"server/api"
	"server/catcherr"
	"server/config"
	"server/directory"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

func init() {
	err := directory.CreateCriticalDirectories()
	catcherr.HandleError(err)
}

func initHandlers(r *mux.Router) {
	// FileServers
	r.PathPrefix(directory.UploadsHTTP).HandlerFunc(uploadsFileServerHandler).Methods(http.MethodGet)

	// API
	api.HandleApi(r)
}

func main() {
	defer catcherr.Recover(`main.main()`)

	r := mux.NewRouter()
	initHandlers(r)

	srv := http.Server{
		Addr:         config.String(`host`),
		Handler:      r,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	catcherr.HandleError(srv.ListenAndServe())
}

func getFileServerHandler(w http.ResponseWriter, r *http.Request, prefix, dir string) http.Handler {
	if strings.HasSuffix(r.URL.Path, directory.Slash) {
		err := errors.New(catcherr.Forbidden.Description)
		catcherr.HandleAndResponse(w, catcherr.Forbidden, err)
	}

	fs := http.FileServer(http.Dir(dir))
	return http.StripPrefix(prefix, fs)
}

func uploadsFileServerHandler(w http.ResponseWriter, r *http.Request) {
	defer catcherr.Recover(`main.uploadsFileServerHandler()`)
	ctx := r.Context()

	uid, err := api.GetUserID(ctx, r)
	catcherr.HandleAndResponse(w, catcherr.Unathorized, err)

	dir := fmt.Sprint(directory.UploadsHTTP, uid, directory.Slash)

	if !strings.HasPrefix(r.URL.Path, dir) {
		err := errors.New(catcherr.Forbidden.Description)
		catcherr.HandleAndResponse(w, catcherr.Forbidden, err)
	}

	handler := getFileServerHandler(w, r, directory.UploadsHTTP, directory.UserUploads())
	handler.ServeHTTP(w, r)
}
