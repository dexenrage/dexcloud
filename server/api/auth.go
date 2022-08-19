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
	"server/apperror"

	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte(`QHhpZGlvCg==`) // UNSAFE

type Claims struct {
	jwt.RegisteredClaims
	Login string `json:"login"`
}

func createToken(acc *Account) (tokenString string, err error) {
	//expirationTime := jwt.NewNumericDate(time.Now().Add(15 * time.Minute))

	claims := &Claims{
		Login:            acc.Login,
		RegisteredClaims: jwt.RegisteredClaims{
			//ExpiresAt: expirationTime,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(jwtKey)
	if err != nil {
		return tokenString, err
	}
	return tokenString, err
}

func checkToken(w http.ResponseWriter, r *http.Request) (err error) {
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			return apperror.ErrUnathorized
		}
		return apperror.ErrBadRequest
	}

	tokenString := c.Value
	claims := &Claims{}

	keyfunc := func(tkn *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	}

	token, err := jwt.ParseWithClaims(tokenString, claims, keyfunc)
	if errors.Is(err, jwt.ErrSignatureInvalid) {
		apperror.ErrUnathorized.Err = jwt.ErrSignatureInvalid
		return apperror.ErrUnathorized
	}
	if err != nil {
		return apperror.ErrUnathorized
	}
	if !token.Valid {
		apperror.ErrUnathorized.Message = "Invalid token"
		return apperror.ErrUnathorized
	}
	return err
}
