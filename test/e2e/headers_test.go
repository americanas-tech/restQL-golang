package e2e

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/b2wdigital/restQL-golang/v4/test"
)

func TestHeadersOnFromStatement(t *testing.T) {
	query := `
from planets
	headers
		X-TID = "abcdef-12345"
`

	planetResponse := `
[{
	"name": "Yavin IV",
	"rotation_period": "24",
	"orbital_period": "4818",
	"diameter": "10200",
	"climate": "temperate, tropical",
	"gravity": "1 standard",
	"terrain": "jungle, rainforests",
	"surface_water": "8",
	"population": "1000",
	"residents": [],
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
		test.Equal(t, r.Method, http.MethodGet)

		xtid := r.Header.Get("X-TID")
		test.Equal(t, xtid, "abcdef-12345")

		w.WriteHeader(200)
		io.WriteString(w, planetResponse)
	})
	mockServer.Start()

	response, err := httpClient.Post(adHocQueryUrl, "text/plain", strings.NewReader(query))
	test.VerifyError(t, err)
	defer response.Body.Close()

	var body map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&body)
	test.VerifyError(t, err)

	test.Equal(t, body, test.Unmarshal(expectedResponse))
}

func TestHeadersOnToStatement(t *testing.T) {
	query := `
to planets
	headers
		X-TID = "abcdef-12345"
`

	planetResponse := `
[{
	"name": "Yavin IV",
	"rotation_period": "24",
	"orbital_period": "4818",
	"diameter": "10200",
	"climate": "temperate, tropical",
	"gravity": "1 standard",
	"terrain": "jungle, rainforests",
	"surface_water": "8",
	"population": "1000",
	"residents": [],
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

		xtid := r.Header.Get("X-TID")
		test.Equal(t, xtid, "abcdef-12345")

		w.WriteHeader(200)
		io.WriteString(w, planetResponse)
	})
	mockServer.Start()

	response, err := httpClient.Post(adHocQueryUrl, "text/plain", strings.NewReader(query))
	test.VerifyError(t, err)
	defer response.Body.Close()

	var body map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&body)
	test.VerifyError(t, err)

	test.Equal(t, body, test.Unmarshal(expectedResponse))
}

func TestHeadersOnIntoStatement(t *testing.T) {
	query := `
into planets
	headers
		X-TID = "abcdef-12345"
`

	planetResponse := `
[{
	"name": "Yavin IV",
	"rotation_period": "24",
	"orbital_period": "4818",
	"diameter": "10200",
	"climate": "temperate, tropical",
	"gravity": "1 standard",
	"terrain": "jungle, rainforests",
	"surface_water": "8",
	"population": "1000",
	"residents": [],
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
		test.Equal(t, r.Method, http.MethodPut)

		xtid := r.Header.Get("X-TID")
		test.Equal(t, xtid, "abcdef-12345")

		w.WriteHeader(200)
		io.WriteString(w, planetResponse)
	})
	mockServer.Start()

	response, err := httpClient.Post(adHocQueryUrl, "text/plain", strings.NewReader(query))
	test.VerifyError(t, err)
	defer response.Body.Close()

	var body map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&body)
	test.VerifyError(t, err)

	test.Equal(t, body, test.Unmarshal(expectedResponse))
}

func TestHeadersOnUpdateStatement(t *testing.T) {
	query := `
update planets
	headers
		X-TID = "abcdef-12345"
`

	planetResponse := `
[{
	"name": "Yavin IV",
	"rotation_period": "24",
	"orbital_period": "4818",
	"diameter": "10200",
	"climate": "temperate, tropical",
	"gravity": "1 standard",
	"terrain": "jungle, rainforests",
	"surface_water": "8",
	"population": "1000",
	"residents": [],
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
		test.Equal(t, r.Method, http.MethodPatch)

		xtid := r.Header.Get("X-TID")
		test.Equal(t, xtid, "abcdef-12345")

		w.WriteHeader(200)
		io.WriteString(w, planetResponse)
	})
	mockServer.Start()

	response, err := httpClient.Post(adHocQueryUrl, "text/plain", strings.NewReader(query))
	test.VerifyError(t, err)
	defer response.Body.Close()

	var body map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&body)
	test.VerifyError(t, err)

	test.Equal(t, body, test.Unmarshal(expectedResponse))
}

func TestHeadersOnDeleteStatement(t *testing.T) {
	query := `
delete planets
	headers
		X-TID = "abcdef-12345"
`

	planetResponse := `
[{
	"name": "Yavin IV",
	"rotation_period": "24",
	"orbital_period": "4818",
	"diameter": "10200",
	"climate": "temperate, tropical",
	"gravity": "1 standard",
	"terrain": "jungle, rainforests",
	"surface_water": "8",
	"population": "1000",
	"residents": [],
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
		test.Equal(t, r.Method, http.MethodDelete)

		xtid := r.Header.Get("X-TID")
		test.Equal(t, xtid, "abcdef-12345")

		w.WriteHeader(200)
		io.WriteString(w, planetResponse)
	})
	mockServer.Start()

	response, err := httpClient.Post(adHocQueryUrl, "text/plain", strings.NewReader(query))
	test.VerifyError(t, err)
	defer response.Body.Close()

	var body map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&body)
	test.VerifyError(t, err)

	test.Equal(t, body, test.Unmarshal(expectedResponse))
}

func TestHeadersFromVariableOnFromStatement(t *testing.T) {
	query := `
from planets
	headers
		X-TID = $xtid
`

	planetResponse := `
[{
	"name": "Yavin IV",
	"rotation_period": "24",
	"orbital_period": "4818",
	"diameter": "10200",
	"climate": "temperate, tropical",
	"gravity": "1 standard",
	"terrain": "jungle, rainforests",
	"surface_water": "8",
	"population": "1000",
	"residents": [],
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
		test.Equal(t, r.Method, http.MethodGet)

		xtid := r.Header.Get("X-TID")
		test.Equal(t, xtid, "abcdef-12345")

		w.WriteHeader(200)
		io.WriteString(w, planetResponse)
	})
	mockServer.Start()

	targetUrl := fmt.Sprintf("%s&xtid=%s", adHocQueryUrl, "abcdef-12345")
	response, err := httpClient.Post(targetUrl, "text/plain", strings.NewReader(query))
	test.VerifyError(t, err)
	defer response.Body.Close()

	var body map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&body)
	test.VerifyError(t, err)

	test.Equal(t, body, test.Unmarshal(expectedResponse))
}

