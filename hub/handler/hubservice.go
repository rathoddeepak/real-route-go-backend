/*
--------------------------------
 Author: Deepak Rathod
--------------------------------
 Date - 05 August 2022
 Phone: +917020814070
 Email: rathoddeepak537@gmail.com 
--------------------------------
 --->   Hub Microservice  <---
--------------------------------
--------------------------------
 --->   Logistics System   <---
--------------------------------
*/

package handler;

import (

	"fmt"

	"context"

	"go-micro.dev/v4/logger"

	pb "hubservice/proto"
	db "hubservice/database"
)

type HubService struct {}
const pointString string = "|ST_GeometryFromText('POINT(%v %v)')|";
//City Related Functions
func (hs *HubService) CreateHub (ctx context.Context, in *pb.CreateHubRequest, out *pb.CreateHubResponse) (error) {
	point := fmt.Sprintf(pointString, in.Lat, in.Lng);
	hubObj := &db.Hub{
		CityId: in.CityId,
		UserId: in.UserId,
		Name: in.Name,
		Geolocation: point,
		Phone: in.Phone,
		Address: in.Address,
		About: in.About,
	}
	hub_id, err := db.InsertHub(hubObj);
	if err != nil {
		return err
	}
	out.HubId = *hub_id;
	logger.Info("Hub Created!");
	return nil;
}

func (hs *HubService) UpdateHub (ctx context.Context, in *pb.UpdateHubRequest, out *pb.UpdateHubResponse) (error) {
	point := fmt.Sprintf(pointString, in.Lat, in.Lng);
	hub := &db.Hub {
		Id: 	in.HubId,
		Name: 	in.Name,
		About: 	in.About,
		Phone: 	in.Phone,
		Address: in.Address,
		Geolocation: point,
		Status: in.Status,
	}
	err := db.UpdateHub(hub);
	if err != nil {
		return err
	}
	out.Status = 200;
	out.Message = "Updated Successfully!";
	return nil;
}

func (hs *HubService) UpdateHubLocation (ctx context.Context, in *pb.UpdateHubLocationRequest, out *pb.UpdateHubResponse) (error) {
	point := fmt.Sprintf(pointString, in.Lat, in.Lng);
	hub := &db.Hub {
		Id: 	in.HubId,
		Geolocation: point,
	}
	err := db.UpdateHubLocaiton(hub);
	if err != nil {
		return err
	}
	out.Status = 200;
	out.Message = "Updated Successfully!";
	return nil;
}

func (hs *HubService) UpdateHubStatus (ctx context.Context, in *pb.UpdateHubStatusRequest, out *pb.UpdateHubResponse) (error) {
	hub := &db.Hub {
		Id: 	in.HubId,
		Status: in.Status,
	}
	err := db.UpdateHubStatus(hub);
	if err != nil {
		return err
	}
	out.Status = in.Status;
	out.Message = "Updated Successfully!";
	return nil;
}

func (hs *HubService) GetHubsOfCity (ctx context.Context, in *pb.GetHubRequest, out *pb.GetHubsResponse) (error) {
	mHubs, err := db.GetHubsOfCity(in.CityId);
	if err != nil {
		return err
	}
	var hubs []*pb.Hub;
	for _, mHub := range *mHubs {
	   hub := makeProtoHub(mHub);
       hubs = append(hubs, hub);
    }
    out.Hubs = hubs;
	return nil;	
}

func (hs *HubService) GetHubsOfUser (ctx context.Context, in *pb.GetHubRequest, out *pb.GetHubsResponse) (error) {
	mHubs, err := db.GetHubsOfUser(in.UserId);
	if err != nil {
		return err
	}
	var hubs []*pb.Hub;
	for _, mHub := range *mHubs {
	   hub := makeProtoHub(mHub);
       hubs = append(hubs, hub);
    }
    out.Hubs = hubs;
	return nil;	
}

func (hs *HubService) GetHubById (ctx context.Context, in *pb.GetHubRequest, out *pb.GetHubResponse) (error) {
	mHub, err := db.GetHubById(in.HubId);
	if err != nil {
		return err
	}
	hub := makeProtoHub(mHub);
    out.Hub = hub;
	return nil;	
}

func (hs *HubService) HubInitData (ctx context.Context, in *pb.HubInitRequest, out *pb.HubInitResponse) (error) {
	var hubs []*pb.Hub;
	var categories []*pb.Category;
	var subcategories []*pb.SubCategory;
	uniqueHubIds := make(map[int64]int);
	uniqueCatIds := make(map[int64]int);
	mHubs, err := db.GetHubsOfCity(in.CityId);	
	if err != nil {
		return err
	}
	for _, mHub := range *mHubs {
	   hub := makeProtoHub(mHub);
       hubs = append(hubs, hub);
       if uniqueHubIds[mHub.Id] == 0 {
       	uniqueHubIds[mHub.Id] = 1;
       }
    }
    for uHubId, _ := range uniqueHubIds {    	
    	mCategories, err := db.GetCategoriesOfHub(uHubId);
		if err != nil {
			continue
		}
		for _, mCategory := range *mCategories {
		   category := makeProtoCategory(mCategory);
	       categories = append(categories, category);
	       if uniqueCatIds[mCategory.Id] == 0 {
	       	uniqueCatIds[mCategory.Id] = 1;
	       }
	    }
    }
    for uCatId, _ := range uniqueCatIds {
    	mSubCategories, err := db.GetSubCategoriesOfCategory(uCatId);
		if err != nil {
			continue
		}		
		for _, mSubCategory := range *mSubCategories {
		   subcategory := makeProtoSubCategory(mSubCategory);
	       subcategories = append(subcategories, subcategory);
	    }
    }
    out.Hubs = hubs;
    out.Categories = categories;
    out.SubCategories = subcategories;
    return nil;
}

func makeProtoHub(mHub *db.Hub) *pb.Hub {
	return &pb.Hub {
	   	Id 	   	: 	mHub.Id,
	   	CityId 	: 	mHub.CityId,
	   	UserId 	: 	mHub.UserId,
	   	Name 	: 	mHub.Name,
	   	About 	: 	mHub.About,
	   	Address : 	mHub.Address,
	   	Phone 	: 	mHub.Phone,
	   	Status 	: 	mHub.Status,
	   	Lat  	: 	mHub.Location.Coordinates[0],
	   	Lng  	: 	mHub.Location.Coordinates[1],
	   	Created : 	mHub.Created,
	}
}