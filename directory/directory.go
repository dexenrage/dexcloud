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
	APIFileServer = `/api/file/`

	APIAuthCheck = `/api/auth/check`
	APIRegister  = `/api/auth/register`
	APILogin     = `/api/auth/login`

	APIFileUpload = `/api/file/upload`
	APIFileList   = `/api/file/list`
	APIFileDelete = `/api/file/delete`

	userDataFolder = `userdata`
)

func init() {
	err := os.Mkdir(userDataFolder, os.ModePerm)
	if errors.Is(err, fs.ErrExist) {
		return
	}
	catcherr.HandleError(err)
}

func UserData() string { return CleanPath(userDataFolder) }

func CleanPath(elem ...string) string {
	return filepath.Clean(filepath.Join(elem...))
}
