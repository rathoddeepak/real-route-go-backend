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

package main

import (
	grpcc "github.com/go-micro/plugins/v4/client/grpc"
	grpcs "github.com/go-micro/plugins/v4/server/grpc"
	
	_ "github.com/go-micro/plugins/v4/registry/kubernetes"

	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"

	"sessionservice/config"
	"sessionservice/handler"

	pb "sessionservice/proto"
)

var (
	service = "sessionservice"
	version = "0.0.1"
)

func main () {
	//Load Config
	if err := config.Load(); err != nil {
		logger.Fatal(err);
	}


	//Initialize Service
	srv := micro.NewService(
		micro.Server(grpcs.NewServer()),
		micro.Client(grpcc.NewClient()),
	)
	srv.Init(
		micro.Name(service),
		micro.Version(version),
		micro.Address(config.Address()),
	)

	//Register Service
	if err := pb.RegisterSessionServiceHandler(srv.Server(), new(handler.SessionService)); err != nil {
		logger.Fatal(err);
	}

	if err := pb.RegisterHealthHandler(srv.Server(), new(handler.Health)); err != nil {
		logger.Fatal(err);
	}

	//Start Service 
	if err := srv.Run(); err != nil {
		logger.Fatal(err);
	}
}