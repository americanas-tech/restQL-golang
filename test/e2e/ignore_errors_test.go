package e2e

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/b2wdigital/restQL-golang/v6/test"
)

func TestReturningMaximumStatusCodeOfStatementsWithoutIgnoreErrors(t *testing.T) {
	query := `
from planets as yavin
	with
		id = 1

from planets as tatooine
	with
		id = 2

from planets as alderaan
	with
		id = 3
	ignore-errors
`

	yavin := `
{
	"name": "Yavin",
	"rotation_period": 24.5,
	"orbital_period": "4818",
	"diameter": "10200",
	"climate": "temperate, tropical",
	"gravity": "1 standard",
	"terrain": { "north": "jungle", "south": "rainforests" },
	"surface_water": "8",
	"population": "1000",
	"residents": ["john", "janne"],
	"films": [1]
}
`

	tatooine := `
{
	"name": "Tatooine",
	"rotation_period": "23",
	"orbital_period": "304",
	"diameter": "10465",
	"climate": "arid",
	"gravity": "1 standard",
	"terrain": "desert",
	"surface_water": "1",
	"population": "200000",
	"residents": ["william", "scarlet"],
	"films": [1, 3, 4, 5, 6]
}
`

	alderaan := `{}`

	expectedResponse := fmt.Sprintf(`
	{
		"yavin": {
			"details": {
				"success": false,
				"status": 404,
				"metadata": {}
			},
			"result": %s
		},
		"tatooine": {
			"details": {
				"success": false,
				"status": 500,
				"metadata": {}
			},
			"result": %s
		},
		"alderaan": {
			"details": {
				"success": false,
				"status": 503,
				"metadata": {"ignore-errors":"ignore"}
			},
			"result": %s
		}
	}`, yavin, tatooine, alderaan)

	mockServer := test.NewMockServer(mockPort)
	defer mockServer.Teardown()

	mockServer.Mux().HandleFunc("/api/planets/1", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		io.WriteString(w, yavin)
	})
	mockServer.Mux().HandleFunc("/api/planets/2", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		io.WriteString(w, tatooine)
	})
	mockServer.Mux().HandleFunc("/api/planets/3", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(503)
		io.WriteString(w, alderaan)
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

func TestReturning2xxWhenStatementWithIgnoreErrorsFail(t *testing.T) {
	query := `
from planets
	with 
		id = 1

from people
	with
		name = planets.leader
	ignore-errors
`

	planetResponse := `
{
	"id": 1,
	"name": "Yavin",
	"rotation_period": 24.5,
	"orbital_period": "4818",
	"diameter": "10200",
	"climate": "temperate, tropical",
	"gravity": "1 standard",
	"terrain": { "north": "jungle", "south": "rainforests" },
	"surface_water": "8",
	"population": "1000",
	"leader": "Yavin King",
	"residents": [1, 2],
	"films": [1]
}
`

	peopleResponse := `{}`

	expectedResponse := fmt.Sprintf(`
	{
		"planets": {
			"details": {
				"success": true,
				"status": 200,
				"metadata": {}
			},
			"result": %s 
		},
		"people": {
			"details": {
				"success": false,
				"status": 404,
				"metadata": {"ignore-errors":"ignore"}
			},
			"result": %s 
		}
	}`, planetResponse, peopleResponse)

	mockServer := test.NewMockServer(mockPort)
	defer mockServer.Teardown()

	mockServer.Mux().HandleFunc("/api/planets/1", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		io.WriteString(w, planetResponse)
	})
	mockServer.Mux().HandleFunc("/api/people/", func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()

		nameParam := params["name"]
		if len(nameParam) == 0 {
			t.Error("got empty nameParam param", nameParam)
			return
		}

		name, err := url.QueryUnescape(nameParam[0])
		test.VerifyError(t, err)
		test.Equal(t, name, "Yavin King")

		w.WriteHeader(404)
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
