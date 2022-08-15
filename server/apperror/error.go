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

package apperror

import (
	"encoding/json"
)

var (
	ErrBadRequest          = NewAppError(nil, "Bad Request", "", "US-000400")
	ErrNotFound            = NewAppError(nil, "Not Found", "", "US-000404")
	ErrUnathorized         = NewAppError(nil, "Unathorized", "", "US-000401")
	ErrInternalServerError = NewAppError(nil, "Internal Server Error", "", "US-000500")
)

type AppError struct {
	Err        error  `json:"-"`
	Message    string `json:"msg,omitempty"`
	DevMessage string `json:"devmsg,omitempty"`
	Code       string `json:"code,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) Marshal() []byte {
	marshal, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return marshal
}

func NewAppError(err error, message, devMessage, code string) *AppError {
	return &AppError{
		Err:        err,
		Message:    message,
		DevMessage: devMessage,
		Code:       code,
	}
}
