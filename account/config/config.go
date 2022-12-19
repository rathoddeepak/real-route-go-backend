/*
--------------------------------
 Author: Deepak Rathod
--------------------------------
 Date - 17 July 2022
 Phone: +917020814070
 Email: rathoddeepak537@gmail.com 
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
	SessionService string
	LogisticsService string
}

var cfg *Config = &Config {
	Port: 1010,
	SessionService: "sessionservice",
	LogisticsService: "logisticsservice",
}

func Address () string {
	return fmt.Sprintf(":%v", cfg.Port);
}

func Get() Config {
	return *cfg
}


func Load () error {
	configor, err := config.NewConfig(config.WithSource(env.NewSource()));
	if err != nil {
		return errors.Wrap(err, "configor.New");
	}
	if err = configor.Load(); err != nil {
		return errors.Wrap(err, "configor.Load");
	}
	if err = configor.Scan(cfg); err != nil {
		return errors.Wrap(err, "configor.Scan");
	}
	return nil;
}