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

package user

import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"server/catcherr"
	"server/database"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

func GeneratePasswordHash(password string) (hash string, err error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashBytes), err
}

func ComparePasswords(
	ctx context.Context,
	wg *sync.WaitGroup,
	login, password string,
	errChan chan<- catcherr.ErrorChan,
) {
	defer wg.Done()
	user, err := database.GetUserInfo(ctx, login)
	catcherr.HandleErrorChannel(errChan, catcherr.InternalServerError, err)

	hashBytes := []byte(user.Password)
	passwordBytes := []byte(password)

	err = bcrypt.CompareHashAndPassword(hashBytes, passwordBytes)
	catcherr.HandleErrorChannel(errChan, catcherr.Unathorized, err)
}

func SaveUploadedFiles(w http.ResponseWriter, files []FileStruct) {
	for _, f := range files {
		path := filepath.Join(f.Directory, f.FileHeader.Filename)

		newFile, err := os.Create(path)
		catcherr.HandleError(w, catcherr.InternalServerError, err)
		defer newFile.Close()

		_, err = io.Copy(newFile, f.File)
		catcherr.HandleError(w, catcherr.InternalServerError, err)
	}
}

func GetFiles(w http.ResponseWriter, dir string) (files []string) {
	dirEntry, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		err = os.Mkdir(dir, os.ModePerm)
	}
	catcherr.HandleError(w, catcherr.InternalServerError, err)

	for _, f := range dirEntry {
		files = append(files, f.Name())
	}
	return files
}
