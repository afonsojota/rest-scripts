package models

type Ticket struct {
	Id         int64  `db:"ticket_id"`
	SolutionId int64  `db:"solution_id"`
	Status     string `db:"status"`
	SiteId     string `db:"site_id"`
	FileId     int32  `db:"file_id"`
}
