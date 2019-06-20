package main

import (
	gorm2 "echoapi/contrib/gorm"
	"echoapi/contrib/prometheus"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

type CustomContext struct {
	echo.Context
	DB *gorm.DB
}

func CustomContextMiddleware(db *gorm.DB) func(h echo.HandlerFunc) echo.HandlerFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &CustomContext{
				Context: c,
				DB:      db,
			}
			return h(cc)
		}
	}
}

func main() {
	config, err := loadConfig()
	if err != nil {
		panic(err)
	}

	db, err := gorm2.ConnectToDB()
	defer db.Close()
	db.LogMode(true)
	db.SingularTable(true)
	db.AutoMigrate(&Todo{})

	e := echo.New()
	//e.Use(middleware.Logger())
	e.Use(CustomContextMiddleware(db))

	if l, ok := e.Logger.(*log.Logger); ok {
		l.SetHeader("${time_rfc3339} ${level}")
		l.SetLevel(log.DEBUG)
	}

	pei := prometheus.NewPrometheusInstrumentation()
	e.Use(pei.PrometheusStatsMiddleware)
	e.GET("/metrics", pei.MetricsEndpoint)

	userResource := UserResource{}

	e.GET("/users", userResource.getUser)
	e.POST("/users", userResource.postUser)

	todoResource := ToDoResource{}

	e.GET("/todos", todoResource.getToDos)
	e.POST("/todos", todoResource.postToDo)

	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%s", config.ServerHost, config.ServerPort)))
}
