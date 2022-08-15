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

	"server/apperror"
	"server/logger"
	"server/user"

	"github.com/gorilla/mux"
)

type Account struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func BuildApi(r *mux.Router) {
	r.HandleFunc(`/api/register`, apperror.Middleware(registerHandler)).Methods("POST")
	r.HandleFunc(`/api/upload`, apperror.Middleware(uploadHandler)).Methods("POST")
	r.HandleFunc(`/api/filelist`, apperror.Middleware(fileListHandler)).Methods("GET")
}

func fileListHandler(w http.ResponseWriter, r *http.Request) (err error) {
	f, err := user.GetFiles(1)
	if err != nil {
		logger.Errorln(err)
		responseInternalServerError(w)
		return err
	}

	data := map[string]interface{}{
		"files": f,
	}

	x, err := json.Marshal(data)
	if err != nil {
		logger.Errorln(err)
		responseInternalServerError(w)
		return err
	}
	responseCustomJSON(w, http.StatusOK, string(x))
	return err
}

var Acl Account

func registerHandler(w http.ResponseWriter, r *http.Request) (err error) {
	bodyBuffer, _ := io.ReadAll(r.Body)
	var acc Account

	err = json.Unmarshal(bodyBuffer, &acc)
	if err != nil {
		logger.Errorln(err)
		responseInternalServerError(w)
		return err
	}

	Acl.Login = acc.Login
	Acl.Password = acc.Password

	//database.RegisterUser(login, pass)

	tkn, err := signIn(&acc)
	if err != nil {
		logger.Errorln(err)
		responseInternalServerError(w)
		return err
	}

	// Temporary
	userID := 1
	dir := fmt.Sprint(`./uploads/`, userID)

	err = os.Mkdir(dir, os.ModePerm)
	if os.IsExist(err) {
		err = nil
	}
	if err != nil {
		logger.Errorln(err)
		responseInternalServerError(w)
		return err
	}
	jsonResp := fmt.Sprintf(`{ "token": "%s" }`, tkn)
	responseCustomJSON(w, http.StatusCreated, jsonResp)
	return err
}

func uploadHandler(w http.ResponseWriter, r *http.Request) (err error) {
	// temporary
	err = checkAuth(w, r)
	if err != nil {
		return err
	}
	const (
		userID        = 1
		createFileDIR = `./uploads/%d/%s`
		redirectPath  = `/profile`
	)

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		logger.Errorln(err)
		responseInternalServerError(w)
		return err
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
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		logger.Errorln(err)
		responseInternalServerError(w)
		return err
	}
	responseOK(w)
	return err
}
