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
	"io"
	"net/http"
	"path/filepath"
	"sync"
	"time"

	"server/catcherr"
	"server/database"
	"server/directory"
	"server/user"

	"github.com/gorilla/mux"
)

func HandleApi(r *mux.Router) {
	r.HandleFunc(directory.ApiCheckAuthHTTP, checkAuthHandler).Methods(http.MethodGet)
	r.HandleFunc(directory.ApiRegisterHTTP, registerHandler).Methods(http.MethodPost)
	r.HandleFunc(directory.ApiLoginHTTP, loginHandler).Methods(http.MethodPost)
	r.HandleFunc(directory.ApiUploadHTTP, uploadHandler).Methods(http.MethodPut)
	r.HandleFunc(directory.ApiFileListHTTP, fileListHandler).Methods(http.MethodGet)
}

func defaultContextTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, 15*time.Second)
}

func getUserDir(userID string) string { return filepath.Join(directory.UserUploads(), userID) }

func fileListHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := defaultContextTimeout(context.Background())
	defer cancel()

	var data fileListStruct
	data.UserID = GetUserID(ctx, w, r)
	data.Files = user.GetFiles(w, getUserDir(data.UserID))

	response.Send(w, responseData{http.StatusOK, data})
}

func checkAuthHandler(w http.ResponseWriter, r *http.Request) {
	defer catcherr.RecoverState(`api.checkAuthHandler`)
	parseToken(w, r)
	response.Send(w, responseData{http.StatusOK, `Authorized`})
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	defer catcherr.RecoverState(`api.registerHandler`)

	ctx, cancel := defaultContextTimeout(context.Background())
	defer cancel()

	bodyBuffer, err := io.ReadAll(r.Body)
	catcherr.HandleError(w, catcherr.InternalServerError, err)

	var acc database.User
	err = json.Unmarshal(bodyBuffer, &acc)
	catcherr.HandleError(w, catcherr.InternalServerError, err)

	var (
		wg        = &sync.WaitGroup{}
		tokenChan = make(chan tokenData, 1)
		errChan   = make(chan catcherr.ErrorChan, 1)
	)
	defer close(tokenChan)
	defer close(errChan)

	wg.Add(2)
	go goRegisterUser(ctx, wg, acc, errChan)
	go goGetTokenData(wg, acc.Login, tokenChan, errChan)

	waitChan := waitGorutines(wg)
	select {
	case <-ctx.Done():
		catcherr.HandleError(w, catcherr.InternalServerError, ctx.Err())
	case <-errChan:
		catcherr.HandleError(w, (<-errChan).CustomError, (<-errChan).Error)
	case <-waitChan:
		response.Send(w, responseData{http.StatusOK, <-tokenChan})
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	defer catcherr.RecoverState(`api.loginHandler`)

	ctx, cancel := defaultContextTimeout(context.Background())
	defer cancel()

	bodyBuffer, err := io.ReadAll(r.Body)
	catcherr.HandleError(w, catcherr.InternalServerError, err)

	var acc database.User
	err = json.Unmarshal(bodyBuffer, &acc)
	catcherr.HandleError(w, catcherr.InternalServerError, err)

	var (
		wg        = &sync.WaitGroup{}
		tokenChan = make(chan tokenData, 1)
		errChan   = make(chan catcherr.ErrorChan)
	)
	defer close(tokenChan)
	defer close(errChan)

	wg.Add(2)
	go user.ComparePasswords(ctx, wg, acc.Login, acc.Password, errChan)
	go goGetTokenData(wg, acc.Login, tokenChan, errChan)

	wgDoneChan := waitGorutines(wg)
	select {
	case <-ctx.Done():
		catcherr.HandleError(w, catcherr.InternalServerError, ctx.Err())
	case data := <-errChan:
		catcherr.HandleError(w, data.CustomError, data.Error)
	case <-wgDoneChan:
		response.Send(w, responseData{http.StatusOK, <-tokenChan})
	}
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
	user.SaveUploadedFile(w, f)
	response.Send(w, responseData{http.StatusOK, `OK`})
}
