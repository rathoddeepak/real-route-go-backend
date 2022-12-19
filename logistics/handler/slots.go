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

	pb "logisticsService/proto"
    db "logisticsService/database"
)

const SLOT_ACTIVE = 0;
const SLOT_INACTIVE = 1;

func (lg *LogisticsService) CreateSlot (ctx context.Context, in *pb.CreateSlotRequest, out *pb.CreateSlotResponse) (error) {
	slot := &db.Slot {
    	CityId  : in.CityId,
    	StartHr : in.StartHr,
    	StartMin: in.StartMin,
    	EndHr   : in.EndHr,
    	EndMin	: in.EndMin, 
    	Title	: in.Title,
    }
    slotId, err := db.InsertSlot(slot);
    if err != nil {
    	return err;
    }
    out.SlotId = *slotId;
    return nil;
}

func (lg *LogisticsService) UpdateSlot (ctx context.Context, in *pb.UpdateSlotRequest, out *pb.UpdateSlotResponse) (error) {
	slot := &db.Slot {
		Id      : in.SlotId,
    	StartHr : in.StartHr,
    	StartMin: in.StartMin,
    	EndHr   : in.EndHr,
    	EndMin	: in.EndMin, 
    	Title	: in.Title,
    	Status	: in.Status,    	
    }
    err := db.UpdateSlot(slot);
    if err != nil {
    	return err;
    }
    out.Status = 200;
    out.Message = UPDATED_MSG;
    return nil;
}

func (lg *LogisticsService) GetCitySlots (ctx context.Context, in *pb.GetSlotRequest, out *pb.GetSlotsResponse) (error) {
	mSlots, err := db.GetSlotsOfCity(in.CityId);
	if err != nil {
		return err
	}
	var slots []*pb.Slot;
	for _, mSlot := range *mSlots {
	   slot := makeProtoSlot(mSlot);
       slots = append(slots, slot);
    }
    out.Slots = slots;
	return nil;	 
}

func (lg *LogisticsService) GetSlotById (ctx context.Context, in *pb.GetSlotRequest, out *pb.GetSlotResponse) (error) {
	mSlot, err := db.GetSlotById(in.SlotId);
	if err != nil {
		return err
	}
	slot := makeProtoSlot(mSlot);
    out.Slot = slot;
	return nil;
}

func (lg *LogisticsService) DeleteSlotById (ctx context.Context, in *pb.GetSlotRequest, out *pb.UpdateSlotResponse) (error) {
	_, err := db.GetSubscriptionsBySlotId(in.SlotId);
	if err == nil {
		return errors.New("Slot is in use!")
	}
	err = db.DeleteSlotById(in.SlotId);
    out.Status = 200;
    out.Message = "Deleted!"
	return nil;
}

func makeProtoSlot (slot *db.Slot) (*pb.Slot){
	return &pb.Slot {
		Id:slot.Id,
		CityId:slot.CityId,
		Title:slot.Title,
		StartHr:slot.StartHr,
		StartMin:slot.StartMin,
		EndHr:slot.EndHr,
		EndMin:slot.EndMin,
		Status:slot.Status,
	}
}