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
	"errors"
	"net/http"
)

type appHandler func(w http.ResponseWriter, r *http.Request) error

func Middleware(h appHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var appErr *CustomError
		if err := h(w, r); err != nil {
			if errors.As(err, &appErr) {
				switch {
				case errors.Is(err, ErrBadRequest):
					response(w, http.StatusBadRequest, ErrBadRequest.Marshal())
					return
				case errors.Is(err, ErrForbidden):
					response(w, http.StatusForbidden, ErrForbidden.Marshal())
					return
				case errors.Is(err, ErrNotFound):
					response(w, http.StatusNotFound, ErrNotFound.Marshal())
					return
				case errors.Is(err, ErrUnathorized):
					response(w, http.StatusUnauthorized, ErrUnathorized.Marshal())
					return
				case errors.Is(err, ErrInternalServerError):
					response(w, http.StatusInternalServerError, ErrInternalServerError.Marshal())
					return
				default:
					err := err.(*CustomError)
					response(w, http.StatusBadRequest, err.Marshal())
					return
				}
			}
			responseTeapot(w, err)
		}
	})
}

func response(w http.ResponseWriter, statusCode int, data []byte) {
	switch {
	case statusCode == 0:
		err := errors.New("statusCode is not set or is equal to zero")
		responseTeapot(w, err)
		return
	case data == nil:
		err := errors.New("response data is nil")
		responseTeapot(w, err)
		return
	}
	w.WriteHeader(statusCode)
	w.Write(data)
}

func responseTeapot(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusTeapot)
	w.Write(systemError(err).Marshal())
}
