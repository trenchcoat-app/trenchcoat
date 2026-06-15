package main

import (
	"flag"
)

var skipMigrations bool

func parseFlags() {
	flag.BoolVar(&skipMigrations, "skip-migrations", false, "skip database migrations on startup")
	flag.BoolVar(&skipMigrations, "M", false, "shorthand for -skip-migrations")

	flag.Parse()
}
