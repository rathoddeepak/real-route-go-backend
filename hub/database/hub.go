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

type Hub struct {
	Id 	   			int64 		`json:id`
	CityId 			int64 		`json:city_id`
	UserId 			int64 		`json:user_id`
	Name  			string  	`json:name`
	Address  		string  	`json:address`
	About  			string  	`json:about`
	Phone  			string  	`json:phone`
	JSONLocation 	[]byte  	`json:json_location`
	Geolocation 	string  	`json:geolocation`
	Location        GeoPoint	`json:location`
	Status 			int32   	`json:status`
	Created   		int64 		`json:created`
}

type GeoPoint struct {
	Type 		string 		`json:type`
	Coordinates []float64   `json:type`
}

func InsertHub(hub *Hub) (*int64, error){
	sql, _, _ := sqlBuilder.Insert(HUB_TABLE).Rows(
		sqlBuilder.Record {
			"city_id" : hub.CityId,
			"user_id" : hub.UserId,
			"name" : hub.Name,
			"about" : hub.About,
			"phone" : hub.Phone,
			"address" : hub.Address,
			"geolocation" : hub.Geolocation,
			"status" : 0,
			"created":time.Now().Unix(),
		},
	).ToSQL();	
	sql = strings.Replace(sql, "'|", "" , 1);
	sql = strings.Replace(sql, "|'", "" , 1);
	sql = strings.ReplaceAll(sql, "''", "'" );
	var hubId int64;
	err := db.QueryRow(sql + " RETURNING id").Scan(&hubId);;	
	if err != nil {
		return nil, err;
	}
	return &hubId, nil
}

func UpdateHub(hub *Hub) (error){	
	sql, _, _ := sqlBuilder.Update(HUB_TABLE).Set(
		sqlBuilder.Record {			
			"name" : hub.Name,
			"about" : hub.About,
			"phone" : hub.Phone,
			"address" : hub.Address,
			"geolocation" : hub.Geolocation,
			"status": hub.Status,
		},
	).Where(sqlBuilder.Ex {"id": hub.Id}).ToSQL();
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

func UpdateHubLocaiton(hub *Hub) (error){	
	sql, _, _ := sqlBuilder.Update(HUB_TABLE).Set(
		sqlBuilder.Record {			
			"geolocation" : hub.Geolocation,
		},
	).Where(sqlBuilder.Ex {"id": hub.Id}).ToSQL();
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

func UpdateHubStatus(hub *Hub) (error){	
	sql, _, _ := sqlBuilder.Update(HUB_TABLE).Set(
		sqlBuilder.Record {			
			"status" : hub.Status,
		},
	).Where(sqlBuilder.Ex {"id": hub.Id}).ToSQL();
	rows, err := db.Query(sql);
	defer rows.Close();
	if err != nil {
		return err;
	}
	return nil;
}

func GetHubsOfCity(city_id int64) (*[]*Hub, error) {
	query := fmt.Sprintf(selectHubSQL, "city_id", city_id);
	var hubs []*Hub;
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
		hub, err := StructHub(rows);
		hubs = append(hubs, hub);
	}
	return &hubs, nil;
}

func GetHubsOfUser(user_id int64) (*[]*Hub, error) {
	query := fmt.Sprintf(selectHubSQL, "user_id", user_id);
	var hubs []*Hub;
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
		hub, err := StructHub(rows);
		hubs = append(hubs, hub);
	}
	return &hubs, nil;
}

func GetHubById(hub_id int64) (*Hub, error) {
	query := fmt.Sprintf(selectHubSQL + " LIMIT 1", "id", hub_id);
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
	hub, err := StructHub(row);
	if err != nil {
		return nil, err;
	}
	return hub, nil;
}

func StructHub (row *sql.Rows) (*Hub, error) {
	var hub Hub;	
	err := row.Scan(
		&hub.Id,
		&hub.UserId,
		&hub.CityId,
		&hub.Name,
		&hub.About,
		&hub.Address,		
		&hub.Phone,
		&hub.Status,
		&hub.JSONLocation,
		&hub.Created,
	);
	err = json.Unmarshal(hub.JSONLocation, &hub.Location);	
	if err != nil {
		return nil, err
	}
	return &hub, nil;
}