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

package directory

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"server/catcherr"
)

const (
	Slash       = `/`
	StaticHTTP  = `/static/`
	UploadsHTTP = `/uploads/`

	ApiCheckAuthHTTP = `/api/auth/check`
	ApiRegisterHTTP  = `/api/auth/register`
	ApiLoginHTTP     = `/api/auth/login`

	ApiUploadFileHTTP = `/api/files/upload`
	ApiFileListHTTP   = `/api/files/list`
	ApiDeleteFileHTTP = `/api/files/delete`

	UserUploadsRoot = "userdata"
)

func CleanPath(elem ...string) string {
	path := filepath.Join(elem...)
	return filepath.Clean(path)
}

func UserUploads() string { return CleanPath(UserUploadsRoot, `uploads`) }

func CreateCriticalDirectories() (err error) {
	defer func() { err = catcherr.RecoverAndReturnError() }()

	directories := [4]string{
		UserUploadsRoot,
		UserUploads(),
	}

	for _, v := range directories {
		err := os.Mkdir(v, os.ModePerm)
		if errors.Is(err, fs.ErrExist) {
			continue
		}
		catcherr.HandleError(err)
	}
	return err
}
