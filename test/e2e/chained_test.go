package e2e

import (
	"encoding/json"
	"fmt"
	"github.com/b2wdigital/restQL-golang/test"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestWithQualifierChainedParamReturningXmlOnFromStatement(t *testing.T) {
	query := `
from planets
	with 
		id = "Yavin"

from people
	with
		name = planets.leader
`

	planetResponse := `
<planet>
   <climate>temperate, tropical</climate>
   <diameter>10200</diameter>
   <films>
      <element>1</element>
   </films>
   <gravity>1 standard</gravity>
   <leader>Yavin King</leader>
   <name>Yavin</name>
   <orbital_period>4818</orbital_period>
   <population>1000</population>
   <residents>
      <element>john</element>
      <element>janne</element>
   </residents>
   <rotation_period>24.5</rotation_period>
   <surface_water>8</surface_water>
   <terrain>
      <north>jungle</north>
      <south>rainforests</south>
   </terrain>
</planet>
`

	peopleResponse := `The request was skipped due to missing { :name } param value`

	expectedResponse := fmt.Sprintf(`
	{
		"planets": {
			"details": {
				"success": false,
				"status": 0,
				"metadata": {}
			},
			"result": "failed to unmarshal response body: outh>rainforests</south>\n   </terrain>\n</planet>\n" 
		},
		"people": {
			"details": {
				"success": false,
				"status": 400,
				"metadata": {}
			},
			"result": "The request was skipped due to missing { :name } param value"
		}
	}`)

	mockServer := test.NewMockServer(mockPort)
	defer mockServer.Teardown()

	mockServer.Mux().HandleFunc("/api/planets/Yavin", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		io.WriteString(w, planetResponse)
	})
	mockServer.Mux().HandleFunc("/api/people/", func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()

		name, err := url.QueryUnescape(params["name"][0])
		test.VerifyError(t, err)
		test.Equal(t, name, "Yavin King")

		w.WriteHeader(200)
		io.WriteString(w, peopleResponse)
	})
	mockServer.Start()

	response, err := httpClient.Post(adHocQueryUrl, "text/plain", strings.NewReader(query))
	test.VerifyError(t, err)
	defer response.Body.Close()

	test.Equal(t, response.StatusCode, 500)

	var body map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&body)
	test.VerifyError(t, err)

	test.Equal(t, body, test.Unmarshal(expectedResponse))
}
