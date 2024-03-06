package model

type HttpRequestMsg struct {
	Uri        string
	Url        string
	Host       string
	Method     string
	Params     map[string]string
	Headers    map[string]string
	Body       string
	RemoteAddr string
	Msg        string
}
