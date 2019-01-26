package main

import (
	"github.com/labstack/echo"
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

func (UserResource) postUser(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		c.Logger().Print(err)
	}
	return c.JSON(http.StatusOK, u)
}

type ToDoResource struct {
}

func (ToDoResource) getToDos(c echo.Context) error {
	cc := c.(*CustomContext)

	todos := make([]Todo, 0)
	cc.DB.Limit(10).Order("id").Find(&todos)

	return c.JSON(http.StatusOK, todos)
}

func (ToDoResource) postToDo(c echo.Context) error {
	cc := c.(*CustomContext)

	todo := new(Todo)
	if err := c.Bind(todo); err != nil {
		c.Logger().Print(err)
		return c.JSON(http.StatusBadRequest, nil)
	}

	cc.DB.Create(&todo)

	return c.JSON(http.StatusOK, todo)
}
