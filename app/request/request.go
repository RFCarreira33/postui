package request

type Model struct {
	url     string
	method  string
	headers map[string]string
	params  string
	body    string
}

func New() *Model {
	return &Model{
		method: "GET",
	}
}

func (m *Model) GetURL() string {
	if m.url[0] == ':' {
		m.url = "http://localhost" + m.url
	}

	return m.url + m.params
}

func (m *Model) GetMethod() string {
	return m.method
}

func (m *Model) GetHeaders() map[string]string {
	return m.headers
}

func (m *Model) SetURL(url string) {
	m.url = url
}

func (m *Model) SetMethod(method string) {
	m.method = method
}

func (m *Model) SetHeaders(headers map[string]string) {
	m.headers = headers
}

func (m *Model) SetParams(params map[string]string) {
	urlEnd := "?"
	for key, val := range params {
		urlEnd += key + "=" + val + "&"
	}
	m.params = urlEnd
}

func (m *Model) SetContentType(s string) {
	switch s {
	case "JSON":
		m.headers["Content-Type"] = "application/json"
	case "Javascript":
		m.headers["Content-Type"] = "application/javascript"
	case "XML":
		m.headers["Content-Type"] = "application/xml"
	case "HTML":
		m.headers["Content-Type"] = "text/html"
	case "Text":
		m.headers["Content-Type"] = "text/plain"
	}
}
