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
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"server/apperror"
	"server/directory"
	"server/user"

	"github.com/gorilla/mux"
)

type Account struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func BuildApi(r *mux.Router) {
	r.HandleFunc(directory.ApiRegisterHTTP, apperror.Middleware(registerHandler)).Methods(http.MethodPost)
	r.HandleFunc(directory.ApiUploadHTTP, apperror.Middleware(uploadHandler)).Methods(http.MethodPost)
	r.HandleFunc(directory.ApiFileListHTTP, apperror.Middleware(fileListHandler)).Methods(http.MethodGet)
}

func fileListHandler(w http.ResponseWriter, r *http.Request) (err error) {
	files, err := user.GetFiles(userDir())
	if err != nil {
		log.Println(err)
		apperror.ErrInternalServerError.Err = err
		return apperror.ErrInternalServerError
	}

	data := map[string]interface{}{
		"files": files,
	}

	x, err := json.Marshal(data)
	if err != nil {
		apperror.ErrInternalServerError.Err = err
		return apperror.ErrInternalServerError
	}
	responseCustomJSON(w, http.StatusOK, string(x))
	return err
}

// temporary
func userDir() string {
	const userID = 1
	folder := strconv.Itoa(userID)
	return filepath.Join(directory.UserUploads(), folder)
}

// temporary
func userFilePath(filename string) string {
	dir := userDir()
	return filepath.Join(dir, filename)
}

func registerHandler(w http.ResponseWriter, r *http.Request) (err error) {
	bodyBuffer, _ := io.ReadAll(r.Body)
	var acc Account

	err = json.Unmarshal(bodyBuffer, &acc)
	if err != nil {
		apperror.ErrInternalServerError.Err = err
		return apperror.ErrInternalServerError
	}

	//database.RegisterUser(login, pass)

	tkn, err := signIn(&acc)
	if err != nil {
		apperror.ErrInternalServerError.Err = err
		return apperror.ErrInternalServerError
	}

	err = os.Mkdir(userDir(), os.ModePerm)
	if os.IsExist(err) {
		err = nil
	}
	if err != nil {
		apperror.ErrInternalServerError.Err = err
		return apperror.ErrInternalServerError
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

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		apperror.ErrInternalServerError.Err = err
		return apperror.ErrInternalServerError
	}
	defer file.Close()

	dst, err := os.Create(userFilePath(fileHeader.Filename))
	if err != nil {
		apperror.ErrInternalServerError.Err = err
		return apperror.ErrInternalServerError
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		apperror.ErrInternalServerError.Err = err
		return apperror.ErrInternalServerError
	}
	responseOK(w)
	return err
}