func TestHeadersWithChainOnFromStatement(t *testing.T) {
	query := `
from planets as first
	with
		id = 1

from planets
	headers
		X-TID = first.name
	with
		id = 2
`

	planetResponse := `
{
	"name": "Yavin IV",
	"rotation_period": "24",
	"orbital_period": "4818",
	"diameter": "10200",
	"climate": "temperate, tropical",
	"gravity": "1 standard",
	"terrain": "jungle, rainforests",
	"surface_water": "8",
	"population": "1000",
	"residents": [],
	"films": [1]
}
`

	expectedResponse := fmt.Sprintf(`
	{
		"first": {
			"details": {
				"success": true,
				"status": 200,
				"metadata": {}
			},
			"result": %s 
		},
		"planets": {
			"details": {
				"success": true,
				"status": 200,
				"metadata": {}
			},
			"result": %s 
		}
	}`, planetResponse, planetResponse)

	mockServer := test.NewMockServer(mockPort)
	defer mockServer.Teardown()

	mockServer.Mux().HandleFunc("/api/planets/1", func(w http.ResponseWriter, r *http.Request) {
		test.Equal(t, r.Method, http.MethodGet)

		w.WriteHeader(200)
		io.WriteString(w, planetResponse)
	})
	mockServer.Mux().HandleFunc("/api/planets/2", func(w http.ResponseWriter, r *http.Request) {
		test.Equal(t, r.Method, http.MethodGet)

		xtid := r.Header.Get("X-TID")
		test.Equal(t, xtid, "Yavin IV")

		w.WriteHeader(200)
		io.WriteString(w, planetResponse)
	})
	mockServer.Start()

	targetUrl := fmt.Sprintf("%s&xtid=%s", adHocQueryUrl, "abcdef-12345")
	response, err := httpClient.Post(targetUrl, "text/plain", strings.NewReader(query))
	test.VerifyError(t, err)
	defer response.Body.Close()

	var body map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&body)
	test.VerifyError(t, err)

	test.Equal(t, body, test.Unmarshal(expectedResponse))
}

func TestForwardingHeadersOnFromStatement(t *testing.T) {
	query := `
from planets
	headers
		X-TID = "abcdef-12345"
`

	planetResponse := `
[{
	"name": "Yavin IV",
	"rotation_period": "24",
	"orbital_period": "4818",
	"diameter": "10200",
	"climate": "temperate, tropical",
	"gravity": "1 standard",
	"terrain": "jungle, rainforests",
	"surface_water": "8",
	"population": "1000",
	"residents": [],
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
		test.Equal(t, r.Method, http.MethodGet)

		xtid := r.Header.Get("X-TID")
		test.Equal(t, xtid, "abcdef-12345")

		custom := r.Header.Get("X-Custom")
		test.Equal(t, custom, "98765")

		w.WriteHeader(200)
		io.WriteString(w, planetResponse)
	})
	mockServer.Start()

	targetUrl := fmt.Sprintf("%s&xtid=%s", adHocQueryUrl, "abcdef-12345")
	req, err := http.NewRequest(http.MethodPost, targetUrl, strings.NewReader(query))
	test.VerifyError(t, err)

	req.Header.Add("Content-Type", "text/plain")
	req.Header.Add("X-Custom", "98765")

	response, err := httpClient.Do(req)
	test.VerifyError(t, err)
	defer response.Body.Close()

	var body map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&body)
	test.VerifyError(t, err)

	test.Equal(t, body, test.Unmarshal(expectedResponse))
}

func TestHeadersShouldBeCaseInsensitiveWhenMergingWithForwardHeaders(t *testing.T) {
	query := `
from planets
	headers
		X-TID = "abcdef-12345"
		accept = "application/json"
`

	planetResponse := `
[{
	"name": "Yavin IV",
	"rotation_period": "24",
	"orbital_period": "4818",
	"diameter": "10200",
	"climate": "temperate, tropical",
	"gravity": "1 standard",
	"terrain": "jungle, rainforests",
	"surface_water": "8",
	"population": "1000",
	"residents": [],
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
		test.Equal(t, r.Method, http.MethodGet)

		xtid := r.Header.Get("X-TID")
		test.Equal(t, xtid, "abcdef-12345")

		accept := r.Header.Get("Accept")
		test.Equal(t, accept, "application/json")

		w.WriteHeader(200)
		io.WriteString(w, planetResponse)
	})
	mockServer.Start()

	request, err := http.NewRequest(http.MethodPost, adHocQueryUrl, strings.NewReader(query))
	test.VerifyError(t, err)
	request.Header.Set("Accept", "*/*")
	request.Header.Set("Content-Type", "text/plain")

	response, err := httpClient.Do(request)
	test.VerifyError(t, err)
	defer response.Body.Close()

	var body map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&body)
	test.VerifyError(t, err)

	test.Equal(t, body, test.Unmarshal(expectedResponse))
}
