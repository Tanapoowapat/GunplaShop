package database

import (
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/Tanapoowapat/GunplaShop/config"
	"github.com/jmoiron/sqlx"
)

func DbConnect(cfg config.IDbConfig) *sqlx.DB {
	db, err := sqlx.Connect("pgx", cfg.Url())
	if err != nil {
		log.Fatalf("Error Fail to connect to database %v", err)
	}
	db.DB.SetMaxOpenConns(cfg.MaxConnection())

	return db
}
