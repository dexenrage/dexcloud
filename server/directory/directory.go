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

import "path/filepath"

const (
	IndexHTTP, Slash = `/`, `/`
	ProfileHTTP      = `/profile`
	StaticHTTP       = `/static/`
	UploadsHTTP      = `/uploads/`

	ApiRegisterHTTP = `/api/register`
	ApiLoginHTTP    = `/api/login`
	ApiUploadHTTP   = `/api/upload`
	ApiFileListHTTP = `/api/filelist`
)

func IndexPage() string { return filepath.Join("web", "index.html") }

func ProfilePage() string { return filepath.Join("web", "profile.html") }

func StaticFiles() string { return filepath.Join(`web`, `static`) }

func UserUploads() string { return filepath.Join(`web`, `uploads`) }
