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
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"server/catcherr"
	"server/database"
	"server/directory"
	"server/user"

	"github.com/gorilla/mux"
)

func HandleApi(ctx context.Context, r *mux.Router) {
	r.HandleFunc(directory.ApiCheckAuthHTTP, checkAuthHandler).Methods(http.MethodGet)
	r.HandleFunc(directory.ApiRegisterHTTP, registerHandler).Methods(http.MethodPost)
	r.HandleFunc(directory.ApiLoginHTTP, loginHandler).Methods(http.MethodPost)
	r.HandleFunc(directory.ApiUploadHTTP, uploadHandler).Methods(http.MethodPost)
	r.HandleFunc(directory.ApiFileListHTTP, fileListHandler).Methods(http.MethodGet)
}

func defaultContextTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5000*time.Second)
}

func getUserDir(userID string) string { return filepath.Join(directory.UserUploads(), userID) }

func fileListHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := defaultContextTimeout()
	defer cancel()

	var data fileListStruct
	data.UserID = GetUserID(ctx, w, r)
	data.Files = user.GetFiles(ctx, w, getUserDir(data.UserID))

	response.Send(ctx, w, responseData{http.StatusOK, data})
}

func checkAuthHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := defaultContextTimeout()
	defer cancel()

	defer catcherr.RecoverState(`api.checkAuthHandler`)
	parseToken(ctx, w, r)
	response.Send(ctx, w, responseData{http.StatusOK, `Authorized`})
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := defaultContextTimeout()
	defer cancel()

	defer catcherr.RecoverState(`api.registerHandler`)

	bodyBuffer, err := io.ReadAll(r.Body)
	catcherr.HandleError(w, catcherr.InternalServerError, err)

	var acc database.User
	err = json.Unmarshal(bodyBuffer, &acc)
	catcherr.HandleError(w, catcherr.InternalServerError, err)

	acc.HashedPassword = user.GeneratePasswordHash(ctx, w, acc.HashedPassword)
	acc = database.RegisterUser(ctx, w, acc)

	userID := fmt.Sprint(acc.ID)

	err = os.Mkdir(getUserDir(userID), os.ModePerm)
	catcherr.HandleError(w, catcherr.InternalServerError, err)

	data := createToken(ctx, w, acc.Login)
	response.Send(ctx, w, responseData{http.StatusOK, data})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := defaultContextTimeout()
	defer cancel()

	defer catcherr.RecoverState(`api.loginHandler`)

	bodyBuffer, err := io.ReadAll(r.Body)
	catcherr.HandleError(w, catcherr.InternalServerError, err)

	var acc database.User
	err = json.Unmarshal(bodyBuffer, &acc)
	catcherr.HandleError(w, catcherr.InternalServerError, err)

	user.CompareLoginCredentials(ctx, w, acc.Login, acc.HashedPassword)

	data := createToken(ctx, w, acc.Login)
	response.Send(ctx, w, responseData{http.StatusOK, data})
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	defer catcherr.RecoverState(`api.uploadHandler`)
	userID := GetUserID(ctx, w, r)

	file, fileHeader, err := r.FormFile("file")
	catcherr.HandleError(w, catcherr.InternalServerError, err)
	defer file.Close()

	f := user.FileStruct{
		Directory:  getUserDir(userID),
		File:       file,
		FileHeader: fileHeader,
	}
	user.SaveUploadedFile(ctx, w, f)
	response.Send(ctx, w, responseData{http.StatusOK, `OK`})
}
