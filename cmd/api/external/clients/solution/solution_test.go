package solution

import (
	_ "github.com/afonsojota/cleancoder_rest-scripts/cmd/api/config"
	"github.com/afonsojota/cleancoder_rest-scripts/cmd/api/utils"
	"github.com/afonsojota/go-afonsojota-toolkit/restful/rest"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

//func TestSolutionHtmlTemplate(t *testing.T) {
//	bytes, err := ioutil.ReadFile("resources/solution.json")
//	if !assert.Nil(t, err, "error reading mock file") {
//		return
//	}
//
//	rest.StartMockupServer()
//	defer rest.StopMockupServer()
//
//	_ = rest.AddMockups(&rest.Mock{
//		URL:          "http://internal-api.afonsojota.com/solutions/47041",
//		HTTPMethod:   http.MethodGet,
//		ReqHeaders:   utils.GetDefaultHeaders(),
//		RespHTTPCode: http.StatusOK,
//		RespBody:     string(bytes),
//	})
//
//	solutionId := int64(47041)
//
//	gateway := NewSolutionGateway()
//	solution, err := gateway.GetSolution(solutionId)
//
//	if !assert.NotNil(t, err, "error from dao") {
//		return
//	}
//
//	template := solution.Versions[0].Contents[0].Body.Template
//	assert.Equal(t, solutionId, solution.Id)
//	assert.Equal(t, "<p><span style=\"font-weight: 400;\">Hola.</span></p>\n<p><span style=\"font-weight: 400;\">Te pedimos disculpas por la demora en responder.&nbsp;</span></p>\n<p><span style=\"font-weight: 400;\">Revisamos la operaci&oacute;n y vimos que el env&iacute;o ya fue entregado a tu comprador.&nbsp;</span></p>\n<p><span style=\"font-weight: 400;\">Desde la secci&oacute;n de &ldquo;Actividad&rdquo; en tu cuenta de Mercado Pago y haciendo clic sobre el cobro, podr&aacute;s ver cu&aacute;ndo tendr&aacute;s el dinero liquidado y disponible para usar.</span></p>\n<p><span style=\"font-weight: 400;\">&iexcl;Hasta pronto!</span></p>", template)
//}

func TestSolutionHtmlTemplateFailError500(t *testing.T) {
	bytes, err := ioutil.ReadFile("/resources/solution.json")
	if !assert.Nil(t, err, "error reading mock file") {
		return
	}

	rest.StartMockupServer()
	defer rest.StopMockupServer()

	_ = rest.AddMockups(&rest.Mock{
		URL:          "http://internal-api.afonsojota.com/solutions/47041",
		HTTPMethod:   http.MethodGet,
		ReqHeaders:   utils.GetDefaultHeaders(),
		RespHTTPCode: http.StatusInternalServerError,
		RespBody:     string(bytes),
	})

	solutionId := int64(47041)

	gateway := NewSolutionGateway()
	_, err = gateway.GetSolution(solutionId)

	assert.NotNil(t, err)
}
