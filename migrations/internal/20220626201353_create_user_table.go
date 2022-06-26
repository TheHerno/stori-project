package main

import (
	"github.com/go-pg/pg/v9/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v2"
)

func init() {
	up := func(db orm.DB) error {
		_, err := db.Exec(`
		CREATE TABLE user (
			UserID serial PRIMARY KEY,
			Name varchar(100) NOT NULL,
			Email varchar(100) NOT NULL,
		)
		`)
		return err
	}

	down := func(db orm.DB) error {
		_, err := db.Exec(`
		DROP TABLE user
		`)
		return err
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20220626201353_create_user_table", up, down, opts)
}
