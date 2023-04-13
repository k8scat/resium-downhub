package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/k8scat/downhub/downloader/wenku"
	"github.com/k8scat/downhub/util"
)

type ParseRequestBody struct {
	Url string `json:"url"`
}

func Wenku(c *gin.Context) {
	var data ParseRequestBody
	if err := c.ShouldBindJSON(&data); err != nil {
		util.Ding(err.Error())
		c.JSON(util.JSONResponse(http.StatusBadRequest, "bad request", nil))
		return
	}
	util.Ding(fmt.Sprintf("Parsing url: %s", data.Url))

	location, err := wenku.GetLocation(data.Url)
	if err != nil {
		util.Ding(err.Error())
		c.JSON(util.JSONResponse(http.StatusInternalServerError, "parse error", nil))
		return
	}
	util.Ding(fmt.Sprintf("parse success: %s", location))
	c.JSON(util.JSONResponse(http.StatusOK, "parse ok", location))
}
