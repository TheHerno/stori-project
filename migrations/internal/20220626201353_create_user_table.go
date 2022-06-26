package main

import (
	"github.com/go-pg/pg/v9/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v2"
)

func init() {
	up := func(db orm.DB) error {
		_, err := db.Exec(`
		CREATE TABLE customer (
			customer_id serial PRIMARY KEY,
			name varchar(100) NOT NULL,
			email varchar(100) NOT NULL,
			created_at timestamp with time zone NOT NULL DEFAULT NOW(),
			updated_at timestamp with time zone NOT NULL DEFAULT NOW(),
			deleted_at timestamp with time zone
		)
		`)
		return err
	}

	down := func(db orm.DB) error {
		_, err := db.Exec(`
		DROP TABLE customer
		`)
		return err
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20220626201353_create_customer_table", up, down, opts)
}
