package microservice

import (
	"flag"
	"strconv"

	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var DefaultDomain string
var DefaultPort = "3000"
var DefaultTimeout = 10

var client *http.Client

func init() {
	DefaultDomain = os.Getenv("GO_ENV")
	var timeout time.Duration
	var seconds int
	var parsErr error
	seconds, parsErr = strconv.Atoi(os.Getenv("MS_DEFAULT_TIMEOUT"))
	if parsErr != nil {
		seconds = DefaultTimeout
	}
	timeout = time.Duration(seconds) * time.Second

	client = &http.Client{
		Timeout: timeout,
	}

	if flag.Lookup("test.v") == nil {
		client.Transport = &http.Transport{
			MaxIdleConns: 20,
		}
	}
}

func Call(ms string, args Args) (*http.Response, error) {
	req, err := Prepare(ms, args)

	if err != nil {
		return nil, err
	} else {
		resp, err := client.Do(req)
		return resp, err
	}
}

func Prepare(ms string, args Args) (*http.Request, error) {
	uurl := os.Getenv(strings.ToUpper(ms) + "_URL")
	if uurl == "" {
		var domain string
		domain = os.Getenv("MS_DOMAIN")
		if domain == "" {
			domain = DefaultDomain
		}
		port := os.Getenv("MS_" + strings.ToUpper(ms) + "_PORT")
		if port == "" {
			port = DefaultPort
		}
		uurl = "http://" + ms + "." + domain + ":" + port
	}
	if args.Path != "" {
		uurl = uurl + "/" + args.Path
	}
	var method string
	if args.Method != "" {
		method = args.Method
	} else {
		method = "GET"
	}
	req, err := http.NewRequest(method, uurl, nil)

	if err != nil {
		return nil, err
	} else {
		query := url.Values{}
		for key, value := range args.Params {
			query.Add(key, value)
		}
		req.URL.RawQuery = query.Encode()
	}

	if len(args.Headers) > 0 {
		for key, value := range args.Headers {
			req.Header.Add(key, value)
		}
	}

	return req, nil
}

// def prepare_microservice_request(microservice, args = {})
//   raise ArgumentError unless microservice
//
//   url = ENV["#{microservice.upcase}_URL"] || "http://#{microservice}.#{ENV['MS_DOMAIN'] || Rails.env}:#{ENV["MS_#{microservice.upcase}_PORT"] || '3000'}"
//   url << "/#{args[:path]}" if args[:path]
//   body = (args[:body].is_a? String) ? args[:body] : args[:body].to_json if args[:body]
//   Typhoeus::Request.new(url, method: args[:method] || :get, headers: MicroserviceDSL.default_headers, body: body, params: args[:params])
// end
//
// def call_microservice(microservice, args={})
//   request = prepare_microservice_request microservice, args
//   response = request.run
//   [response.body, response.headers['content-type'], response.code]
// end
