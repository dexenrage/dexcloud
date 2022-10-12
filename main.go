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
	"net/http"
	"server/api"
	"server/catcherr"
	"server/config"
	"server/directory"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

func initHandlers(r *mux.Router) {
	// FileServer
	r.PathPrefix(directory.APIFileServer).HandlerFunc(fileServer).Methods(http.MethodGet)

	// API
	api.Handle(r)
}

func main() {
	defer catcherr.Recover(`main.main()`)
	r := mux.NewRouter()
	initHandlers(r)

	var (
		host    = config.String(config.Host)
		timeout = 15 * time.Second
	)

	srv := http.Server{
		Addr:         host,
		Handler:      r,
		WriteTimeout: timeout,
		ReadTimeout:  timeout,
	}
	catcherr.HandleError(srv.ListenAndServe())
}

// ToDo: Auth check
func fileServer(w http.ResponseWriter, r *http.Request) {
	defer catcherr.Recover(`main.fileServer()`)

	if strings.HasSuffix(r.URL.Path, `/`) {
		err := errors.New(catcherr.Forbidden.Description)
		catcherr.HandleAndResponse(w, catcherr.Forbidden, err)
	}

	fs := http.FileServer(http.Dir(directory.APIFileServer))
	fs.ServeHTTP(w, r)
}
