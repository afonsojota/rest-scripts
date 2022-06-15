package services

import (
	"bytes"
	"errors"
	"io/ioutil"
	"testing"

	"github.com/golang/mock/gomock"
	mock_dao "github.com/afonsojota/cleancoder_rest-scripts/cmd/api/dao/mocks"
	mock_ticket "github.com/afonsojota/cleancoder_rest-scripts/cmd/api/external/clients/ticket/mocks"
	mock_contact "github.com/afonsojota/cleancoder_rest-scripts/cmd/api/external/clients/contact/mocks"
	mock_solution "github.com/afonsojota/cleancoder_rest-scripts/cmd/api/external/clients/solution/mocks"
	"github.com/afonsojota/cleancoder_rest-scripts/cmd/api/models"
	"github.com/afonsojota/cleancoder_rest-scripts/cmd/api/utils"
	"github.com/stretchr/testify/assert"
)

func TestCloseTicket(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockticket := mock_ticket.NewMockGateway(ctrl)
	mockRepository := mock_dao.NewMockRepository(ctrl)
	service := NewBulkService(mockRepository, nil, nil, mockticket)
	Ticket := &models.Ticket{Id: 19}
	ownerID := "afonso"

	t.Run("Success", func(t *testing.T) {
		mockticket.EXPECT().CloseTicket(gomock.Eq(Ticket), gomock.Eq(ownerID)).Return(nil)
		mockRepository.EXPECT().Changetickettatus(gomock.Eq(Ticket)).Return(nil)
		err := service.(*bulkService).closeTicket(Ticket, ownerID)
		assert.Nil(t, err)
	})
	t.Run("Fail - Database", func(t *testing.T) {
		mockticket.EXPECT().CloseTicket(gomock.Eq(Ticket), gomock.Eq(ownerID)).Return(nil)
		mockRepository.EXPECT().Changetickettatus(gomock.Eq(Ticket)).Return(errors.New("Error DB"))
		err := service.(*bulkService).closeTicket(Ticket, ownerID)
		assert.NotNil(t, err)
	})
	t.Run("Fail - Gateway", func(t *testing.T) {
		mockticket.EXPECT().CloseTicket(gomock.Eq(Ticket), gomock.Eq(ownerID)).Return(errors.New("Error Gateway"))
		err := service.(*bulkService).closeTicket(Ticket, ownerID)
		assert.NotNil(t, err)
	})

}

func TestSaveError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockticket := mock_ticket.NewMockGateway(ctrl)
	mockRepository := mock_dao.NewMockRepository(ctrl)
	service := NewBulkService(mockRepository, nil, nil, mockticket)

	Ticket := &models.Ticket{Id: 19, Status: "ERROR"}
	sendErrors := &models.SendErros{TicketId: Ticket.Id, Error: "Error SendTicket"}

	mockRepository.EXPECT().InsertSendErrors(gomock.Eq(sendErrors)).Return(nil)
	mockRepository.EXPECT().Changetickettatus(gomock.Eq(Ticket)).Return(nil)

	ticketResp := service.(*bulkService).saveError(Ticket, errors.New(sendErrors.Error))

	assert.Equal(t, string(utils.ERROR), ticketResp.Status)

}

