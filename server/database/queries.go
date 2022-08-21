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

const (
	createTableQuery = `
		CREATE TABLE IF NOT EXISTS users(
			id serial PRIMARY KEY,
			login text,
			password text
		);
	`
	regUserQuery           = `INSERT INTO users(login, password) VALUES ('%s', '%s');`
	getHashedPasswordQuery = `SELECT password FROM users WHERE login='%s';`
	getUserIDQuery         = `SELECT id FROM users WHERE login='%s';`
)
