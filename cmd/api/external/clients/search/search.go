package search

import (
	"context"
	"errors"
	"fmt"
	"github.com/afonsojota/cleancoder_rest-scripts/cmd/api/utils"
	"time"

	"github.com/afonsojota/cleancoder_rest-scripts/cmd/api/models"
	"github.com/afonsojota/go-afonsojota-toolkit/restful/rest"
	"github.com/spf13/viper"
)

//Gateway for contacts API
//go:generate mockgen -source=./search.go -destination=./mocks/search_mock.go
type Gateway interface {
	Indexer(c *models.Ticket, owner string) error
}

type searchGateway struct {
	rb rest.RequestBuilder
}

func NewSearchGateway() Gateway {
	return &searchGateway{
		rb: rest.RequestBuilder{
			BaseURL:     viper.Get("host.default").(string),
			ContentType: rest.JSON,
			CustomPool: &rest.CustomPool{
				MaxIdleConnsPerHost: 100,
			},
			Timeout: 5 * time.Second,
		},
	}
}

func (gtw *searchGateway) Indexer(c *models.Ticket, owner string) error {
	uri := fmt.Sprintf("%s/%s%d", viper.GetString("search.path"), "indexer?ids=", c.Id)

	data := map[string]interface{}{}

	response := gtw.rb.Post(uri, data, rest.Headers(utils.GetDefaultHeaders()), rest.Context(context.Background()))

	if response.Err != nil || !(int(response.StatusCode/100) == 2) {
		return buildErrorMessage(response)
	}
	return nil
}

func buildErrorMessage(response *rest.Response) error {
	return errors.New(response.String())
}
