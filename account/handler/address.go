/*
--------------------------------
 Author: Deepak Rathod
--------------------------------
 Date - 18 August 2022
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

	pb "accountservice/proto"
	db "accountservice/database"
)

const pointString string = "|ST_GeometryFromText('POINT(%v %v)')|";
//Address Related Functions
func (as *AccountService) CreateAddress (ctx context.Context, in *pb.CreateAddressRequest, out *pb.CreateAddressResponse) (error) {
	point := fmt.Sprintf(pointString, in.Lat, in.Lng);
	addressObj := &db.Address {
		Address: in.Address,
		UserId: in.UserId,
		Landmark: in.Landmark,
		Flat: in.Flat,
		Geolocation: point,
		Type: in.Type,		
	}
	address_id, err := db.InsertAddress(addressObj);
	if err != nil {
		return err
	}
	out.AddressId = *address_id;
	return nil;
}

func (as *AccountService) UpdateAddress (ctx context.Context, in *pb.UpdateAddressRequest, out *pb.UpdateAddressResponse) (error) {
	point := fmt.Sprintf(pointString, in.Lat, in.Lng);
	addressObj := &db.Address {
		Id: 	in.AddressId,
		Address: in.Address,
		Status: in.Status,
		Landmark: in.Landmark,
		Flat: in.Flat,
		Geolocation: point,
		Type: in.Type,
	}
	err := db.UpdateAddress(addressObj);
	if err != nil {
		return err
	}
	out.Status = 200;
	out.Message = "Updated Successfully!";
	return nil;
}

func (as *AccountService) UpdateAddressLocation (ctx context.Context, in *pb.UpdateAddressLocationRequest, out *pb.UpdateAddressResponse) (error) {
	point := fmt.Sprintf(pointString, in.Lat, in.Lng);
	address := &db.Address {
		Id: 	in.AddressId,
		Geolocation: point,
	}
	err := db.UpdateAddressLocaiton(address);
	if err != nil {
		return err
	}
	out.Status = 200;
	out.Message = "Updated Successfully!";
	return nil;
}

func (as *AccountService) UpdateAddressStatus (ctx context.Context, in *pb.UpdateAddressStatusRequest, out *pb.UpdateAddressResponse) (error) {
	address := &db.Address {
		Id: 	in.AddressId,
		Status: in.Status,
	}
	err := db.UpdateAddressStatus(address);
	if err != nil {
		return err
	}
	out.Status = in.Status;
	out.Message = "Updated Successfully!";
	return nil;
}

func (as *AccountService) GetAddressesOfUser (ctx context.Context, in *pb.GetAddressRequest, out *pb.GetAddressesResponse) (error) {
	mAddresses, err := db.GetAddressesOfUser(in.UserId);
	if err != nil {
		return err
	}
	var addresses []*pb.Address;
	for _, mAddress := range *mAddresses {
	   address := makeProtoAddress(mAddress);
       addresses = append(addresses, address);
    }
    out.Addresses = addresses;
	return nil;	
}

func (as *AccountService) GetAddressById (ctx context.Context, in *pb.GetAddressRequest, out *pb.GetAddressResponse) (error) {
	mAddress, err := db.GetAddressById(in.AddressId);
	if err != nil {
		return err
	}
	address := makeProtoAddress(mAddress);
    out.Address = address;
	return nil;	
}

func makeProtoAddress(mAddress *db.Address) *pb.Address {
	return &pb.Address {
	   	Id 	   		: 	mAddress.Id,
	   	UserId 		: 	mAddress.UserId,
	   	Address 	: 	mAddress.Address,
	   	Landmark 	: 	mAddress.Landmark,
	   	Flat 		: 	mAddress.Flat,
	   	Type 		: 	mAddress.Type,
	   	Status 		: 	mAddress.Status,
	   	Lat  		: 	mAddress.Location.Coordinates[0],
	   	Lng  		: 	mAddress.Location.Coordinates[1],
	   	Created 	: 	mAddress.Created,
	}
}