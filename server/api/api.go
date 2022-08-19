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

	"server/apperror"
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
	r.HandleFunc(directory.ApiRegisterHTTP, apperror.Middleware(registerHandler)).Methods(http.MethodPost)
	r.HandleFunc(directory.ApiLoginHTTP, apperror.Middleware(loginHandler)).Methods(http.MethodPost)
	r.HandleFunc(directory.ApiUploadHTTP, apperror.Middleware(uploadHandler)).Methods(http.MethodPost)
	r.HandleFunc(directory.ApiFileListHTTP, apperror.Middleware(fileListHandler)).Methods(http.MethodGet)
}

func fileListHandler(w http.ResponseWriter, r *http.Request) (err error) {
	userID, err := r.Cookie("userid")
	if err != nil {
		if err == http.ErrNoCookie {
			return apperror.ErrUnathorized
		}
		return apperror.ErrBadRequest
	}

	files, err := user.GetFiles(userDir(userID.Value))
	if err != nil {
		return apperror.ErrInternalServerError
	}

	data := map[string]interface{}{
		"files": files,
	}

	x, err := json.Marshal(data)
	if err != nil {
		return apperror.ErrInternalServerError
	}
	responseCustomJSON(w, http.StatusOK, x)
	return err
}

func userDir(userID string) string {
	return filepath.Join(directory.UserUploads(), userID)
}

func registerHandler(w http.ResponseWriter, r *http.Request) (err error) {
	bodyBuffer, err := io.ReadAll(r.Body)
	if err != nil {
		return apperror.ErrInternalServerError
	}

	var acc Account
	err = json.Unmarshal(bodyBuffer, &acc)
	if err != nil {
		return apperror.ErrInternalServerError
	}

	acc.UserID, err = database.RegisterUser(acc.Login, acc.Password)
	if err != nil {
		return apperror.ErrInternalServerError
	}

	err = os.Mkdir(userDir(acc.UserID), os.ModePerm)
	if err != nil {
		return apperror.ErrInternalServerError
	}

	token, err := createToken(&acc)
	if err != nil {
		return err
	}

	dataMap := map[string]string{
		"userid": acc.UserID,
		"token":  token,
	}

	data, err := json.Marshal(dataMap)
	if err != nil {
		return apperror.ErrInternalServerError
	}
	responseCustomJSON(w, http.StatusOK, data)
	return err
}

func loginHandler(w http.ResponseWriter, r *http.Request) (err error) {
	bodyBuffer, err := io.ReadAll(r.Body)
	if err != nil {
		return apperror.ErrInternalServerError
	}

	var acc Account
	err = json.Unmarshal(bodyBuffer, &acc)
	if err != nil {
		return apperror.ErrInternalServerError
	}

	err = user.CompareLoginData(acc.Login, acc.Password)
	if err != nil {
		return err
	}

	acc.UserID, err = database.GetUserID(acc.Login)
	if err != nil {
		return err
	}

	token, err := createToken(&acc)
	if err != nil {
		return err
	}

	dataMap := map[string]string{
		"userid": acc.UserID,
		"token":  token,
	}

	data, err := json.Marshal(dataMap)
	if err != nil {
		return apperror.ErrInternalServerError
	}
	responseCustomJSON(w, http.StatusOK, data)
	return err
}

func uploadHandler(w http.ResponseWriter, r *http.Request) (err error) {
	err = checkToken(w, r)
	if err != nil {
		return err
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		return apperror.ErrInternalServerError
	}
	defer file.Close()

	userID, err := r.Cookie("userid")
	if err != nil {
		if err == http.ErrNoCookie {
			return apperror.ErrUnathorized
		}
		return apperror.ErrBadRequest
	}

	path := filepath.Join(userDir(userID.Value), fileHeader.Filename)
	dst, err := os.Create(path)
	if err != nil {
		return apperror.ErrInternalServerError
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return apperror.ErrInternalServerError
	}
	responseOK(w)
	return err
}
