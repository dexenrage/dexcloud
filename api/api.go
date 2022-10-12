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

	"server/auth"
	"server/catcherr"
	"server/database"
	"server/directory"
	"server/response"
	"server/user"

	"github.com/gorilla/mux"
)

func Handle(r *mux.Router) {
	// Auth
	r.HandleFunc(directory.APIAuthCheck, authCheckFunc).Methods(http.MethodGet)
	r.HandleFunc(directory.APIRegister, registerFunc).Methods(http.MethodPost)
	r.HandleFunc(directory.APILogin, loginFunc).Methods(http.MethodPost)

	// Files
	r.HandleFunc(directory.APIFileUpload, fileUploadFunc).Methods(http.MethodPut)
	r.HandleFunc(directory.APIFileDelete, fileDeleteFunc).Methods(http.MethodDelete)
	r.HandleFunc(directory.APIFileList, fileListFunc).Methods(http.MethodGet)
}

func fileListFunc(w http.ResponseWriter, r *http.Request) {
	defer catcherr.Recover(`api.fileListFunc()`)
	ctx := r.Context()

	login, err := auth.GetLoginFromCookie(r)
	catcherr.HandleAndResponse(w, catcherr.Unathorized, err)

	err = auth.VerifyUser(r, login)
	catcherr.HandleAndResponse(w, catcherr.Unathorized, err)

	files, err := database.GetFileList(ctx, login)
	catcherr.HandleAndResponse(w, catcherr.InternalServerError, err)

	fileList, err := json.Marshal(files)
	catcherr.HandleAndResponse(w, catcherr.InternalServerError, err)

	err = response.Send(w, response.Data{StatusCode: http.StatusOK, Data: fileList})
	catcherr.HandleError(err)
}

func authCheckFunc(w http.ResponseWriter, r *http.Request) {
	defer catcherr.Recover(`api.authCheckFunc()`)

	login, err := auth.GetLoginFromCookie(r)
	catcherr.HandleAndResponse(w, catcherr.Unathorized, err)

	err = auth.VerifyUser(r, login)
	catcherr.HandleAndResponse(w, catcherr.Unathorized, err)

	statusText := http.StatusText(http.StatusOK)
	err = response.Send(w, response.Data{
		StatusCode: http.StatusOK,
		Data:       statusText,
	})
	catcherr.HandleError(err)
}

func registerFunc(w http.ResponseWriter, r *http.Request) {
	defer catcherr.Recover(`api.registerFunc()`)
	ctx := r.Context()

	bodyBuffer, err := io.ReadAll(r.Body)
	catcherr.HandleAndResponse(w, catcherr.InternalServerError, err)

	var u database.User
	err = json.Unmarshal(bodyBuffer, &u)
	catcherr.HandleAndResponse(w, catcherr.InternalServerError, err)

	token, err := user.Register(ctx, u)
	catcherr.HandleAndResponse(w, catcherr.InternalServerError, err)

	err = response.Send(w, response.Data{StatusCode: http.StatusOK, Data: token})
	catcherr.HandleError(err)
}

func loginFunc(w http.ResponseWriter, r *http.Request) {
	defer catcherr.Recover(`api.loginFunc()`)
	ctx := r.Context()

	bodyBuffer, err := io.ReadAll(r.Body)
	catcherr.HandleAndResponse(w, catcherr.InternalServerError, err)

	var u database.User
	err = json.Unmarshal(bodyBuffer, &u)
	catcherr.HandleAndResponse(w, catcherr.InternalServerError, err)

	token, err := user.Login(ctx, u)
	catcherr.HandleAndResponse(w, catcherr.Unathorized, err)

	err = response.Send(w, response.Data{StatusCode: http.StatusOK, Data: token})
	catcherr.HandleError(err)
}

func fileUploadFunc(w http.ResponseWriter, r *http.Request) {
	defer catcherr.Recover(`api.fileUploadFunc()`)
	ctx := r.Context()

	login, err := auth.GetLoginFromCookie(r)
	catcherr.HandleAndResponse(w, catcherr.Unathorized, err)

	err = auth.VerifyUser(r, login)
	catcherr.HandleAndResponse(w, catcherr.Unathorized, err)

	err = r.ParseMultipartForm(32 << 20) // 32 MB
	catcherr.HandleAndResponse(w, catcherr.InternalServerError, err)

	fileList := r.MultipartForm.File[`file`]
	for _, fileHeader := range fileList {
		checksum, err := user.SaveFile(ctx, fileHeader)
		catcherr.HandleAndResponse(w, catcherr.InternalServerError, err)

		// ToDo: delete the file if catch an error
		err = database.SaveFileInfo(ctx, login, fileHeader.Filename, checksum)
		catcherr.HandleAndResponse(w, catcherr.InternalServerError, err)
	}

	statusText := http.StatusText(http.StatusOK)
	err = response.Send(w, response.Data{StatusCode: http.StatusOK, Data: statusText})
	catcherr.HandleError(err)
}

func fileDeleteFunc(w http.ResponseWriter, r *http.Request) {
	defer catcherr.Recover(`api.fileDeleteFunc()`)
	ctx := r.Context()

	login, err := auth.GetLoginFromCookie(r)
	catcherr.HandleAndResponse(w, catcherr.Unathorized, err)

	err = auth.VerifyUser(r, login)
	catcherr.HandleAndResponse(w, catcherr.Unathorized, err)

	files, err := database.GetFileList(ctx, login)
	catcherr.HandleAndResponse(w, catcherr.InternalServerError, err)

	bodyBuffer, err := io.ReadAll(r.Body)
	catcherr.HandleAndResponse(w, catcherr.InternalServerError, err)

	var file database.File
	err = json.Unmarshal(bodyBuffer, &file)
	catcherr.HandleAndResponse(w, catcherr.InternalServerError, err)

	for _, v := range files {
		if v.Checksum == file.Checksum {
			err = user.RemoveFile(ctx, v.Checksum)
			catcherr.HandleAndResponse(w, catcherr.InternalServerError, err)

			err = database.RemoveFileInfo(ctx, login, v.Checksum)
			catcherr.HandleAndResponse(w, catcherr.InternalServerError, err)
			break
		}
	}
	statusText := http.StatusText(http.StatusOK)
	err = response.Send(w, response.Data{StatusCode: http.StatusOK, Data: statusText})
	catcherr.HandleError(err)
}
