package e2e

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/b2wdigital/restQL-golang/v6/test"
)

const savedQueryUrl = "http://localhost:9000/run-query/test/variable-resolution/1?tenant=DEFAULT"

func TestVariableResolutionUsingQueryParametersOnFromStatement(t *testing.T) {
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

		test.Equal(t, params["residents"], []string{"john", "janne"})

		name, err := url.QueryUnescape(params["name"][0])
		test.VerifyError(t, err)
		test.Equal(t, name, "Yavin")

		w.WriteHeader(200)
		io.WriteString(w, planetResponse)
	})
	mockServer.Start()

	targetUrl := savedQueryUrl + "&name=Yavin&residents=john&residents=janne"
	response, err := httpClient.Get(targetUrl)
	test.VerifyError(t, err)
	defer response.Body.Close()

	test.Equal(t, response.StatusCode, 200)

	var result map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&result)
	test.VerifyError(t, err)

	test.Equal(t, result, test.Unmarshal(expectedResponse))
}

func TestVariableResolutionUsingBodyOnFromStatement(t *testing.T) {
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

		test.Equal(t, params["residents"], []string{"john", "janne"})

		name, err := url.QueryUnescape(params["name"][0])
		test.VerifyError(t, err)
		test.Equal(t, name, "Yavin")

		w.WriteHeader(200)
		io.WriteString(w, planetResponse)
	})
	mockServer.Start()

	body := `{ "name": "Yavin", "residents": ["john", "janne"] }`

	response, err := httpClient.Post(savedQueryUrl, "application/json", strings.NewReader(body))
	test.VerifyError(t, err)
	defer response.Body.Close()

	test.Equal(t, response.StatusCode, 200)

	var result map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&result)
	test.VerifyError(t, err)

	test.Equal(t, result, test.Unmarshal(expectedResponse))
}

func TestVariableResolutionUsingHeadersOnFromStatement(t *testing.T) {
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

		name, err := url.QueryUnescape(params["name"][0])
		test.VerifyError(t, err)
		test.Equal(t, name, "Yavin")

		w.WriteHeader(200)
		io.WriteString(w, planetResponse)
	})
	mockServer.Start()

	request, err := http.NewRequest(http.MethodGet, savedQueryUrl, ioutil.NopCloser(strings.NewReader("")))
	test.VerifyError(t, err)

	request.Header["name"] = []string{"Yavin"}

	response, err := httpClient.Do(request)
	test.VerifyError(t, err)
	defer response.Body.Close()

	test.Equal(t, response.StatusCode, 200)

	var result map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&result)
	test.VerifyError(t, err)

	test.Equal(t, result, test.Unmarshal(expectedResponse))
}

func TestVariableResolutionUsingQueryParametersOnToStatement(t *testing.T) {
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
		test.Equal(t, r.Method, http.MethodPost)

		b, err := io.ReadAll(r.Body)
		test.VerifyError(t, err)

		body := string(b)

		test.NotEqual(t, body, "")
		test.Equal(t, test.Unmarshal(body), test.Unmarshal(`{"planets": [{"name": "Yavin"}], "climate": "temperate"}`))

		w.WriteHeader(200)
		io.WriteString(w, planetResponse)
	})
	mockServer.Start()

	targetUrl := `http://localhost:9000/run-query/test/variable-resolution-on-to-statement/1?tenant=DEFAULT&body=[{"name": "Yavin"}]&climate=temperate`
	response, err := httpClient.Get(targetUrl)
	test.VerifyError(t, err)
	defer response.Body.Close()

	test.Equal(t, response.StatusCode, 200)

	var result map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&result)
	test.VerifyError(t, err)

	test.Equal(t, result, test.Unmarshal(expectedResponse))
}
