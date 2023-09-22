package request

type Model struct {
	url     string
	method  string
	headers map[string]string
	params  string
}

func New() *Model {
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	return &Model{
		url:     "",
		method:  "GET",
		headers: headers,
	}
}

func (m *Model) GetURL() string {
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
