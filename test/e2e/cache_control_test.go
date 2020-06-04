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

func TestGlobalCacheControlOnFromStatement(t *testing.T) {
	query := `
use max-age 60

from planets
	with 
		id = "Yavin"
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

	mockServer.Mux().HandleFunc("/api/planets/Yavin", func(w http.ResponseWriter, r *http.Request) {
		test.Equal(t, r.Method, http.MethodGet)

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

	cacheControl := response.Header.Get("Cache-Control")
	test.Equal(t, cacheControl, "max-age=60")

}

func TestResourceCacheControlOnFromStatement(t *testing.T) {
	query := `
from planets
	max-age 60
	with 
		id = "Yavin"
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

	mockServer.Mux().HandleFunc("/api/planets/Yavin", func(w http.ResponseWriter, r *http.Request) {
		test.Equal(t, r.Method, http.MethodGet)

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

	cacheControl := response.Header.Get("Cache-Control")
	test.Equal(t, cacheControl, "max-age=60")

}

func TestMinimumQueryDefinedCacheControlOnFromStatement(t *testing.T) {
	query := `
use max-age 80

from planets
	max-age 60
	with 
		id = "Yavin"

from planets
	max-age 40
	with 
		id = "Yavin"
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

	mockServer.Mux().HandleFunc("/api/planets/Yavin", func(w http.ResponseWriter, r *http.Request) {
		test.Equal(t, r.Method, http.MethodGet)

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

	cacheControl := response.Header.Get("Cache-Control")
	test.Equal(t, cacheControl, "max-age=40")

}

func TestMinimumBetweenQueryDefinedCacheControlAndUpstreamCacheControlOnFromStatement(t *testing.T) {
	query := `
use max-age 80

from planets
	max-age 60
	with 
		id = "Yavin"

from planets
	max-age 40
	with 
		id = "Yavin"
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

	mockServer.Mux().HandleFunc("/api/planets/Yavin", func(w http.ResponseWriter, r *http.Request) {
		test.Equal(t, r.Method, http.MethodGet)

		w.Header().Add("Cache-Control", "max-age=10")
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

	cacheControl := response.Header.Get("Cache-Control")
	test.Equal(t, cacheControl, "max-age=10")

}

func TestUpstreamNoCacheOnFromStatement(t *testing.T) {
	query := `
use max-age 80

from planets
	max-age 60
	with 
		id = "Yavin"

from planets
	max-age 40
	with 
		id = "Yavin"
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

	mockServer.Mux().HandleFunc("/api/planets/Yavin", func(w http.ResponseWriter, r *http.Request) {
		test.Equal(t, r.Method, http.MethodGet)

		w.Header().Add("Cache-Control", "no-cache")
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

	cacheControl := response.Header.Get("Cache-Control")
	test.Equal(t, cacheControl, "no-cache")

}
