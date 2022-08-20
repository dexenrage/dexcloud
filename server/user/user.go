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
	"net/http"
	"os"
	"server/catcherr"
	"server/database"

	"golang.org/x/crypto/bcrypt"
)

func GeneratePasswordHash(w http.ResponseWriter, password string) (hash string) {
	defer catcherr.RecoverState(`user.GeneratePasswordHash`)

	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	catcherr.HandleError(w, catcherr.InternalServerError, err)

	return string(hashBytes)
}

func CompareLoginCredentials(w http.ResponseWriter, login, password string) {
	defer catcherr.RecoverState(`user.CompareLoginData`)
	hash := database.GetHashedPassword(w, login)

	hashBytes := []byte(hash)
	passwordBytes := []byte(password)

	err := bcrypt.CompareHashAndPassword(hashBytes, passwordBytes)
	catcherr.HandleError(w, catcherr.Unathorized, err)
}

func GetFiles(w http.ResponseWriter, dir string) (files []string) {
	defer catcherr.RecoverState(`user.GetFiles`)

	dirEntry, err := os.ReadDir(dir)
	catcherr.HandleError(w, catcherr.InternalServerError, err)

	for _, f := range dirEntry {
		files = append(files, f.Name())
	}
	return files
}
