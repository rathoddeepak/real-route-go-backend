/*
--------------------------------
 Author: Deepak Rathod
--------------------------------
 Date - 04 August 2022
 Phone: +917020814070
 Email: rathoddeepak537@gmail.com 
--------------------------------
 --->   City Microservice  <---
--------------------------------
--------------------------------
 --->   Logistics System   <---
--------------------------------
*/

package config

import (
	"fmt"
	"github.com/pkg/errors"
	"go-micro.dev/v4/config"
	"go-micro.dev/v4/config/source/env"
)

type Config struct {
	Port int
}

var cfg *Config = &Config {
	Port: 4040,
}

func Address () string {
	return fmt.Sprintf(":%v", cfg.Port);
}

func Load () error {
	configor, err := config.NewConfig(config.WithSource(env.NewSource()));

	if err != nil {
		return errors.Wrap(err, "configor.New");
	}

	if configor.Load(); err != nil {
		return errors.Wrap(err, "configor.Load");
	}

	if configor.Scan(cfg); err != nil {
		return errors.Wrap(err, "configor.Scan");
	}

	return nil;
}
