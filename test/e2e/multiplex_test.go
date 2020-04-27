package e2e

import (
	"encoding/json"
	"fmt"
	"github.com/b2wdigital/restQL-golang/test"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestMultiplexingWithStaticParameterOnFromStatement(t *testing.T) {
	query := `
from planets
	with
		id = [1, 2]
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

	expectedResponse := fmt.Sprintf(`
	{
		"planets": {
			"details": [
				{
					"success": true,
					"status": 200,
					"metadata": {}
				},
				{
					"success": true,
					"status": 200,
					"metadata": {}
				}
			],
			"result": [%s,%s]
		}
	}`, yavin, tatooine)

	mockServer := test.NewMockServer(mockPort)
	defer mockServer.Teardown()

	mockServer.Mux().HandleFunc("/api/planets/1", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		io.WriteString(w, yavin)
	})
	mockServer.Mux().HandleFunc("/api/planets/2", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		io.WriteString(w, tatooine)
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

func TestMultiplexingWithChainedParameterOnFromStatement(t *testing.T) {
	query := `
from planets
	with 
		id = 1

from people
	with
		id = planets.residents
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

	john := `
{
	"id": 1,
	"name": "john",
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

	janne := `
{
	"id": 2,
	"name": "janne",
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
		"planets": {
			"details": {
				"success": true,
				"status": 200,
				"metadata": {}
			},
			"result": %s 
		},
		"people": {
			"details": [
				{
					"success": true,
					"status": 200,
					"metadata": {}
				},
				{
					"success": true,
					"status": 200,
					"metadata": {}
				}
			],
			"result": [%s, %s]
		}
	}`, planetResponse, john, janne)

	mockServer := test.NewMockServer(mockPort)
	defer mockServer.Teardown()

	mockServer.Mux().HandleFunc("/api/planets/1", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		io.WriteString(w, planetResponse)
	})
	mockServer.Mux().HandleFunc("/api/people/1", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		io.WriteString(w, john)
	})
	mockServer.Mux().HandleFunc("/api/people/2", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		io.WriteString(w, janne)
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

func TestMultiplexingWithStaticObjectParameterOnFromStatement(t *testing.T) {
	query := `
from planets
	with 
		id = 1

from people
	with
		profile = { name: planets.residents }
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
	"residents": ["john", "janne"],
	"films": [1]
}
`

	john := `
{
	"id": 1,
	"name": "john",
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

	janne := `
{
	"id": 2,
	"name": "janne",
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
		"planets": {
			"details": {
				"success": true,
				"status": 200,
				"metadata": {}
			},
			"result": %s 
		},
		"people": {
			"details": [ 
				{
					"success": true,
					"status": 200,
					"metadata": {}
				},
				{
					"success": true,
					"status": 200,
					"metadata": {}
				}
			],
			"result": [%s, %s]
		}
	}`, planetResponse, john, janne)

	mockServer := test.NewMockServer(mockPort)
	defer mockServer.Teardown()

	mockServer.Mux().HandleFunc("/api/planets/1", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		io.WriteString(w, planetResponse)
	})
	mockServer.Mux().HandleFunc("/api/people/", func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		profileData := params["profile"][0]

		var profile map[string]string
		err := json.Unmarshal([]byte(profileData), &profile)
		test.VerifyError(t, err)

		if profile["name"] == "john" {
			w.WriteHeader(200)
			io.WriteString(w, john)
			return
		}

		if profile["name"] == "janne" {
			w.WriteHeader(200)
			io.WriteString(w, janne)
			return
		}
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

func TestFlatteningWithStaticParameterOnFromStatement(t *testing.T) {
	query := `
from planets
	with
		residents = ["john", "janne"] -> flatten
`

	planetResponse := `
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

	expectedResponse := fmt.Sprintf(`
	{
		"planets": {
			"details": {
				"success": true,
				"status": 200,
				"metadata": {}
			},
			"result": %s
		}
	}`, planetResponse)

	mockServer := test.NewMockServer(mockPort)
	defer mockServer.Teardown()

	mockServer.Mux().HandleFunc("/api/planets/", func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()

		test.Equal(t, params["residents"], []string{"john", "janne"})

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

func TestFlatteningWithChainedParameterOnFromStatement(t *testing.T) {
	query := `
from planets
	with 
		id = 1

from people
	with
		name = planets.residents -> flatten
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
	"residents": ["john", "janne"],
	"films": [1]
}
`

	peopleResponse := `[
{
	"id": 1,
	"name": "john",
	"height": "172",
	"mass": "77",
	"hair_color": "blond",
	"skin_color": "fair",
	"eye_color": "blue",
	"birth_year": "19BBY",
	"gender": "male",
	"homeworld": 1,
	"films": [1, 2, 3, 6]
},
{
	"id": 2,
	"name": "janne",
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
]
`

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
				"success": true,
				"status": 200,
				"metadata": {}
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

		test.Equal(t, params["name"], []string{"john", "janne"})

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

func TestFlatteningWithStaticObjectParameterOnFromStatement(t *testing.T) {
	query := `
from planets
	with 
		id = 1

from people
	with
		profile = { name: planets.residents } -> flatten
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
	"residents": ["john", "janne"],
	"films": [1]
}
`

	peopleResponse := ` [
{
	"id": 1,
	"name": "john",
	"height": "172",
	"mass": "77",
	"hair_color": "blond",
	"skin_color": "fair",
	"eye_color": "blue",
	"birth_year": "19BBY",
	"gender": "male",
	"homeworld": 1,
	"films": [1, 2, 3, 6]
},
{
	"id": 2,
	"name": "janne",
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
]
`

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
			"details":{
				"success": true,
				"status": 200,
				"metadata": {}
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
		profileData := params["profile"][0]

		test.Equal(t, profileData, `{"name":["john","janne"]}`)

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
