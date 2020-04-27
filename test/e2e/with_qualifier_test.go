package e2e

import (
	"encoding/json"
	"fmt"
	"github.com/b2wdigital/restQL-golang/test"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestWithQualifierOnFromStatement(t *testing.T) {
	query := `
from planets
	with 
		name = "Yavin IV"
		population = 1000
		residents = ["john", "janne"] -> flatten
		rotation_period = 24.5
		terrain = { "north": "jungle", "south": "rainforests" }
`

	planetResponse := `
[{
	"name": "Yavin IV",
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
}]
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

		test.Equal(t, params["name"][0], "Yavin IV")
		test.Equal(t, params["population"][0], "1000")
		test.Equal(t, params["rotation_period"][0], "24.5")
		test.Equal(t, params["terrain"][0], `{"north":"jungle","south":"rainforests"}`)

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

func TestWithQualifierOnPostStatement(t *testing.T) {
	query := `
to planets
	with 
		name = "Yavin IV"
		population = 1000
		rotation_period = 24.5
		terrain = { "north": "jungle", "south": "rainforests" }
		residents = ["john", "janne"] -> flatten
`

	planetResponse := `
{
	"name": "Yavin IV",
	"rotation_period": 24.5,
	"orbital_period": "4818",
	"diameter": "10200",
	"climate": "temperate, tropical",
	"gravity": "1 standard",
	"terrain": { "north": "jungle", "south": "rainforests" },
	"surface_water": "8",
	"population": 1000,
	"residents": ["john", "janne"],
	"films": [1]
}
`

	expectedResponse := fmt.Sprintf(`
	{
		"planets": {
			"details": {
				"success": true,
				"status": 201,
				"metadata": {}
			},
			"result": %s 
		}
	}`, planetResponse)

	mockServer := test.NewMockServer(mockPort)
	defer mockServer.Teardown()

	mockServer.Mux().HandleFunc("/api/planets/", func(w http.ResponseWriter, r *http.Request) {
		test.Equal(t, r.Method, http.MethodPost)

		b, err := ioutil.ReadAll(r.Body)
		test.VerifyError(t, err)

		body := string(b)
		test.NotEqual(t, body, "")

		expectedBody := `
{
	"name": "Yavin IV",
	"population": 1000,
	"residents": ["john", "janne"],
	"rotation_period": 24.5,
	"terrain": { "north": "jungle", "south": "rainforests" }
}
`

		test.Equal(t, test.Unmarshal(expectedBody), test.Unmarshal(body))

		w.WriteHeader(201)
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

func TestWithQualifierOnPutStatement(t *testing.T) {
	query := `
into planets
	with 
		name = "Yavin IV"
		population = 1000
		rotation_period = 24.5
		terrain = { "north": "jungle", "south": "rainforests" }
		residents = ["john", "janne"] -> flatten
`

	planetResponse := `
{
	"name": "Yavin IV",
	"rotation_period": 24.5,
	"orbital_period": "4818",
	"diameter": "10200",
	"climate": "temperate, tropical",
	"gravity": "1 standard",
	"terrain": { "north": "jungle", "south": "rainforests" },
	"surface_water": "8",
	"population": 1000,
	"residents": ["john", "janne"],
	"films": [1]
}
`

	expectedResponse := fmt.Sprintf(`
	{
		"planets": {
			"details": {
				"success": true,
				"status": 201,
				"metadata": {}
			},
			"result": %s 
		}
	}`, planetResponse)

	mockServer := test.NewMockServer(mockPort)
	defer mockServer.Teardown()

	mockServer.Mux().HandleFunc("/api/planets/", func(w http.ResponseWriter, r *http.Request) {
		test.Equal(t, r.Method, http.MethodPut)

		b, err := ioutil.ReadAll(r.Body)
		test.VerifyError(t, err)

		body := string(b)
		test.NotEqual(t, body, "")

		expectedBody := `
{
	"name": "Yavin IV",
	"population": 1000,
	"residents": ["john", "janne"],
	"rotation_period": 24.5,
	"terrain": { "north": "jungle", "south": "rainforests" }
}
`

		test.Equal(t, test.Unmarshal(expectedBody), test.Unmarshal(body))

		w.WriteHeader(201)
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

func TestWithQualifierOnPatchStatement(t *testing.T) {
	query := `
update planets
	with 
		name = "Yavin IV"
		population = 1000
		rotation_period = 24.5
		terrain = { "north": "jungle", "south": "rainforests" }
		residents = ["john", "janne"] -> flatten
`

	planetResponse := `
{
	"name": "Yavin IV",
	"rotation_period": 24.5,
	"orbital_period": "4818",
	"diameter": "10200",
	"climate": "temperate, tropical",
	"gravity": "1 standard",
	"terrain": { "north": "jungle", "south": "rainforests" },
	"surface_water": "8",
	"population": 1000,
	"residents": ["john", "janne"],
	"films": [1]
}
`

	expectedResponse := fmt.Sprintf(`
	{
		"planets": {
			"details": {
				"success": true,
				"status": 201,
				"metadata": {}
			},
			"result": %s 
		}
	}`, planetResponse)

	mockServer := test.NewMockServer(mockPort)
	defer mockServer.Teardown()

	mockServer.Mux().HandleFunc("/api/planets/", func(w http.ResponseWriter, r *http.Request) {
		test.Equal(t, r.Method, http.MethodPatch)

		b, err := ioutil.ReadAll(r.Body)
		test.VerifyError(t, err)

		body := string(b)
		test.NotEqual(t, body, "")

		expectedBody := `
{
	"name": "Yavin IV",
	"population": 1000,
	"residents": ["john", "janne"],
	"rotation_period": 24.5,
	"terrain": { "north": "jungle", "south": "rainforests" }
}
`

		test.Equal(t, test.Unmarshal(expectedBody), test.Unmarshal(body))

		w.WriteHeader(201)
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

func TestWithQualifierOnDeleteStatement(t *testing.T) {
	query := `
delete planets
	with 
		name = "Yavin IV"
		population = 1000
		residents = ["john", "janne"] -> flatten
		rotation_period = 24.5
		terrain = { "north": "jungle", "south": "rainforests" }
`

	planetResponse := `{}`

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

		test.Equal(t, params["name"][0], "Yavin IV")
		test.Equal(t, params["population"][0], "1000")
		test.Equal(t, params["rotation_period"][0], "24.5")
		test.Equal(t, params["terrain"][0], `{"north":"jungle","south":"rainforests"}`)

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

func TestWithQualifierVariableParamOnFromStatement(t *testing.T) {
	query := `
from planets
	with 
		name = $name
		population = 1000
		residents = ["john", "janne"] -> flatten
		rotation_period = 24.5
		terrain = { "north": "jungle", "south": "rainforests" }
`

	planetResponse := `
[{
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
}]
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

		test.Equal(t, params["name"][0], "Yavin")
		test.Equal(t, params["population"][0], "1000")
		test.Equal(t, params["rotation_period"][0], "24.5")
		test.Equal(t, params["terrain"][0], `{"north":"jungle","south":"rainforests"}`)

		test.Equal(t, params["residents"], []string{"john", "janne"})

		w.WriteHeader(200)
		io.WriteString(w, planetResponse)
	})
	mockServer.Start()

	response, err := httpClient.Post(adHocQueryUrl+`&name=Yavin`, "text/plain", strings.NewReader(query))
	test.VerifyError(t, err)
	defer response.Body.Close()

	test.Equal(t, response.StatusCode, 200)

	var body map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&body)
	test.VerifyError(t, err)

	test.Equal(t, body, test.Unmarshal(expectedResponse))
}

func TestWithQualifierChainedParamOnFromStatement(t *testing.T) {
	query := `
from planets
	with 
		name = "Yavin"
		population = 1000
		residents = ["john", "janne"] -> flatten
		rotation_period = 24.5
		terrain = { "north": "jungle", "south": "rainforests" }

from people
	with
		name = planets.leader
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
	"leader": "Yavin King",
	"residents": ["john", "janne"],
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

	mockServer.Mux().HandleFunc("/api/planets/", func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()

		test.Equal(t, params["name"][0], "Yavin")
		test.Equal(t, params["population"][0], "1000")
		test.Equal(t, params["rotation_period"][0], "24.5")
		test.Equal(t, params["terrain"][0], `{"north":"jungle","south":"rainforests"}`)

		test.Equal(t, params["residents"], []string{"john", "janne"})

		w.WriteHeader(200)
		io.WriteString(w, planetResponse)
	})
	mockServer.Mux().HandleFunc("/api/people/", func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()

		test.Equal(t, params["name"][0], "Yavin King")

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
