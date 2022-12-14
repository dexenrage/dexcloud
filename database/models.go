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

import "github.com/uptrace/bun"

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`
	ID            int64  `bun:"id,pk,autoincrement"`
	Login         string `bun:"login,notnull" json:"login"`
	Password      string `bun:"password,notnull" json:"password"`
}

type File struct {
	bun.BaseModel `bun:"table:files,alias:f"`
	UserID        int64  `bun:"uid,notnull"`
	Name          string `bun:"name,notnull" json:"filename"`
	Checksum      string `bun:"checksum,notnull" json:"checksum"`
}
