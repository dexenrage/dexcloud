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
	"net/http"
	"os"

	"server/logger"
	"server/user"

	"github.com/gorilla/mux"
)

type Account struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func BuildApi(r *mux.Router) {
	r.HandleFunc(`/api/register`, registerHandler).Methods("POST")
	r.HandleFunc(`/api/upload`, uploadHandler).Methods("POST")
	r.HandleFunc(`/api/filelist`, fileListHandler).Methods("GET")
}

func fileListHandler(w http.ResponseWriter, r *http.Request) {
	f, err := user.GetFiles(1)
	if err != nil {
		logger.Errorln(err)
		responseInternalServerError(w)
		return
	}

	data := map[string]interface{}{
		"files": f,
	}

	x, err := json.Marshal(data)
	if err != nil {
		logger.Errorln(err)
		responseInternalServerError(w)
		return
	}
	responseCustomJSON(w, http.StatusOK, string(x))
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	bodyBuffer, _ := io.ReadAll(r.Body)
	var acc Account

	err := json.Unmarshal(bodyBuffer, &acc)
	if err != nil {
		logger.Errorln(err)
		responseInternalServerError(w)
		return
	}

	//database.RegisterUser(login, pass)

	// Temporary
	userID := 1
	dir := fmt.Sprint(`./uploads/`, userID)

	err = os.Mkdir(dir, os.ModePerm)
	if err != nil {
		logger.Errorln(err)
		responseInternalServerError(w)
		return
	}
	responseCreated(w)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	const (
		userID        = 1
		createFileDIR = `./uploads/%d/%s`
		redirectPath  = `/profile`
	)

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		logger.Errorln(err)
		responseInternalServerError(w)
		return
	}
	defer file.Close()

	var (
		filename = fileHeader.Filename
		filepath = fmt.Sprintf(createFileDIR, userID, filename)
	)

	dst, err := os.Create(filepath)
	if err != nil {
		logger.Errorln(err)
		responseInternalServerError(w)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		logger.Errorln(err)
		responseInternalServerError(w)
		return
	}
	responseOK(w)
}
