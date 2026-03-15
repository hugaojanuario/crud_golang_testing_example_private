package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/hugaojanuario/crud_golang_testing_example_private/internal/config"
	"github.com/hugaojanuario/crud_golang_testing_example_private/internal/database/postgre"
	"github.com/hugaojanuario/crud_golang_testing_example_private/internal/handler"
	"github.com/hugaojanuario/crud_golang_testing_example_private/internal/repository"
	"github.com/hugaojanuario/crud_golang_testing_example_private/internal/service"
)

func main() {
	fmt.Println("iniciando a api . . .")

	cfg := config.NewConfig()

	db, err := postgre.NewConn(postgre.Config{
		DbHost:     cfg.DbHost,
		DBPort:     cfg.DBPort,
		DbUser:     cfg.DbUser,
		DbPassword: cfg.DbPassword,
		DbName:     cfg.DbName,
		SslMode:    cfg.SslMode,
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	repo := repository.NewRepostory(db)
	service := service.NewService(repo)
	h := handler.NewHandler(service)

	r := gin.Default()

	handler.NewConnection(r, h)

	log.Printf("servidor rodando na porta %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf(err.Error())
	}
}
