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

	"encoding/json"

	pb "cityservice/proto"
	db "cityservice/database"
)

//Fence Related Functions
func (cs *CityService) CreateFence (ctx context.Context, in *pb.CreateFenceRequest, out *pb.CreateFenceResponse) (error) {
	fence := &db.Fence{
		CityId: in.CityId,
		HubId: in.HubId,
		Name: in.Name,
	}
	fence_id, err := db.InsertFence(fence, in.Coords);
	if err != nil {
		return err
	}
	out.FenceId = fence_id;
	return nil;
}

func (cs *CityService) UpdateFence (ctx context.Context, in *pb.UpdateFenceRequest, out *pb.UpdateResponse) (error) {
	err := db.UpdateFence(in.FenceId, in.Name, in.Status, in.Polygon);
	if err != nil {
		return err
	}
	out.Status = in.Status;
	out.Message = "Updated Successfully!";
	return nil;
}

func (cs *CityService) GetFenceById (ctx context.Context, in *pb.GetFenceRequest, out *pb.GetFenceResponse) (error) {
	fence, err := db.GetFenceById(in.FenceId);
	if err != nil {
		return err
	}
	res := makeProtoFence(fence);
	if res == nil {
		return errors.New("Fence Not Found!");
	}
	out.Fence = res;
	return nil;
}

func (cs *CityService) GetFenceByGeoPoint (ctx context.Context, in *pb.GetFenceRequest, out *pb.GetFenceResponse) (error) {
	fence, err := db.GetFenceByGeoPoint(in.Lat, in.Lng);
	if err != nil {
		return err
	}	
	res := makeProtoFence(fence);
	if res == nil {
		return errors.New("Fence Not Found!");
	}
	out.Fence = res;
	return nil;
}

func (cs *CityService) GetFenceByGeoPointAndId (ctx context.Context, in *pb.GetFenceRequest, out *pb.GetFenceResponse) (error) {
	fence, err := db.GetFenceByGeoPointAndId(in.HubId, in.Lat, in.Lng);
	if err != nil {
		return err
	}	
	res := makeProtoFence(fence);
	if res == nil {
		return errors.New("Fence Not Found!");
	}
	out.Fence = res;
	return nil;
}

func (cs *CityService) GetFences (ctx context.Context, in *pb.GetFenceRequest, out *pb.GetFencesResponse) (error) {
	mfences, err := db.GetCityFences(in.CityId);
	if err != nil {
		return err
	}
	var fences []*pb.Fence;
	for _, mFence := range *mfences {
		fence := makeProtoFence(mFence);
		if fence == nil {
			continue
		}
        fences = append(fences, fence);
    }
    out.Fences = fences;
	return nil;	
}

func (cs *CityService) GetHubFences (ctx context.Context, in *pb.GetFenceRequest, out *pb.GetFencesResponse) (error) {
	mfences, err := db.GetHubFences(in.HubId);
	if err != nil {
		return err
	}
	var fences []*pb.Fence;
	for _, mFence := range *mfences {
		fence := makeProtoFence(mFence);
		if fence == nil {
			continue
		}
        fences = append(fences, fence);
    }
    out.Fences = fences;
	return nil;	
}

func makeProtoFence(mFence *db.Fence) *pb.Fence {
	polygon, err := json.Marshal(mFence.Polygon);
	if err != nil {
		return nil;
	}
	return &pb.Fence {	   	
		Id:mFence.Id,
		CityId:mFence.CityId,
		HubId:mFence.HubId,
		Name:mFence.Name,
		Status:mFence.Status,
		Polygon:string(polygon),
		Created:mFence.Created,
	}
}