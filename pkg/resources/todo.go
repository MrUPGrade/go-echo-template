package resources

import (
	"echoapi/pkg/db/models"
	"echoapi/pkg/echo/context"
	"github.com/labstack/echo"
	"net/http"
)

type ToDoResource struct {
}

func (ToDoResource) Get(c echo.Context) error {
	cc := c.(*context.CustomContext)

	todos := make([]models.Todo, 0)
	cc.DB.Limit(10).Order("id").Find(&todos)

	return c.JSON(http.StatusOK, todos)
}

func (ToDoResource) Post(c echo.Context) error {
	cc := c.(*context.CustomContext)

	todo := new(models.Todo)
	if err := c.Bind(todo); err != nil {
		c.Logger().Print(err)
		return c.JSON(http.StatusBadRequest, err)
	}

	cc.DB.Create(&todo)

	return c.JSON(http.StatusOK, todo)
}
