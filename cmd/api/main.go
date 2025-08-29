package main

import (
	"log"

	"github.com/aryadisastra/main-app-be/internal/config"
	"github.com/aryadisastra/main-app-be/internal/db"
	"github.com/aryadisastra/main-app-be/internal/router"
	_ "github.com/aryadisastra/main-app-be/internal/swagger"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	cfg := config.Load()
	gdb := db.Open(cfg.DBDsn)

	r := router.New(gdb, cfg.JWTSecret)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Printf("logistics service listening on :%s", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
