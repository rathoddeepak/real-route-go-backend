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
	"flag"
	"net/http"
	"fmt"
	"context"
	
	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"	
	"google.golang.org/grpc"
	"github.com/rs/cors"


	accountService   "justify_backend/proto/account"
	sessionService   "justify_backend/proto/session"
	cityService      "justify_backend/proto/city"
	hubService       "justify_backend/proto/hub"
	logisticsService "justify_backend/proto/logistics"

	"justify_backend/gateway"
)

var (
	accountEndpoint		= flag.String("account_ep",   "localhost:1010", "go.micro.srv.accountservice address")
	sessionEndpoint 	= flag.String("session_ep",   "localhost:2020", "go.micro.srv.sessionservice address")
	cityEndpoint 		= flag.String("city_ep", 	  "localhost:4040", "go.micro.srv.cityservice address")
	hubEndpoint 		= flag.String("hub_ep", 	  "localhost:5050", "go.micro.srv.hubservice address")
	logisticsEndpoint 	= flag.String("logistics_ep", "localhost:6060", "go.micro.srv.logisticsservice address")
)

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	endpoints := &gateway.EndPoints {
		AccountEndpoint   : accountEndpoint,
		SessionEndpoint   : sessionEndpoint,
		CityEndpoint 	  : cityEndpoint,
		HubEndpoint 	  : hubEndpoint,
		LogisticsEndpoint : logisticsEndpoint,
	}

	err := accountService.RegisterAccountServiceHandlerFromEndpoint(ctx, mux, *accountEndpoint, opts)
	if err != nil {
		return err
	}
	err = sessionService.RegisterSessionServiceHandlerFromEndpoint(ctx, mux, *sessionEndpoint, opts)
	if err != nil {
		return err
	}
	err = cityService.RegisterCityServiceHandlerFromEndpoint(ctx, mux, *cityEndpoint, opts)
	if err != nil {
		return err
	}
	err = hubService.RegisterHubServiceHandlerFromEndpoint(ctx, mux, *hubEndpoint, opts)
	if err != nil {
		return err
	}
	err = logisticsService.RegisterLogisticsServiceHandlerFromEndpoint(ctx, mux, *logisticsEndpoint, opts)
	if err != nil {
		return err
	}
	//For Handling Mix API Call
	err = gateway.SetupClients(endpoints, mux);

	if err != nil {
		return err
	}	
	fmt.Println("Listening on Port 8080");
	//For Handling Static files
	go gateway.ServeStatic();

	handler := cors.AllowAll().Handler(mux);
	return http.ListenAndServe(":8080", handler);
}

func main() {
	flag.Parse()

	defer glog.Flush()

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}