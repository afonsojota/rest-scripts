package contact

import (
	"net/http"
	"testing"

	_ "github.com/afonsojota/cleancoder_rest-scripts/cmd/api/config"
	"github.com/afonsojota/cleancoder_rest-scripts/cmd/api/models"
	"github.com/afonsojota/cleancoder_rest-scripts/cmd/api/utils"
	"github.com/afonsojota/go-afonsojota-toolkit/restful/rest"
	"github.com/stretchr/testify/assert"
)

func TestSendContact(t *testing.T) {
	rest.StartMockupServer()
	defer rest.StopMockupServer()

	_ = rest.AddMockups(&rest.Mock{
		URL:          "http://internal-api.afonsojota.com/contacts/dev/agent",
		HTTPMethod:   http.MethodPost,
		ReqHeaders:   utils.GetDefaultHeaders(),
		RespHTTPCode: http.StatusOK,
	})

	ticketContact := new(models.Ticket)
	ticketContact.Id = int64(123)
	ticketContact.SolutionId = int64(10)
	ticketContact.SiteId = "MLA"

	contact := new(models.Contact)
	models.NewContact(ticketContact.Id, ticketContact.SiteId, ticketContact.SolutionId)

	template := "email body"
	contactGateway := NewContactGateway()
	err := contactGateway.SendContact(contact, template, "admin")

	assert.NotNil(t, err)
}

func TestSendContactFailError500(t *testing.T) {
	rest.StartMockupServer()
	defer rest.StopMockupServer()

	_ = rest.AddMockups(&rest.Mock{
		URL:          "http://internal-api.afonsojota.com/contacts/dev/agent",
		HTTPMethod:   http.MethodPost,
		ReqHeaders:   utils.GetDefaultHeaders(),
		RespHTTPCode: http.StatusInternalServerError,
	})

	ticketContact := new(models.Ticket)
	ticketContact.Id = int64(123)
	ticketContact.SolutionId = int64(10)
	ticketContact.SiteId = "MLB"

	contact := models.NewContact(ticketContact.Id, ticketContact.SiteId, ticketContact.SolutionId)

	template := "email body"
	contactGateway := NewContactGateway()
	err := contactGateway.SendContact(contact, template, "admin")

	assert.NotNil(t, err)
}
