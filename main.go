package main

import (
	"fmt"
	"os"

	"github.com/Tanapoowapat/GunplaShop/config"
	"github.com/Tanapoowapat/GunplaShop/pkg/database"
)

func envPath() string {
	if len(os.Args) == 1 {
		return ".env"
	} else {
		return os.Args[1]
	}
}

func main() {
	// LoadConfig(envPath())
	cfg := config.LoadConfig(envPath())

	db := database.DbConnect(cfg.Db())
	defer db.Close()

	fmt.Println(db)

}
