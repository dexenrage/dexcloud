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
	"net/url"
	"server/catcherr"
	"server/config"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

var db *bun.DB

func init() {
	ctx := context.Background()

	var (
		host     = config.String(`db_host`)
		username = config.String(`db_user`)
		password = config.String(`db_pass`)
		sslmode  = config.String(`db_sslmode`)
		dbName   = config.String(`db_name`)

		user = url.UserPassword(username, password)
	)

	dsn := url.URL{
		Scheme: `postgres`,
		Host:   host,
		User:   user,
		Path:   dbName,
	}
	{
		q := dsn.Query()
		q.Add(`sslmode`, sslmode)
		dsn.RawQuery = q.Encode()
	}

	// Open a PostgreSQL database.
	pgdb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn.String())))

	// Create a Bun db on top of it.
	db = bun.NewDB(pgdb, pgdialect.New())

	// Print all queries to stdout.
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	_, err := db.NewCreateTable().Model((*User)(nil)).IfNotExists().Exec(ctx)
	catcherr.HandleError(err)
}

func RegisterUser(ctx context.Context, u User) (user User, err error) {
	defer func() { err = catcherr.RecoverAndReturnError() }()

	_, err = db.NewInsert().Model(&u).Exec(ctx)
	catcherr.HandleError(err)

	return GetUserInfo(ctx, u.Login)
}

func GetUserInfo(ctx context.Context, login string) (user User, err error) {
	u := new(User)
	err = db.NewSelect().Model(u).Where(`login = ?`, login).Scan(ctx)
	return *u, err
}
