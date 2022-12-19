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
	_ "github.com/lib/pq"
	"go-micro.dev/v4/logger"
	"database/sql"
	"fmt"
)

const CREDITED = true;
const DEBITED = false;

const USER_TABLE string = "users";
const COMPANY_TABLE string = "companies";
const ADRESS_TABLE string = "addresses"
const TXN_TABLE string = "txns";
const COMPANY_BILLING_TABLE string = "company_billing"

const errText string = "No Records!";
const RESPONSE_PARSE_ERROR string = "Unable to parse data!"
const selectAddressSQL string = "SELECT id, user_id, address, landmark, flat, created, status, ST_AsGeoJSON(geolocation)::jsonb FROM addresses WHERE (%v = %v)";
var db *sql.DB;

var psqlInfo string = fmt.Sprintf(
	"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
    "localhost",5432,
    "postgres",
    "Rathod@123",
    "rr",
)
//Vehicle Types Constants

const VEHICLE_BYCYCLE = 0;
const VEHICLE_SCOOTY = 1;
const VEHICLE_BIKE = 2;
const VEHICLE_CAR = 3;
const VEHICLE_XUV = 4;
const VEHICLE_BUS = 5;
const VEHICLE_TRUCK = 6;

func init () {
	logger.Info("Opening Database Connection!");
	mDB, err := sql.Open("postgres", psqlInfo);
	if err != nil {
		logger.Fatal(err);
	}
	err = mDB.Ping();
	db = mDB;
	if err != nil {
		logger.Fatal(err);
	}else{
		logger.Info("Ping Successfull!");
	}
}

func GetCompanyColumns () ([]interface{}) {
	id := "id"
	name := "name"
	email := "email"
	contact := "contact"
	role := "role"
	created := "created"		
	settings := "settings"
	columns := []interface{}{id,name,email,contact,role,created,settings}
	return columns;
}

func GetCompanyBillingColumns () ([]interface{}) {
	id := "id"
	rp_subscription_id := "rp_subscription_id"
	rp_plan_id := "rp_plan_id"
	amount := "amount"
	start_at := "start_at"
	end_at := "end_at"		
	status := "status"
	columns := []interface{}{id,rp_subscription_id,rp_plan_id,amount,start_at,end_at,status}
	return columns;
}