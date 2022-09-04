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
	IndexHTTP, Slash = `/`, `/`
	RegisterHTTP     = `/register`
	LoginHTTP        = `/login`
	ProfileHTTP      = `/profile`
	StaticHTTP       = `/static/`
	UploadsHTTP      = `/uploads/`

	ApiCheckAuthHTTP = `/api/checkauth`
	ApiRegisterHTTP  = `/api/register`
	ApiLoginHTTP     = `/api/login`
	ApiUploadHTTP    = `/api/upload`
	ApiFileListHTTP  = `/api/filelist`

	StaticFilesRoot = "web"
	UserUploadsRoot = "userdata"
)

func CleanPath(elem ...string) string {
	path := filepath.Join(elem...)
	return filepath.Clean(path)
}

func IndexPage() string { return CleanPath(StaticFilesRoot, "index.html") }

func RegisterPage() string { return CleanPath(StaticFilesRoot, "register.html") }

func LoginPage() string { return CleanPath(StaticFilesRoot, "login.html") }

func ProfilePage() string { return CleanPath(StaticFilesRoot, "profile.html") }

func StaticFiles() string { return CleanPath(StaticFilesRoot, `static`) }

func UserUploads() string { return CleanPath(UserUploadsRoot, `uploads`) }

func CreateCriticalDirectories() (err error) {
	defer func() { err = catcherr.RecoverAndReturnError() }()

	directories := [4]string{
		StaticFilesRoot,
		StaticFiles(),
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
