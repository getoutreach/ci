// Copyright 2021 Outreach Corporation. All Rights Reserved.

// Description: Implements the "ghaccesstoken token" command. This command
// returns a token from a pool of Github Apps, ensuring that the token is good
// for _at least_ 5 API requests. It is the caller's responsibility to ensure
// that, when the token is rate limited, a new token is requested or to provide
// enough Github Applications that the pool is not easily exhausted.
package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/getoutreach/ci/internal/github"
	"github.com/getoutreach/gobox/pkg/cfg"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// filterPrefix iterates over a set of strings and returns only
// those with the provided prefix
func filterPrefix(prefix string, strs []string) []string {
	filtered := make([]string, 0)
	for _, s := range strs {
		if strings.HasPrefix(s, prefix) {
			filtered = append(filtered, s)
		}
	}
	return filtered
}

// NewTokenCmd creates a new "ghaccesstoken token" command
func NewTokenCmd(log logrus.FieldLogger) *cli.Command {
	return &cli.Command{
		Name:        "token",
		Description: "Returns a valid, non-ratelimited token for use with the Github API",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "env-prefix",
				EnvVars: []string{"GHACCESSTOKEN_ENV_PREFIX"},
				Usage:   "Prefix to use when looking for credentials. Environment variable format is: appID,installID,AccessToken_OR_PEM",
				Value:   "GHACCESSTOKEN_GHAPP",
			},
		},
		Action: func(c *cli.Context) error {
			envPrefix := c.String("env-prefix")
			creds := make([]*github.Credential, 0)

			envVars := filterPrefix(envPrefix, os.Environ())
			sort.Strings(envVars)

			// Discover all env vars matching the specified prefix
			for _, env := range envVars {
				// remove the value
				env = strings.Split(env, "=")[0]

				envV := os.Getenv(env)
				if envV == "" {
					return fmt.Errorf("env '%s' was empty", env)
				}

				spl := strings.SplitN(envV, ",", 3)
				if len(spl) != 3 {
					return fmt.Errorf("processing env '%s': expected format appID,creds but didn't match", env)
				}

				appID := 0
				if spl[0] != "" {
					var err error
					appID, err = strconv.Atoi(spl[0])
					if err != nil {
						return errors.Wrapf(err, "processing env '%s': appID wasn't able to be turned into an int", env)
					}
				}

				installID := 0
				if spl[1] != "" {
					var err error
					installID, err = strconv.Atoi(spl[1])
					if err != nil {
						return errors.Wrapf(err, "processing env '%s': installID wasn't able to be turned into an int", env)
					}
				}

				log.Printf("Using credentials from environment variable: %s", env)

				// Create the cred. If appID is it's zero value then we assume that this is a PAT.
				cred := &github.Credential{
					Name: env,
				}

				tokenOrPem := spl[2]
				b, err := base64.StdEncoding.DecodeString(tokenOrPem)
				if err != nil {
					return errors.Wrap(err, "failed to decode access_token/pem as base64")
				}
				tokenOrPem = string(b)

				if appID != 0 {
					cred.AppID = &appID
					cred.InstallID = &installID
					cred.PEM = cfg.SecretData(tokenOrPem)
				} else {
					cred.AccessToken = cfg.SecretData(tokenOrPem)
				}

				creds = append(creds, cred)
			}

			t, err := github.GetToken(context.Background(), creds, log)
			if err != nil {
				return err
			}

			fmt.Println(string(t))
			return nil
		},
	}
}
