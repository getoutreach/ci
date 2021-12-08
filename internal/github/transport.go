// Copyright 2021 Outreach Corporation. All Rights Reserved.

// Description: Implements a static token transport for use in
// rate limit validation
package github

import (
	"net/http"

	"github.com/getoutreach/gobox/pkg/cfg"
)

// staticTokenTransport is a small http.Roundtripper that injects
// a static token into all HTTP requests made with it.
type staticTokenTransport struct {
	token cfg.SecretData
}

// RoundTrip implements http.RoundTripper interface.
func (t *staticTokenTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+string(t.token))

	resp, err := http.DefaultTransport.RoundTrip(req)
	return resp, err
}
