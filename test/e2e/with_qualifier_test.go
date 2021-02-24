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

	"github.com/b2wdigital/restQL-golang/v5/test"
)

func TestWithQualifierOnFromStatement(t *testing.T) {
	query := `
from planets
	with
		planets = "5m"
		name = "Yavin IV"
		population = 1000
		residents = ["john", "janne"] -> no-multiplex
		rotation_period = 24.5
		terrain = { "north": "jungle", "south": "rainforests" }
		hot = true
		timestamp = null

from planets as second
	with
		to = planets.from
		planets = "5m"
		name = "Yavin IV"
		population = 1000
		residents = ["john", "janne"] -> no-multiplex
		rotation_period = 24.5
		terrain = { "north": "jungle", "south": "rainforests" }
		hot = true
		timestamp = null

from planets as third
	with
		to = second.from
		planets = "5m"
		name = "Yavin IV"
		population = 1000
		residents = ["john", "janne"] -> no-multiplex
		rotation_period = 24.5
		terrain = { "north": "jungle", "south": "rainforests" }
		hot = true
		timestamp = null
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
	"population": "1000",
	"residents": ["john", "janne"],
	"films": [1],
	"from": "sector A"
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
		"second": {
			"details": {
				"success": true,
				"status": 200,
				"metadata": {}
			},
			"result": %s
		},
		"third": {
			"details": {
				"success": true,
				"status": 200,
				"metadata": {}
			},
			"result": %s
		}
	}`, planetResponse, planetResponse, planetResponse)

	mockServer := test.NewMockServer(mockPort)
	defer mockServer.Teardown()

	mockServer.Mux().HandleFunc("/api/planets/", func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()

		var expectedTimestamp []string

		test.Equal(t, params["timestamp"], expectedTimestamp)
		test.Equal(t, params["planets"][0], "5m")
		test.Equal(t, params["hot"][0], "true")
		test.Equal(t, params["population"][0], "1000")
		test.Equal(t, params["rotation_period"][0], "24.5")
		test.Equal(t, params["residents"], []string{"john", "janne"})

		name, err := url.QueryUnescape(params["name"][0])
		test.VerifyError(t, err)
		test.Equal(t, name, "Yavin IV")

		terrain, err := url.QueryUnescape(params["terrain"][0])
		test.VerifyError(t, err)
		test.Equal(t, terrain, `{"north":"jungle","south":"rainforests"}`)

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

func TestWithQualifierOnToStatement(t *testing.T) {
	query := `
to planets
	with 
		name = "Yavin IV"
		population = 1000
		rotation_period = 24.5
		terrain = { "north": "jungle", "south": "rainforests" }
		residents = ["john", "janne"] -> no-multiplex
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
		residents = ["john", "janne"] -> no-multiplex
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
		residents = ["john", "janne"] -> no-multiplex
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
		residents = ["john", "janne"] -> no-multiplex
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

		test.Equal(t, params["population"][0], "1000")
		test.Equal(t, params["rotation_period"][0], "24.5")
		test.Equal(t, params["residents"], []string{"john", "janne"})

		name, err := url.QueryUnescape(params["name"][0])
		test.VerifyError(t, err)
		test.Equal(t, name, "Yavin IV")

		terrain, err := url.QueryUnescape(params["terrain"][0])
		test.VerifyError(t, err)
		test.Equal(t, terrain, `{"north":"jungle","south":"rainforests"}`)

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
		residents = ["john", "janne"] -> no-multiplex
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

		test.Equal(t, params["population"][0], "1000")
		test.Equal(t, params["rotation_period"][0], "24.5")
		test.Equal(t, params["residents"], []string{"john", "janne"})

		name, err := url.QueryUnescape(params["name"][0])
		test.VerifyError(t, err)
		test.Equal(t, name, "Yavin")

		terrain, err := url.QueryUnescape(params["terrain"][0])
		test.VerifyError(t, err)
		test.Equal(t, terrain, `{"north":"jungle","south":"rainforests"}`)

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
		residents = ["john", "janne"] -> no-multiplex
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

		test.Equal(t, params["population"][0], "1000")
		test.Equal(t, params["rotation_period"][0], "24.5")
		test.Equal(t, params["residents"], []string{"john", "janne"})

		name, err := url.QueryUnescape(params["name"][0])
		test.VerifyError(t, err)
		test.Equal(t, name, "Yavin")

		terrain, err := url.QueryUnescape(params["terrain"][0])
		test.VerifyError(t, err)
		test.Equal(t, terrain, `{"north":"jungle","south":"rainforests"}`)

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

func TestWithQualifierEncodersOnFromStatement(t *testing.T) {
	query := `
from planets
	with 
		name = "Yavin IV" -> base64
		population = 1000
		leaders = [["King", "Queen"], ["Governor"]] -> flatten -> no-multiplex
		residents = ["john", "janne"] -> json
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
	"population": "1000",
	"residents": ["john", "janne"],
	"films": [1, 2, 3]
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

		test.Equal(t, params["population"][0], "1000")
		test.Equal(t, params["leaders"], []string{"King", "Queen", "Governor"})

		name, err := url.QueryUnescape(params["name"][0])
		test.VerifyError(t, err)
		test.Equal(t, name, "WWF2aW4gSVY=")

		residents, err := url.QueryUnescape(params["residents"][0])
		test.VerifyError(t, err)
		test.Equal(t, residents, `["john","janne"]`)

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

func TestWithQualifierDynamicBodyOnToStatement(t *testing.T) {
	query := `
to planets
	with
		$planet
`

	planet := `
{
	"name": "Yavin",
	"rotation_period": 24.5,
	"orbital_period": "4818",
	"diameter": "10200",
	"climate": "temperate, tropical",
	"gravity": "1 standard",
	"terrain": { "north": "jungle", "south": "rainforests" },
	"surface_water": "8",
	"population": 1000,
	"residents": ["john", "janne"]
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
	}`, planet)

	mockServer := test.NewMockServer(mockPort)
	defer mockServer.Teardown()

	mockServer.Mux().HandleFunc("/api/planets/", func(w http.ResponseWriter, r *http.Request) {
		test.Equal(t, r.Method, http.MethodPost)

		b, err := ioutil.ReadAll(r.Body)
		test.VerifyError(t, err)

		body := string(b)
		test.NotEqual(t, body, "")

		test.Equal(t, test.Unmarshal(body), test.Unmarshal(removeWhitespaces(planet)))

		w.WriteHeader(201)
		io.WriteString(w, planet)
	})
	mockServer.Start()

	target := fmt.Sprintf("%s&planet=%s", adHocQueryUrl, removeWhitespaces(planet))
	response, err := httpClient.Post(target, "text/plain", strings.NewReader(query))
	test.VerifyError(t, err)
	defer response.Body.Close()

	test.Equal(t, response.StatusCode, 200)

	var body map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&body)
	test.VerifyError(t, err)

	test.Equal(t, body, test.Unmarshal(expectedResponse))
}

