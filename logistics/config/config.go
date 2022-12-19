/*
--------------------------------
 Author: Deepak Rathod
--------------------------------
 Date - 18 August 2022
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
	Port 			 int
	LogisticsService string
	AccountService 	 string
	HubService 		 string
	CityService 	 string
}

var cfg *Config = &Config {
	Port            : 6060,
	LogisticsService: "logisticsservice",
	AccountService  : "accountservice",
	HubService      : "hubservice",
	CityService     : "cityservice",
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