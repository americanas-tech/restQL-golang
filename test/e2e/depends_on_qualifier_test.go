package e2e

import (
	"encoding/json"
	"fmt"
	"github.com/b2wdigital/restQL-golang/v6/test"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestDependsOnQualifierOnFromStatement(t *testing.T) {
	query := `
from planets as first
	with
		id = 1

from planets as second
	depends-on first
	with
		id = 2
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
	"films": [1],
	"from": "sector A"
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
		"second": {
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

	planetOneWasCalled := false
	mockServer.Mux().HandleFunc("/api/planets/1", func(w http.ResponseWriter, r *http.Request) {
		planetOneWasCalled = true

		w.WriteHeader(200)
		io.WriteString(w, planetResponse)
	})
	mockServer.Mux().HandleFunc("/api/planets/2", func(w http.ResponseWriter, r *http.Request) {
		test.Equal(t, planetOneWasCalled, true)

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

func TestDependsOnQualifierOnToStatement(t *testing.T) {
	query := `
from planets as first
	with
		id = 1

to planets as second
	depends-on first
	with
		id = 2
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
	"films": [1],
	"from": "sector A"
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
		"second": {
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

	planetOneWasCalled := false
	mockServer.Mux().HandleFunc("/api/planets/1", func(w http.ResponseWriter, r *http.Request) {
		planetOneWasCalled = true

		w.WriteHeader(200)
		io.WriteString(w, planetResponse)
	})
	mockServer.Mux().HandleFunc("/api/planets/2", func(w http.ResponseWriter, r *http.Request) {
		test.Equal(t, r.Method, http.MethodPost)
		test.Equal(t, planetOneWasCalled, true)

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

func TestDependsOnQualifierOnIntoStatement(t *testing.T) {
	query := `
from planets as first
	with
		id = 1

into planets as second
	depends-on first
	with
		id = 2
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
	"films": [1],
	"from": "sector A"
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
		"second": {
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

	planetOneWasCalled := false
	mockServer.Mux().HandleFunc("/api/planets/1", func(w http.ResponseWriter, r *http.Request) {
		planetOneWasCalled = true

		w.WriteHeader(200)
		io.WriteString(w, planetResponse)
	})
	mockServer.Mux().HandleFunc("/api/planets/2", func(w http.ResponseWriter, r *http.Request) {
		test.Equal(t, r.Method, http.MethodPut)
		test.Equal(t, planetOneWasCalled, true)

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

func TestDependsOnQualifierOnUpdateStatement(t *testing.T) {
	query := `
from planets as first
	with
		id = 1

update planets as second
	depends-on first
	with
		id = 2
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
	"films": [1],
	"from": "sector A"
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
		"second": {
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

	planetOneWasCalled := false
	mockServer.Mux().HandleFunc("/api/planets/1", func(w http.ResponseWriter, r *http.Request) {
		planetOneWasCalled = true

		w.WriteHeader(200)
		io.WriteString(w, planetResponse)
	})
	mockServer.Mux().HandleFunc("/api/planets/2", func(w http.ResponseWriter, r *http.Request) {
		test.Equal(t, r.Method, http.MethodPatch)
		test.Equal(t, planetOneWasCalled, true)

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

func TestDependsOnQualifierOnDeleteStatement(t *testing.T) {
	query := `
from planets as first
	with
		id = 1

delete planets as second
	depends-on first
	with
		id = 2
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
	"films": [1],
	"from": "sector A"
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
		"second": {
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

	planetOneWasCalled := false
	mockServer.Mux().HandleFunc("/api/planets/1", func(w http.ResponseWriter, r *http.Request) {
		planetOneWasCalled = true

		w.WriteHeader(200)
		io.WriteString(w, planetResponse)
	})
	mockServer.Mux().HandleFunc("/api/planets/2", func(w http.ResponseWriter, r *http.Request) {
		test.Equal(t, r.Method, http.MethodDelete)
		test.Equal(t, planetOneWasCalled, true)

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

func TestDependsOnQualifierOnMultiplexedFromStatement(t *testing.T) {
	query := `
from planets as first
	with
		id = 1

from planets as second
	depends-on first
	with
		id = [2, 3, 4]
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
	"films": [1],
	"from": "sector A"
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
		"second": {
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
				},
				{
					"success": true,
					"status": 200,
					"metadata": {}
				}
			],
			"result": [%s, %s, %s]
		}
	}`, planetResponse, planetResponse, planetResponse, planetResponse)

	mockServer := test.NewMockServer(mockPort)
	defer mockServer.Teardown()

	planetOneWasCalled := false
	mockServer.Mux().HandleFunc("/api/planets/1", func(w http.ResponseWriter, r *http.Request) {
		planetOneWasCalled = true

		w.WriteHeader(200)
		io.WriteString(w, planetResponse)
	})
	mockServer.Mux().HandleFunc("/api/planets/2", func(w http.ResponseWriter, r *http.Request) {
		test.Equal(t, planetOneWasCalled, true)

		w.WriteHeader(200)
		io.WriteString(w, planetResponse)
	})
	mockServer.Mux().HandleFunc("/api/planets/3", func(w http.ResponseWriter, r *http.Request) {
		test.Equal(t, planetOneWasCalled, true)

		w.WriteHeader(200)
		io.WriteString(w, planetResponse)
	})
	mockServer.Mux().HandleFunc("/api/planets/4", func(w http.ResponseWriter, r *http.Request) {
		test.Equal(t, planetOneWasCalled, true)

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

func TestDependsOnQualifierOnMultiplexedTargetStatement(t *testing.T) {
	query := `
from planets as first
	with
		id = [1, 2, 3]

from planets as second
	depends-on first
	with
		id = 4
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
	"films": [1],
	"from": "sector A"
}
`

	expectedResponse := fmt.Sprintf(`
	{
		"second" : {
			"details": {
				"success": true,
				"status": 200,
				"metadata": {}
			},
			"result": %s
		},
		"first": {
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
				},
				{
					"success": true,
					"status": 200,
					"metadata": {}
				}
			],
			"result": [%s, %s, %s]
		}
	}`, planetResponse, planetResponse, planetResponse, planetResponse)

	mockServer := test.NewMockServer(mockPort)
	defer mockServer.Teardown()

	planetOneWasCalled := false
	planetTwoWasCalled := false
	planetThreeWasCalled := false
	mockServer.Mux().HandleFunc("/api/planets/1", func(w http.ResponseWriter, r *http.Request) {
		planetOneWasCalled = true

		w.WriteHeader(200)
		io.WriteString(w, planetResponse)
	})
	mockServer.Mux().HandleFunc("/api/planets/2", func(w http.ResponseWriter, r *http.Request) {
		planetTwoWasCalled = true

		w.WriteHeader(200)
		io.WriteString(w, planetResponse)
	})
	mockServer.Mux().HandleFunc("/api/planets/3", func(w http.ResponseWriter, r *http.Request) {
		planetThreeWasCalled = true

		w.WriteHeader(200)
		io.WriteString(w, planetResponse)
	})
	mockServer.Mux().HandleFunc("/api/planets/4", func(w http.ResponseWriter, r *http.Request) {
		test.Equal(t, planetOneWasCalled, true)
		test.Equal(t, planetTwoWasCalled, true)
		test.Equal(t, planetThreeWasCalled, true)

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

func TestDependsOnQualifierOnWhenTargetHasIgnoreErrors(t *testing.T) {
	query := `
from planets as first
	with
		id = 1
	ignore-errors

from planets as second
	depends-on first
	with
		id = 2
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
	"films": [1],
	"from": "sector A"
}
`

	expectedResponse := fmt.Sprintf(`
	{
		"first": {
			"details": {
				"success": false,
				"status": 408,
				"metadata": {"ignore-errors": "ignore"}
			},
			"result": %s
		},
		"second": {
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

	planetOneWasCalled := false
	mockServer.Mux().HandleFunc("/api/planets/1", func(w http.ResponseWriter, r *http.Request) {
		planetOneWasCalled = true

		w.WriteHeader(408)
		io.WriteString(w, planetResponse)
	})
	mockServer.Mux().HandleFunc("/api/planets/2", func(w http.ResponseWriter, r *http.Request) {
		test.Equal(t, planetOneWasCalled, true)

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

func TestDependsOnQualifierOnBlocksExecutionIfTargetFail(t *testing.T) {
	query := `
from planets as first
	with
		id = 1

from planets as second
	depends-on first
	with
		id = 2
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
	"films": [1],
	"from": "sector A"
}
`

	expectedResponse := fmt.Sprintf(`
	{
		"first": {
			"details": {
				"success": false,
				"status": 408,
				"metadata": {}
			},
			"result": %s
		},
		"second": {
			"details": {
				"success": false,
				"status": 400,
				"metadata": {}
			},
			"result": "The request was skipped due to unresolved dependency { first }"
		}
	}`, planetResponse)

	mockServer := test.NewMockServer(mockPort)
	defer mockServer.Teardown()

	planetOneWasCalled := false
	mockServer.Mux().HandleFunc("/api/planets/1", func(w http.ResponseWriter, r *http.Request) {
		planetOneWasCalled = true

		w.WriteHeader(408)
		io.WriteString(w, planetResponse)
	})
	mockServer.Mux().HandleFunc("/api/planets/2", func(w http.ResponseWriter, r *http.Request) {
		t.Error("second statement should not have been called")
	})
	mockServer.Start()

	response, err := httpClient.Post(adHocQueryUrl, "text/plain", strings.NewReader(query))
	test.VerifyError(t, err)
	defer response.Body.Close()

	test.Equal(t, response.StatusCode, 408)

	var body map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&body)
	test.VerifyError(t, err)

	test.Equal(t, body, test.Unmarshal(expectedResponse))
	test.Equal(t, planetOneWasCalled, true)
}
