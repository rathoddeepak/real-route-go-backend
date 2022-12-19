/*
--------------------------------
 Author: Deepak Rathod
--------------------------------
 Date - 05 September 2022
 Phone: +917020814070
 Email: rathoddeepak537@gmail.com 
--------------------------------
 --->   Logistics System   <---
--------------------------------
        RealRoute Code:
--------------------------------
*/
package database;

import (
    "fmt"

    "errors"
    "strings"
    "encoding/json"
    _ "github.com/lib/pq"   
    "database/sql"  

    sqlBuilder "github.com/doug-martin/goqu/v9"
)

type Location struct {
    Id              int64       `json:"id"`
    CityId          int64       `json:"city_id"`
    Name            string      `json:"name"`
    Address         string      `json:"address"`
    Contact         string      `json:"contact"`
    JSONLocation    []byte      `json:json_location`
    Location        GeoPoint    `json:location`
}

const location_geostring_from string = "\"ST_AsGeoJSON(geolocation)::jsonb\"";
const location_geostring_to string = "ST_AsGeoJSON(geolocation)::jsonb";


func InsertLocation(point string, location *Location) (*int64, error){
    sql, _, _ := sqlBuilder.Insert(LOCATION_TABLE).Rows(
        sqlBuilder.Record {
            "name"    : location.Name,
            "city_id" : location.CityId,
            "address" : location.Address,
            "contact" : location.Contact,
            "geolocation": point,
        },
    ).Returning("id").ToSQL();
    var location_id int64;
    sql = strings.Replace(sql, "'|", "" , 1);
    sql = strings.Replace(sql, "|'", "" , 1);
    sql = strings.ReplaceAll(sql, "''", "'" );
    err := db.QueryRow(sql).Scan(&location_id);
    if err != nil {
        return nil, err;
    }
    return &location_id, nil
}

func UpdateLocation(geolocation string, location *Location) (error){   
    sql, _, _ := sqlBuilder.Update(LOCATION_TABLE).Set(
        sqlBuilder.Record {
            "name"        : location.Name,
            "address"     : location.Address,
            "contact"     : location.Contact,
            "geolocation" : geolocation,
        },
    ).Where(sqlBuilder.Ex {"id": location.Id}).ToSQL();
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

func GetLocationById(location_id int64) (*Location, error){  
    query, _, _ := sqlBuilder.Select(GetLocationColumns()...).From(LOCATION_TABLE).Where (
        sqlBuilder.Ex{"id": location_id},
    ).ToSQL();
    query = strings.Replace(query, location_geostring_from, location_geostring_to, 1);
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
    location, err := StructLocation(row);
    if err != nil {
        return nil, err;
    }
    return location, nil;
}

func GetLocationsOfCity (city_id int64) (*[]*Location, error){
    query, _, _ := sqlBuilder.Select(GetLocationColumns()...).From(LOCATION_TABLE).Where (
        sqlBuilder.Ex{"city_id": city_id},
    ).ToSQL();
    query = strings.Replace(query, location_geostring_from, location_geostring_to, 1);
    rows, err := db.Query(query);
    if err != nil {
        return nil, err;
    }
    defer rows.Close();
    var locations []*Location;
    for rows.Next() {
        if err := rows.Err(); err != nil {
            fmt.Println(err);            
            continue
        }
        location, err := StructLocation(rows);
        if err == nil {         
            locations = append(locations, location);
        }else{          
            fmt.Println(err);      
            continue
        }
    }   
    return &locations, nil
}

func StructLocation (row *sql.Rows) (*Location, error) {
    var location Location;
    err := row.Scan(    
        &location.Id,  
        &location.CityId,  
        &location.Name,
        &location.Address,
        &location.Contact,
        &location.JSONLocation,
    );
    if err != nil {
        return nil, err
    }
    err = json.Unmarshal(location.JSONLocation, &location.Location);  
    if err != nil {
        return nil, err
    }
    return &location, nil;
}