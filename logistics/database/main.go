/*
--------------------------------
 Author: Deepak Rathod
--------------------------------
 Date - 18 August 2022
 Phone: +917020814070
 Email: rathoddeepak537@gmail.com 
--------------------------------
 --->   Logistics Service  <---
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

	"go-micro.dev/v4/broker"
)

const SUBSCRIPTION_TABLE string = "subscriptions";
const ROUTE_NODE_TABLE string = "route_nodes";
const ROUTE_TABLE string = "routes";
const SLOT_TABLE string = "slots";
const AGENT_TABLE string = "agents";
const TASK_TABLE string = "tasks";
const TASK_POINT_TABLE string = "task_points";
const VEHICLE_TABLE string = "vehicles";

const TASK_CREATED = 0;
const TASK_ASSIGNED = 1;
const TASK_ACCEPTED = 2;
const TASK_STARTED = 3;
const TASK_COMPLETED = 4;
const TASK_CUSTOMER_CANCEL = 5;
const TASK_AGENT_CANCEL = 6;
const TASK_SYSTEM_CANCEL_PAID = 7;
const TASK_SYSTEM_CANCEL_NOT_PAID = 8;
const TASK_TRANSFERRED = 9;

const errText string = "No Records!";
const SUBSCRIPTION_DISPATCH = 0;
var db *sql.DB;

const HUB_NODE = 0;
const SUBSCRIPTION_NODE = 1;
const LOCATION_NODE = 2;
const ADDRESS_NODE = 3;

const NODE_ACTION_PICKIUP = 0;
const NODE_ACTION_DELIVERY = 1;
const NODE_ACTION_FIELD_TASK = 2;

//Agent Status
const AGENT_OFFLINE= 0;
const AGENT_ONLINE = 1;
const AGENT_BUSY = 2;
const AGENT_ONLINE_BUSY = 3; //For Filtering Purpose Only

const TASK_STATUS_AGENT_PENDING = 19;
const TASK_STATUS_COMPLETED = 20;
const TASK_STATUS_PENDING = 21;

const REALTIME_TOPIC = "go.micro.topic.logi_live"

var psqlInfo string = fmt.Sprintf(
	"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
    "localhost",5432,
    "postgres",
    "Rathod@123",
    "rr",
)

//if Auto Cancel is on for task then task will be cancelled
//when home page of real route app is loaded
const TASK_AUTO_CANCEL_ACTIVE = 1;
const TASK_AUTO_CANCEL_INACTIVE = 0;

//Real time events
const (
    EVENT_TASK_CREATED                     = "task_created";
    EVENT_TASK_ASSIGNED                    = "task_assigned";
    EVENT_TASK_ACCEPTED                    = "task_accepted";
    EVENT_TASK_STARTED                     = "task_started";
    EVENT_TASK_COMPLETED                   = "task_completed";
    EVENT_TASK_FULL_COMPLETED              = "task_full_completed";
    EVENT_TASK_CUSTOMER_CANCEL             = "task_customer_cancel"
    EVENT_TASK_AGENT_CANCEL                = "task_agent_cancel"
    EVENT_TASK_SYSTEM_CANCEL_PAID          = "task_cancel_paid"
    EVENT_TASK_SYSTEM_CANCEL_NOT_PAID      = "task_cancel_not_paid"
    EVENT_TASK_TRANSFERRED			       = "task_transferred"
    EVENT_TASK_AGENT_CHANGE			       = "task_agent_change"

    EVENT_AGENT_LOCATION_UPDATE            = "agent_location_update"
    EVENT_AGENT_STATUS_UPDATE              = "agent_status_update"
)

//Realtime live updates
type TaskUpdate_Realtime struct {
    Event          string       `json:"event,omitempty"`
    Task_id        int64        `json:"id,omitempty"`
    Status         int32        `json:"status,omitempty"`
    Point_id       int64        `json:"point_id,omitempty"`
    City_id        int64        `json:"city_id"`
    Agent_id       int64        `json:"agent_id,omitempty"`
    Agent_name     string       `json:"agent_name,omitempty"`
    Agent_phone    string       `json:"agent_phone,omitempty"`
    Agent_status   string       `json:"agent_status,omitempty"`
    Agent_lat      float64      `json:"agent_lat,omitempty"`
    Agent_lng      float64      `json:"agent_lng,omitempty"`
}


//RealRoute Code: Start
const LOCATION_TABLE string = "locations";
//RealRoute Code: End

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

func PublishMessage(body []byte) {
    err := broker.Publish(REALTIME_TOPIC, &broker.Message{Body: body}); 
    if err != nil {
        logger.Info(err);
    }
}

func GetColumns () ([]interface{}) {
	idCol := "id"
	userIdCol := "user_id"
	product_idCol := "product_id"
	start_fromCol := "start_from"
	dlv_insCol := "dlv_ins"
	dlv_daysCol := "dlv_days"
	sunCol := "sun"
	monCol := "mon"
	tueCol := "tue"
	wedCol := "wed"
	thuCol := "thu"
	friCol := "fri"
	satCol := "sat"
	patternCol := "pattern"
	createdCol := "created"
	statusCol := "status"
	address_idCol := "address_id"
	slot_idCol := "slot_id"
	city_idCol := "city_id"
	hub_idCol := "hub_id"
	expEndCol := "exp_end"
	columns := []interface{}{idCol,userIdCol,product_idCol,start_fromCol,dlv_insCol,dlv_daysCol,sunCol,monCol,tueCol,wedCol,thuCol,friCol,satCol,patternCol,createdCol,statusCol,address_idCol,slot_idCol,city_idCol,hub_idCol,expEndCol}
	return columns;
}

func GetRouteNodeColumns () ([]interface{}) {
	idCol := "id"
	RouteId := "route_id"
	ActionType := "action_type"
	NodeType := "node_type"
	NodeId := "node_id"
	columns := []interface{}{idCol,RouteId,ActionType,NodeType,NodeId}
	return columns;
}		

func GetRouteColumns () ([]interface{}) {
	idCol := "id"
	hubIdCol := "hub_id"
	cityIdCol := "city_id"
	agentIdCol := "agent_id"
	vehicleIdCol := "vehicle_id"
	name := "name"
	dispatchTime := "dispatch_time"
	startTime := "start_time"
	dispatchType := "dispatch_type"		
	columns := []interface{}{idCol,hubIdCol,cityIdCol,agentIdCol,vehicleIdCol,name,dispatchTime,startTime,dispatchType}
	return columns;
}

func GetAgentColumns () ([]interface{}) {
	idCol := "id"
	cityIdCol := "city_id"
	name := "name"
	phone := "phone"
	passcode := "passcode"
	avatar := "avatar"
	locationCol := "ST_AsGeoJSON(location)::jsonb"
	statusCol := "status"
		
	columns := []interface{}{idCol,cityIdCol,name,phone,passcode,avatar,locationCol,statusCol}
	return columns;
}

func GetTaskColumns () ([]interface{}) {
	idCol := "id";
	city_idCol := "city_id";
	route_idCol := "route_id";
	hub_idCol := "hub_id";
	user_idCol := "user_id";
	agent_idCol := "agent_id";
	vehicle_idCol := "vehicle_id";
	start_afterCol := "start_after";
	end_afterCol := "end_after";
	visible_afterCol := "visible_after";
	statusCol := "status";
	createdCol := "created";
	autoCancel := "auto_cancel";
	return []interface{}{idCol,city_idCol,route_idCol,hub_idCol,user_idCol,agent_idCol,vehicle_idCol,start_afterCol,end_afterCol,visible_afterCol,statusCol,createdCol,autoCancel}
}

func GetTaskPointColumns () ([]interface{}) {
	idCol := "id";
	task_idCol := "task_id";
	hub_idCol := "hub_id";
	user_idCol := "user_id";
	agent_idCol := "agent_id";
	subscription_idCol := "subscription_id";
	dependent_idCol := "dependent_id";
	task_type := "task_type";	
	name := "name"
	contact := "contact"
	address := "address"
	lat := "lat"
	lng := "lng"
	statusCol := "status";
	createdCol := "created";

	return []interface{}{idCol,task_idCol,hub_idCol,user_idCol,agent_idCol,subscription_idCol,dependent_idCol,task_type,name,contact,address,lat,lng,statusCol,createdCol}
}


//RealRoute Code: Start
func GetLocationColumns () ([]interface{}) {
	idCol := "id"
	cityIdCol := "city_id"
	name := "name"
	address := "address"
	contact := "contact"
	locationCol := "ST_AsGeoJSON(geolocation)::jsonb"
		
	columns := []interface{}{idCol,cityIdCol,name,address,contact,locationCol}
	return columns;
}
//RealRoute Code: End

func GetVehicleColumns () ([]interface{}) {
	id := "id"
	cityId := "city_id"
	companyId := "company_id"
	name := "name"
	photo := "photo"
	vehicleType := "type"
	locationCol := "ST_AsGeoJSON(location)::jsonb"
	status := "status"
	created := "created"		
	columns := []interface{}{id,cityId,companyId,name,photo,vehicleType,locationCol,status,created}
	return columns;
}