func TestWithQualifierMultiplexedDynamicBodyOnToStatement(t *testing.T) {
	query := `
to planets
	with
		$planet
`

	yavin := `
{
	"name": "Yavin",
	"rotation_period": 24.5,
	"orbital_period": "4818",
	"diameter": "10200",
	"climate": "temperate",
	"gravity": "standard",
	"terrain": { "north": "jungle", "south": "rainforests" },
	"surface_water": "8",
	"population": 1000,
	"residents": ["john", "janne"]
}`

	tatooine := `
{
	"name": "Tatooine",
	"rotation_period": "23",
	"orbital_period": "304",
	"diameter": "10465",
	"climate": "arid",
	"gravity": "standard",
	"terrain": "desert",
	"surface_water": "1",
	"population": "200000",
	"residents": ["william", "scarlet"]
}`

	expectedResponse := fmt.Sprintf(`
	{
		"planets": {
			"details": [
				{
					"success": true,
					"status": 201,
					"metadata": {}
				},
				{
					"success": true,
					"status": 201,
					"metadata": {}
				}
			],
			"result": [%s, %s] 
		}
	}`, yavin, tatooine)

	mockServer := test.NewMockServer(mockPort)
	defer mockServer.Teardown()

	mockServer.Mux().HandleFunc("/api/planets/", func(w http.ResponseWriter, r *http.Request) {
		test.Equal(t, r.Method, http.MethodPost)

		b, err := ioutil.ReadAll(r.Body)
		test.VerifyError(t, err)

		test.NotEqual(t, string(b), "")
		body := test.Unmarshal(string(b)).(map[string]interface{})

		if body["name"] == "Yavin" {
			test.Equal(t, body, test.Unmarshal(yavin))

			w.WriteHeader(201)
			io.WriteString(w, yavin)
			return
		}

		if body["name"] == "Tatooine" {
			test.Equal(t, body, test.Unmarshal(tatooine))

			w.WriteHeader(201)
			io.WriteString(w, tatooine)
		}
	})
	mockServer.Start()

	target := fmt.Sprintf("%s&planet=%s&planet=%s", adHocQueryUrl, removeWhitespaces(yavin), removeWhitespaces(tatooine))
	response, err := httpClient.Post(target, "text/plain", strings.NewReader(query))
	test.VerifyError(t, err)
	defer response.Body.Close()

	test.Equal(t, response.StatusCode, 200)

	var body map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&body)
	test.VerifyError(t, err)

	test.Equal(t, body, test.Unmarshal(expectedResponse))
}

