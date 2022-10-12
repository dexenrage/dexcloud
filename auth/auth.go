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

package auth

import (
	"context"
	"net/http"
	"server/config"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type (
	jwtClaims struct {
		jwt.RegisteredClaims
		Login string `jsons:"login"`
	}

	Token struct {
		Login   string `json:"login"`
		Token   string `json:"token"`
		Expires string `json:"expires"`
	}
)

/*
 * ToDo: Refresh token
 */

func CreateToken(ctx context.Context, login string) (t Token, err error) {
	expirationTime := jwt.NewNumericDate(time.Now().Add(15 * time.Minute))
	claims := &jwtClaims{
		Login: login,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: expirationTime,
		},
	}
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := tkn.SignedString(config.Bytes(config.JWTKey))
	if err != nil {
		return Token{}, err
	}

	expires := expirationTime.UTC().Format(http.TimeFormat)
	t = Token{
		Login:   login,
		Token:   token,
		Expires: expires,
	}
	return t, nil
}

func GetLoginFromCookie(r *http.Request) (string, error) {
	lc, err := r.Cookie(`login`)
	if err != nil {
		return ``, err
	}
	return lc.Value, nil
}

func VerifyUser(r *http.Request, login string) error {
	tokenCookie, err := r.Cookie(`token`)
	if err != nil {
		return err
	}

	var (
		claims  = &jwtClaims{}
		key     = config.Bytes(config.JWTKey)
		keyfunc = func(tkn *jwt.Token) (any, error) { return key, nil }
	)

	token, err := jwt.ParseWithClaims(tokenCookie.Value, claims, keyfunc)
	if err != nil {
		return err
	}

	switch {
	case claims.Login != login:
		return jwt.ErrSignatureInvalid
	case !token.Valid:
		return jwt.ErrSignatureInvalid
	}
	return nil
}
