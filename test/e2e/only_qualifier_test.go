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

func TestOnlyQualifierOnFromStatement(t *testing.T) {
	query := `
from planets
	with id = 1
		only
			id
			name
			gravity
			residents -> matches("^j")
`

	planetResponse := `
{
	"id": 1,
	"name": "Yavin IV",
	"rotation_period": 24.5,
	"orbital_period": "4818",
	"diameter": "10200",
	"climate": "temperate, tropical",
	"gravity": "1 standard",
	"terrain": { "north": "jungle", "south": "rainforests" },
	"surface_water": "8",
	"population": "1000",
	"residents": ["john", "janne", "kyle"],
	"films": [1]
}
`

	expectedResponse := `
	{
		"planets": {
			"details": {
				"success": true,
				"status": 200,
				"metadata": {}
			},
			"result": {
				"id": 1,
				"name": "Yavin IV",
				"gravity": "1 standard",
				"residents": ["john", "janne"]
			}
		}
	}`

	mockServer := test.NewMockServer(mockPort)
	defer mockServer.Teardown()

	mockServer.Mux().HandleFunc("/api/planets/1", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		io.WriteString(w, planetResponse)
	})
	mockServer.Start()

	response, err := httpClient.Post(adHocQueryUrl, "text/plain", strings.NewReader(query))
	test.VerifyError(t, err)
	defer response.Body.Close()

	test.Equal(t, response.StatusCode, 200)

	var body map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&body)
	test.VerifyError(t, err)

	test.Equal(t, body, test.Unmarshal(expectedResponse))
}

func TestHiddenQualifierOnFromStatement(t *testing.T) {
	query := `
from planets
	with 
		id = 1
	hidden

from people
	with
		name = planets.leader
`

	planetResponse := `
{
	"id": 1,
	"name": "Yavin IV",
	"rotation_period": 24.5,
	"orbital_period": "4818",
	"diameter": "10200",
	"climate": "temperate, tropical",
	"gravity": "1 standard",
	"terrain": { "north": "jungle", "south": "rainforests" },
	"surface_water": "8",
	"population": "1000",
	"leader": "Yavin King",
	"residents": ["john", "janne", "kyle"],
	"films": [1]
}
`

	peopleResponse := `
{
	"name": "Yavin King",
	"height": "172",
	"mass": "77",
	"hair_color": "blond",
	"skin_color": "fair",
	"eye_color": "blue",
	"birth_year": "19BBY",
	"gender": "male",
	"homeworld": 1,
	"films": [1, 2, 3, 6]
}
`

	expectedResponse := fmt.Sprintf(`
	{
		"people": {
			"details": {
				"success": true,
				"status": 200,
				"metadata": {}
			},
			"result": %s
		}
	}`, peopleResponse)

	mockServer := test.NewMockServer(mockPort)
	defer mockServer.Teardown()

	mockServer.Mux().HandleFunc("/api/planets/1", func(w http.ResponseWriter, r *http.Request) {
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

	test.Equal(t, response.StatusCode, 200)

	var body map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&body)
	test.VerifyError(t, err)

	test.Equal(t, body, test.Unmarshal(expectedResponse))
}
