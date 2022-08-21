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
	"errors"
	"net/http"
	"server/catcherr"
	"server/database"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte(`QHhpZGlvCg==`) // UNSAFE

type Claims struct {
	jwt.RegisteredClaims
	Login string `json:"login"`
}

func createToken(w http.ResponseWriter, login string) (token, expires string) {
	defer catcherr.RecoverState(`api.createToken`)
	expirationTime := jwt.NewNumericDate(time.Now().Add(15 * time.Minute))

	claims := &Claims{
		Login: login,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: expirationTime,
		},
	}
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := tkn.SignedString(jwtKey)
	catcherr.HandleError(w, catcherr.InternalServerError, err)

	expires = expirationTime.UTC().Format(http.TimeFormat)
	return token, expires
}

func getCookie(w http.ResponseWriter, r *http.Request, name string) (cookie string) {
	c, err := r.Cookie(name)
	if errors.Is(err, http.ErrNoCookie) {
		catcherr.HandleError(w, catcherr.Unathorized, err)
	}
	catcherr.HandleError(w, catcherr.BadRequest, err)
	return c.Value
}

func parseToken(w http.ResponseWriter, r *http.Request) {
	tokenCookie := getCookie(w, r, `token`)
	var (
		claims  = &Claims{}
		keyfunc = func(tkn *jwt.Token) (interface{}, error) { return jwtKey, nil }
	)

	token, err := jwt.ParseWithClaims(tokenCookie, claims, keyfunc)
	if errors.Is(err, jwt.ErrSignatureInvalid) {
		catcherr.HandleError(w, catcherr.Unathorized, err)
	}
	catcherr.HandleError(w, catcherr.BadRequest, err)

	if login := getCookie(w, r, `login`); claims.Login != login {
		err := errors.New(`Invalid login`)
		catcherr.HandleError(w, catcherr.Unathorized, err)
	}

	if !token.Valid {
		err := errors.New(`Invalid token`)
		catcherr.HandleError(w, catcherr.Unathorized, err)
	}
}

func GetUserID(w http.ResponseWriter, r *http.Request) (userID string) {
	defer catcherr.RecoverState(`api.GetUserID`)
	parseToken(w, r)
	login := getCookie(w, r, `login`)
	return database.GetUserID(w, login)
}
