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
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
	"errors"

	sqlBuilder "github.com/doug-martin/goqu/v9"
)

const TABLE string = "cities";
const errText string = "No Records!";

type City struct {
	Id 			 int64 		`json:id`
	CompanyId  	 int64 		`json:company_id`
	Status  	 int32 		`json:status`
	Name 	     string		`json:name`
	Created 	 int64 		`json:created`
}

//Related Subscribe Page Where You Can Show Number of Days
type DeliveryDays struct {
	Days []struct {
		DayNo  int64  `json:"day_no"`
		Active bool   `json:"active"`
	} `json:"days"`
}
func (a DeliveryDays) Value() (driver.Value, error) {
    return json.Marshal(a)
}
func (a *DeliveryDays) Scan(value interface{}) error {
    b, ok := value.([]byte)
    if !ok {
        return errors.New("type assertion to []byte failed")
    }

    return json.Unmarshal(b, &a)
}


func (city *City) toJSON() (string, error) {
	res, err := json.Marshal(city);
	if err != nil {
		return "", err
	}
	return string(res), nil;
}

func InsertCity (city *City) (int64, error){
	timeUnix := time.Now().Unix();
	const sql string = "Insert into cities (company_id, name, status, created) values (%v, '%s', %v, %v) RETURNING id";		
	var lastInsertId int64;
	err := db.QueryRow(fmt.Sprintf(sql, city.CompanyId, city.Name, 0, timeUnix)).Scan(&lastInsertId);
	logger.Info(fmt.Sprintf(sql, city.CompanyId, city.Name, 0, timeUnix));
	if err != nil {
		return 0, err
	}
	return lastInsertId, nil;
}

func UpdateCity(city_id int64, name string, status int32) error {
	query, _, err := sqlBuilder.Update(TABLE).Set(
		sqlBuilder.Record {
			"status":status,
			"name": name,
		},
	).Where(
		sqlBuilder.Ex {"id": city_id},
	).ToSQL();
	if err != nil {
		logger.Fatal(err);
	}
	rows, err := db.Query(query);
	defer rows.Close();
	if err != nil {
		return err;
	}
	return nil;
}

func GetCityById (city_id int64) (*City, error){
	query, _, _ := sqlBuilder.From(TABLE).Where (
		sqlBuilder.Ex{"id": city_id},
	).Limit(1).ToSQL();
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

	city, err := StructCity(rows);
	if err != nil {
		return nil, err
	}
	return city, nil
}

func GetCompanyCities (company_id int64) (*[]City, error){
	query, _, _ := sqlBuilder.From(TABLE).Where (
		sqlBuilder.Ex{"company_id": company_id},
	).ToSQL();
	rows, err := db.Query(query);
	if err != nil {
		return nil, err;
	}
	defer rows.Close();
	var cities []City;
	for rows.Next() {
		if err := rows.Err(); err != nil {			
			continue
		}
		city, err := StructCity(rows);
		if err == nil {			
			cities = append(cities, *city);
		}else{
			logger.Fatal(err);
			continue
		}
	}	
	return &cities, nil
}

//City Settings
func SetSubDeliveryDays(city_id int64, dlvDays *DeliveryDays) (int64, error){	
	id, _, _ := GetSubDeliveryDays(city_id);
	if id == nil {
		query, _, err := sqlBuilder.Insert(SETTING_TABLE).Rows(sqlBuilder.Record {
			"data":dlvDays,"city_id":city_id,"key": CITY_DLV_DAYS,
		}).ToSQL();
		if err != nil {
			return 0, err;
		}
		var lastInsertId int64;
		err = db.QueryRow(query + " RETURNING id").Scan(&lastInsertId);
		if err != nil {
			return 0, err
		}
		return lastInsertId, nil;
	}else{
		query, _, err := sqlBuilder.Update(SETTING_TABLE).Set(sqlBuilder.Record{"data":dlvDays}).Where(sqlBuilder.Ex{"id": *id}).ToSQL();
		if err != nil {
			return 0, err
		}
		rows, err := db.Query(query);
		defer rows.Close();
		if err != nil {
			return 0, err;
		}
		return *id, nil;
	}
}

func GetSubDeliveryDays (city_id int64) (*int64, *DeliveryDays, error){
	query, _, _ := sqlBuilder.From(SETTING_TABLE).Select("id", "data").Where (
		sqlBuilder.Ex{"city_id": city_id},
		sqlBuilder.Ex{"key": CITY_DLV_DAYS},		
	).Limit(1).ToSQL();
	rows, err := db.Query(query);
	if err != nil {
		return nil, nil, err;
	}
	defer rows.Close();
	present := rows.Next();
	if present == false {
		return nil, nil, errors.New(errText);
	}
	if err := rows.Err(); err != nil {
		return nil, nil, err;
	}
	return toDeliveryDays(rows)
}

func toDeliveryDays (row *sql.Rows) (*int64, *DeliveryDays, error) {
	var id int64;
	var dlvDays DeliveryDays;
	err := row.Scan(
		&id,
		&dlvDays,
	);
	if err != nil {
		return nil, nil, err
	}
	return &id, &dlvDays, nil;
}

func StructCity (row *sql.Rows) (*City, error) {
	var city City;
	err := row.Scan(
		&city.Id,
		&city.CompanyId,
		&city.Name, 
		&city.Status, 
		&city.Created,
	);
	if err != nil {
		return nil, err
	}
	return &city, nil;
}