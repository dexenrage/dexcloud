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
	"database/sql"
	"log"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

var db *bun.DB

func init() {
	ctx := context.Background()

	// Open a PostgreSQL database.
	dsn := "postgres://postgres:123456@localhost:5432/users?sslmode=disable"
	pgdb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	// Create a Bun db on top of it.
	db = bun.NewDB(pgdb, pgdialect.New())

	// Print all queries to stdout.
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	_, err := db.NewCreateTable().Model((*User)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		log.Panicln(err)
	}
}

func RegisterUser(ctx context.Context, u User) (user User, err error) {
	_, err = db.NewInsert().Model(&u).Exec(ctx)
	if err != nil {
		return user, err
	}
	return GetUserInfo(ctx, u.Login)
}

func GetUserInfo(ctx context.Context, login string) (user User, err error) {
	u := new(User)
	err = db.NewSelect().Model(u).Where(`login = ?`, login).Scan(ctx)
	return *u, err
}
