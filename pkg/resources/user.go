package resources

import (
	"github.com/labstack/echo"
	"net/http"
)

type User struct {
	Name string `json:"name"`
}

type UserResource struct {
}

func (UserResource) Get(c echo.Context) error {
	return c.JSON(http.StatusOK, User{"Stanisław"})
}

func (UserResource) Post(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		c.Logger().Print(err)
	}
	return c.JSON(http.StatusOK, u)
}