func TestWithQualifierFlattenedDynamicBodyOnToStatement(t *testing.T) {
	query := `
to planets
	with
		$planet -> no-multiplex
`

	yavin := `
{
	"name": "Yavin",
	"rotation_period": 24.5,
	"orbital_period": "4818",
	"diameter": "10200",
	"climate": "temperate",
	"gravity": "standard",
	"terrain": { "north": "jungle", "south": "rainforests" },
	"surface_water": "8",
	"population": 1000,
	"residents": ["john", "janne"]
}`

	tatooine := `
{
	"name": "Tatooine",
	"rotation_period": "23",
	"orbital_period": "304",
	"diameter": "10465",
	"climate": "arid",
	"gravity": "standard",
	"terrain": "desert",
	"surface_water": "1",
	"population": "200000",
	"residents": ["william", "scarlet"]
}`

	expectedResponse := fmt.Sprintf(`
	{
		"planets": {
			"details": {
				"success": true,
				"status": 201,
				"metadata": {}
			},
			"result": [%s, %s] 
		}
	}`, yavin, tatooine)

	mockServer := test.NewMockServer(mockPort)
	defer mockServer.Teardown()

	mockServer.Mux().HandleFunc("/api/planets/", func(w http.ResponseWriter, r *http.Request) {
		test.Equal(t, r.Method, http.MethodPost)

		b, err := ioutil.ReadAll(r.Body)
		test.VerifyError(t, err)

		body := string(b)

		expectedBody := fmt.Sprintf(`[%s, %s]`, yavin, tatooine)

		test.NotEqual(t, body, "")
		test.Equal(t, test.Unmarshal(body), test.Unmarshal(expectedBody))

		response := fmt.Sprintf(`[%s, %s]`, yavin, tatooine)

		w.WriteHeader(201)
		io.WriteString(w, response)
	})
	mockServer.Start()

	target := fmt.Sprintf("%s&planet=%s&planet=%s", adHocQueryUrl, removeWhitespaces(yavin), removeWhitespaces(tatooine))
	response, err := httpClient.Post(target, "text/plain", strings.NewReader(query))
	test.VerifyError(t, err)
	defer response.Body.Close()

	test.Equal(t, response.StatusCode, 200)

	var body map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&body)
	test.VerifyError(t, err)

	test.Equal(t, body, test.Unmarshal(expectedResponse))
}

