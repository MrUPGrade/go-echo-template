package context

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
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
