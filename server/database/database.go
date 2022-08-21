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
	"net/http"
	"server/catcherr"
	"server/logger"

	"github.com/jackc/pgx/v4"
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

func RegisterUser(w http.ResponseWriter, login, hashedPassword string) (userID string) {
	defer catcherr.RecoverState(`database.RegisterUser`)

	conn, err := connect()
	catcherr.HandleError(w, catcherr.InternalServerError, err)
	defer conn.Close(context.Background())

	query := fmt.Sprintf(regUserQuery, login, hashedPassword)

	resp, err := conn.Query(context.Background(), query)
	catcherr.HandleError(w, catcherr.InternalServerError, err)
	defer resp.Close()

	userID = GetUserID(w, login)
	catcherr.HandleError(w, catcherr.InternalServerError, err)
	return userID
}

func GetHashedPassword(w http.ResponseWriter, login string) (hash string) {
	defer catcherr.RecoverState(`database.GetUserPasswordHash`)

	conn, err := connect()
	catcherr.HandleError(w, catcherr.InternalServerError, err)
	defer conn.Close(context.Background())

	query := fmt.Sprintf(getHashedPasswordQuery, login)

	rows, err := conn.Query(context.Background(), query)
	catcherr.HandleError(w, catcherr.InternalServerError, err)
	defer rows.Close()

	for rows.Next() {
		values, err := rows.Values()
		catcherr.HandleError(w, catcherr.InternalServerError, err)
		hash = values[0].(string)
	}
	return hash
}

func GetUserID(w http.ResponseWriter, login string) (userID string) {
	defer catcherr.RecoverState(`database.GetUserID`)

	conn, err := connect()
	catcherr.HandleError(w, catcherr.InternalServerError, err)
	defer conn.Close(context.Background())

	query := fmt.Sprintf(getUserIDQuery, login)

	rows, err := conn.Query(context.Background(), query)
	catcherr.HandleError(w, catcherr.InternalServerError, err)
	defer rows.Close()

	for rows.Next() {
		values, err := rows.Values()
		catcherr.HandleError(w, catcherr.InternalServerError, err)
		userID = fmt.Sprint(values[0].(int32))
	}
	return userID
}
