package main

import (
	dbConfig "echoapi/pkg/db/config"
	"echoapi/pkg/db/models"
	echoConfig "echoapi/pkg/echo/config"
	"echoapi/pkg/echo/context"
	"echoapi/pkg/prometheus"
	"echoapi/pkg/resources"
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

func registerRoutes(e *echo.Echo) {
	pei := prometheus.NewPrometheusInstrumentation()
	e.Use(pei.PrometheusStatsMiddleware)
	e.GET("/metrics", pei.MetricsEndpoint)

	userResource := resources.UserResource{}

	e.GET("/user", userResource.Get)
	e.POST("/user", userResource.Post)

	todoResource := resources.ToDoResource{}

	e.GET("/todo", todoResource.Get)
	e.POST("/todo", todoResource.Post)
}

func main() {
	config, err := echoConfig.LoadConfig()
	if err != nil {
		panic(err)
	}

	db, err := dbConfig.ConnectToDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_ = db.LogMode(true)
	db.SingularTable(true)
	_ = db.AutoMigrate(&models.Todo{})

	e := echo.New()
	e.Use(context.CustomContextMiddleware(db))

	registerRoutes(e)

	if l, ok := e.Logger.(*log.Logger); ok {
		l.SetHeader("${time_rfc3339} ${level}")
		l.SetLevel(log.DEBUG)
	}

	echoConfigString := fmt.Sprintf("%s:%s", config.ServerHost, config.ServerPort)
	err = e.Start(echoConfigString)
	if err != nil {
		e.Logger.Fatal(err)
	}
}