func TestWithQualifierDynamicBodyAsBase64OnToStatement(t *testing.T) {
	query := `
to planets
	with
		$planet -> base64
`

	planet := `
{
	"name": "Yavin",
	"rotation_period": 24.5,
	"orbital_period": "4818",
	"diameter": "10200",
	"climate": "temperate, tropical",
	"gravity": "1 standard",
	"terrain": { "north": "jungle", "south": "rainforests" },
	"surface_water": "8",
	"population": 1000,
	"residents": ["john", "janne"]
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
	}`, planet)

	mockServer := test.NewMockServer(mockPort)
	defer mockServer.Teardown()

	mockServer.Mux().HandleFunc("/api/planets/", func(w http.ResponseWriter, r *http.Request) {
		test.Equal(t, r.Method, http.MethodPost)

		b, err := ioutil.ReadAll(r.Body)
		test.VerifyError(t, err)

		body := string(b)
		test.Equal(t, body, "bWFwW2NsaW1hdGU6dGVtcGVyYXRlLHRyb3BpY2FsIGRpYW1ldGVyOjEwMjAwIGdyYXZpdHk6MXN0YW5kYXJkIG5hbWU6WWF2aW4gb3JiaXRhbF9wZXJpb2Q6NDgxOCBwb3B1bGF0aW9uOjEwMDAgcmVzaWRlbnRzOltqb2huIGphbm5lXSByb3RhdGlvbl9wZXJpb2Q6MjQuNSBzdXJmYWNlX3dhdGVyOjggdGVycmFpbjptYXBbbm9ydGg6anVuZ2xlIHNvdXRoOnJhaW5mb3Jlc3RzXV0=")

		w.WriteHeader(201)
		io.WriteString(w, planet)
	})
	mockServer.Start()

	target := fmt.Sprintf("%s&planet=%s", adHocQueryUrl, removeWhitespaces(planet))
	response, err := httpClient.Post(target, "text/plain", strings.NewReader(query))
	test.VerifyError(t, err)
	defer response.Body.Close()

	test.Equal(t, response.StatusCode, 200)

	var body map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&body)
	test.VerifyError(t, err)

	test.Equal(t, body, test.Unmarshal(expectedResponse))
}

func TestWithQualifierDynamicBodyWithJsonEncoderOnToStatement(t *testing.T) {
	query := `
to planets
	with
		$planet -> json
`

	planet := `
{
	"name": "Yavin",
	"rotation_period": 24.5,
	"orbital_period": "4818",
	"diameter": "10200",
	"climate": "temperate, tropical",
	"gravity": "1 standard",
	"terrain": { "north": "jungle", "south": "rainforests" },
	"surface_water": "8",
	"population": 1000,
	"residents": ["john", "janne"]
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
	}`, planet)

	mockServer := test.NewMockServer(mockPort)
	defer mockServer.Teardown()

	mockServer.Mux().HandleFunc("/api/planets/", func(w http.ResponseWriter, r *http.Request) {
		test.Equal(t, r.Method, http.MethodPost)

		b, err := ioutil.ReadAll(r.Body)
		test.VerifyError(t, err)

		body := string(b)
		test.NotEqual(t, body, "")

		test.Equal(t, test.Unmarshal(body), test.Unmarshal(removeWhitespaces(planet)))

		w.WriteHeader(201)
		io.WriteString(w, planet)
	})
	mockServer.Start()

	target := fmt.Sprintf("%s&planet=%s", adHocQueryUrl, removeWhitespaces(planet))
	response, err := httpClient.Post(target, "text/plain", strings.NewReader(query))
	test.VerifyError(t, err)
	defer response.Body.Close()

	test.Equal(t, response.StatusCode, 200)

	var body map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&body)
	test.VerifyError(t, err)

	test.Equal(t, body, test.Unmarshal(expectedResponse))
}

