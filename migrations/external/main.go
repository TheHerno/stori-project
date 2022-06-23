package main

import (
	"crypto/tls"
	"log"
	"os"
	"stori-service/src/libs/env"

	"github.com/go-pg/pg/v9"
	migrations "github.com/robinjoseph08/go-pg-migrations/v2"
)

const directory = "migrations/external"

func main() {
	dbHost := env.StoriServicePostgresqlHost
	dbPort := env.StoriServicePostgresqlPort
	dbName := env.StoriServicePostgresqlName
	if env.AppEnv == "testing" {
		dbName = env.StoriServicePostgresqlNameTest
	}
	dbUser := env.StoriServicePostgresqlUsername
	dbPassword := env.StoriServicePostgresqlPassword
	dbSSLMode := env.StoriServicePostgresqlSSLMode

	options := &pg.Options{
		Addr:     dbHost + ":" + dbPort,
		User:     dbUser,
		Database: dbName,
		Password: dbPassword,
	}
	if dbSSLMode != "disable" {
		options.TLSConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	db := pg.Connect(options)

	err := migrations.Run(db, directory, os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}
