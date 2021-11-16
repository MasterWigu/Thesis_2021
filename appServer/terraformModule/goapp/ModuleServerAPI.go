package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func declareAPIEndpoint() *gin.Engine {
	r := gin.Default()

	r.POST("/plan", planActionHandler)
	r.POST("/execute", execActionHandler)
	r.POST("/execute/:actionId/confirm", confirmActionHandler)

	r.POST("/actions/:actionId/ledgerId/:ledgerId", informIdHandler)



	return r
}

func startListening(r *gin.Engine) {
	r.RunTLS(":8060", "/certs/cert.pem", "/certs/key.pem")
}

func planActionHandler(c *gin.Context) {

	formFile, err := c.FormFile("toolResources")
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	resp, err := processor.planAction(formFile)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSONP(http.StatusOK, resp)
}

func execActionHandler(c *gin.Context) {
	formFile, err := c.FormFile("toolResources")
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	resp, err := processor.execAction(formFile)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSONP(http.StatusOK, resp)
}

func confirmActionHandler(c *gin.Context) {
	actionId := c.Param("actionId")
	if actionId == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	resp, idFound, err := processor.ConfirmAction(actionId)
	if !idFound {
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSONP(http.StatusOK, resp)
}

func informIdHandler(c *gin.Context) {
	actionId := c.Param("actionId")
	if actionId == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	ledgerId := c.Param("ledgerId")
	if ledgerId == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	err := processor.BindLedgerId(actionId, ledgerId)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}