func TestProcessTicket(t *testing.T) {
	file, err := ioutil.ReadFile("resources/ticket_to_answer.csv")
	buffer := bytes.NewBuffer(file)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mock_dao.NewMockRepository(ctrl)
	mockSolution := mock_solution.NewMockGateway(ctrl)
	mockContact := mock_contact.NewMockGateway(ctrl)
	mockticket := mock_ticket.NewMockGateway(ctrl)

	service := NewBulkService(mockRepository, mockSolution, mockContact, mockticket)

	template := "body template"
	body := &models.Body{Template: template}
	content := models.Content{Channel: "DEFAULT", Body: *body}
	contents := []models.Content{content}
	version := models.Version{Sites: []string{"MLB"}, Contents: contents}
	solution := &models.Solution{Versions: []models.Version{version}}

	mockRepository.EXPECT().InsertInputFile(gomock.Any()).Return(int32(1), nil)
	mockSolution.EXPECT().GetSolution(gomock.Any()).Return(solution, nil)
	mockContact.EXPECT().SendContact(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	mockRepository.EXPECT().ChangeInputFileStatus(gomock.Any()).Return(nil)
	mockRepository.EXPECT().InsertTicket(gomock.Any()).Return(nil)
	mockticket.EXPECT().CloseTicket(gomock.Any(), gomock.Any()).Return(nil)
	mockRepository.EXPECT().Changetickettatus(gomock.Any()).Return(nil)

	inputFile := &models.InputFile{Owner: "admin"}
	inputFile.Data = *buffer

	err = service.Processticket(inputFile)

	assert.Nil(t, err)
}

func TestProcessChangeQueueticket(t *testing.T) {
	file, err := ioutil.ReadFile("resources/ticket_change_queue.csv")
	buffer := bytes.NewBuffer(file)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mock_dao.NewMockRepository(ctrl)
	mockticket := mock_ticket.NewMockGateway(ctrl)

	service := NewBulkService(mockRepository, nil, nil, mockticket)

	mockticket.EXPECT().ChangeQueue(gomock.Any(), gomock.Any()).Return(nil)
	mockRepository.EXPECT().InsertInputFile(gomock.Any()).Return(int32(1), nil)
	mockRepository.EXPECT().InsertTicket(gomock.Any()).Return(nil)
	mockRepository.EXPECT().Changetickettatus(gomock.Any()).Return(nil)
	mockRepository.EXPECT().ChangeInputFileStatus(gomock.Any()).Return(nil)

	inputFile := &models.InputFile{Owner: "admin"}
	inputFile.Data = *buffer

	err = service.ProcessChangeQueueticket(inputFile)

	assert.Nil(t, err)
}

func TestProcessTicketInputFileError(t *testing.T) {
	file, err := ioutil.ReadFile("resources/ticket_to_answer.csv")
	buffer := bytes.NewBuffer(file)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mock_dao.NewMockRepository(ctrl)

	service := NewBulkService(mockRepository, nil, nil, nil)

	inputFile := &models.InputFile{Owner: "admin"}
	inputFile.Data = *buffer

	mockRepository.EXPECT().InsertInputFile(gomock.Any()).Return(int32(1), errors.New("Error InsertInputFile"))
	err = service.Processticket(inputFile)
	assert.NotNil(t, err)

	assert.True(t, true)
}

func TestProcessTicketInsertError(t *testing.T) {
	file, err := ioutil.ReadFile("resources/ticket_to_answer.csv")
	buffer := bytes.NewBuffer(file)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mock_dao.NewMockRepository(ctrl)

	service := NewBulkService(mockRepository, nil, nil, nil)

	inputFile := &models.InputFile{Owner: "admin"}
	inputFile.Data = *buffer

	mockRepository.EXPECT().InsertInputFile(gomock.Any()).Return(int32(1), nil)
	mockRepository.EXPECT().InsertTicket(gomock.Any()).Return(errors.New("Error InsertTicket"))
	mockRepository.EXPECT().InsertSendErrors(gomock.Any()).Return(nil)
	mockRepository.EXPECT().Changetickettatus(gomock.Any()).Return(nil)
	mockRepository.EXPECT().ChangeInputFileStatus(gomock.Any()).Return(nil)

	err = service.Processticket(inputFile)
	assert.Nil(t, err)

	assert.True(t, true)
}

func TestProcessticketendContactError(t *testing.T) {
	file, _ := ioutil.ReadFile("resources/ticket_to_answer.csv")
	buffer := bytes.NewBuffer(file)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mock_dao.NewMockRepository(ctrl)
	mockContact := mock_contact.NewMockGateway(ctrl)
	mockSolution := mock_solution.NewMockGateway(ctrl)

	service := NewBulkService(mockRepository, mockSolution, mockContact, nil)

	template := "body template"
	body := &models.Body{Template: template}
	content := models.Content{Channel: "DEFAULT", Body: *body}
	contents := []models.Content{content}
	version := models.Version{Sites: []string{"MLA"}, Contents: contents}
	solution := &models.Solution{Versions: []models.Version{version}}

	inputFile := &models.InputFile{Owner: "admin"}
	inputFile.Data = *buffer

	mockRepository.EXPECT().InsertInputFile(gomock.Any()).Return(int32(1), nil)
	mockRepository.EXPECT().InsertTicket(gomock.Any()).Return(nil)
	mockSolution.EXPECT().GetSolution(gomock.Any()).Return(solution, nil)
	mockContact.EXPECT().SendContact(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("Error SendContact"))
	mockRepository.EXPECT().InsertSendErrors(gomock.Any()).Return(nil)
	mockRepository.EXPECT().Changetickettatus(gomock.Any()).Return(nil)
	mockRepository.EXPECT().ChangeInputFileStatus(gomock.Any()).Return(nil)

	service.Processticket(inputFile)

	assert.True(t, true)
}

func TestProcessticketolutionError(t *testing.T) {
	file, _ := ioutil.ReadFile("resources/ticket_to_answer.csv")
	buffer := bytes.NewBuffer(file)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mock_dao.NewMockRepository(ctrl)
	mockContact := mock_contact.NewMockGateway(ctrl)
	mockSolution := mock_solution.NewMockGateway(ctrl)

	service := NewBulkService(mockRepository, mockSolution, mockContact, nil)

	template := "body template"
	body := &models.Body{Template: template}
	content := models.Content{Channel: "DEFAULT", Body: *body}
	contents := []models.Content{content}
	version := models.Version{Sites: []string{"MLB"}, Contents: contents}
	solution := &models.Solution{Versions: []models.Version{version}}

	inputFile := &models.InputFile{Owner: "admin"}
	inputFile.Data = *buffer

	mockRepository.EXPECT().InsertInputFile(gomock.Any()).Return(int32(1), nil)
	mockRepository.EXPECT().InsertTicket(gomock.Any()).Return(nil)
	mockSolution.EXPECT().GetSolution(gomock.Any()).Return(solution, errors.New("Error GetSolution"))
	mockRepository.EXPECT().InsertSendErrors(gomock.Any()).Return(nil)
	mockRepository.EXPECT().Changetickettatus(gomock.Any()).Return(nil)
	mockRepository.EXPECT().ChangeInputFileStatus(gomock.Any()).Return(nil)

	err := service.Processticket(inputFile)

	assert.Nil(t, err)
}

func TestProcessTicketCloseError(t *testing.T) {
	file, _ := ioutil.ReadFile("resources/ticket_to_answer.csv")
	buffer := bytes.NewBuffer(file)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mock_dao.NewMockRepository(ctrl)
	mockContact := mock_contact.NewMockGateway(ctrl)
	mockSolution := mock_solution.NewMockGateway(ctrl)

	mockticket := mock_ticket.NewMockGateway(ctrl)

	service := NewBulkService(mockRepository, mockSolution, mockContact, mockticket)

	template := "body template"
	body := &models.Body{Template: template}
	content := models.Content{Channel: "DEFAULT", Body: *body}
	contents := []models.Content{content}
	version := models.Version{Sites: []string{"MLB"}, Contents: contents}
	solution := &models.Solution{Versions: []models.Version{version}}

	inputFile := &models.InputFile{Owner: "admin"}
	inputFile.Data = *buffer

	mockRepository.EXPECT().InsertInputFile(gomock.Any()).Return(int32(1), nil)
	mockRepository.EXPECT().InsertTicket(gomock.Any()).Return(nil)
	mockSolution.EXPECT().GetSolution(gomock.Any()).Return(solution, nil)
	mockContact.EXPECT().SendContact(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	mockRepository.EXPECT().InsertSendErrors(gomock.Any()).Return(nil)
	mockRepository.EXPECT().Changetickettatus(gomock.Any()).Return(nil)
	mockRepository.EXPECT().ChangeInputFileStatus(gomock.Any()).Return(nil)
	mockticket.EXPECT().CloseTicket(gomock.Any(), gomock.Any()).Return(errors.New("Error GetSolution"))

	err := service.Processticket(inputFile)

	assert.Nil(t, err)
}
