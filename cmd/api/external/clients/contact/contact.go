package contact

import (
	"context"
	"errors"
	"fmt"
	"github.com/afonsojota/cleancoder_rest-scripts/cmd/api/utils"
	"strings"
	"time"

	"github.com/afonsojota/cleancoder_rest-scripts/cmd/api/models"
	"github.com/afonsojota/go-afonsojota-toolkit/restful/rest"
	"github.com/spf13/viper"
)

//Gateway for contacts API
//go:generate mockgen -source=./contact.go -destination=./mocks/contact_mock.go
type Gateway interface {
	SendContact(c *models.Contact, template string, owner string) error
}

type contactGateway struct {
	rb rest.RequestBuilder
}

func NewContactGateway() Gateway {
	return &contactGateway{
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

func (gtw *contactGateway) SendContact(c *models.Contact, template string, owner string) error {
	template = strings.Replace(template, "\n", `\n`, -1)

	uri := fmt.Sprintf("%s/%s", viper.GetString("contacts.path"), "agent")

	data := map[string]interface{}{
		"ticket_id": c.TicketId,
		"agent": map[string]interface{}{
			"name": "CX_SYSTEM_APPLICATION", "is_new_ticket": true,
		},
		"carrier":           "mail",
		"parent_comment":    1,
		"skip_notification": false,
		"is_classification": false,
		"solution_id":       c.SolutionId,
		"type":              "agent",
		"body_text":         nil,
		"body_html":         strings.ReplaceAll(template, "\\n", ""),
		"subject":           c.Subject,
		"emails_cc":         []string{},
		"emails_bcc":        []string{},
		"email_to":          []string{},
	}

	response := gtw.rb.Post(uri, data, rest.Headers(utils.GetDefaultHeaders()), rest.Context(context.Background()))

	if response.Err != nil || !(int(response.StatusCode/100) == 2) {
		return buildErrorMessage(response)
	}
	return nil
}

func buildErrorMessage(response *rest.Response) error {
	return errors.New(response.String())
}
