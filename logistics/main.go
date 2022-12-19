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

package main

import (
	
	grpcc "github.com/go-micro/plugins/v4/client/grpc"
	grpcs "github.com/go-micro/plugins/v4/server/grpc"
	
	_ "github.com/go-micro/plugins/v4/registry/kubernetes"

	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"

	"logisticsService/config"
	"logisticsService/handler"
	pb "logisticsService/proto"


	"go-micro.dev/v4/broker"
	"go-micro.dev/v4/cmd"
)

var (
	service = "logisticsservice"
	version = "latest"
)

func main () {
	//Load Configurations
	if err := config.Load(); err != nil {
		logger.Fatal(err);
	}


	//Create Service
	srv := micro.NewService (
		micro.Server(grpcs.NewServer()),
		micro.Client(grpcc.NewClient()),
	);

	srv.Init(
		micro.Name(service),
		micro.Version(version),
		micro.Address(config.Address()),
	);

	//Register Handler
	cfg, client := config.Get(), srv.Client();
	logistics := &handler.LogisticsService {
		HubService: pb.NewHubService(cfg.HubService, client),
		AccountService: pb.NewAccountService(cfg.AccountService, client),
		CityService: pb.NewCityService(cfg.CityService, client),
	}
	if err := pb.RegisterLogisticsServiceHandler(srv.Server(), logistics); err != nil {
		logger.Fatal(err);
	}
	if err := pb.RegisterHealthHandler(srv.Server(), new(handler.Health)); err != nil {
		logger.Fatal(err);
	}

	//Kafka Setup
	cmd.Init();

	if err := broker.Init(); err != nil {
		logger.Fatal(err);
	}

	if err := broker.Connect(); err != nil {
		logger.Fatal(err); 
	}

	//Run Service
	if err := srv.Run(); err != nil {
		logger.Fatal(err);
	}

}