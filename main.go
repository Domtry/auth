package main

import (
	"auth/api"
	"auth/config"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title           CMagic Auth
// @version         1.0.0
// @description     Service d'authentification.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8000
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {

	server := echo.New()

	server.Use(middleware.Recover())
	server.Use(middleware.Logger())
	server.Use(middleware.CORS())

	dbConfig, err := config.LoadDBonfig()
	if err != nil {
		panic("File not found")
	}

	db, err := config.GetDB(dbConfig)
	if err == nil {
		err = config.CreateUpdateTable(db)
		if err != nil {
			return
		}
		err = config.InitSystem(db)
		if err != nil {
			return
		}
	}

	api.GlobalSetup(server, db)

	server.GET("/swagger/*", echoSwagger.WrapHandler)
	err = server.Start(":8000")
	if err != nil {
		return
	}
}
