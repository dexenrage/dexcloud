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
	"net/http"
	"os"
	"server/catcherr"
	"server/database"
	"server/user"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte(`QHhpZGlvCg==`) // UNSAFE

func createToken(login string) (data tokenData, err error) {
	expirationTime := jwt.NewNumericDate(time.Now().Add(15 * time.Minute))

	claims := &jwtClaims{
		Login: login,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: expirationTime,
		},
	}
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := tkn.SignedString(jwtKey)
	expires := expirationTime.UTC().Format(http.TimeFormat)

	return tokenData{
		Login:   login,
		Token:   token,
		Expires: expires,
	}, err

}

func getCookie(r *http.Request, name string) (value string, err error) {
	c, err := r.Cookie(name)
	return c.Value, err
}

func parseToken(r *http.Request, loginCookie string) (err error) {
	defer func() { err = catcherr.RecoverAndReturnError() }()

	tokenCookie, err := getCookie(r, `token`)
	catcherr.HandleError(err)

	var (
		claims  = &jwtClaims{}
		keyfunc = func(tkn *jwt.Token) (interface{}, error) { return jwtKey, nil }
	)

	token, err := jwt.ParseWithClaims(tokenCookie, claims, keyfunc)
	catcherr.HandleError(err)

	switch {
	case claims.Login != loginCookie:
		catcherr.HandleError(jwt.ErrSignatureInvalid)
	case !token.Valid:
		catcherr.HandleError(jwt.ErrSignatureInvalid)
	}
	return err
}

func checkAuth(r *http.Request) (err error) {
	defer func() { err = catcherr.RecoverAndReturnError() }()

	login, err := getCookie(r, `login`)
	catcherr.HandleError(err)

	return parseToken(r, login)
}

func registerUser(ctx context.Context, acc database.User) (err error) {
	defer func() { err = catcherr.RecoverAndReturnError() }()

	acc.Password, err = user.GeneratePasswordHash(acc.Password)
	catcherr.HandleError(err)

	acc, err = database.RegisterUser(ctx, acc)
	catcherr.HandleError(err)

	var (
		userID = fmt.Sprint(acc.ID)
		dir    = getUserDir(userID)
	)

	return os.Mkdir(dir, os.ModePerm)
}

func GetUserID(ctx context.Context, r *http.Request) (userID string, err error) {
	defer func() { err = catcherr.RecoverAndReturnError() }()

	login, err := getCookie(r, `login`)
	catcherr.HandleError(err)

	err = parseToken(r, login)
	catcherr.HandleError(err)

	user, err := database.GetUserInfo(ctx, login)
	catcherr.HandleError(err)

	return fmt.Sprint(user.ID), err
}
