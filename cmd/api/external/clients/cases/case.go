package ticket

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	_ "github.com/afonsojota/cleancoder_rest-scripts/cmd/api/config"
	"github.com/afonsojota/cleancoder_rest-scripts/cmd/api/models"
	"github.com/afonsojota/cleancoder_rest-scripts/cmd/api/utils"
	"github.com/afonsojota/go-afonsojota-toolkit/restful/rest"
	"github.com/spf13/viper"
)

const ReasonId = 1185

//go:generate mockgen -source=./ticket.go -destination=./mocks/ticket_mock.go
type Gateway interface {
	CloseTicket(c *models.Ticket, owner string) error
	ChangeQueue(c *models.Ticket, owner string) error
	OpenTicket(c *models.Ticket, owner string) error
}

type ticketGateway struct {
	rb rest.RequestBuilder
}

func NewTicketGateway() Gateway {
	return &ticketGateway{
		rb: rest.RequestBuilder{
			BaseURL:     viper.Get("host.default").(string),
			ContentType: rest.JSON,
			CustomPool: &rest.CustomPool{
				MaxIdleConnsPerHost: 100,
			},
			Timeout: 2 * time.Second,
		},
	}
}

func (gtw *ticketGateway) CloseTicket(c *models.Ticket, owner string) error {

	data := map[string]interface{}{}

	if c.SolutionId == 0 {
		data = map[string]interface{}{
			"reason_id": ReasonId,
		}
	}

	uri := fmt.Sprintf("%s/%d/%s", viper.GetString("ticket.path"), c.Id, "close")
	response := gtw.rb.Put(uri, data, rest.Headers(utils.GetDefaultHeaders()), rest.Context(context.Background()))

	if response.Err != nil || !(int(response.StatusCode/100) == 2) {
		return buildErrorMessage(response)
	}
	return nil
}

func (gtw *ticketGateway) ChangeQueue(c *models.Ticket, owner string) error {

	newQueueId, _ := strconv.Atoi(c.SiteId)

	body := map[string]interface{}{
		"queue_id": newQueueId,
	}

	uri := fmt.Sprintf("%s/%d/%s", viper.GetString("ticket.path"), c.Id, "change_queue")
	response := gtw.rb.Put(uri, body, rest.Headers(utils.GetDefaultHeaders()), rest.Context(context.Background()))

	if response.Err != nil || !(int(response.StatusCode/100) == 2) {
		return buildErrorMessage(response)
	}
	return nil
}

func (gtw *ticketGateway) OpenTicket(c *models.Ticket, owner string) error {

	data := map[string]interface{}{}

	uri := fmt.Sprintf("%s/%d/%s", viper.GetString("ticket.path"), c.Id, "open")
	response := gtw.rb.Put(uri, data, rest.Headers(utils.GetDefaultHeaders()), rest.Context(context.Background()))

	if response.Err != nil || !(int(response.StatusCode/100) == 2) {
		return buildErrorMessage(response)
	}
	return nil
}

func buildErrorMessage(response *rest.Response) error {
	return errors.New(response.String())
}
