/*
--------------------------------
 Author: Deepak Rathod
--------------------------------
 Date - 04 August 2022
 Phone: +917020814070
 Email: rathoddeepak537@gmail.com 
--------------------------------
 --->   Hub Microservice   <---
--------------------------------
 --->      Products        <---
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

const errText string = "No Records";

const HUB_TABLE string = "hubs";
const CATEGORY_TABLE string = "categories";
const SUBCAT_TABLE string = "sub_categories";
const PRODUCT_TABLE string = "products";
const PRODUCT_IMAGE_FOLDER = "uploads/images/%s/%s.%s";
const selectCatSQL = "SELECT id,hub_id, name, status,image,blurhash FROM categories WHERE (%v = %v)";
const selectSubCatSQL = "SELECT id,main_cat_id,name,status,image,blurhash FROM sub_categories WHERE (%v = %v)";
const selectHubSQL = "SELECT id, user_id, city_id, name, about, address, phone, status, ST_AsGeoJSON(geolocation)::jsonb, created FROM hubs WHERE (%v = %v)";
const selectProductSQL = "select id,name,category_id,sub_category_id,hub_id,price,higher_price,qty,max_limit,status,image,blurhash,city_id from products WHERE (%v = %v)"

var psqlInfo string = fmt.Sprintf(
	"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",

    "localhost",5432,
    "postgres",
    "Rathod@123",
    "rr",
);

var db *sql.DB;

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