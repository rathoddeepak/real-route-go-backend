/*
--------------------------------
 Author: Deepak Rathod
--------------------------------
 Date - 04 August 2022
 Phone: +917020814070
 Email: rathoddeepak537@gmail.com 
--------------------------------
 --->   City Microservice  <---
--------------------------------
 --->   Logistics System   <---
--------------------------------
*/

package handler;

import (
	"context"
	"errors"

	"go-micro.dev/v4/logger"

	pb "cityservice/proto"
	db "cityservice/database"
)

type CityService struct {}

//City Related Functions
func (city *CityService) CreateCity (ctx context.Context, in *pb.CreateCityRequest, out *pb.CreateCityResponse) (error) {
	if(in.CompanyId == 0){
		return errors.New("Company Id is required!");
	}
	cityObj := &db.City{
		CompanyId: in.CompanyId,
		Name: in.Name,
	}
	logger.Info(in.CompanyId)
	city_id, err := db.InsertCity(cityObj);
	if err != nil {
		return err
	}
	out.CityId = city_id;
	
	logger.Info("City Created!");

	return nil;
}

func (city *CityService) UpdateCity (ctx context.Context, in *pb.UpdateCityRequest, out *pb.UpdateResponse) (error) {
	err := db.UpdateCity(in.CityId, in.Name, in.Status);
	if err != nil {
		return err
	}
	out.Status = in.Status;
	out.Message = "Updated Successfully!";
	return nil;
}

func (city *CityService) GetCities (ctx context.Context, in *pb.GetCityRequest, out *pb.GetCitiesResponse) (error) {
	mCities, err := db.GetCompanyCities(in.CompanyId);
	if err != nil {
		return err
	}
	var cities []*pb.City;
	for _, mCity := range *mCities {
	   city := pb.City {
	   	Id : mCity.Id,
	   	Name : mCity.Name,
	   	CompanyId : mCity.CompanyId,
	   	Status : mCity.Status,
	   	Created: mCity.Created,
       }
       cities = append(cities, &city);
    }
    out.Cities = cities;
	return nil;	
}

func (cs *CityService) GetCityById (ctx context.Context, in *pb.GetCityRequest, out *pb.GetCityResponse) (error) {
	city, err := db.GetCityById(in.CityId);
	if err != nil {
		return err
	}
	out.City = &pb.City {
		Id: city.Id,
		CompanyId: city.CompanyId,
		Name: city.Name,
		Status: city.Status,
		Created: city.Created,
	}
	return nil;
}

func (cs *CityService) GetCityByGeoPoint (ctx context.Context, in *pb.GetCityRequest, out *pb.GetCityResponse) (error) {
	fence, err := db.GetFenceByGeoPoint(in.Lat, in.Lng);
	if err != nil {
		return err
	}
	city, err := db.GetCityById(fence.CityId);
	if err != nil {
		return err
	}
	out.City = &pb.City {
		Id: city.Id,
		CompanyId: city.CompanyId,
		Name: city.Name,
		Status: city.Status,
		Created: city.Created,
	}
	return nil;
}