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

package config

import (
	"server/catcherr"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
)

const (
	JWTKey = `jwt_key`
	Host   = `host`

	DBHost     = `db_host`
	DBUser     = `db_user`
	DBPassword = `db_pass`
	DBSSLMode  = `db_sslmode`
	DBName     = `db_name`
)

var cfg *koanf.Koanf

func init() {
	cfg = koanf.New(`.`)
	err := cfg.Load(file.Provider(`config.yml`), yaml.Parser())
	catcherr.HandleError(err)
}

func String(path string) string { return cfg.String(path) }
func Bytes(path string) []byte  { return cfg.Bytes(path) }
