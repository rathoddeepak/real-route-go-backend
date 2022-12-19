/*
--------------------------------
 Author: Deepak Rathod
--------------------------------
 Date - 20 August 2022
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
	"errors"
	"context"
	"time"

	pb "logisticsService/proto"
    db "logisticsService/database"
)
const INVALID_LINK_ERR_MSG = "Invalid Link!";

func (lg *LogisticsService) GetTrackingData (ctx context.Context, in *pb.GetTrackingDataRequest, out *pb.GetTrackingDataResponse) (error) {
	task, err := db.GetTaskById(in.TaskId);
	if err != nil {
		return errors.New(INVALID_LINK_ERR_MSG);
	}

	t1 := time.Unix(task.Created, 0)
    t2 := time.Now();
    //TODO: Make Expire For Link Customized
	if t1.Sub(t2).Hours() > 30 {
		return errors.New("Link Expired");
	} 

	routeTask, err := db.GetRouteTaskById(in.TaskId);
	if err != nil {
		return errors.New(INVALID_LINK_ERR_MSG)
	}

	cityObject, err := lg.CityService.GetCityById(ctx, &pb.GetCityRequest {
		CityId : task.CityId,
	});
	if err != nil {
		return errors.New(INVALID_LINK_ERR_MSG)
	}

	cpy, err := lg.AccountService.GetCompany(ctx, &pb.GetCompanyRequest {
		CompanyId : cityObject.City.CompanyId,
	});
	if err != nil {
		return errors.New(INVALID_LINK_ERR_MSG)
	}
	out.Task = routeTask;
	out.CityId = task.CityId;
	out.Company = &pb.GetTrackingDataResponse_TrackingCompany {
		Name : cpy.Company.Name,
		Contact : cpy.Company.Contact,
	}
    return nil;
}