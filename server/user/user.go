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
	"os"
	"server/apperror"
	"server/database"

	"golang.org/x/crypto/bcrypt"
)

func CompareLoginData(login, password string) (err error) {
	dbLogin, dbPassword, err := database.CheckUserData(login, password)
	if err != nil {
		return apperror.ErrInternalServerError
	}

	if login != dbLogin {
		return apperror.ErrUnathorized
	}

	dbPasswordBytes := []byte(dbPassword)
	passwordBytes := []byte(password)

	err = bcrypt.CompareHashAndPassword(dbPasswordBytes, passwordBytes)
	if err != nil {
		return apperror.ErrUnathorized
	}
	return err
}

func GetFiles(dir string) ([]string, error) {
	var files []string

	dirEntry, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, f := range dirEntry {
		files = append(files, f.Name())
	}
	return files, err
}
