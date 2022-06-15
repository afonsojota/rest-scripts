package ticket

import (
	"net/http"
	"testing"

	"github.com/afonsojota/cleancoder_rest-scripts/cmd/api/models"
	"github.com/afonsojota/cleancoder_rest-scripts/cmd/api/utils"
	"github.com/afonsojota/go-afonsojota-toolkit/restful/rest"
	"github.com/stretchr/testify/assert"
)

func TestCloseTicket(t *testing.T) {
	rest.StartMockupServer()
	defer rest.StopMockupServer()

	_ = rest.AddMockups(&rest.Mock{
		URL:          "http://internal-api.afonsojota.com/ticket/dev/123/close",
		HTTPMethod:   http.MethodPut,
		ReqHeaders:   utils.GetDefaultHeaders(),
		RespHTTPCode: http.StatusOK,
	})

	ticketContact := new(models.Ticket)
	ticketContact.Id = int64(123)
	ticketContact.SolutionId = int64(10)
	ticketContact.SiteId = "MLA"

	contactGateway := NewTicketGateway()
	err := contactGateway.CloseTicket(ticketContact, "admin")

	assert.NotNil(t, err)
}

func TestCloseTicketFailError(t *testing.T) {
	rest.StartMockupServer()
	defer rest.StopMockupServer()

	_ = rest.AddMockups(&rest.Mock{
		URL:          "http://internal-api.afonsojota.com/ticket/dev/123/close",
		HTTPMethod:   http.MethodPut,
		ReqHeaders:   utils.GetDefaultHeaders(),
		RespHTTPCode: http.StatusInternalServerError,
	})

	ticketContact := new(models.Ticket)
	ticketContact.Id = int64(123)
	ticketContact.SolutionId = int64(0)
	ticketContact.SiteId = "MLA"

	contactGateway := NewTicketGateway()
	err := contactGateway.CloseTicket(ticketContact, "admin")

	assert.NotNil(t, err)
}

func TestChangeQueue(t *testing.T) {
	rest.StartMockupServer()
	defer rest.StopMockupServer()

	_ = rest.AddMockups(&rest.Mock{
		URL:          "http://internal-api.afonsojota.com/ticket/dev/123/change_queue",
		HTTPMethod:   http.MethodPut,
		ReqHeaders:   utils.GetDefaultHeaders(),
		RespHTTPCode: http.StatusOK,
	})

	ticketContact := new(models.Ticket)
	ticketContact.Id = int64(123)
	ticketContact.SiteId = "1233"

	contactGateway := NewTicketGateway()
	err := contactGateway.ChangeQueue(ticketContact, "admin")

	assert.NotNil(t, err)
}
