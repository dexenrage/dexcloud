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
	UserID   string `json:"userid"`
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

/*func customResponse(w http.ResponseWriter, status int, data []byte) {
	defer catcherr.RecoverState(`api.customResponse`)

	w.Header().Set(defaultResponseType, defaultResponseValue)
	w.WriteHeader(status)

	_, err := w.Write(data)
	catcherr.HandleError(w, catcherr.InternalServerError, err)
}*/

func fileListHandler(w http.ResponseWriter, r *http.Request) {
	var (
		userID = GetUserID(w, r)
		files  = user.GetFiles(w, getUserDir(userID))
		data   = map[string]interface{}{
			"userid": userID,
			"files":  files,
		}
	)
	response.Send(w, http.StatusOK, data)
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

	token, expires := createToken(w, acc.Login)
	data := map[string]string{
		"login":   acc.Login,
		"token":   token,
		"expires": expires,
	}
	response.Send(w, http.StatusOK, data)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	defer catcherr.RecoverState(`api.loginHandler`)

	bodyBuffer, err := io.ReadAll(r.Body)
	catcherr.HandleError(w, catcherr.InternalServerError, err)

	var acc Account
	err = json.Unmarshal(bodyBuffer, &acc)
	catcherr.HandleError(w, catcherr.InternalServerError, err)

	user.CompareLoginCredentials(w, acc.Login, acc.Password)

	token, expires := createToken(w, acc.Login)
	data := map[string]string{
		"login":   acc.Login,
		"token":   token,
		"expires": expires,
	}
	response.Send(w, http.StatusOK, data)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	defer catcherr.RecoverState(`api.uploadHandler`)
	userID := GetUserID(w, r)

	file, fileHeader, err := r.FormFile("file")
	catcherr.HandleError(w, catcherr.InternalServerError, err)
	defer file.Close()

	f := user.FileStruct{
		Directory:  getUserDir(userID),
		File:       file,
		FileHeader: fileHeader,
	}
	user.SaveUploadedFile(w, f)
	response.Send(w, http.StatusOK, `OK`)
}
