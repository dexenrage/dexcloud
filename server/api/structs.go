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

import "github.com/golang-jwt/jwt/v4"

type (
	account struct {
		UserID   string `json:"userid"`
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	fileListStruct struct {
		UserID string   `json:"userid"`
		Files  []string `json:"files,omitempty"`
	}

	jwtClaims struct {
		jwt.RegisteredClaims
		Login string `json:"login"`
	}

	tokenStruct struct {
		Login   string `json:"login"`
		Token   string `json:"token"`
		Expires string `json:"expires"`
	}

	responseData struct {
		StatusCode int         `json:"status"`
		Data       interface{} `json:"data"`
	}
)
