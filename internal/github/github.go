package github

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/getoutreach/gobox/pkg/cfg"
	gc "github.com/google/go-github/v34/github"
	"github.com/pkg/errors"
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

	resetAtDur := time.Until(rl.Core.Reset.Time)
	fmt.Fprintf(os.Stderr, "Credential Status for %s: %d calls remaining for %s\n", c.Name, rl.Core.Remaining, resetAtDur)

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
func GetToken(ctx context.Context, creds []*Credential) (cfg.SecretData, error) {
	errs := make([]error, 0)
	for _, c := range creds {
		t, err := c.GetToken(ctx)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		return t, err
	}

	return "", fmt.Errorf("failed to find non-ratelimited token: %v", errs)
}
