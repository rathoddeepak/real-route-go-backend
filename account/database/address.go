/*
--------------------------------
 Author: Deepak Rathod
--------------------------------
 Date - 17 August 2022
 Phone: +917020814070
 Email: rathoddeepak537@gmail.com 
--------------------------------
 --->   Logistics System   <---
--------------------------------
*/
package database;

import (
	"fmt"
	"time"
	"strings"
	"errors"
	"encoding/json"
	_ "github.com/lib/pq"	
	"database/sql"
	sqlBuilder "github.com/doug-martin/goqu/v9"
)

type Address struct {
	Id 	   			int64 			 `json:"id"`
	UserId 			int64 			 `json:"user_id"`
	Address  		string  		 `json:"address"`

	NullLandmark    sql.NullString   	 `json:"landmark"`
	NullFlat  	    sql.NullString   	 `json:"flat"`

	Landmark  		string   	 	 `json:"landmark"`
	Flat  			string   	 	 `json:"flat"`

	Type  			int32  	    	 `json:"type"`
	JSONLocation 	[]byte  		 `json:"json_location"`
	Geolocation 	string  		 `json:"geolocation"`
	Location        GeoPoint		 `json:"location"`
	Created   		int64 			 `json:"created"`
	Status   		int32 			 `json:"status"`
}

type GeoPoint struct {
	Type 		string 		`json:type`
	Coordinates []float64   `json:type`
}

func InsertAddress(address *Address) (*int64, error){
	record := sqlBuilder.Record {
		"user_id" : address.UserId,
		"address" : address.Address,
		"type" : address.Type,
		"created":time.Now().Unix(),
		"geolocation" : address.Geolocation,
		"status":1,
	}
	if address.Landmark != "" {
		record["landmark"] = address.Landmark;
	}
	if address.Flat != "" {
		record["flat"] = address.Flat;
	}
	sql, _, _ := sqlBuilder.Insert(ADRESS_TABLE).Rows(record).ToSQL();	
	sql = strings.Replace(sql, "'|", "" , 1);
	sql = strings.Replace(sql, "|'", "" , 1);
	sql = strings.ReplaceAll(sql, "''", "'" );
	var addressId int64;
	err := db.QueryRow(sql + " RETURNING id").Scan(&addressId);
	if err != nil {
		return nil, err;
	}
	return &addressId, nil
}

func UpdateAddress(address *Address) (error){
	record := sqlBuilder.Record {			
		"address" : address.Address,
		"type" : address.Type,
		"status":address.Status,
		"geolocation" : address.Geolocation,
	}
	if address.Landmark != "" {
		record["landmark"] = address.Landmark;
	}
	if address.Flat != "" {
		record["flat"] = address.Flat;
	}
	sql, _, _ := sqlBuilder.Update(ADRESS_TABLE).Set(record).Where(sqlBuilder.Ex {"id": address.Id}).ToSQL();
	sql = strings.Replace(sql, "'|", "" , 1);
	sql = strings.Replace(sql, "|'", "" , 1);
	sql = strings.ReplaceAll(sql, "''", "'" );
	rows, err := db.Query(sql);
	defer rows.Close();

	if err != nil {
		return err;
	}

	return nil;
}

func UpdateAddressLocaiton(address *Address) (error){	
	sql, _, _ := sqlBuilder.Update(ADRESS_TABLE).Set(
		sqlBuilder.Record {			
			"geolocation" : address.Geolocation,
		},
	).Where(sqlBuilder.Ex {"id": address.Id}).ToSQL();
	sql = strings.Replace(sql, "'|", "" , 1);
	sql = strings.Replace(sql, "|'", "" , 1);
	sql = strings.ReplaceAll(sql, "''", "'" );
	rows, err := db.Query(sql);
	defer rows.Close();
	if err != nil {
		return err;
	}
	return nil;
}

func UpdateAddressStatus(address *Address) (error){	
	sql, _, _ := sqlBuilder.Update(ADRESS_TABLE).Set(
		sqlBuilder.Record {			
			"status" : address.Status,
		},
	).Where(sqlBuilder.Ex {"id": address.Id}).ToSQL();
	rows, err := db.Query(sql);
	defer rows.Close();
	if err != nil {
		return err;
	}
	return nil;
}

func GetAddressesOfUser(user_id int64) (*[]*Address, error) {
	query := fmt.Sprintf(selectAddressSQL, "user_id", user_id);
	var addresses []*Address;
	rows, err := db.Query(query);
	if err != nil {
		return nil, err;
	}
	defer rows.Close();
	for rows.Next() {		
		err := rows.Err();
		if err != nil {
			continue
		}
		address, err := StructAddress(rows);
		addresses = append(addresses, address);
	}
	return &addresses, nil;
}

func GetAddressById(address_id int64) (*Address, error) {
	query := fmt.Sprintf(selectAddressSQL + " LIMIT 1", "id", address_id);
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
	address, err := StructAddress(row);
	if err != nil {
		return nil, err;
	}
	return address, nil;
}

func StructAddress (row *sql.Rows) (*Address, error) {
	var address Address;
	err := row.Scan(
		&address.Id,
		&address.UserId,
		&address.Address,
		&address.NullLandmark,
		&address.NullFlat,
		&address.Created,		
		&address.Status,
		&address.JSONLocation,
	);
	address.Landmark = address.NullLandmark.String
	address.Flat = address.NullFlat.String
	err = json.Unmarshal(address.JSONLocation, &address.Location);	
	if err != nil {
		return nil, err
	}
	return &address, nil;
}