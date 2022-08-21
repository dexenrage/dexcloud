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

package catcherr

import "net/http"

type (
	Error struct {
		Error interface{} `json:"error"`
	}

	Response struct {
		StatusCode int         `json:"status"`
		Data       interface{} `json:"data"`
	}

	CustomError struct {
		StatusCode  int
		Description string
	}
)

var (
	BadRequest          CustomError
	Unathorized         CustomError
	Forbidden           CustomError
	NotFound            CustomError
	InternalServerError CustomError
)

func (e *CustomError) BadRequest() {
	e.StatusCode = http.StatusBadRequest
	e.Description = `Bad Request`
}

func (e *CustomError) Unathorized() {
	e.StatusCode = http.StatusUnauthorized
	e.Description = `Unathorized`
}

func (e *CustomError) Forbidden() {
	e.StatusCode = http.StatusForbidden
	e.Description = `Forbidden`
}

func (e *CustomError) NotFound() {
	e.StatusCode = http.StatusNotFound
	e.Description = `Not Found`
}

func (e *CustomError) InternalServerError() {
	e.StatusCode = http.StatusInternalServerError
	e.Description = `Internal Server Error`
}
