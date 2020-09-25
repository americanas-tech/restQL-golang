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

	"github.com/b2wdigital/restQL-golang/v4/test"
)

func TestChainedParamOnFromStatement(t *testing.T) {
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

func TestChainInterruptionWhenResourceFailOnFromStatement(t *testing.T) {
	query := `
from planets
	with 
		name = "Yavin"

from people
	with
		name = planets.leader
`

	expectedResponse := `
	{
		"planets": {
			"details": {
				"success": false,
				"status": 500,
				"metadata": {}
			},
			"result": {} 
		},
		"people": {
			"details": {
				"success": false,
				"status": 400,
				"metadata": {}
			},
			"result": "The request was skipped due to missing { :name } param value"
		}
	}`

	mockServer := test.NewMockServer(mockPort)
	defer mockServer.Teardown()

	mockServer.Mux().HandleFunc("/api/planets/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		io.WriteString(w, "{}")
	})
	mockServer.Mux().HandleFunc("/api/people/", func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("Resource /api/people called. Expected no call.")
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

func TestChainedParamOnToStatement(t *testing.T) {
	query := `
from planets
	with 
		id = "Yavin"

to people
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
				"status": 201,
				"metadata": {}
			},
			"result": %s
		}
	}`, planetResponse, peopleResponse)

	mockServer := test.NewMockServer(mockPort)
	defer mockServer.Teardown()

	mockServer.Mux().HandleFunc("/api/planets/Yavin", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		io.WriteString(w, planetResponse)
	})
	mockServer.Mux().HandleFunc("/api/people/", func(w http.ResponseWriter, r *http.Request) {
		test.Equal(t, r.Method, http.MethodPost)

		b, err := ioutil.ReadAll(r.Body)
		test.VerifyError(t, err)

		test.NotEqual(t, string(b), "")
		body := test.Unmarshal(string(b))

		test.Equal(t, body, test.Unmarshal(`{"name": "Yavin King"}`))

		w.WriteHeader(201)
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

func TestChainedParamOnIntoStatement(t *testing.T) {
	query := `
from planets
	with 
		id = "Yavin"

into people
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
				"status": 201,
				"metadata": {}
			},
			"result": %s
		}
	}`, planetResponse, peopleResponse)

	mockServer := test.NewMockServer(mockPort)
	defer mockServer.Teardown()

	mockServer.Mux().HandleFunc("/api/planets/Yavin", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		io.WriteString(w, planetResponse)
	})
	mockServer.Mux().HandleFunc("/api/people/", func(w http.ResponseWriter, r *http.Request) {
		test.Equal(t, r.Method, http.MethodPut)

		b, err := ioutil.ReadAll(r.Body)
		test.VerifyError(t, err)

		test.NotEqual(t, string(b), "")
		body := test.Unmarshal(string(b))

		test.Equal(t, body, test.Unmarshal(`{"name": "Yavin King"}`))

		w.WriteHeader(201)
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

func TestChainedParamOnUpdateStatement(t *testing.T) {
	query := `
from planets
	with 
		id = "Yavin"

update people
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
				"status": 201,
				"metadata": {}
			},
			"result": %s
		}
	}`, planetResponse, peopleResponse)

	mockServer := test.NewMockServer(mockPort)
	defer mockServer.Teardown()

	mockServer.Mux().HandleFunc("/api/planets/Yavin", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		io.WriteString(w, planetResponse)
	})
	mockServer.Mux().HandleFunc("/api/people/", func(w http.ResponseWriter, r *http.Request) {
		test.Equal(t, r.Method, http.MethodPatch)

		b, err := ioutil.ReadAll(r.Body)
		test.VerifyError(t, err)

		test.NotEqual(t, string(b), "")
		body := test.Unmarshal(string(b))

		test.Equal(t, body, test.Unmarshal(`{"name": "Yavin King"}`))

		w.WriteHeader(201)
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

func TestChainedParamOnDeleteStatement(t *testing.T) {
	query := `
from planets
	with 
		id = "Yavin"

delete people
	with
		id = planets.leader
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
	"leader": "King",
	"residents": ["john", "janne"],
	"films": [1]
}
`

	peopleResponse := `
{
	"name": "King",
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
				"status": 201,
				"metadata": {}
			},
			"result": %s
		}
	}`, planetResponse, peopleResponse)

	mockServer := test.NewMockServer(mockPort)
	defer mockServer.Teardown()

	mockServer.Mux().HandleFunc("/api/planets/Yavin", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		io.WriteString(w, planetResponse)
	})
	mockServer.Mux().HandleFunc("/api/people/King", func(w http.ResponseWriter, r *http.Request) {
		test.Equal(t, r.Method, http.MethodDelete)

		w.WriteHeader(201)
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

func TestChainedParamReturningXmlOnFromStatement(t *testing.T) {
	query := `
from planets
	with 
		id = "Yavin"

from people
	with
		name = planets.leader
`

	planetResponse := `
<planet>
   <climate>temperate, tropical</climate>
   <diameter>10200</diameter>
   <films>
      <element>1</element>
   </films>
   <gravity>1 standard</gravity>
   <leader>Yavin King</leader>
   <name>Yavin</name>
   <orbital_period>4818</orbital_period>
   <population>1000</population>
   <residents>
      <element>john</element>
      <element>janne</element>
   </residents>
   <rotation_period>24.5</rotation_period>
   <surface_water>8</surface_water>
   <terrain>
      <north>jungle</north>
      <south>rainforests</south>
   </terrain>
</planet>
`

	peopleResponse := `
[
	{
		"name": "Luke Skywalker",
		"height": "172",
		"mass": "77",
		"hair_color": "blond",
		"skin_color": "fair",
		"eye_color": "blue",
		"birth_year": "19BBY",
		"gender": "male",
		"homeworld": "http://swapi.dev/api/planets/1/",
		"species": [],
		"vehicles": [
			"http://swapi.dev/api/vehicles/14/",
			"http://swapi.dev/api/vehicles/30/"
		],
		"starships": [
			"http://swapi.dev/api/starships/12/",
			"http://swapi.dev/api/starships/22/"
		]
	},
	{
		"name": "C-3PO",
		"height": "167",
		"mass": "75",
		"hair_color": "n/a",
		"skin_color": "gold",
		"eye_color": "yellow",
		"birth_year": "112BBY",
		"gender": "n/a",
		"homeworld": "http://swapi.dev/api/planets/1/",
		"species": ["http://swapi.dev/api/species/2/"],
		"vehicles": [],
		"starships": []
	},
	{
		"name": "R2-D2",
		"height": "96",
		"mass": "32",
		"hair_color": "n/a",
		"skin_color": "white, blue",
		"eye_color": "red",
		"birth_year": "33BBY",
		"gender": "n/a",
		"homeworld": "http://swapi.dev/api/planets/8/",
		"species": ["http://swapi.dev/api/species/2/"],
		"vehicles": [],
		"starships": []
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
			"result": "\n<planet>\n   <climate>temperate, tropical</climate>\n   <diameter>10200</diameter>\n   <films>\n      <element>1</element>\n   </films>\n   <gravity>1 standard</gravity>\n   <leader>Yavin King</leader>\n   <name>Yavin</name>\n   <orbital_period>4818</orbital_period>\n   <population>1000</population>\n   <residents>\n      <element>john</element>\n      <element>janne</element>\n   </residents>\n   <rotation_period>24.5</rotation_period>\n   <surface_water>8</surface_water>\n   <terrain>\n      <north>jungle</north>\n      <south>rainforests</south>\n   </terrain>\n</planet>\n" 
		},
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

	mockServer.Mux().HandleFunc("/api/planets/Yavin", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		io.WriteString(w, planetResponse)
	})
	mockServer.Mux().HandleFunc("/api/people/", func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		var expectedNameParam []string
		test.Equal(t, params["name"], expectedNameParam)

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

func TestChainedParamReturningStringOnFromStatement(t *testing.T) {
	query := `
from planets as yavinAuth
	with 
		id = "Yavin"

from planets as tatooineAuth
	with 
		id = "Tatooine"

from people
	with
		name = tatooineAuth.leader
`

	yavinResponse := "a1HO1SDrIJ"
	tatooineResponse := ""
	peopleResponse := `
[
	{
		"name": "Luke Skywalker",
		"height": "172",
		"mass": "77",
		"hair_color": "blond",
		"skin_color": "fair",
		"eye_color": "blue",
		"birth_year": "19BBY",
		"gender": "male",
		"homeworld": "http://swapi.dev/api/planets/1/",
		"species": [],
		"vehicles": [
			"http://swapi.dev/api/vehicles/14/",
			"http://swapi.dev/api/vehicles/30/"
		],
		"starships": [
			"http://swapi.dev/api/starships/12/",
			"http://swapi.dev/api/starships/22/"
		]
	},
	{
		"name": "C-3PO",
		"height": "167",
		"mass": "75",
		"hair_color": "n/a",
		"skin_color": "gold",
		"eye_color": "yellow",
		"birth_year": "112BBY",
		"gender": "n/a",
		"homeworld": "http://swapi.dev/api/planets/1/",
		"species": ["http://swapi.dev/api/species/2/"],
		"vehicles": [],
		"starships": []
	},
	{
		"name": "R2-D2",
		"height": "96",
		"mass": "32",
		"hair_color": "n/a",
		"skin_color": "white, blue",
		"eye_color": "red",
		"birth_year": "33BBY",
		"gender": "n/a",
		"homeworld": "http://swapi.dev/api/planets/8/",
		"species": ["http://swapi.dev/api/species/2/"],
		"vehicles": [],
		"starships": []
	}
]
`

	expectedResponse := fmt.Sprintf(`
	{
		"yavinAuth": {
			"details": {
				"success": true,
				"status": 200,
				"metadata": {}
			},
			"result": "%s" 
		},
		"tatooineAuth": {
			"details": {
				"success": true,
				"status": 200,
				"metadata": {}
			}
		},
		"people": {
			"details": {
				"success": true,
				"status": 200,
				"metadata": {}
			},
			"result": %s
		}
	}`, yavinResponse, peopleResponse)

	mockServer := test.NewMockServer(mockPort)
	defer mockServer.Teardown()

	mockServer.Mux().HandleFunc("/api/planets/Yavin", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		io.WriteString(w, yavinResponse)
	})
	mockServer.Mux().HandleFunc("/api/planets/Tatooine", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		io.WriteString(w, tatooineResponse)
	})
	mockServer.Mux().HandleFunc("/api/people/", func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		var expectedNameParam []string
		test.Equal(t, params["name"], expectedNameParam)

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

func TestChainedParameterUsingKeyword(t *testing.T) {
	query := `
from planets as withLeader
	with 
		name = "Yavin"
		population = 1000
		residents = ["john", "janne"] -> no-multiplex
		rotation_period = 24.5
		terrain = { "north": "jungle", "south": "rainforests" }

from people
	with
		name = withLeader.with.leader
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
	"with": { "leader": "Yavin King" },
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
		"withLeader": {
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

	//bytes, err := ioutil.ReadAll(response.Body)
	//test.VerifyError(t, err)

	//fmt.Printf("response : %s\n", string(bytes))
	test.Equal(t, response.StatusCode, 200)

	var body map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&body)
	test.VerifyError(t, err)

	test.Equal(t, body, test.Unmarshal(expectedResponse))
}

func TestChainedParamTargetingUnknownStatementOnFromStatement(t *testing.T) {
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
		name = foo.leader
`

	expectedResponse := `
	{
		"error": "chained parameter targeting unknown statement : foo.leader"
	}`

	response, err := httpClient.Post(adHocQueryUrl, "text/plain", strings.NewReader(query))
	test.VerifyError(t, err)
	defer response.Body.Close()

	test.Equal(t, response.StatusCode, 400)

	var body map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&body)
	test.VerifyError(t, err)

	test.Equal(t, body, test.Unmarshal(expectedResponse))
}

func TestChainedParamWithVariableOnBracketsOnFromStatement(t *testing.T) {
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
		favoriteLandscape = planets.terrain[$direction]
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

		test.Equal(t, params["favoriteLandscape"][0], "jungle")

		w.WriteHeader(200)
		io.WriteString(w, peopleResponse)
	})
	mockServer.Start()

	target := adHocQueryUrl + "&direction=north"
	response, err := httpClient.Post(target, "text/plain", strings.NewReader(query))
	test.VerifyError(t, err)
	defer response.Body.Close()

	test.Equal(t, response.StatusCode, 200)

	var body map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&body)
	test.VerifyError(t, err)

	test.Equal(t, body, test.Unmarshal(expectedResponse))
}
