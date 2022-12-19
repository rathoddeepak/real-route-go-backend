/*
--------------------------------
 Author: Deepak Rathod
--------------------------------
 Date - 23 August 2022
 Phone: +917020814070
 Email: rathoddeepak537@gmail.com 
--------------------------------
 --->   Logistics System   <---
--------------------------------
*/
package database;

import (
	"fmt"
	"errors"
	_ "github.com/lib/pq"	
	"database/sql"

	sqlBuilder "github.com/doug-martin/goqu/v9"
)

type Route struct {
	Id 	   			int64 		`json:"id"`
	HubId 			int64 		`json:"hub_id"`
	CityId 			int64 		`json:"city_id"`
	AgentId 		int64 		`json:"agent_id"`
	VehicleId 		int64 		`json:"vehicle_id"`
	Name 	   		string 		`json:"name"`
	StartTime       int64 		`json:"start_time"`
	DispatchTime    int64 		`json:"dispatch_time"`
	DispatchType    int32 		`json:"dispatch_type"`
}

func InsertRoute(route *Route) (*int64, error){
	sql, _, _ := sqlBuilder.Insert(ROUTE_TABLE).Rows(
		sqlBuilder.Record {
			"name"    : route.Name,
			"hub_id"  : route.HubId,
			"city_id"  : route.CityId,
			"start_time" : 0,
			"vehicle_id" : route.VehicleId,
			"dispatch_time" : 0,
			"agent_id" : route.AgentId,
			"dispatch_type" : SUBSCRIPTION_DISPATCH,
		},
	).Returning("id").ToSQL();
	var route_id int64;
	err := db.QueryRow(sql).Scan(&route_id);
	if err != nil {
		return nil, err;
	}
	return &route_id, nil
}

func UpdateRoute(rotue *Route) (error){	
	sql, _, _ := sqlBuilder.Update(ROUTE_TABLE).Set(
		sqlBuilder.Record {
			"name"          : rotue.Name,
			"start_time"    : rotue.StartTime,
			"dispatch_time" : rotue.DispatchTime,
		},
	).Where(sqlBuilder.Ex {"id": rotue.Id}).ToSQL();
	rows, err := db.Query(sql);
	defer rows.Close();
	if err != nil {
		fmt.Println(err);
		return err;
	}
	return nil;
}

func AssignRouteAgent(route_id int64, agent_id int64) (error){	
	sql, _, _ := sqlBuilder.Update(ROUTE_TABLE).Set(
		sqlBuilder.Record {		
			"agent_id"  : agent_id,
		},
	).Where(sqlBuilder.Ex {"id": route_id}).ToSQL();
	rows, err := db.Query(sql);
	defer rows.Close();
	if err != nil {
		return err;
	}
	return nil;
}

func AssignRouteVehicle(route_id int64, vehicle_id int64) (error){	
	sql, _, _ := sqlBuilder.Update(ROUTE_TABLE).Set(
		sqlBuilder.Record {		
			"vehicle_id"  : vehicle_id,
		},
	).Where(sqlBuilder.Ex {"id": route_id}).ToSQL();
	rows, err := db.Query(sql);
	defer rows.Close();
	if err != nil {
		return err;
	}
	return nil;
}

func GetRouteById(route_id int64) (*Route, error){	
	query, _, _ := sqlBuilder.Select(GetRouteColumns()...).From(ROUTE_TABLE).Where (
		sqlBuilder.Ex{"id": route_id},
	).ToSQL();
	row, err := db.Query(query);
	if err != nil {
		return nil, err;
	}
	defer row.Close();
	present := row.Next();
	if present == false {
		return nil, errors.New(errText);
	}
	if err := row.Err(); err != nil {
		return nil, err;
	}
	route, err := StructRoute(row);
	if err != nil {
		return nil, err;
	}
	return route, nil;
}

func GetRoutesOfHub (hub_id int64) (*[]*Route, error){
	query, _, _ := sqlBuilder.Select(GetRouteColumns()...).From(ROUTE_TABLE).Where (
		sqlBuilder.Ex{"hub_id": hub_id},
	).ToSQL();
	rows, err := db.Query(query);
	if err != nil {
		return nil, err;
	}
	defer rows.Close();
	var routes []*Route;
	for rows.Next() {
		if err := rows.Err(); err != nil {			
			continue
		}
		route, err := StructRoute(rows);
		if err == nil {			
			routes = append(routes, route);
		}else{
			fmt.Println(err);
			continue
		}
	}	
	return &routes, nil
}

func GetRoutesOfCity (city_id int64) (*[]*Route, error){
	query, _, _ := sqlBuilder.Select(GetRouteColumns()...).From(ROUTE_TABLE).Where (
		sqlBuilder.Ex{"city_id": city_id},
	).ToSQL();
	rows, err := db.Query(query);
	if err != nil {
		return nil, err;
	}
	defer rows.Close();
	var routes []*Route;
	for rows.Next() {
		if err := rows.Err(); err != nil {			
			continue
		}
		route, err := StructRoute(rows);
		if err == nil {			
			routes = append(routes, route);
		}else{
			fmt.Println(err);
			continue
		}
	}	
	return &routes, nil
}

func StructRoute (row *sql.Rows) (*Route, error) {
	var route Route;
	err := row.Scan(	
		&route.Id,	
		&route.HubId,
		&route.CityId,
		&route.AgentId,
		&route.VehicleId,
		&route.Name,
		&route.DispatchTime,	
		&route.StartTime,
		&route.DispatchType,
	);
	if err != nil {
		return nil, err
	}
	return &route, nil;
}