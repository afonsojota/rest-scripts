package dao

import (
	"github.com/afonsojota/cleancoder_rest-scripts/cmd/api/models"
	"github.com/afonsojota/cleancoder_rest-scripts/cmd/api/utils"
	"github.com/afonsojota/go-afonsojota-toolkit/goutils/logger"
	"upper.io/db.v3/lib/sqlbuilder"
)

//go:generate mockgen -source=./repository.go -destination=./mocks/repository_mock.go
type Repository interface {
	InsertTicket(c *models.Ticket) error
	Changetickettatus(c *models.Ticket) error
	InsertInputFile(i *models.InputFile) (int32, error)
	ChangeInputFileStatus(i *models.InputFile) error
	InsertSendErrors(s *models.SendErros) error
}

type bulkMySqlRepository struct {
	session sqlbuilder.Database
}

func NewBulkRepository(session sqlbuilder.Database) Repository {
	return &bulkMySqlRepository{
		session,
	}
}

func (repo *bulkMySqlRepository) InsertTicket(c *models.Ticket) error {
	c.Status = string(utils.RUNNING)
	_, err := repo.session.Collection("send_ticket").Insert(c)
	return err
}

func (repo *bulkMySqlRepository) Changetickettatus(c *models.Ticket) error {
	res := repo.session.Collection("send_ticket").Find("ticket_id", c.Id)
	var sendTicket models.Ticket

	if res.Err() != nil {
		return res.Err()
	}

	if err := res.One(&sendTicket); err != nil {
		logger.Error("Close ticket error", err)
		return err
	}

	sendTicket.SolutionId = c.SolutionId
	sendTicket.Status = c.Status

	return res.Update(sendTicket)
}

func (repo *bulkMySqlRepository) InsertInputFile(file *models.InputFile) (int32, error) {

	file.Status = string(utils.PENDING)

	_, err := repo.session.Collection("input_file").Insert(file)

	res := repo.session.Collection("input_file").Find("name", file.Name)
	var inputFile models.InputFile
	if err = res.One(&inputFile); err != nil {
		return 0, err
	}

	return inputFile.Id, nil
}

func (repo *bulkMySqlRepository) ChangeInputFileStatus(file *models.InputFile) error {

	res := repo.session.Collection("input_file").Find("id", file.Id)

	if res.Err() != nil {
		return res.Err()
	}
	var inputFile models.InputFile
	if err := res.One(&inputFile); err != nil {
		logger.Error("Error changing input file status", err)
		return err
	}

	inputFile.Status = file.Status

	return res.Update(inputFile)
}

func (repo *bulkMySqlRepository) InsertSendErrors(s *models.SendErros) error {
	_, err := repo.session.Collection("send_errors").Insert(s)
	return err
}
