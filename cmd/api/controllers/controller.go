package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/afonsojota/cleancoder_rest-scripts/cmd/api/services"
	"github.com/afonsojota/cleancoder_rest-scripts/cmd/api/services/parser"
)

type BulkController struct {
	service services.Service
}

func NewBulkController(service services.Service) *BulkController {
	return &BulkController{
		service,
	}
}

func (controller *BulkController) SetupRoutes(router *gin.Engine) {
	router.POST("/bulk/ticket/close", controller.AnswerAndClose)
	router.POST("/bulk/ticket/change_queue", controller.ChangeQueue)
	router.POST("/bulk/ticket/open_ticket", controller.OpenTicket)
	router.POST("/bulk/ticket/search/index", controller.SearchIndex)
}

func (controller *BulkController) ChangeQueue(c *gin.Context) {
	adminID := c.GetHeader("X-Admin-Id")
	file, err := parser.GetFile(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"This file could not be processed: ": err.Error()})
	} else {

		file.Owner = adminID

		err = controller.service.ProcessChangeQueueticket(file)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"This file could not be processed: ": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "processed"})
		}

	}
}

func (controller *BulkController) AnswerAndClose(c *gin.Context) {

	adminID := c.GetHeader("X-Admin-Id")
	file, err := parser.GetFile(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"This file could not be processed: ": err.Error()})
	} else {

		file.Owner = adminID

		err = controller.service.Processticket(file)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"This file could not be processed: ": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "processed"})
		}

	}
}

func (controller *BulkController) OpenTicket(c *gin.Context) {

	adminID := c.GetHeader("X-Admin-Id")
	file, err := parser.GetFile(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"This file could not be processed: ": err.Error()})
	} else {

		file.Owner = adminID

		err = controller.service.OpenTicket(file)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"This file could not be processed: ": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "processed"})
		}

	}
}

func (controller *BulkController) SearchIndex(c *gin.Context) {

	adminID := c.GetHeader("X-Admin-Id")
	file, err := parser.GetFile(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"This file could not be processed: ": err.Error()})
	} else {

		file.Owner = adminID

		err = controller.service.SearchIndex(file)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"This file could not be processed: ": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "processed"})
		}

	}
}