func TestWithQualifierMultiplexedDynamicBodyAndMultiplexPathParamOnToStatement(t *testing.T) {
	query := `
to planets
	with
		$planet
		id = ["Yavin", "Tatooine"]
`

	yavin := `
{
	"name": "Yavin",
	"rotation_period": 24.5,
	"orbital_period": "4818",
	"diameter": "10200",
	"climate": "temperate",
	"gravity": "standard",
	"terrain": { "north": "jungle", "south": "rainforests" },
	"surface_water": "8",
	"population": 1000,
	"residents": ["john", "janne"]
}`

	tatooine := `
{
	"name": "Tatooine",
	"rotation_period": "23",
	"orbital_period": "304",
	"diameter": "10465",
	"climate": "arid",
	"gravity": "standard",
	"terrain": "desert",
	"surface_water": "1",
	"population": "200000",
	"residents": ["william", "scarlet"]
}`

	expectedResponse := fmt.Sprintf(`
	{
		"planets": {
			"details": [
				{
					"success": true,
					"status": 201,
					"metadata": {}
				},
				{
					"success": true,
					"status": 201,
					"metadata": {}
				}
			],
			"result": [%s, %s] 
		}
	}`, yavin, tatooine)

	mockServer := test.NewMockServer(mockPort)
	defer mockServer.Teardown()

	mockServer.Mux().HandleFunc("/api/planets/Yavin", func(w http.ResponseWriter, r *http.Request) {
		test.Equal(t, r.Method, http.MethodPost)

		b, err := ioutil.ReadAll(r.Body)
		test.VerifyError(t, err)

		body := string(b)

		test.NotEqual(t, body, "")
		test.Equal(t, test.Unmarshal(body), test.Unmarshal(yavin))

		w.WriteHeader(201)
		io.WriteString(w, yavin)
	})
	mockServer.Mux().HandleFunc("/api/planets/Tatooine", func(w http.ResponseWriter, r *http.Request) {
		test.Equal(t, r.Method, http.MethodPost)

		b, err := ioutil.ReadAll(r.Body)
		test.VerifyError(t, err)

		body := string(b)

		test.NotEqual(t, body, "")
		test.Equal(t, test.Unmarshal(body), test.Unmarshal(tatooine))

		w.WriteHeader(201)
		io.WriteString(w, tatooine)
	})
	mockServer.Start()

	target := fmt.Sprintf("%s&planet=%s&planet=%s", adHocQueryUrl, removeWhitespaces(yavin), removeWhitespaces(tatooine))
	response, err := httpClient.Post(target, "text/plain", strings.NewReader(query))
	test.VerifyError(t, err)
	defer response.Body.Close()

	test.Equal(t, response.StatusCode, 200)

	var body map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&body)
	test.VerifyError(t, err)

	test.Equal(t, body, test.Unmarshal(expectedResponse))
}

func TestWithQualifierSendingQueryParametersOnToStatement(t *testing.T) {
	query := `
to starships
	with 
		id = $shipIds
		name = $shipName
`

	starshipResponse := `
{
	"id": 1,
	"name": "X-wing",
	"model": "T-65 X-wing",
	"manufacturer": "Incom Corporation",
	"cost_in_credits": "149999",
	"length": "12.5",
	"max_atmosphering_speed": "1050",
	"crew": "1",
	"passengers": "0",
	"cargo_capacity": "110",
	"consumables": "1 week",
	"hyperdrive_rating": "1.0",
	"MGLT": "100",
	"starship_class": "Starfighter"
}
`

	expectedResponse := fmt.Sprintf(`
	{
		"starships": {
			"details": {
				"success": true,
				"status": 201,
				"metadata": {}
			},
			"result": %s 
		}
	}`, starshipResponse)

	mockServer := test.NewMockServer(mockPort)
	defer mockServer.Teardown()

	mockServer.Mux().HandleFunc("/api/starships", func(w http.ResponseWriter, r *http.Request) {
		test.Equal(t, r.Method, http.MethodPost)

		params := r.URL.Query()
		test.Equal(t, params["id"][0], "1")
		test.Equal(t, params["name"][0], "X-wing")

		w.WriteHeader(201)
		io.WriteString(w, starshipResponse)
	})
	mockServer.Start()

	target := adHocQueryUrl + `&shipIds=1&shipName=X-wing`
	response, err := httpClient.Post(target, "text/plain", strings.NewReader(query))
	test.VerifyError(t, err)
	defer response.Body.Close()

	test.Equal(t, response.StatusCode, 200)

	var body map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&body)
	test.VerifyError(t, err)

	test.Equal(t, body, test.Unmarshal(expectedResponse))
}

