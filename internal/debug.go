package internal

import (
	"log"
	"net/http"
	"net/http/httputil"
	"os"
)

func NewDebugTransport(innerTransport http.RoundTripper) http.RoundTripper {
	return &DebugTransport{
		transport: innerTransport,
	}
}

type DebugTransport struct {
	transport http.RoundTripper
}

func (c *DebugTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	if os.Getenv("SB_DEBUG") != "" {
		logRequest(request)
	}
	response, err := c.transport.RoundTrip(request)
	if os.Getenv("SB_DEBUG") != "" {
		logResponse(response, err)
	}
	return response, err
}

const logRequestTemplate = `DEBUG:
---[ REQUEST ]--------------------------------------------------------
%s
----------------------------------------------------------------------
`

const logResponseTemplate = `DEBUG:
---[ RESPONSE ]-------------------------------------------------------
%s
----------------------------------------------------------------------
`

func logRequest(r *http.Request) {
	body, err := httputil.DumpRequestOut(r, true)
	if err != nil {
		return
	}
	log.Printf(logRequestTemplate, body)
}

func logResponse(r *http.Response, err error) {
	if err != nil {
		log.Printf(logResponseTemplate, err)
		return
	}
	body, err := httputil.DumpResponse(r, true)
	if err != nil {
		return
	}
	log.Printf(logResponseTemplate, body)
}
