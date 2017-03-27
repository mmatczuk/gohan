// Copyright (C) 2015 NTT Innovation Institute, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package migration

import (
	"database/sql"
	"fmt"
	"os"
	"path"

	"github.com/pressly/goose"

	"github.com/cloudwan/gohan/log"
	"github.com/cloudwan/gohan/util"
)

var logger = log.NewLogger()

func LoadConfig(configFile string) (err error) {
	config := util.GetConfig()

	err = config.ReadConfig(configFile)

	if err != nil {
		fmt.Printf("error: failed to load config: %s\n", err.Error())
		return
	}

	err = os.Chdir(path.Dir(configFile))

	if err != nil {
		fmt.Printf("error: chdir() failed: %s\n", err.Error())
		return
	}

	return
}

func readGooseConfig() (dbType, dbConnection, migrationsPath string) {
	config := util.GetConfig()
	dbType = config.GetString("database/type", "sqlite3")
	dbConnection = config.GetString("database/connection", "")
	migrationsPath = config.GetString("database/migrations", "etc/db/migrations")
	return
}

func EnsureVersion() error {
	dbType, dbConnection, migrationsPath := readGooseConfig()

	if err := goose.SetDialect(dbType); err != nil {
		return fmt.Errorf("migration: failed to set goose dialect: %s", err)
	}

	db, err := sql.Open(dbType, dbConnection)

	if err != nil {
		return fmt.Errorf("migration: failed to open db: %s", err)
	}

	v, err := goose.EnsureDBVersion(db)
	if err != nil {
		return fmt.Errorf("migration: failed to ensure db version: %s", err)
	}
	logger.Info("migration path: %q, db version: %d", migrationsPath, v)

	ms, err := goose.CollectMigrations(migrationsPath, 0, goose.MaxVersion)
	if err != nil {
		return fmt.Errorf("migration: failed to list migrations: %s", err)
	}
	m, err := ms.Last()
	if err != nil {
		return fmt.Errorf("migration: failed to get last migration: %s", err)
	}
	if m.Version != v {
		return fmt.Errorf("migration: version mismatch db version=%d, last migration=%d", v, m.Version)
	}

	return nil
}

func Help() {
	fmt.Println("missing subcommand: help, up, up-by-one, up-to, create, down, down-to, redo, status, version")
}

func Run(subcmd string, args []string) {
	dbType, dbConnection, migrationsPath := readGooseConfig()

	if err := goose.SetDialect(dbType); err != nil {
		fmt.Printf("error: failed to set goose dialect: %s\n", err.Error())
		return
	}

	db, err := sql.Open(dbType, dbConnection)

	if err != nil {
		fmt.Printf("error: failed to open db: %s\n", err.Error())
		return
	}

	err = goose.Run(subcmd, db, migrationsPath, args...)
	if err != nil {
		fmt.Println(err)
	}
}
