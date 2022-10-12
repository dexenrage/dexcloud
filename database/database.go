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
		host     = config.String(config.DBHost)
		username = config.String(config.DBUser)
		password = config.String(config.DBPassword)
		sslmode  = config.String(config.DBSSLMode)
		dbName   = config.String(config.DBName)

		user = url.UserPassword(username, password)
	)

	const scheme = `postgres`
	dsn := url.URL{
		Scheme: scheme,
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

	_, err = db.NewCreateTable().Model((*File)(nil)).IfNotExists().Exec(ctx)
	catcherr.HandleError(err)
}

func RegisterUser(ctx context.Context, u User) (user User, err error) {
	defer func() { err = catcherr.RecoverAndReturnError() }()
	_, err = db.NewInsert().Model(&u).Exec(ctx)
	if err != nil {
		return User{}, err
	}
	return GetUser(ctx, u.Login)
}

func GetUser(ctx context.Context, login string) (user User, err error) {
	u := new(User)
	err = db.NewSelect().Model(u).Where(`login = ?`, login).Scan(ctx)
	return *u, err
}

func GetFileList(ctx context.Context, login string) (files []File, err error) {
	u, err := GetUser(ctx, login)
	if err != nil {
		return nil, err
	}

	err = db.NewSelect().Model(&files).Where(`f.uid = ?`, u.ID).Scan(ctx)
	return files, err
}

func SaveFileInfo(ctx context.Context, login, filename, checksum string) error {
	u, err := GetUser(ctx, login)
	if err != nil {
		return err
	}

	f := new(File)
	f.UserID = u.ID
	f.Name = filename
	f.Checksum = checksum

	_, err = db.NewInsert().Model(f).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func RemoveFileInfo(ctx context.Context, login, checksum string) error {
	u, err := GetUser(ctx, login)
	if err != nil {
		return err
	}

	f := new(File)
	_, err = db.NewDelete().Model(f).Where(`uid = ?`, u.ID).
		Where(`checksum = ?`, checksum).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
