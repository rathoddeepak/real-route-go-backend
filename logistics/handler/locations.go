/*
--------------------------------
 Author: Deepak Rathod
--------------------------------
 Date - 05 September 2022
 Phone: +917020814070
 Email: rathoddeepak537@gmail.com 
--------------------------------
 --->   Logistics System   <---
--------------------------------
        RealRoute Code:
--------------------------------
*/

package handler;

import (
	"fmt"
	"context"

	pb "logisticsService/proto"
    	db "logisticsService/database"
)

func (lg *LogisticsService) CreateLocation (ctx context.Context, in *pb.CreateLocationRequest, out *pb.CreateLocationResponse) (error) {
    location := &db.Location {
    	CityId   : in.CityId,
    	Name     : in.Name,
    	Address  : in.Address,
    	Contact  : in.Contact,
    }
    point := fmt.Sprintf(PointString, in.Lat, in.Lng);
    locationId, err := db.InsertLocation(point, location);
    if err != nil {
    	return err;
    }
    out.LocationId = *locationId;
    return nil;
}

func (lg *LogisticsService) UpdateLocation (ctx context.Context, in *pb.UpdateLocationRequest, out *pb.UpdateLocationResponse) (error) {
	point := fmt.Sprintf(PointString, in.Lat, in.Lng);
	location := &db.Location {
    	Id  	 : in.LocationId,
    	Name     : in.Name,
    	Address  : in.Address,
    	Contact  : in.Contact,
    }
    err := db.UpdateLocation(point, location);
    if err != nil {
    	return err;
    }
    out.Status = 200;
    out.Message = UPDATED_MSG;
    return nil;
}

func (lg *LogisticsService) CityLocations (ctx context.Context, in *pb.GetLocationRequest, out *pb.GetLocationsResponse) (error) {
	mLocations, err := db.GetLocationsOfCity(in.CityId);
	if err != nil {
		return err
	}
	var locations []*pb.Location;
	for _, mLocation := range *mLocations {
	   location := makeProtoLocation(mLocation);
       locations = append(locations, location);
    }
    out.Locations = locations;
	return nil;	 
}

func (lg *LogisticsService) GetLocationById (ctx context.Context, in *pb.GetLocationRequest, out *pb.GetLocationResponse) (error) {
	mLocation, err := db.GetLocationById(in.LocationId);
	if err != nil {
		return err
	}
	location := makeProtoLocation(mLocation);
    out.Location = location;
	return nil;
}

func makeProtoLocation (location *db.Location) (*pb.Location){
	return &pb.Location {
		Id       : location.Id,
		CityId   : location.CityId,
		Name     : location.Name,
		Address  : location.Address,
		Contact  : location.Contact,
		Lat   	 : location.Location.Coordinates[0],
		Lng      : location.Location.Coordinates[1],
	}
}