package main

import (
	"github.com/afonsojota/cleancoder_rest-scripts/cmd/api/config"
	"github.com/afonsojota/cleancoder_rest-scripts/cmd/api/controllers"
	"github.com/afonsojota/cleancoder_rest-scripts/cmd/api/dao"
	"github.com/afonsojota/cleancoder_rest-scripts/cmd/api/external/clients/ticket"
	"github.com/afonsojota/cleancoder_rest-scripts/cmd/api/external/clients/contact"
	"github.com/afonsojota/cleancoder_rest-scripts/cmd/api/external/clients/search"
	"github.com/afonsojota/cleancoder_rest-scripts/cmd/api/external/clients/solution"
	"github.com/afonsojota/cleancoder_rest-scripts/cmd/api/services"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/afonsojota/go-afonsojota-toolkit/gingonic/mlhandlers"
	"github.com/afonsojota/go-afonsojota-toolkit/goutils/logger"
	"github.com/newrelic/go-agent/_integrations/nrgin/v1"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := run(port); err != nil {
		logger.Errorf("error running server", err)
	}
}

func run(port string) error {
	router := mlhandlers.DefaultafonsojotaRouter()
	router.GET("/ping", func(c *gin.Context) {
		if txn := nrgin.Transaction(c); txn != nil {
			_ = txn.Ignore()
		}

		c.String(http.StatusOK, "pong")
	})

	repository := dao.NewBulkRepository(config.GetDatabase())
	controller := controllers.NewBulkController(
		services.NewBulkService(
			repository,
			solution.NewSolutionGateway(),
			contact.NewContactGateway(),
			ticket.NewTicketGateway(),
			search.NewSearchGateway(),
		),
	)
	controller.SetupRoutes(router)

	return router.Run(":" + port)
}
