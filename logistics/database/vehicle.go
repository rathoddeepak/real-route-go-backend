/*
--------------------------------
 Author: Deepak Rathod
--------------------------------
 Date - 16 September 2022
 Phone: +917020814070
 Email: rathoddeepak537@gmail.com 
--------------------------------
 --->   Logistics System   <---
--------------------------------
*/

package database

import (
	"fmt"
	_ "github.com/lib/pq"
	"database/sql"
	"errors"
	"encoding/json"

	"strings"
	"time"

	sqlBuilder "github.com/doug-martin/goqu/v9"
	pb "logisticsService/proto"
)

//TODO: Shift to live database!
func GetVehicleTypes () *[]*pb.VehicleType {
	const VEHICLE_BICYCLE = 1;	
	const VEHICLE_SCOOTY = 2;
	const VEHICLE_BIKE = 3;
	const VEHICLE_MOTORBIKE = 4;
	const VEHICLE_CAR = 5;
	const VEHICLE_SUV = 6;
	const VEHICLE_BUS = 7;
	const VEHICLE_TRUCK = 8;
	var vehicleTypes []*pb.VehicleType;
	vehicleTypes = append(vehicleTypes, &pb.VehicleType {
		Id: VEHICLE_BICYCLE,
		Name: "Bicycle",
	});
	vehicleTypes = append(vehicleTypes, &pb.VehicleType {
		Id: VEHICLE_MOTORBIKE,
		Name: "MotorBike",
	});
	vehicleTypes = append(vehicleTypes, &pb.VehicleType {
		Id: VEHICLE_CAR,
		Name: "Car",
	});
	vehicleTypes = append(vehicleTypes, &pb.VehicleType {
		Id: VEHICLE_SUV,
		Name: "XUV",
	});
	vehicleTypes = append(vehicleTypes, &pb.VehicleType {
		Id: VEHICLE_BUS,
		Name: "BUS",
	});
	vehicleTypes = append(vehicleTypes, &pb.VehicleType {
		Id: VEHICLE_TRUCK,
		Name: "Truck",
	});
	return &vehicleTypes;
}

func InsertVehicle (in *pb.CreateVehicleRequest) (*int64, error) {
	point := fmt.Sprintf(PointString, 0, 0);
	record := sqlBuilder.Record {
		"company_id" : in.CompanyId,
		"city_id" : in.CityId,
		"name" : in.Name,
		"type" : in.Type,		
		"created":time.Now().Unix(),
		"status":1,
		"location" : point,
	}
	sql, _, _ := sqlBuilder.Insert(VEHICLE_TABLE).Rows(record).Returning("id").ToSQL();	
	var vehicleId int64;
	sql = strings.Replace(sql, "'|", "" , 1);
	sql = strings.Replace(sql, "|'", "" , 1);
	sql = strings.ReplaceAll(sql, "''", "'" );
	err := db.QueryRow(sql).Scan(&vehicleId);
	if err != nil {
		return nil, err;
	}
	return &vehicleId, nil
}

func UpdateVehicle (in *pb.UpdateVehicleRequest) (error) {
	sql, _, _ := sqlBuilder.Update(VEHICLE_TABLE).Set(
		sqlBuilder.Record {
			"name" : in.Name,
			"type" : in.Type,
		},
	).Where(sqlBuilder.Ex {"id": in.VehicleId}).ToSQL();
	rows, err := db.Query(sql);
	defer rows.Close();
	if err != nil {
		fmt.Println(err);
		return err;
	}
	return nil;
}

func UpdateVehicleLocation(geolocation string, vehicle_id int64) (error){	
	sql, _, _ := sqlBuilder.Update(VEHICLE_TABLE).Set(
		sqlBuilder.Record {			
			"location" : geolocation,
		},
	).Where(sqlBuilder.Ex {"id": vehicle_id}).ToSQL();
	sql = strings.Replace(sql, "'|", "" , 1);
	sql = strings.Replace(sql, "|'", "" , 1);
	sql = strings.ReplaceAll(sql, "''", "'" );
	rows, err := db.Query(sql);	
	if err != nil {
		return err;
	}else{
		defer rows.Close();
	}
	return nil;
}

func UpdateVehicleImage (in *pb.UpdateVehicleRequest) (error) {
	sql, _, _ := sqlBuilder.Update(VEHICLE_TABLE).Set(
		sqlBuilder.Record {
			"photo" : in.Image,
		},
	).Where(sqlBuilder.Ex {"id": in.VehicleId}).ToSQL();
	rows, err := db.Query(sql);
	defer rows.Close();
	if err != nil {
		fmt.Println(err);
		return err;
	}
	return nil;
}

func GetVehiclesOfCity (city_id int64) (*[]*pb.Vehicle, error){
	query, _, _ := sqlBuilder.Select(GetVehicleColumns()...).From(VEHICLE_TABLE).Where (
		sqlBuilder.Ex{"city_id": city_id},
	).ToSQL();
	query = strings.Replace(query, geostring_from, geostring_to, 1);
	rows, err := db.Query(query);
	if err != nil {
		return nil, err;
	}
	defer rows.Close();
	var vehicles []*pb.Vehicle;
	for rows.Next() {
		if err := rows.Err(); err != nil {			
			continue
		}
		vehicle, err := StructVehicle(rows);
		if err == nil {			
			vehicles = append(vehicles, vehicle);
		}else{			
			continue
		}
	}	
	return &vehicles, nil
}

func GetVehicleById (vehicle_id int64) (*pb.Vehicle, error){
	query, _, _ := sqlBuilder.From(VEHICLE_TABLE).Select(GetVehicleColumns()...).Where (
		sqlBuilder.Ex{"id": vehicle_id},
	).Limit(1).ToSQL();
	query = strings.Replace(query, geostring_from, geostring_to, 1);
	rows, err := db.Query(query);
	if err != nil {
		return nil, err;
	}
	defer rows.Close();
	present := rows.Next();
	if present == false {
		return nil, errors.New(errText);
	}
	if err := rows.Err(); err != nil {
		return nil, err;
	}

	vehicle, err := StructVehicle(rows);
	if err != nil {
		return nil, err
	}
	return vehicle, nil
}

func StructVehicle (row *sql.Rows) (*pb.Vehicle, error) {
	var vehicle pb.Vehicle;
	var location GeoPoint;
	var JSONLocation []byte;
	var nullPhoto sql.NullString;
	err := row.Scan(	
		&vehicle.Id,	
		&vehicle.CityId,	
		&vehicle.CompanyId,
		&vehicle.Name,
		&nullPhoto,
		&vehicle.Type,
		&JSONLocation,
		&vehicle.Status,
		&vehicle.Created,
	);
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(JSONLocation, &location);
	if err != nil {
		return nil, err
	}
	vehicle.Photo = nullPhoto.String;
	vehicle.Lat = location.Coordinates[0];
	vehicle.Lng = location.Coordinates[1];
	return &vehicle, nil;
}