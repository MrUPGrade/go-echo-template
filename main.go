package main

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"net/http"
)

type User struct {
	Name string `json:"name"`
}

type UserResource struct {
}

func (UserResource) getUser(c echo.Context) error {
	return c.JSON(http.StatusOK, User{"Stanisław"})
}

func (UserResource) postUser(c echo.Context) (err error) {
	u := new(User)
	if err = c.Bind(u); err != nil {
		c.Logger().Print(err)
	}
	return c.JSON(http.StatusOK, u)
}

func main() {
	config, err := loadConfig()
	if err != nil {
		panic(err)
	}

	e := echo.New()
	if l, ok := e.Logger.(*log.Logger); ok {
		l.SetHeader("${time_rfc3339} ${level}")
		l.SetLevel(log.DEBUG)
	}

	userResource := UserResource{}

	e.GET("/users", userResource.getUser)
	e.POST("/users", userResource.postUser)
	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%s", config.ServerHost, config.ServerPort)))
}
