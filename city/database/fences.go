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
--------------------------------
 --->   Logistics System   <---
--------------------------------
*/

package database

import (
	_ "github.com/lib/pq"
	"go-micro.dev/v4/logger"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
	"errors"

	_ "github.com/doug-martin/goqu/v9"
)

const FENCE_TABLE string = "fences";
const polyJSONString string = "'{\"type\":\"Polygon\",\"coordinates\":[%v]}'";
const pointString = "ST_GeomFromText('Point(%v %v)')";
const selectROW = "id, city_id, hub_id, name, status, ST_AsGeoJSON(polygon)::jsonb, created";

type Polygon struct {
	Type 		string		   `json:type`
	Coordinates [][][]float64  `json:coordinates`
}

type Fence struct {
	Id 			  int64 		`json:id`
	CityId  	  int64 		`json:city_id`
	HubId 	      int64 		`json:hub_id`
	Name 	 	  string		`json:name`
	Polygon 	  Polygon		`json:polygon`
	PolygonString []byte		`json:polygonString`
	Status 	      int32 		`json:status`
	Created 	  int64 		`json:created`	
}

func (fence *Fence) toJSON() (string, error) {
	res, err := json.Marshal(fence);
	if err != nil {
		return "", err
	}
	return string(res), nil;
}

func InsertFence (fence *Fence, coords string) (int64, error){
	timeUnix := time.Now().Unix();
	polygon := fmt.Sprintf(polyJSONString, coords);
	const sql string = "Insert into fences (city_id, hub_id, name, status, created, polygon) values (%v, %v, '%v', %v, %v, %v) RETURNING id";
	var lastInsertId int64;
	fmt.Println(fmt.Sprintf(sql, fence.CityId, fence.HubId, fence.Name, 0, timeUnix, polygon));
	err := db.QueryRow(fmt.Sprintf(sql, fence.CityId, fence.HubId, fence.Name, 0, timeUnix, polygon)).Scan(&lastInsertId);
	if err != nil {
		return 0, err
	}
	return lastInsertId, nil;
}

func UpdateFence(fence_id int64, name string, status int32, coords string) error {
	polygon := fmt.Sprintf(polyJSONString, coords);
	const sql string = "update fences set name = '%v', status = %v, polygon = %v where id = %v";
	rows, err := db.Query(fmt.Sprintf(sql, name, status, polygon, fence_id));
	defer rows.Close();
	if err != nil {
		return err;
	}
	return nil;
}

func GetFenceById (fence_id int64) (*Fence, error){
	const query string = "select %v from fences where id = %v";
	rows, err := db.Query(fmt.Sprintf(query, selectROW, fence_id));
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
	fence, err := StructFence(rows);
	if err != nil {
		return nil, err
	}
	return fence, nil
}

func GetFenceByGeoPoint (lat float64, lng float64) (*Fence, error){
	point := fmt.Sprintf(pointString, lng, lat);
	const query string = "SELECT %v FROM fences WHERE ST_Contains(polygon, %v)";
	rows, err := db.Query(fmt.Sprintf(query, selectROW, point));	
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
	fence, err := StructFence(rows);
	if err != nil {
		return nil, err
	}
	return fence, nil
}

func GetFenceByGeoPointAndId (hub_id int64, lat float64, lng float64) (*Fence, error) {
	point := fmt.Sprintf(pointString, lng, lat);
	const query string = "SELECT %v FROM fences WHERE ST_Contains(polygon, %v) And hub_id = %v";
	rows, err := db.Query(fmt.Sprintf(query, selectROW, point, hub_id));
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
	fence, err := StructFence(rows);
	if err != nil {
		return nil, err
	}
	return fence, nil
}

func GetCityFences (city_id int64) (*[]*Fence, error){
	const query string = "select %v from fences where city_id = %v";
	rows, err := db.Query(fmt.Sprintf(query, selectROW, city_id));
	if err != nil {
		return nil, err;
	}
	defer rows.Close();
	var fences []*Fence;
	for rows.Next() {
		if err := rows.Err(); err != nil {			
			continue
		}
		fence, err := StructFence(rows);
		if err == nil {			
			fences = append(fences, fence);
		}else{
			logger.Info(err);
			continue
		}
	}	
	return &fences, nil
}

func GetHubFences (hub_id int64) (*[]*Fence, error){
	const query string = "select %v from fences where hub_id = %v";
	rows, err := db.Query(fmt.Sprintf(query, selectROW, hub_id));
	if err != nil {
		return nil, err;
	}
	defer rows.Close();
	var fences []*Fence;
	for rows.Next() {
		if err := rows.Err(); err != nil {			
			continue
		}
		fence, err := StructFence(rows);
		if err == nil {			
			fences = append(fences, fence);
		}else{
			logger.Info(err);
			continue
		}
	}	
	return &fences, nil
}

func StructFence (row *sql.Rows) (*Fence, error) {
	var fence Fence;
	err := row.Scan(
		&fence.Id,
		&fence.CityId,
		&fence.HubId,
		&fence.Name,
		&fence.Status,
		&fence.PolygonString,
		&fence.Created,
	);
	err = json.Unmarshal(fence.PolygonString, &fence.Polygon);
	if err != nil {
		return nil, err;
	}
	if err != nil {
		return nil, err
	}
	return &fence, nil;
}