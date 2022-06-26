package main

import (
	"github.com/go-pg/pg/v9/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v2"
)

func init() {
	up := func(db orm.DB) error {
		/*
			I seed the user table here because user managment is not on the code challenge.
		*/
		_, err := db.Exec(`
			INSERT INTO customer (customer_id, name, email) VALUES
			(1, 'Pepe Perez', 'hernanlistort.i@gmail.com'),
			(2, 'Juan Perez', 'hernanlistor.ti@gmail.com')
		`)
		return err
	}

	down := func(db orm.DB) error {
		_, err := db.Exec(`
			DELETE FROM customer WHERE customer_id IN (1, 2)
		`)
		return err
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20220626201946_seed_table_customer", up, down, opts)
}
