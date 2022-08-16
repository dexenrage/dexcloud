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

package main

import "path/filepath"

const (
	staticHTTP  = `/static`
	uploadsHTTP = `/uploads`
	profileHTTP = `/profile`
)

func indexPagePath() string {
	return filepath.Join("web", "index.html")
}

func profilePagePath() string {
	return filepath.Join("web", "profile.html")
}

func staticFiles() string {
	return filepath.Join(`web`, `static`)
}

func uploadsFiles() string {
	return filepath.Clean(`./uploads`)
}
