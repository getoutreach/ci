package github

import (
	"net/http"

	"github.com/getoutreach/gobox/pkg/cfg"
)

type staticTokenTransport struct {
	token cfg.SecretData
}

// RoundTrip implements http.RoundTripper interface.
func (t *staticTokenTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+string(t.token))

	resp, err := http.DefaultTransport.RoundTrip(req)
	return resp, err
}
