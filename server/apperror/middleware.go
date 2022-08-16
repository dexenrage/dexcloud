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
					w.WriteHeader(http.StatusBadRequest)
					w.Write(ErrBadRequest.Marshal())
					return
				case errors.Is(err, ErrNotFound):
					w.WriteHeader(http.StatusNotFound)
					w.Write(ErrNotFound.Marshal())
					return
				case errors.Is(err, ErrUnathorized):
					w.WriteHeader(http.StatusUnauthorized)
					w.Write(ErrUnathorized.Marshal())
					return
				case errors.Is(err, ErrInternalServerError):
					w.WriteHeader(http.StatusInternalServerError)
					w.Write(ErrInternalServerError.Marshal())
					return
				default:
					err := err.(*CustomError)
					w.WriteHeader(http.StatusBadRequest)
					w.Write(err.Marshal())
					return
				}
			}
			w.WriteHeader(http.StatusTeapot)
			w.Write(systemError(err).Marshal())
		}
	})
}
