package ai

import "net/http"

type mockTransport struct {
	response *http.Response
	err      error
}

func (t *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.response, t.err
}