func TestWithQualifierUsingParameterAsBodyOnToStatement(t *testing.T) {
	query := `
to planets
	with 
		body = [{ "north": "jungle", "south": "rainforests" }] -> as-body
		id = 1
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

	mockServer.Mux().HandleFunc("/api/planets/1", func(w http.ResponseWriter, r *http.Request) {
		test.Equal(t, r.Method, http.MethodPost)

		b, err := ioutil.ReadAll(r.Body)
		test.VerifyError(t, err)

		body := string(b)
		test.NotEqual(t, body, "")

		expectedBody := `[{ "north": "jungle", "south": "rainforests" }]`

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

func TestExplodeWithStaticObjectParameterWithListInsideOnToStatement(t *testing.T) {
	query := `
from planets
	with 
		id = 1

to people with 
	profiles = [{homeworld: "Yavin", friends: ["mark", "janne"]},
				{homeworld: "Yavin", friends: ["john", "anne"]}] -> no-multiplex
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

	people := `
{
	"inserted": 2
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
	}`, planetResponse, people)

	mockServer := test.NewMockServer(mockPort)
	defer mockServer.Teardown()

	mockServer.Mux().HandleFunc("/api/planets/1", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		io.WriteString(w, planetResponse)
	})
	mockServer.Mux().HandleFunc("/api/people/", func(w http.ResponseWriter, r *http.Request) {
		test.Equal(t, r.Method, http.MethodPost)
		expectedBody := map[string]interface{}{
			"profiles": []interface{}{
				[]interface{}{
					map[string]interface{}{"homeworld": "Yavin", "friends": "mark"},
					map[string]interface{}{"homeworld": "Yavin", "friends": "janne"},
				},
				[]interface{}{
					map[string]interface{}{"homeworld": "Yavin", "friends": "john"},
					map[string]interface{}{"homeworld": "Yavin", "friends": "anne"},
				},
			},
		}

		var body map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&body)
		test.VerifyError(t, err)

		test.Equal(t, body, expectedBody)

		io.WriteString(w, people)
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

func TestNoExplodeWithStaticObjectParameterWithListInsideOnToStatement(t *testing.T) {
	query := `
from planets
	with 
		id = 1

to people with 
	profiles = [{homeworld: "Yavin", friends: ["mark", "janne"]},
				{homeworld: "Yavin", friends: ["john", "anne"]}] -> no-explode -> no-multiplex
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

	people := `
{
	"inserted": 2
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
	}`, planetResponse, people)

	mockServer := test.NewMockServer(mockPort)
	defer mockServer.Teardown()

	mockServer.Mux().HandleFunc("/api/planets/1", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		io.WriteString(w, planetResponse)
	})
	mockServer.Mux().HandleFunc("/api/people/", func(w http.ResponseWriter, r *http.Request) {
		test.Equal(t, r.Method, http.MethodPost)
		expectedBody := map[string]interface{}{
			"profiles": []interface{}{
				map[string]interface{}{"homeworld": "Yavin", "friends": []interface{}{"mark", "janne"}},
				map[string]interface{}{"homeworld": "Yavin", "friends": []interface{}{"john", "anne"}},
			},
		}

		var body map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&body)
		test.VerifyError(t, err)

		test.Equal(t, body, expectedBody)

		io.WriteString(w, people)
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

func removeWhitespaces(s string) string {
	return strings.Join(strings.Fields(s), "")
}
