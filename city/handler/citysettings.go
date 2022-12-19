/*
--------------------------------
 Author: Deepak Rathod
--------------------------------
 Date - 18 August 2022
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

//City Settings
func (city *CityService) SetDeliveryDays (ctx context.Context, in *pb.SetDeliveryDaysRequest, out *pb.UpdateResponse) (error) {
	if(in.CityId == 0){
		return errors.New("City Id is required!");
	}
	var dlvDays db.DeliveryDays;
	err := json.Unmarshal([]byte(in.Data), &dlvDays);
	if err != nil {
		return err;
	}
	_, err = db.SetSubDeliveryDays(in.CityId, &dlvDays);
	if err != nil {
		return err
	}
	out.Status = 200;
	out.Message = "Data Updated!";
	return nil;
}

func (city *CityService) GetDeliveryDays (ctx context.Context, in *pb.GetCityRequest, out *pb.GetDeliveryDaysResponse) (error) {
	id, dlvDays, err := db.GetSubDeliveryDays(in.CityId);
	if err != nil {
		return err
	}
	data, err := json.Marshal(dlvDays)
    if err != nil {
    	return err;
    }
    out.Id = *id;
    out.Data = string(data);
	return nil;
}