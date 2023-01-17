package main

import (
	"GoBBS/config"
	"GoBBS/domain/service"
	"GoBBS/interface/handler"
	"GoBBS/usecase"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	env := config.GetEnv()
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?parseTime=true",
		env.DBUser,
		env.DBPassword,
		env.DBHost,
		env.DBName,
	)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	usecase := usecase.NewUserUseCase(db, service.NewUserServiceFactory())
	handler.NewUserHandler(usecase).RegistUserHandlerFunc()

	log.Fatal(http.ListenAndServe(":8100", nil))
}
