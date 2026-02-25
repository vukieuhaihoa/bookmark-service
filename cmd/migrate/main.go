package main

import (
	"github.com/vukieuhaihoa/bookmark-service/internal/infrastructure"
)

func main() {
	_ = infrastructure.CreateSQLDBAndMigration()
}
