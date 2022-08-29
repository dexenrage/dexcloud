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
	"context"
	"fmt"
	"os"
	"server/catcherr"
	"server/database"
	"server/user"
	"sync"
)

func waitGorutines(wg *sync.WaitGroup) <-chan bool {
	wgDoneChan := make(chan bool)
	go func() {
		wg.Wait()
		close(wgDoneChan)
	}()
	return wgDoneChan
}

func goRegisterUser(
	ctx context.Context,
	wg *sync.WaitGroup,
	acc database.User,
	errChan chan<- catcherr.ErrorChan,
) {
	defer wg.Done()
	var err error

	acc.Password, err = user.GeneratePasswordHash(acc.Password)
	catcherr.HandleErrorChannel(errChan, catcherr.InternalServerError, err)

	acc, err = database.RegisterUser(ctx, acc)
	catcherr.HandleErrorChannel(errChan, catcherr.InternalServerError, err)

	userID := fmt.Sprint(acc.ID)
	err = os.Mkdir(getUserDir(userID), os.ModePerm)
	catcherr.HandleErrorChannel(errChan, catcherr.InternalServerError, err)
}

func goGetTokenData(
	wg *sync.WaitGroup,
	login string,
	tokenChan chan<- tokenData,
	errChan chan<- catcherr.ErrorChan,
) {
	defer wg.Done()
	data, err := createToken(login)
	catcherr.HandleErrorChannel(errChan, catcherr.Unathorized, err)
	tokenChan <- data
}
