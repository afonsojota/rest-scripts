package solution

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/afonsojota/cleancoder_rest-scripts/cmd/api/models"
	"github.com/afonsojota/cleancoder_rest-scripts/cmd/api/utils"
	"time"

	"github.com/afonsojota/go-afonsojota-toolkit/restful/rest"
	"github.com/spf13/viper"
)

//go:generate mockgen -source=./solution.go -destination=./mocks/solution_mock.go
type Gateway interface {
	GetSolution(id int64) (*models.Solution, error)
}

type solutionGateway struct {
	rb rest.RequestBuilder
}

func NewSolutionGateway() Gateway {
	return &solutionGateway{
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

func (gtw *solutionGateway) GetSolution(id int64) (*models.Solution, error) {

	solution := new(models.Solution)

	uri := fmt.Sprintf("%s/%d", viper.GetString("solutions.path"), id)
	resp := gtw.rb.Get(uri, rest.Headers(utils.GetDefaultHeaders()), rest.Context(context.Background()))

	if resp.Err != nil || !(int(resp.StatusCode/100) == 2) {
		return nil, buildErrorMessage(resp)
	}

	err := json.Unmarshal(resp.Bytes(), &solution)

	return solution, err
}

func buildErrorMessage(response *rest.Response) error {
	return errors.New(response.String())
}
