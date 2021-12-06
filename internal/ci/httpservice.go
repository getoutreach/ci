// Copyright 2021 Outreach Corporation. All Rights Reserved.

// Description: This file exposes the private HTTP service for ci.
// Managed: true

package ci //nolint:revive // Why: This nolint is here just in case your project name contains any of [-_].

import (
	"context"
	"fmt"
	"net/http"

	"github.com/getoutreach/httpx/pkg/handlers"
	// Place any extra imports for your service code here
	///Block(imports)
	///EndBlock(imports)
)

// HTTPService handles internal http requests
type HTTPService struct {
	handlers.Service
}

// Run is the entrypoint for the HTTPService serviceActivity.
func (s *HTTPService) Run(ctx context.Context, config *Config) error {
	// create a http handler (handlers.Service does metrics, health etc)
	///Block(privatehandler)
	s.App = http.NotFoundHandler()
	///EndBlock(privatehandler)
	return s.Service.Run(ctx, fmt.Sprintf("%s:%d", config.ListenHost, config.HTTPPort))
}
