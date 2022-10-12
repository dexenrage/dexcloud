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

package user

import (
	"context"
	"crypto/sha256"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"server/auth"
	"server/database"
	"server/directory"

	"golang.org/x/crypto/bcrypt"
)

func Register(ctx context.Context, u database.User) (token auth.Token, err error) {
	u.Password, err = generatePasswordHash(ctx, u.Password)
	if err != nil {
		return auth.Token{}, err
	}

	u, err = database.RegisterUser(ctx, u)
	if err != nil {
		return auth.Token{}, err
	}

	token, err = auth.CreateToken(ctx, u.Login)
	if err != nil {
		return auth.Token{}, err
	}
	return token, nil
}

func Login(ctx context.Context, u database.User) (token auth.Token, err error) {
	userInfo, err := database.GetUser(ctx, u.Login)
	if err != nil {
		return auth.Token{}, err
	}

	err = comparsePasswords(ctx, u.Password, userInfo.Password)
	if err != nil {
		return auth.Token{}, err
	}

	token, err = auth.CreateToken(ctx, u.Login)
	if err != nil {
		return auth.Token{}, err
	}
	return token, nil
}

func generatePasswordHash(ctx context.Context, password string) (hash string, err error) {
	if ctx.Err() != nil {
		return ``, err
	}
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ``, err
	}
	hash = string(b)
	return hash, nil
}

func comparsePasswords(ctx context.Context, password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

type FileData struct {
	File       multipart.File
	FileHeader *multipart.FileHeader
}

func SaveFile(ctx context.Context, f *multipart.FileHeader) (sha256sum string, err error) {
	if ctx.Err() != nil {
		return ``, ctx.Err()
	}

	file, err := f.Open()
	if err != nil {
		return ``, err
	}
	defer file.Close()

	checksum := sha256.New()
	_, err = io.Copy(checksum, file)
	if err != nil {
		return ``, err
	}

	sha256sum = string(checksum.Sum(nil))
	return sha256sum, nil
}

func RemoveFile(ctx context.Context, checksum string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	path := filepath.Join(directory.UserData(), checksum)
	if err := os.Remove(path); err != nil {
		return err
	}
	return nil
}
