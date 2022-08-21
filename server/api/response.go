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

package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"server/catcherr"
)

type CustomResponse struct {
	Status int
	Data   interface{}
}

const (
	defaultResponseType  = "Content-Type"
	defaultResponseValue = "application/json"
)

var response CustomResponse

func (resp *CustomResponse) Send(w http.ResponseWriter, status int, data interface{}) {
	defer catcherr.RecoverState(`api.response.Send`)
	w.Header().Set(defaultResponseType, defaultResponseValue)

	resp.Status = status
	resp.Data = data

	if _, ok := resp.Data.(map[string]interface{}); !ok {
		if _, ok := resp.Data.(map[string]string); !ok {
			resp.Data = map[string]string{
				"data": fmt.Sprint(resp.Data),
			}
		}
	}

	jsonData, err := json.Marshal(resp.Data)
	catcherr.HandleError(w, catcherr.InternalServerError, err)

	_, err = w.Write(jsonData)
	catcherr.HandleError(w, catcherr.InternalServerError, err)
}
