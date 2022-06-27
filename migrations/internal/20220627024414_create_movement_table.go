package main

import (
	"github.com/go-pg/pg/v9/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v2"
)

func init() {
	up := func(db orm.DB) error {
		_, err := db.Exec(`
		CREATE TABLE movement (
			movement_id serial PRIMARY KEY,
			customer_id int NOT NULL,
			quantity int NOT NULL,
			available int NOT NULL,
			type int NOT NULL,
			date timestamp with time zone NOT NULL,
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

	migrations.Register("20220627024414_create_movement_table", up, down, opts)
}
