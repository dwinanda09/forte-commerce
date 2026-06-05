package middleware

import (
	"github.com/labstack/echo/v4"
)

func SellerOnly(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role, _ := c.Get("role").(string)
		if role != "seller" {
			return c.JSON(403, map[string]string{
				"error": "Seller access required",
			})
		}
		return next(c)
	}
}
