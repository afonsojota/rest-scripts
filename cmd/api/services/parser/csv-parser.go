package parser

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/afonsojota/cleancoder_rest-scripts/cmd/api/models"
	"github.com/afonsojota/go-afonsojota-toolkit/goutils/logger"
	"io"
	"strconv"
	"strings"
	"time"
)

func GetFile(c *gin.Context) (*models.InputFile, error) {
	var buf bytes.Buffer
	file, header, _ := c.Request.FormFile("file")
	defer file.Close()

	t := time.Now()
	formatted := fmt.Sprintf("_%d%02d%02d_%02d%02d%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())

	filename := strings.Split(header.Filename, ".")[0] + formatted
	fileFormat := strings.Split(header.Filename, ".")[1]

	if fileFormat != "csv" {
		return nil, errors.New("Invalid file format")
	}

	logger.Infof("File name %s", filename)

	io.Copy(&buf, file)
	inputFile := new(models.InputFile)
	inputFile.Name = filename
	inputFile.Data = buf

	return inputFile, nil
}

func CsvToTicket(file *models.InputFile) ([]*models.Ticket, error) {
	reader := bytes.NewReader(file.Data.Bytes())
	csv := csv.NewReader(reader)

	ticket := []*models.Ticket{}
	lineHeader := true

	for {

		record, err := csv.Read()
		if err == io.EOF {
			break
		}

		if lineHeader {
			lineHeader = false
			continue
		}

		ticketId, err := strconv.ParseInt(record[0], 10, 64)
		if err != nil {
			return nil, err
		}

		siteId := record[1]
		var solutionId int64 = 0
		if len(record) > 2 {
			var value int64
			value, err = strconv.ParseInt(record[2], 10, 32)
			if err != nil {
				return nil, err
			}
			solutionId = value
		}
		ticket = append(ticket, &models.Ticket{
			Id:         ticketId,
			SiteId:     siteId,
			SolutionId: solutionId,
			FileId:     file.Id,
		})
	}
	return ticket, nil
}

func CsvToOpen(file *models.InputFile) ([]*models.Ticket, error) {
	reader := bytes.NewReader(file.Data.Bytes())
	csv := csv.NewReader(reader)

	ticket := []*models.Ticket{}
	lineHeader := true

	for {

		record, err := csv.Read()
		if err == io.EOF {
			break
		}

		if lineHeader {
			lineHeader = false
			continue
		}

		ticketId, err := strconv.ParseInt(record[0], 10, 64)
		if err != nil {
			return nil, err
		}

		ticket = append(ticket, &models.Ticket{
			Id:     ticketId,
			FileId: file.Id,
		})
	}
	return ticket, nil
}
