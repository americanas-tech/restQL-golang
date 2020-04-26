package e2e

import "net/http"

var httpClient = &http.Client{}

const mockPort = 65000
const adHocQueryUrl = "http://localhost:9000/run-query?tenant=DEFAULT"
