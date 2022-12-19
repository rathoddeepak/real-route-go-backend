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
	"fmt"
)

var psqlInfo string = fmt.Sprintf(
	"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",

    "localhost",5432,
    "postgres",
    "Rathod@123",
    "rr",
);
var db *sql.DB;

//City Setting Constants
const SETTING_TABLE = "settings";
const CITY_DLV_DAYS int = 0;

func init() {
	logger.Info("Opening Database Connection!");
	mDB, err := sql.Open("postgres", psqlInfo);
	if err != nil {
		logger.Fatal(err);
	}
	err = mDB.Ping();
	if err != nil {
		logger.Fatal(err);
	}else{
		db = mDB;
		logger.Info("Ping Successfull!");		
	}
}