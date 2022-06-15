package models

import "fmt"

type Contact struct {
	Subject    string
	Body       string
	SolutionId int64
	SiteId     string
	TicketId     int64
}

func NewContact(ticketId int64, siteId string, solutionId int64) *Contact {
	c := new(Contact)
	c.TicketId = ticketId
	c.SiteId = siteId
	c.SolutionId = solutionId
	if c.SiteId == "MLB" {
		c.Subject = fmt.Sprintf("Resposta a sua consulta [Ticket: %d]", ticketId)
	} else {
		c.Subject = fmt.Sprintf("Respuesta a tu consulta [Ticket: %d]", ticketId)
	}
	return c
}
