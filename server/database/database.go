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

package database

import (
	"context"
	"fmt"
	"math/rand"
	"server/logger"
	"strconv"

	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	conn, err := connect()
	if err != nil {
		logger.Panicln(err)
	}
	defer conn.Close(context.Background())

	resp, err := conn.Query(context.Background(), createTableQuery)
	if err != nil {
		logger.Panicln(err)
	}
	defer resp.Close()
}

func connect() (*pgx.Conn, error) {
	const auth = "user=postgres password=123456 host=localhost port=5432 dbname=users"

	cs, err := pgx.ParseConfig(auth)
	if err != nil {
		return nil, err
	}

	conn, err := pgx.ConnectConfig(context.Background(), cs)
	if err != nil {
		return nil, err
	}
	return conn, err
}

func RegisterUser(login, pass string) error {
	conn, err := connect()
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	pass = string(hash)
	token := strconv.Itoa(rand.Int()) // temporary (!!!)
	query := fmt.Sprintf(regUserQuery, login, pass, token)

	resp, err := conn.Query(context.Background(), query)
	if err != nil {
		return err
	}
	defer resp.Close()
	return err
}
