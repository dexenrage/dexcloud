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

import (
	"encoding/json"
	"fmt"
	"net/http"
	"server/logger"
)

func init() {
	BadRequest.BadRequest()
	Unathorized.Unathorized()
	Forbidden.Forbidden()
	NotFound.NotFound()
	InternalServerError.InternalServerError()
}

func HandleError(w http.ResponseWriter, c CustomError, err error) {
	if err != nil {
		sendErrorData(w, c.StatusCode, c.Description)
		panic(err)
	}
}

func RecoverState(sender string) {
	if r := recover(); r != nil {
		logError(sender, r)
	}
}

func logError(sender string, data interface{}) {
	const tmpl = `[ Sender: %s ]: %v `
	logger.Errorln(fmt.Errorf(tmpl, sender, data))
}

func sendErrorData(w http.ResponseWriter, statusCode int, data interface{}) {
	var resp Response
	resp.StatusCode = statusCode
	resp.Response = data

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		logError(`ResponseWriter`, fmt.Sprint(`Can't send data to user. Reason: `, err))
	}
}
