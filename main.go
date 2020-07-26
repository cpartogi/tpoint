package main

import (
	"github.com/cpartogi/tpoint/controllers"
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()

	e.GET("/point/:member_id", controllers.GetPoint)
	e.POST("/point", controllers.AddPoint)
	e.PUT("/point", controllers.UpdatePoint)

	e.Logger.Fatal(e.Start(":1234"))
}
