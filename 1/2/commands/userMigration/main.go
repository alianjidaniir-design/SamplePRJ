package main

import (
	"flag"
	"fmt"
	"log"

	mysqlDataSource "Fiber/API/2/models/student/datasourse/mySqlDS"
)

func main() {
	envCfg, err := mysqlDataSource.LoadConfigFromEnv()
	if err != nil {
		log.Fatalf("[user-migration] load env failed: %v", err)
	}

	dsn := flag.String("dsn", envCfg.DSN, "MySQL DSN")
	table := flag.String("table", envCfg.StudentTableName, "students table name")
	flag.Parse()

	if *dsn == "" {
		log.Fatal("[user-migration] missing DSN: set MYSQL_DSN or pass --dsn")
	}

	cfg := envCfg
	cfg.DSN = *dsn
	cfg.StudentTableName = *table

	if err := mysqlDataSource.ValidateTableName(cfg.StudentTableName); err != nil {
		log.Fatalf("[user-migration] invalid table name: %v", err)
	}

	db, err := mysqlDataSource.Open(cfg)
	if err != nil {
		log.Fatalf("[user-migration] connect mysql failed: %v", err)
	}
	defer db.Close()

	if err := mysqlDataSource.EnsureTaskTable(db, cfg.StudentTableName); err != nil {
		log.Fatalf("[user-migration] create table failed: %v", err)
	}

	fmt.Printf("[user-migration] table is ready: %s\n", cfg.StudentTableName)
}
