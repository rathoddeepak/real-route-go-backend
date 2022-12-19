/*
--------------------------------
 Author: Deepak Rathod
--------------------------------
 Date - 04 August 2022
 Phone: +917020814070
 Email: rathoddeepak537@gmail.com 
--------------------------------
 --->   Hub Microservice   <---
--------------------------------
--------------------------------
 --->   Logistics System   <---
--------------------------------
*/

package main

import (
	
	grpcc "github.com/go-micro/plugins/v4/client/grpc"
	grpcs "github.com/go-micro/plugins/v4/server/grpc"

	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"

	"hubservice/config"
	"hubservice/handler"

	pb "hubservice/proto"
)

var (
	service = "hubservice"
	version = "lastest"
)

func main () {
	//Load Configurations
	if err := config.Load(); err != nil {
		logger.Fatal(err);
	}

    //Create Service
    srv := micro.NewService(
		micro.Server(grpcs.NewServer()),
		micro.Client(grpcc.NewClient()),
	)
	srv.Init(
		micro.Name(service),
		micro.Version(version),
		micro.Address(config.Address()),
	)

	//Register Handler
	if err := pb.RegisterHubServiceHandler(srv.Server(), new(handler.HubService)); err != nil {
		logger.Fatal(err);
	}
	if err := pb.RegisterHealthHandler(srv.Server(), new(handler.Health)); err != nil {
		logger.Fatal(err);
	}
	
	//Run Service
	if err := srv.Run(); err != nil {
		logger.Fatal(err);
	}
}