package microservice_test

import (
	"github.com/michele/go.microservice"

	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrepare(t *testing.T) {
	microservice.DefaultDomain = "development"
	microservice.DefaultPort = "3000"

	ms := "microservice"
	args := microservice.Args{
		Path:    "ping",
		Body:    `{"data": "test"}`,
		Method:  "POST",
		Params:  map[string]string{"foo": "bar", "more": "less"},
		Headers: map[string]string{"Content-Type": "application/json", "Accept": "text/plain"},
	}

	req, err := microservice.Prepare(ms, args)

	assert.Nil(t, err)

	assert.Equal(t, "http", req.URL.Scheme)
	assert.Equal(t, ms+"."+microservice.DefaultDomain+":"+microservice.DefaultPort, req.URL.Host)
	assert.Equal(t, "/"+args.Path, req.URL.Path)

	query := url.Values{}
	for key, value := range args.Params {
		query.Add(key, value)
	}
	rawQuery := query.Encode()

	assert.Equal(t, rawQuery, req.URL.RawQuery)

	for k, v := range args.Headers {
		assert.Equal(t, v, req.Header.Get(k))
	}
}
