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
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"server/catcherr"

	"golang.org/x/crypto/bcrypt"
)

func GeneratePasswordHash(password string) (hash string, err error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b), err
}

func ComparsePasswords(ctx context.Context, password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func SaveUploadedFiles(files []FileStruct) (err error) {
	defer func() { err = catcherr.RecoverAndReturnError() }()
	for _, f := range files {
		path := filepath.Join(f.Directory, f.FileHeader.Filename)

		newFile, err := os.Create(path)
		catcherr.HandleError(err)
		defer newFile.Close()

		_, err = io.Copy(newFile, f.File)
		catcherr.HandleError(err)
	}
	return err
}

func GetFiles(dir string) (files []string, err error) {
	defer func() { err = catcherr.RecoverAndReturnError() }()

	entry, err := os.ReadDir(dir)
	if errors.Is(err, fs.ErrNotExist) {
		err = os.Mkdir(dir, os.ModePerm)
	}
	catcherr.HandleError(err)

	for _, f := range entry {
		files = append(files, f.Name())
	}
	return files, err
}
