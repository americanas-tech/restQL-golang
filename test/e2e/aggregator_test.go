package e2e

import (
	"encoding/json"
	"github.com/b2wdigital/restQL-golang/v4/test"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestInOperatorOnFromStatement(t *testing.T) {
	query := `
from planets
	with 
		id = "Yavin"

from planets as sats in planets.sats
	with 
		id = planets.satellites
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
	"satellites": ["Havoc"],
	"films": [1]
}
`

	satellitiesResponse := `
{
	"name": "Havoc",
	"rotation_period": 24.5,
	"orbital_period": "4818",
	"diameter": "10200",
	"climate": "frost",
	"gravity": "1 standard",
	"terrain": { "north": "desert", "south": "desert" },
	"surface_water": "8",
	"population": "0"
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
				"satellites": ["Havoc"],
				"sats": [{
					"name": "Havoc",
					"rotation_period": 24.5,
					"orbital_period": "4818",
					"diameter": "10200",
					"climate": "frost",
					"gravity": "1 standard",
					"terrain": { "north": "desert", "south": "desert" },
					"surface_water": "8",
					"population": "0"
				}],
				"films": [1]
			}
		},
		"sats":{
			"details": [{
				"success": true,
				"status": 200,
				"metadata": {}
			}]
		}
	}`

	mockServer := test.NewMockServer(mockPort)
	defer mockServer.Teardown()

	mockServer.Mux().HandleFunc("/api/planets/Yavin", func(w http.ResponseWriter, r *http.Request) {
		test.Equal(t, r.Method, http.MethodGet)

		w.WriteHeader(200)
		io.WriteString(w, planetResponse)
	})
	mockServer.Mux().HandleFunc("/api/planets/Havoc", func(w http.ResponseWriter, r *http.Request) {
		test.Equal(t, r.Method, http.MethodGet)

		w.WriteHeader(200)
		io.WriteString(w, satellitiesResponse)
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

func TestInOperatorOnToStatement(t *testing.T) {
	query := `
from planets
	with 
		id = "Yavin"

to planets as sats in planets.sats
	with 
		id = planets.satellites
		main_body = planets.name
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
	"satellites": ["Havoc"],
	"films": [1]
}
`

	satellitiesResponse := `
{
	"name": "Havoc",
	"rotation_period": 24.5,
	"orbital_period": "4818",
	"diameter": "10200",
	"climate": "frost",
	"gravity": "1 standard",
	"terrain": { "north": "desert", "south": "desert" },
	"surface_water": "8",
	"population": "0"
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
				"satellites": ["Havoc"],
				"sats": [{
					"name": "Havoc",
					"rotation_period": 24.5,
					"orbital_period": "4818",
					"diameter": "10200",
					"climate": "frost",
					"gravity": "1 standard",
					"terrain": { "north": "desert", "south": "desert" },
					"surface_water": "8",
					"population": "0"
				}],
				"films": [1]
			}
		},
		"sats":{
			"details": [{
				"success": true,
				"status": 200,
				"metadata": {}
			}]
		}
	}`

	mockServer := test.NewMockServer(mockPort)
	defer mockServer.Teardown()

	mockServer.Mux().HandleFunc("/api/planets/Yavin", func(w http.ResponseWriter, r *http.Request) {
		test.Equal(t, r.Method, http.MethodGet)

		w.WriteHeader(200)
		io.WriteString(w, planetResponse)
	})
	mockServer.Mux().HandleFunc("/api/planets/Havoc", func(w http.ResponseWriter, r *http.Request) {
		test.Equal(t, r.Method, http.MethodPost)

		b, err := ioutil.ReadAll(r.Body)
		test.VerifyError(t, err)

		test.NotEqual(t, string(b), "")
		body := test.Unmarshal(string(b))

		test.Equal(t, body, test.Unmarshal(`{"main_body": "Yavin"}`))

		w.WriteHeader(200)
		io.WriteString(w, satellitiesResponse)
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

func TestInOperatorOnIntoStatement(t *testing.T) {
	query := `
from planets
	with 
		id = "Yavin"

into planets as sats in planets.sats
	with 
		id = planets.satellites
		main_body = planets.name
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
	"satellites": ["Havoc"],
	"films": [1]
}
`

	satellitiesResponse := `
{
	"name": "Havoc",
	"rotation_period": 24.5,
	"orbital_period": "4818",
	"diameter": "10200",
	"climate": "frost",
	"gravity": "1 standard",
	"terrain": { "north": "desert", "south": "desert" },
	"surface_water": "8",
	"population": "0"
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
				"satellites": ["Havoc"],
				"sats": [{
					"name": "Havoc",
					"rotation_period": 24.5,
					"orbital_period": "4818",
					"diameter": "10200",
					"climate": "frost",
					"gravity": "1 standard",
					"terrain": { "north": "desert", "south": "desert" },
					"surface_water": "8",
					"population": "0"
				}],
				"films": [1]
			}
		},
		"sats":{
			"details": [{
				"success": true,
				"status": 200,
				"metadata": {}
			}]
		}
	}`

	mockServer := test.NewMockServer(mockPort)
	defer mockServer.Teardown()

	mockServer.Mux().HandleFunc("/api/planets/Yavin", func(w http.ResponseWriter, r *http.Request) {
		test.Equal(t, r.Method, http.MethodGet)

		w.WriteHeader(200)
		io.WriteString(w, planetResponse)
	})
	mockServer.Mux().HandleFunc("/api/planets/Havoc", func(w http.ResponseWriter, r *http.Request) {
		test.Equal(t, r.Method, http.MethodPut)

		b, err := ioutil.ReadAll(r.Body)
		test.VerifyError(t, err)

		test.NotEqual(t, string(b), "")
		body := test.Unmarshal(string(b))

		test.Equal(t, body, test.Unmarshal(`{"main_body": "Yavin"}`))

		w.WriteHeader(200)
		io.WriteString(w, satellitiesResponse)
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

func TestInOperatorOnUpdateStatement(t *testing.T) {
	query := `
from planets
	with 
		id = "Yavin"

update planets as sats in planets.sats
	with 
		id = planets.satellites
		main_body = planets.name
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
	"satellites": ["Havoc"],
	"films": [1]
}
`

	satellitiesResponse := `
{
	"name": "Havoc",
	"rotation_period": 24.5,
	"orbital_period": "4818",
	"diameter": "10200",
	"climate": "frost",
	"gravity": "1 standard",
	"terrain": { "north": "desert", "south": "desert" },
	"surface_water": "8",
	"population": "0"
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
				"satellites": ["Havoc"],
				"sats": [{
					"name": "Havoc",
					"rotation_period": 24.5,
					"orbital_period": "4818",
					"diameter": "10200",
					"climate": "frost",
					"gravity": "1 standard",
					"terrain": { "north": "desert", "south": "desert" },
					"surface_water": "8",
					"population": "0"
				}],
				"films": [1]
			}
		},
		"sats":{
			"details": [{
				"success": true,
				"status": 200,
				"metadata": {}
			}]
		}
	}`

	mockServer := test.NewMockServer(mockPort)
	defer mockServer.Teardown()

	mockServer.Mux().HandleFunc("/api/planets/Yavin", func(w http.ResponseWriter, r *http.Request) {
		test.Equal(t, r.Method, http.MethodGet)

		w.WriteHeader(200)
		io.WriteString(w, planetResponse)
	})
	mockServer.Mux().HandleFunc("/api/planets/Havoc", func(w http.ResponseWriter, r *http.Request) {
		test.Equal(t, r.Method, http.MethodPatch)

		b, err := ioutil.ReadAll(r.Body)
		test.VerifyError(t, err)

		test.NotEqual(t, string(b), "")
		body := test.Unmarshal(string(b))

		test.Equal(t, body, test.Unmarshal(`{"main_body": "Yavin"}`))

		w.WriteHeader(200)
		io.WriteString(w, satellitiesResponse)
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

func TestInOperatorOnDeleteStatement(t *testing.T) {
	query := `
from planets
	with 
		id = "Yavin"

delete planets as sats in planets.sats
	with 
		id = planets.satellites
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
	"satellites": ["Havoc"],
	"films": [1]
}
`

	satellitiesResponse := `
{
	"name": "Havoc",
	"rotation_period": 24.5,
	"orbital_period": "4818",
	"diameter": "10200",
	"climate": "frost",
	"gravity": "1 standard",
	"terrain": { "north": "desert", "south": "desert" },
	"surface_water": "8",
	"population": "0"
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
				"satellites": ["Havoc"],
				"sats": [{
					"name": "Havoc",
					"rotation_period": 24.5,
					"orbital_period": "4818",
					"diameter": "10200",
					"climate": "frost",
					"gravity": "1 standard",
					"terrain": { "north": "desert", "south": "desert" },
					"surface_water": "8",
					"population": "0"
				}],
				"films": [1]
			}
		},
		"sats":{
			"details": [{
				"success": true,
				"status": 200,
				"metadata": {}
			}]
		}
	}`

	mockServer := test.NewMockServer(mockPort)
	defer mockServer.Teardown()

	mockServer.Mux().HandleFunc("/api/planets/Yavin", func(w http.ResponseWriter, r *http.Request) {
		test.Equal(t, r.Method, http.MethodGet)

		w.WriteHeader(200)
		io.WriteString(w, planetResponse)
	})
	mockServer.Mux().HandleFunc("/api/planets/Havoc", func(w http.ResponseWriter, r *http.Request) {
		test.Equal(t, r.Method, http.MethodDelete)

		w.WriteHeader(200)
		io.WriteString(w, satellitiesResponse)
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

func TestInOperatorAndHiddenFilterOnFromStatement(t *testing.T) {
	query := `
from planets
	with 
		id = "Yavin"

from planets as planetSats in planets.satellitesDetails
	with 
		id = planets.satellites
	hidden
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
	"satellites": ["Havoc"],
	"films": [1]
}
`

	satellitiesResponse := `
{
	"name": "Havoc",
	"rotation_period": 24.5,
	"orbital_period": "4818",
	"diameter": "10200",
	"climate": "frost",
	"gravity": "1 standard",
	"terrain": { "north": "desert", "south": "desert" },
	"surface_water": "8",
	"population": "0"
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
				"satellites": ["Havoc"],
				"satellitesDetails": [{
					"name": "Havoc",
					"rotation_period": 24.5,
					"orbital_period": "4818",
					"diameter": "10200",
					"climate": "frost",
					"gravity": "1 standard",
					"terrain": { "north": "desert", "south": "desert" },
					"surface_water": "8",
					"population": "0"
				}],
				"films": [1]
			}
		}
	}`

	mockServer := test.NewMockServer(mockPort)
	defer mockServer.Teardown()

	mockServer.Mux().HandleFunc("/api/planets/Yavin", func(w http.ResponseWriter, r *http.Request) {
		test.Equal(t, r.Method, http.MethodGet)

		w.WriteHeader(200)
		io.WriteString(w, planetResponse)
	})
	mockServer.Mux().HandleFunc("/api/planets/Havoc", func(w http.ResponseWriter, r *http.Request) {
		test.Equal(t, r.Method, http.MethodGet)

		w.WriteHeader(200)
		io.WriteString(w, satellitiesResponse)
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

func TestInOperatorAndOnlyFilterOnFromStatement(t *testing.T) {
	query := `
from planets
	with 
		id = "Yavin"

from planets as sats in planets.satellitesDetails
	with 
		id = planets.satellites
	only
		name
		rotation_period
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
	"satellites": ["Havoc"],
	"films": [1]
}
`

	satellitiesResponse := `
{
	"name": "Havoc",
	"rotation_period": 24.5,
	"orbital_period": "4818",
	"diameter": "10200",
	"climate": "frost",
	"gravity": "1 standard",
	"terrain": { "north": "desert", "south": "desert" },
	"surface_water": "8",
	"population": "0"
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
				"satellites": ["Havoc"],
				"satellitesDetails": [{
					"name": "Havoc",
					"rotation_period": 24.5
				}],
				"films": [1]
			}
		},
		"sats":{
			"details": [{
				"success": true,
				"status": 200,
				"metadata": {}
			}]
		}
	}`

	mockServer := test.NewMockServer(mockPort)
	defer mockServer.Teardown()

	mockServer.Mux().HandleFunc("/api/planets/Yavin", func(w http.ResponseWriter, r *http.Request) {
		test.Equal(t, r.Method, http.MethodGet)

		w.WriteHeader(200)
		io.WriteString(w, planetResponse)
	})
	mockServer.Mux().HandleFunc("/api/planets/Havoc", func(w http.ResponseWriter, r *http.Request) {
		test.Equal(t, r.Method, http.MethodGet)

		w.WriteHeader(200)
		io.WriteString(w, satellitiesResponse)
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
