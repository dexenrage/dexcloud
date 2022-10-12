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
	"errors"
	"fmt"
	"log"
	"net/http"
)

func init() {
	Unathorized.Unathorized()
	Forbidden.Forbidden()
	InternalServerError.InternalServerError()
}

func HandleError(err any) {
	if err != nil {
		panic(err)
	}
}

func HandleAndResponse(w http.ResponseWriter, c CustomError, err any) {
	if err != nil {
		sendErrorData(w, c)
		panic(err)
	}
}

func RecoverAndReturnError() (err error) {
	if r := recover(); r != nil {
		return errors.New(fmt.Sprint(r))
	}
	return nil
}

func Recover(sender string) {
	if r := recover(); r != nil {
		logError(sender, r)
	}
}

func logError(sender string, data any) {
	const template = `[ Sender: %s ]: %v `
	log.Printf(template, sender, data)
}

func sendErrorData(w http.ResponseWriter, c CustomError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(c.StatusCode)

	err := json.NewEncoder(w).Encode(c)
	if err != nil {
		const sender = `catcherr.sendErrorData()`
		const template = `Can't send data to user. Reason: %v`
		logError(sender, fmt.Sprintf(template, err))
	}
}
