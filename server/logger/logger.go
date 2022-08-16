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

package logger

import (
	"errors"
	"fmt"

	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

func init() {
	lg, err := zap.NewDevelopment()
	if err != nil {
		msg := `Failed to initialize Zap logger: %v`
		msg = fmt.Sprintf(msg, err)
		err = errors.New(msg)
		panic(err)
	}
	logger = lg.Sugar()
	defer Sync()
}

func Sync()             { logger.Sync() }
func Errorln(err error) { logger.Errorln(err) }
func Panicln(err error) { logger.Panicln(err) }

func Fatalln(err error) {
	logger.Sync()
	logger.Fatalln(err)
}
