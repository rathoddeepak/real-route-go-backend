/*
--------------------------------
 Author: Deepak Rathod
--------------------------------
 Date - 22 August 2022
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
	"fmt"
	"context"

	pb "logisticsService/proto"
    db "logisticsService/database"
)

func (lg *LogisticsService) CreateVehicle (ctx context.Context, in *pb.CreateVehicleRequest, out *pb.CreateVehicleResponse) (error) {
    vehicleId, err := db.InsertVehicle(in);
    if err != nil {
    	return err;
    }
    out.VehicleId = *vehicleId;
    return nil;
}

func (lg *LogisticsService) UpdateVehicle (ctx context.Context, in *pb.UpdateVehicleRequest, out *pb.UpdateVehicleResponse) (error) {
    err := db.UpdateVehicle(in);
    if err != nil {
    	return err;
    }
    out.Status = 200;
    out.Message = UPDATED_MSG;
    return nil;
}

func (lg *LogisticsService) UpdateVehicleImage (ctx context.Context, in *pb.UpdateVehicleRequest, out *pb.UpdateVehicleResponse) (error) {
	//TODO delete previous image
    err := db.UpdateVehicleImage(in);
    if err != nil {
    	return err;
    }
    out.Status = 200;
    out.Message = UPDATED_MSG;
    return nil;
}

func (lg *LogisticsService) UpdateVehicleLocation (ctx context.Context, in *pb.UpdateVehicleRequest, out *pb.UpdateVehicleResponse) (error) {
	point := fmt.Sprintf(PointString, in.Lat, in.Lng);
    err := db.UpdateVehicleLocation(point, in.VehicleId);
    if err != nil {
    	return err;
    }    
    out.Status = 200;
    out.Message = UPDATED_MSG;
    return nil;
}

func (lg *LogisticsService) GetVehicle (ctx context.Context, in *pb.GetVehicleRequest, out *pb.GetVehicleResponse) (error) {
	vehicle, err := db.GetVehicleById(in.VehicleId);
	if err != nil {
		return err
	}
    out.Vehicle = vehicle;
	return nil;	 
}

func (lg *LogisticsService) GetVehiclesOfCity (ctx context.Context, in *pb.GetVehicleRequest, out *pb.GetVehiclesResponse) (error) {
	vehicles, err := db.GetVehiclesOfCity(in.CityId);
	if err != nil {
		return err
	}
    out.Vehicles = *vehicles;
	return nil;	 
}