// Copyright 2021 Outreach Corporation. All Rights Reserved.

// Description: This file is the focal point of configuration that needs passed
// to various parts of the service.
// Managed: true

package ci //nolint:revive // Why: This nolint is here just in case your project name contains any of [-_].

import (
	"context"
	"os"

	"github.com/getoutreach/gobox/pkg/cfg"
	"github.com/getoutreach/gobox/pkg/events"
	"github.com/getoutreach/gobox/pkg/log"
)

// Config tracks config needed for ci
type Config struct {
	ListenHost string `yaml:"ListenHost"`
	HTTPPort   int    `yaml:"HTTPPort"`
	///Block(config)
	///EndBlock(config)
}

// MarshalLog can be used to write config to log
func (c *Config) MarshalLog(addfield func(key string, value interface{})) {
	///Block(marshalconfig)
	///EndBlock(marshalconfig)
}

// LoadConfig returns a new Config type that has been loaded in accordance to the environment
// that the service was deployed in, with all necessary tweaks made before returning.
func LoadConfig(ctx context.Context) *Config { //nolint: funlen // Why: This function is long for extensibility reasons since it is generated by bootstrap.
	// NOTE: Defaults should generally be set in the config
	// override jsonnet file: deployments/ci/ci.config.jsonnet
	c := Config{
		// Defaults to [::]/0.0.0.0 which will broadcast to all reachable
		// IPs on a server on the given port for the respective service.
		ListenHost: "",
		HTTPPort:   8000,
		///Block(defconfig)
		///EndBlock(defconfig)
	}

	// Attempt to load a local config file on top of the defaults
	if err := cfg.Load("ci.yaml", &c); os.IsNotExist(err) {
		log.Info(ctx, "No configuration file detected. Using default settings")
	} else if err != nil {
		log.Error(ctx, "Failed to load configuration file, will use default settings", events.NewErrorInfo(err))
	}

	// Do any necessary tweaks/augmentations to your configuration here
	///Block(configtweak)
	///EndBlock(configtweak)

	log.Info(ctx, "Configuration data of the application:\n", &c)

	return &c
}