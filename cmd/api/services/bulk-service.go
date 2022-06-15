package services

import (
	"github.com/afonsojota/cleancoder_rest-scripts/cmd/api/dao"
	"github.com/afonsojota/cleancoder_rest-scripts/cmd/api/external/clients/ticket"
	"github.com/afonsojota/cleancoder_rest-scripts/cmd/api/external/clients/contact"
	"github.com/afonsojota/cleancoder_rest-scripts/cmd/api/external/clients/search"
	"github.com/afonsojota/cleancoder_rest-scripts/cmd/api/external/clients/solution"
	"github.com/afonsojota/cleancoder_rest-scripts/cmd/api/models"
	"github.com/afonsojota/cleancoder_rest-scripts/cmd/api/services/parser"
	"github.com/afonsojota/cleancoder_rest-scripts/cmd/api/utils"
	"github.com/afonsojota/go-afonsojota-toolkit/goutils/logger"
)

type Service interface {
	Processticket(file *models.InputFile) error
	ProcessChangeQueueticket(file *models.InputFile) error
	OpenTicket(file *models.InputFile) error
	SearchIndex(file *models.InputFile) error
}

type bulkService struct {
	repository  dao.Repository
	solutionGTW solution.Gateway
	contactGTW  contact.Gateway
	ticketGTW    ticket.Gateway
	searchGTW	search.Gateway
}

func NewBulkService(repository dao.Repository, solutionGTW solution.Gateway, contactGTW contact.Gateway, ticketGTW ticket.Gateway, searchGTW search.Gateway) Service {
	return &bulkService{
		repository:  repository,
		solutionGTW: solutionGTW,
		contactGTW:  contactGTW,
		ticketGTW:    ticketGTW,
		searchGTW:	searchGTW,
	}
}

func (bulk *bulkService) OpenTicket(file *models.InputFile) error {
	ticket, err := parser.CsvToOpen(file)
	if err != nil {
		return err
	}

	var fileId int32
	fileId, err = bulk.repository.InsertInputFile(file)
	file.Id = fileId

	if err != nil {
		return err
	}

	for _, c := range ticket {
		c.FileId = fileId
		if err := bulk.repository.InsertTicket(c); err != nil {
			bulk.saveError(c, err)
			continue
		}

		logger.Infof("Id: %d", c.Id)

		err := bulk.openTicket(c, file.Owner)
		if err != nil {
			bulk.saveError(c, err)
			continue
		}
	}

	file.Status = string(utils.FINISHED)
	return bulk.repository.ChangeInputFileStatus(file)
}

func (bulk *bulkService) SearchIndex(file *models.InputFile) error {
	ticket, err := parser.CsvToOpen(file)
	if err != nil {
		return err
	}

	var fileId int32
	fileId, err = bulk.repository.InsertInputFile(file)
	file.Id = fileId

	if err != nil {
		return err
	}

	for _, c := range ticket {
		c.FileId = fileId
		if err := bulk.repository.InsertTicket(c); err != nil {
			bulk.saveError(c, err)
			continue
		}

		logger.Infof("Id: %d", c.Id)

		err := bulk.searchIndex(c, file.Owner)
		if err != nil {
			bulk.saveError(c, err)
			continue
		}
	}

	file.Status = string(utils.FINISHED)
	return bulk.repository.ChangeInputFileStatus(file)
}

func (bulk *bulkService) ProcessChangeQueueticket(file *models.InputFile) error {
	ticket, err := parser.CsvToTicket(file)
	if err != nil {
		return err
	}

	var fileId int32
	fileId, _ = bulk.repository.InsertInputFile(file)
	file.Id = fileId

	for _, c := range ticket {
		c.FileId = fileId
		if err := bulk.repository.InsertTicket(c); err != nil {
			bulk.saveError(c, err)
			continue
		}

		logger.Infof("Id: %d| QueueId: %s\n", c.Id, c.SiteId)

		err := bulk.changeQueue(c, file.Owner)
		if err != nil {
			bulk.saveError(c, err)
			continue
		}
	}

	file.Status = string(utils.FINISHED)
	return bulk.repository.ChangeInputFileStatus(file)

}

func (bulk *bulkService) Processticket(file *models.InputFile) error {
	ticket, err := parser.CsvToTicket(file)
	if err != nil {
		return err
	}

	var fileId int32
	fileId, err = bulk.repository.InsertInputFile(file)
	file.Id = fileId

	if err != nil {
		return err
	}

	for _, c := range ticket {
		c.FileId = fileId
		if err := bulk.repository.InsertTicket(c); err != nil {
			bulk.saveError(c, err)
			continue
		}

		if c.SolutionId > 0 {
			solution, err := bulk.solutionGTW.GetSolution(c.SolutionId)
			if err != nil {
				bulk.saveError(c, err)
				continue
			}

			template := bulk.getBodyContent(solution.Versions, c.SiteId)

			contact := models.NewContact(c.Id, c.SiteId, c.SolutionId)
			err = bulk.contactGTW.SendContact(contact, template, file.Owner)
			if err != nil {
				bulk.saveError(c, err)
				continue
			}
		}

		logger.Infof("Id: %d| SolutionId: %d| SiteId: %s\n", c.Id, c.SolutionId, c.SiteId)

		err := bulk.closeTicket(c, file.Owner)
		if err != nil {
			bulk.saveError(c, err)
			continue
		}
	}

	file.Status = string(utils.FINISHED)
	return bulk.repository.ChangeInputFileStatus(file)

}

func (bulk *bulkService) getBodyContent(versions []models.Version, siteId string) string {
	for _, v := range versions {
		for _, s := range v.Sites {
			if s == siteId {
				for _, c := range v.Contents {
					if c.Channel == "DEFAULT" {
						return c.Body.Template
					}
				}
			}
		}
	}
	return ""
}

func (bulk *bulkService) saveError(c *models.Ticket, err error) *models.Ticket {
	sendError := new(models.SendErros)
	sendError.TicketId = c.Id
	sendError.Error = err.Error()
	bulk.repository.InsertSendErrors(sendError)
	c.Status = string(utils.ERROR)
	bulk.repository.Changetickettatus(c)
	return c
}

func (bulk *bulkService) closeTicket(c *models.Ticket, owner string) error {
	if err := bulk.ticketGTW.CloseTicket(c, owner); err != nil {
		return err
	}
	c.Status = string(utils.FINISHED)
	return bulk.repository.Changetickettatus(c)
}

func (bulk *bulkService) changeQueue(c *models.Ticket, owner string) error {
	if err := bulk.ticketGTW.ChangeQueue(c, owner); err != nil {
		return err
	}
	c.Status = string(utils.FINISHED)
	return bulk.repository.Changetickettatus(c)
}

func (bulk *bulkService) openTicket(c *models.Ticket, owner string) error {
	if err := bulk.ticketGTW.OpenTicket(c, owner); err != nil {
		return err
	}
	c.Status = string(utils.FINISHED)
	return bulk.repository.Changetickettatus(c)
}

func (bulk *bulkService) searchIndex(c *models.Ticket, owner string) error {
	if err := bulk.searchGTW.Indexer(c, owner); err != nil {
		return err
	}
	c.Status = string(utils.FINISHED)
	return bulk.repository.Changetickettatus(c)
}