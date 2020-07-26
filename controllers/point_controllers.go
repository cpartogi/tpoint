package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/cpartogi/tpoint/config"
	"github.com/cpartogi/tpoint/models"
	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
)

func GetPoint(c echo.Context) error {
	member_id, _ := strconv.Atoi(c.Param("member_id"))

	//ambil header
	api_key := c.Request().Header.Get("x-api-key")
	conf := config.GetConfig()
	if api_key != conf.API_KEY {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "invalid api key"})
	}

	result, err := models.GetPoint(member_id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func AddPoint(c echo.Context) error {

	p := new(models.Pointlog)
	if err := c.Bind(p); err != nil {
		return err
	}

	member_id := p.Member_id
	point_type := p.Point_type
	point_desc := p.Point_desc
	point_amount := p.Point_amount
	created_by := p.Created_by

	//ambil header
	api_key := c.Request().Header.Get("x-api-key")
	conf := config.GetConfig()
	if api_key != conf.API_KEY {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "invalid api key"})
	}

	if member_id == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "member_id is empty"})
	}

	if point_type == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "point_type is empty"})
	}

	if point_type == 1 && point_amount < 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "point amount for point_type 1 must more than 0"})
	}

	if point_type == 2 && point_amount > 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "point amount for point_type 2 must less than 0"})
	}

	if point_amount == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "point_amount is empty"})
	}

	if point_desc == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "point_desc is empty"})
	}

	if created_by == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "created_by is empty"})
	}

	var insert_id uuid.UUID
	insert_id = uuid.NewV4()

	var currentTime time.Time

	currentTime = time.Now()

	created_at := currentTime.Format("2006-01-02 15:04:05.000000")

	result, err := models.AddPoint(insert_id, member_id, point_type, point_desc, point_amount, created_by, created_at)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

func UpdatePoint(c echo.Context) error {
	p := new(models.Pointlog)
	if err := c.Bind(p); err != nil {
		return err
	}
	id := p.Id
	member_id := p.Member_id
	point_type := p.Point_type
	point_desc := p.Point_desc
	point_amount := p.Point_amount
	created_by := p.Created_by
	created_at := p.Created_at

	result, err := models.UpdatePoint(id, member_id, point_type, point_desc, point_amount, created_by, created_at)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)

}
