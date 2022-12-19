/*
--------------------------------
 Author: Deepak Rathod
--------------------------------
 Date - 23 August 2022
 Phone: +917020814070
 Email: rathoddeepak537@gmail.com 
--------------------------------
 ---> Account Microservice <---
--------------------------------
--------------------------------
 --->   Logistics System   <---
--------------------------------
*/

package handler;

import (
	"context"
	pb "logisticsService/proto"
    db "logisticsService/database"
)

func (lg *LogisticsService) CreateRoute (ctx context.Context, in *pb.CreateRouteRequest, out *pb.CreateRouteResponse) (error) {
	route := &db.Route {
    	HubId     : in.HubId,
    	CityId    : in.CityId,
    	Name   	  : in.Name,
    	AgentId   : in.AgentId,
    	VehicleId : in.VehicleId,
    }
    routeId, err := db.InsertRoute(route);
    if err != nil {
    	return err;
    }
    out.RouteId = *routeId;
    return nil;
}

func (lg *LogisticsService) UpdateRoute (ctx context.Context, in *pb.UpdateRouteRequest, out *pb.UpdateRouteResponse) (error) {
	route := &db.Route {
		Id           : in.RouteId,
    	Name         : in.Name,
    	DispatchTime : in.DispatchTime,
    	StartTime    : in.StartTime,
    }
    err := db.UpdateRoute(route);
    if err != nil {
    	return err;
    }
    out.Status = 200;
    out.Message = UPDATED_MSG;
    return nil;
}

func (lg *LogisticsService) AssignRouteAgent (ctx context.Context, in *pb.AssignRouteAgentRequest, out *pb.UpdateRouteResponse) (error) {
    err := db.AssignRouteAgent(in.RouteId, in.AgentId);
    if err != nil {
    	return err;
    }
    out.Status = 200;
    out.Message = UPDATED_MSG;
    return nil;
}

func (lg *LogisticsService) AssignRouteVehicle (ctx context.Context, in *pb.AssignRouteVehicleRequest, out *pb.UpdateRouteResponse) (error) {
    err := db.AssignRouteVehicle(in.RouteId, in.VehicleId);
    if err != nil {
    	return err;
    }
    out.Status = 200;
    out.Message = UPDATED_MSG;
    return nil;
}

func (lg *LogisticsService) GetRoutesOfHub (ctx context.Context, in *pb.GetRouteRequest, out *pb.GetRoutesResponse) (error) {
	mRoutes, err := db.GetRoutesOfHub(in.HubId);
	if err != nil {
		return err
	}
	var routes []*pb.Route;
	for _, mRoute := range *mRoutes {
	   route := makeProtoRoute(mRoute);
	   if route.AgentId != 0 {
	   	agent, err := db.GetAgentById(route.AgentId);
	   	if err == nil {
	   		route.AgentName = agent.Name;
	   	}
	   }
	   if route.VehicleId != 0 {
	   	vehicle, err := db.GetVehicleById(route.VehicleId);
	   	if err == nil {
	   		route.VehicleName = vehicle.Name;
	   	}
	   }
       routes = append(routes, route);
    }
    out.Routes = routes;
	return nil;	 
}

func (lg *LogisticsService) GetRoutesOfCity (ctx context.Context, in *pb.GetRouteRequest, out *pb.GetRoutesResponse) (error) {
	mRoutes, err := db.GetRoutesOfCity(in.CityId);
	if err != nil {
		return err
	}
	var routes []*pb.Route;
	for _, mRoute := range *mRoutes {
	   
	   route := makeProtoRoute(mRoute);	   
	   if route.AgentId != 0 {
	   	agent, err := db.GetAgentById(route.AgentId);
	   	if err == nil {
	   		route.AgentName = agent.Name;
	   	}
	   }

	   if route.VehicleId != 0 {
	   	vehicle, err := db.GetVehicleById(route.VehicleId);
	   	if err == nil {
	   		route.VehicleName = vehicle.Name;
	   	}
	   }


       routes = append(routes, route);
    }
    out.Routes = routes;
	return nil;	 
}

func (lg *LogisticsService) GetRouteById (ctx context.Context, in *pb.GetRouteRequest, out *pb.GetRouteResponse) (error) {
	mRoute, err := db.GetRouteById(in.RouteId);
	if err != nil {
		return err
	}
	rotue := makeProtoRoute(mRoute);
    out.Route = rotue;
	return nil;
}

func makeProtoRoute (route *db.Route) (*pb.Route){
	return &pb.Route {
		Id   		  : route.Id,
		HubId 		  : route.HubId,
		CityId 		  : route.CityId,
		Name 		  : route.Name,
		DispatchTime  : route.DispatchTime,
		StartTime	  : route.StartTime,
		AgentId       : route.AgentId,
		VehicleId     : route.VehicleId,
		DispatchType  : route.DispatchType,
	}
}