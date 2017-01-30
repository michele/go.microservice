package microservice

type Args struct {
	Path    string
	Body    string
	Method  string
	Params  map[string]string
	Headers map[string]string
}
