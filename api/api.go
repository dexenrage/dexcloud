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
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"server/catcherr"
	"server/database"
	"server/directory"
	"server/user"

	"github.com/gorilla/mux"
)

func HandleApi(r *mux.Router) {
	// Auth
	r.HandleFunc(directory.ApiCheckAuthHTTP, authCheckHandler).Methods(http.MethodGet)
	r.HandleFunc(directory.ApiRegisterHTTP, registerHandler).Methods(http.MethodPost)
	r.HandleFunc(directory.ApiLoginHTTP, loginHandler).Methods(http.MethodPost)

	// Files
	r.HandleFunc(directory.ApiUploadFileHTTP, fileUploadHandler).Methods(http.MethodPut)
	r.HandleFunc(directory.ApiFileListHTTP, fileListHandler).Methods(http.MethodGet)
	r.HandleFunc(directory.ApiDeleteFileHTTP, fileDeletionHandler).Methods(http.MethodDelete)
}

func getUserDir(userID string) string {
	return directory.CleanPath(directory.UserUploads(), userID)
}

func fileListHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	uid, err := GetUserID(ctx, r)
	catcherr.HandleAndResponse(w, catcherr.Unathorized, err)

	dir := getUserDir(uid)

	files, err := user.GetFiles(dir)
	catcherr.HandleAndResponse(w, catcherr.InternalServerError, err)

	data := fileListStruct{
		UserID: uid,
		Files:  files,
	}

	err = response.Send(w, responseData{http.StatusOK, data})
	catcherr.HandleError(err)
}

func authCheckHandler(w http.ResponseWriter, r *http.Request) {
	defer catcherr.Recover(`api.checkAuthHandler()`)

	_, err := checkAuth(r)
	catcherr.HandleAndResponse(w, catcherr.Unathorized, err)

	statusText := http.StatusText(http.StatusOK)
	err = response.Send(w, responseData{http.StatusOK, statusText})
	catcherr.HandleError(err)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	defer catcherr.Recover(`api.registerHandler()`)
	ctx := r.Context()

	bodyBuffer, err := io.ReadAll(r.Body)
	catcherr.HandleAndResponse(w, catcherr.InternalServerError, err)

	var acc database.User
	err = json.Unmarshal(bodyBuffer, &acc)
	catcherr.HandleAndResponse(w, catcherr.InternalServerError, err)

	err = registerUser(ctx, acc)
	catcherr.HandleAndResponse(w, catcherr.InternalServerError, err)

	data, err := createToken(acc.Login)
	catcherr.HandleAndResponse(w, catcherr.Unathorized, err)

	err = response.Send(w, responseData{http.StatusOK, data})
	catcherr.HandleError(err)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	defer catcherr.Recover(`api.loginHandler()`)
	ctx := r.Context()

	bodyBuffer, err := io.ReadAll(r.Body)
	catcherr.HandleAndResponse(w, catcherr.InternalServerError, err)

	var acc database.User
	err = json.Unmarshal(bodyBuffer, &acc)
	catcherr.HandleAndResponse(w, catcherr.InternalServerError, err)

	userInfo, err := database.GetUserInfo(ctx, acc.Login)
	catcherr.HandleAndResponse(w, catcherr.InternalServerError, err)

	err = user.ComparsePasswords(ctx, acc.Password, userInfo.Password)
	catcherr.HandleAndResponse(w, catcherr.Unathorized, err)

	data, err := createToken(acc.Login)
	catcherr.HandleAndResponse(w, catcherr.Unathorized, err)

	err = response.Send(w, responseData{http.StatusOK, data})
	catcherr.HandleAndResponse(w, catcherr.InternalServerError, err)
}

func fileUploadHandler(w http.ResponseWriter, r *http.Request) {
	defer catcherr.Recover(`api.uploadHandler()`)
	ctx := r.Context()

	uid, err := GetUserID(ctx, r)
	catcherr.HandleAndResponse(w, catcherr.Unathorized, err)

	err = r.ParseMultipartForm(32 << 20)
	catcherr.HandleAndResponse(w, catcherr.InternalServerError, err)

	destination := getUserDir(uid)

	var (
		file     multipart.File
		fileData []user.FileStruct
		fileList = r.MultipartForm.File["file"]
	)

	for _, fileHeader := range fileList {
		file, err = fileHeader.Open()

		// Save successfully uploaded files if catch an error
		if err != nil {
			saveErr := user.SaveUploadedFiles(fileData)
			catcherr.HandleAndResponse(w, catcherr.InternalServerError, saveErr)
			catcherr.HandleAndResponse(w, catcherr.InternalServerError, err)
		}
		defer file.Close()

		fileData = append(
			fileData,
			user.FileStruct{
				Directory:  destination,
				File:       file,
				FileHeader: fileHeader,
			},
		)

	}
	err = user.SaveUploadedFiles(fileData)
	catcherr.HandleAndResponse(w, catcherr.InternalServerError, err)

	statusText := http.StatusText(http.StatusOK)
	err = response.Send(w, responseData{http.StatusOK, statusText})
	catcherr.HandleError(err)
}

func fileDeletionHandler(w http.ResponseWriter, r *http.Request) {
	defer catcherr.Recover(`api.fileDeletionHandler()`)

	ctx := r.Context()

	uid, err := GetUserID(ctx, r)
	catcherr.HandleAndResponse(w, catcherr.Unathorized, err)

	userDirectory := getUserDir(uid)

	bodyBuffer, err := io.ReadAll(r.Body)
	catcherr.HandleAndResponse(w, catcherr.InternalServerError, err)

	var fileList fileListStruct
	err = json.Unmarshal(bodyBuffer, &fileList)
	catcherr.HandleAndResponse(w, catcherr.InternalServerError, err)

	for _, v := range fileList.Files {
		path := filepath.Join(userDirectory, v)
		err = os.Remove(path)
		catcherr.HandleAndResponse(w, catcherr.InternalServerError, err)
	}

	statusText := http.StatusText(http.StatusOK)
	err = response.Send(w, responseData{http.StatusOK, statusText})
	catcherr.HandleError(err)
}
