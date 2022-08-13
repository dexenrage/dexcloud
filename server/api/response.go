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
	"io"
	"net/http"
)

const (
	defaultType  = "Content-Type"
	defaultValue = "application/json"
)

func responseCustomJSON(w http.ResponseWriter, status int, msg string) {
	w.Header().Set(defaultType, defaultValue)
	w.WriteHeader(status)
	io.WriteString(w, msg)
}

func responseOK(w http.ResponseWriter) {
	w.Header().Set(defaultType, defaultValue)
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, `{ "message": "OK" }`)
}

func responseCreated(w http.ResponseWriter) {
	w.Header().Set(defaultType, defaultValue)
	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, `{ "message": "Created" }`)
}

func responseInternalServerError(w http.ResponseWriter) {
	w.Header().Set(defaultType, defaultValue)
	w.WriteHeader(http.StatusInternalServerError)
	io.WriteString(w, `{ "message": "Internal Server Error" }`)
}
