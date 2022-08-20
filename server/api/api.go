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
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"server/catcherr"
	"server/database"
	"server/directory"
	"server/user"

	"github.com/gorilla/mux"
)

type Account struct {
	UserID   string `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

func HandleApi(r *mux.Router) {
	r.HandleFunc(directory.ApiRegisterHTTP, registerHandler).Methods(http.MethodPost)
	r.HandleFunc(directory.ApiLoginHTTP, loginHandler).Methods(http.MethodPost)
	r.HandleFunc(directory.ApiUploadHTTP, uploadHandler).Methods(http.MethodPost)
	r.HandleFunc(directory.ApiFileListHTTP, fileListHandler).Methods(http.MethodGet)
}

func getUserDir(userID string) string { return filepath.Join(directory.UserUploads(), userID) }

const (
	defaultResponseType  = "Content-Type"
	defaultResponseValue = "application/json"
)

func customResponse(w http.ResponseWriter, status int, data []byte) {
	defer catcherr.RecoverState(`api.customResponse`)

	w.Header().Set(defaultResponseType, defaultResponseValue)
	w.WriteHeader(status)

	_, err := w.Write(data)
	catcherr.HandleError(w, catcherr.InternalServerError, err)
}

func fileListHandler(w http.ResponseWriter, r *http.Request) {
	defer catcherr.RecoverState(`api.fileListHandler`)

	userID, err := r.Cookie("userid")
	if errors.Is(err, http.ErrNoCookie) {
		catcherr.HandleError(w, catcherr.Unathorized, err)
	}
	catcherr.HandleError(w, catcherr.BadRequest, err)

	files := user.GetFiles(w, getUserDir(userID.Value))
	catcherr.HandleError(w, catcherr.InternalServerError, err)

	data := map[string]interface{}{
		"files": files,
	}

	x, err := json.Marshal(data)
	catcherr.HandleError(w, catcherr.InternalServerError, err)
	customResponse(w, http.StatusOK, x)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	defer catcherr.RecoverState(`api.registerHandler`)

	bodyBuffer, err := io.ReadAll(r.Body)
	catcherr.HandleError(w, catcherr.InternalServerError, err)

	var acc Account
	err = json.Unmarshal(bodyBuffer, &acc)
	catcherr.HandleError(w, catcherr.InternalServerError, err)

	acc.Password = user.GeneratePasswordHash(w, acc.Password)
	acc.UserID = database.RegisterUser(w, acc.Login, acc.Password)

	err = os.Mkdir(getUserDir(acc.UserID), os.ModePerm)
	catcherr.HandleError(w, catcherr.InternalServerError, err)

	token, expiresAt := createToken(w, acc.Login)

	dataMap := map[string]string{
		"userid":  acc.UserID,
		"token":   token,
		"expires": expiresAt,
	}

	data, err := json.Marshal(dataMap)
	catcherr.HandleError(w, catcherr.InternalServerError, err)
	customResponse(w, http.StatusOK, data)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	defer catcherr.RecoverState(`api.loginHandler`)

	bodyBuffer, err := io.ReadAll(r.Body)
	catcherr.HandleError(w, catcherr.InternalServerError, err)

	var acc Account
	err = json.Unmarshal(bodyBuffer, &acc)
	catcherr.HandleError(w, catcherr.InternalServerError, err)

	user.CompareLoginCredentials(w, acc.Login, acc.Password)
	acc.UserID = database.GetUserID(w, acc.Login)

	token, expiresAt := createToken(w, acc.Login)

	dataMap := map[string]string{
		"userid":  acc.UserID,
		"token":   token,
		"expires": expiresAt,
	}

	data, err := json.Marshal(dataMap)
	catcherr.HandleError(w, catcherr.Unathorized, err)

	customResponse(w, http.StatusOK, data)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	defer catcherr.RecoverState(`api.uploadHandler`)
	parseToken(w, r)

	file, fileHeader, err := r.FormFile("file")
	catcherr.HandleError(w, catcherr.InternalServerError, err)
	defer file.Close()

	userID, err := r.Cookie("userid")
	if errors.Is(err, http.ErrNoCookie) {
		catcherr.HandleError(w, catcherr.Unathorized, err)
	}
	catcherr.HandleError(w, catcherr.BadRequest, err)

	path := filepath.Join(getUserDir(userID.Value), fileHeader.Filename)

	dst, err := os.Create(path)
	catcherr.HandleError(w, catcherr.InternalServerError, err)
	defer dst.Close()

	_, err = io.Copy(dst, file)
	catcherr.HandleError(w, catcherr.InternalServerError, err)
	customResponse(w, http.StatusOK, []byte(`OK`))
}
