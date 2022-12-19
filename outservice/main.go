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

	"context"
	"encoding/json"

	"outservice/config"
	"outservice/handler"
	pb "outservice/proto"

	"go-micro.dev/v4"
	"go-micro.dev/v4/cmd"
	"go-micro.dev/v4/broker"
	"go-micro.dev/v4/logger"

	_ "github.com/go-micro/plugins/v4/broker/kafka"
)

var (
	smsTopic = "go.micro.topic.sendotp"
	notificationTopic = "go.micro.topic.sendnotification"
)

const (
	service = "outservice"
	version = "latest"
)

func main () {
	//Load Config
	if err := config.Load(); err != nil {
		logger.Fatal(err);
	}

	//Register Handler
	srv := micro.NewService(
		micro.Server(grpcs.NewServer()),
		micro.Client(grpcc.NewClient()),
	)
	srv.Init(
		micro.Name(service),
		micro.Version(version),
		micro.Address(config.Address()),
	)
	outService := &handler.OutService {}
	if err := pb.RegisterOutServiceHandler(srv.Server(), outService); err != nil {
		logger.Fatal(err);
	}
	if err := pb.RegisterHealthHandler(srv.Server(), new(handler.Health)); err != nil {
		logger.Fatal(err);
	}

	cmd.Init();

	if err := broker.Init(); err != nil {
		logger.Fatal(err);
	}

	if err := broker.Connect(); err != nil {
		logger.Fatal(err); 
	}

	_, err := broker.Subscribe(smsTopic, func(p broker.Event) error {
		var req pb.SendSMSRequest;
		err := json.Unmarshal(p.Message().Body, &req);
		if err != nil {
			logger.Fatal(err);
		}else{
			outService.SendSMS(context.Background(), &req, nil);
		}		
		return nil
	});

	_, err = broker.Subscribe(notificationTopic, func(p broker.Event) error {
		var req pb.SendNotificaitonRequest;
		err := json.Unmarshal(p.Message().Body, &req);
		if err != nil {
			logger.Fatal(err);
		}else{
			outService.SendMobileNotification(context.Background(), &req, nil);
		}		
		return nil
	});

	if err != nil {
		logger.Fatal(err);
	}

	if err := srv.Run(); err != nil {
		logger.Fatal(err);
	}
}