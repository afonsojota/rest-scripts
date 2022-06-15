package models

type SendErros struct {
	TicketId int64  `db:"ticket_id"`
	Error  string `db:"error"`
}
