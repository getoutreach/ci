// Copyright 2021 Outreach Corporation. All Rights Reserved.

// Description: Implements the core Github credential pooling logic

// Package github contains all logic for implementing Github credential pooling
package github

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/getoutreach/gobox/pkg/cfg"
	gc "github.com/google/go-github/v34/github"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Credential is a type of credential to use when talking to Github
type Credential struct {
	// Name is an optional field to supply to make it easier to identify a
	// set of credentials
	Name string

	// AppID is the ID of the Github App that this credential is tied to.
	// If this is set then PEM must also be set.
	AppID *int

	// InstallID is the install ID of the AppID provided above. Must be set
	// when AppID is set.
	InstallID *int

	// PEM is the private key of the Github App that this credential is
	// tied to. This should not be provided if AppID is 0.
	PEM cfg.SecretData

	// AccessToken is a Github Personal Access Token to be used. This must
	// be set if AppID is not set.
	AccessToken cfg.SecretData
}

// getToken returns a valid token from a Credential or an error if one
// can't be obtained
func (c *Credential) getToken(ctx context.Context) (cfg.SecretData, error) {
	if c.AppID == nil {
		return c.AccessToken, nil
	}

	atsp, err := ghinstallation.NewAppsTransport(http.DefaultTransport, int64(*c.AppID), []byte(string(c.PEM)))
	if err != nil {
		return "", errors.Wrap(err, "failed to create github app client transport")
	}

	// Create a Token for use as our Github App
	token, _, err := gc.NewClient(&http.Client{Transport: atsp}).
		Apps.CreateInstallationToken(ctx, int64(*c.InstallID), &gc.InstallationTokenOptions{})
	if err != nil {
		return "", errors.Wrap(err, "failed to create access token from Github App")
	}
	if token == nil {
		return "", fmt.Errorf("malformed response from Github")
	}

	return cfg.SecretData(*token.Token), nil
}

// GetToken returns a valid token from a credential and ensures it's not rate-limited
// if it is, or one cannot be obtained, then instead an error is returned.
func (c *Credential) GetToken(ctx context.Context) (cfg.SecretData, error) {
	t, err := c.getToken(ctx)
	if err != nil {
		return "", err
	}

	cli := gc.NewClient(&http.Client{
		Transport: &staticTokenTransport{token: t},
	})

	rl, _, err := cli.RateLimits(ctx)
	if err != nil {
		return "", err
	}

	if rl.Core == nil {
		return "", fmt.Errorf("failed to calculate rate limit")
	}

	if rl.Core.Remaining < 10 {
		return "", fmt.Errorf("token close to be, or is, rate limited")
	}

	// Now use the token against a known working endpoint because /rate_limits
	// doesn't always work for Github Apps
	_, _, err = cli.Organizations.Get(ctx, "getoutreach")
	if err != nil {
		return "", errors.Wrap(err, "failed to test token")
	}

	return t, nil
}

// GetToken returns a Github Token that is not rate-limited from
// a pool of Github Apps and Access Tokens.
func GetToken(ctx context.Context, creds []*Credential, logger logrus.FieldLogger) (cfg.SecretData, error) {
	indexes := rand.Perm(len(creds))
	for _, i := range indexes {
		c := creds[i]
		log := logger.WithField("name", c.Name)

		t, err := c.GetToken(ctx)
		if err != nil {
			log.WithError(err).Warn("failed to use credential")
			continue
		}

		log.Info("selected token")
		return t, err
	}

	return "", fmt.Errorf("failed to find non-ratelimited token")
}
