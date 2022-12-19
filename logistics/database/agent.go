/*
--------------------------------
 Author: Deepak Rathod
--------------------------------
 Date - 22 August 2022
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
	"strings"
	"encoding/json"
	_ "github.com/lib/pq"	
	"database/sql"	

	pb "logisticsService/proto"
	sqlBuilder "github.com/doug-martin/goqu/v9"
)

type Agent struct {
	Id 	   			int64 		`json:"id"`
	CityId 			int64 		`json:"city_id"`
	Name 	   		string 		`json:"name"`
	Phone 	   		string 		`json:"phone"`
	Passcode 	   	string 		`json:"passcode"`	
	Avatar 	   		string 		`json:"avatar"`	
	JSONLocation 	[]byte  	`json:json_location`
	Location        GeoPoint	`json:location`
	Status   		int32 		`json:"status"`
}

const PointString string = "|ST_GeometryFromText('POINT(%v %v)')|";
const geostring_from string = "\"ST_AsGeoJSON(location)::jsonb\"";
const geostring_to string = "ST_AsGeoJSON(location)::jsonb";

func InsertAgent(agent *Agent) (*int64, error){
	point := fmt.Sprintf(PointString, 0, 0);
	sql, _, _ := sqlBuilder.Insert(AGENT_TABLE).Rows(
		sqlBuilder.Record {
			"name"    : agent.Name,
			"city_id" : agent.CityId,
			"phone"   : agent.Phone,
			"passcode": agent.Passcode,
			"location": point,
			"status"  : AGENT_OFFLINE,
		},
	).Returning("id").ToSQL();
	var agent_id int64;
	sql = strings.Replace(sql, "'|", "" , 1);
	sql = strings.Replace(sql, "|'", "" , 1);
	sql = strings.ReplaceAll(sql, "''", "'" );
	err := db.QueryRow(sql).Scan(&agent_id);
	if err != nil {
		return nil, err;
	}
	return &agent_id, nil
}

func UpdateAgent(agent *Agent) (error){	
	sql, _, _ := sqlBuilder.Update(AGENT_TABLE).Set(
		sqlBuilder.Record {
			"name"    : agent.Name,
			"passcode": agent.Passcode,
		},
	).Where(sqlBuilder.Ex {"id": agent.Id}).ToSQL();
	rows, err := db.Query(sql);
	defer rows.Close();
	if err != nil {
		fmt.Println(err);
		return err;
	}
	return nil;
}

func UpdateAgentStatus(agent *Agent) (error){	
	sql, _, _ := sqlBuilder.Update(AGENT_TABLE).Set(
		sqlBuilder.Record {		
			"status"  : agent.Status,
		},
	).Where(sqlBuilder.Ex {"id": agent.Id}).ToSQL();
	rows, err := db.Query(sql);	
	if err != nil {
		return err;
	}else{
		defer rows.Close();
	}
	return nil;
}

func UpdateAgentAvatar(agent *Agent) (error){	
	sql, _, _ := sqlBuilder.Update(AGENT_TABLE).Set(
		sqlBuilder.Record {		
			"avatar"  : agent.Avatar,
		},
	).Where(sqlBuilder.Ex {"id": agent.Id}).ToSQL();
	rows, err := db.Query(sql);	
	if err != nil {
		return err;
	}else{
		defer rows.Close();
	}
	return nil;
}

func UpdateAgentLocation(geolocation string, agent_id int64) (error){	
	sql, _, _ := sqlBuilder.Update(AGENT_TABLE).Set(
		sqlBuilder.Record {			
			"location" : geolocation,
		},
	).Where(sqlBuilder.Ex {"id": agent_id}).ToSQL();
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

func GetAgentById(agent_id int64) (*Agent, error){	
	query, _, _ := sqlBuilder.Select(GetAgentColumns()...).From(AGENT_TABLE).Where (
		sqlBuilder.Ex{"id": agent_id},
	).ToSQL();
	query = strings.Replace(query, geostring_from, geostring_to, 1);
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
	agent, err := StructAgent(row);
	if err != nil {
		return nil, err;
	}
	return agent, nil;
}

func VerifyAgent(phone string, passcode string) (*Agent, error){	
	query, _, _ := sqlBuilder.Select(GetAgentColumns()...).From(AGENT_TABLE).Where (
		sqlBuilder.Ex{"phone": phone},
		sqlBuilder.Ex{"passcode": passcode},
	).ToSQL();
	query = strings.Replace(query, geostring_from, geostring_to, 1);
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
	agent, err := StructAgent(row);
	if err != nil {
		return nil, err;
	}
	return agent, nil;
}

func GetAgentByPhone(city_id int64, phone string) (*Agent, error){	
	query, _, _ := sqlBuilder.Select(GetAgentColumns()...).From(AGENT_TABLE).Where (
		sqlBuilder.Ex{"phone": phone},
		sqlBuilder.Ex{"city_id": city_id},
	).ToSQL();
	query = strings.Replace(query, geostring_from, geostring_to, 1);
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
	agent, err := StructAgent(row);
	if err != nil {
		return nil, err;
	}
	return agent, nil;
}

func GetAgentsOfCity (city_id int64) (*[]*Agent, error){
	query, _, _ := sqlBuilder.Select(GetAgentColumns()...).From(AGENT_TABLE).Where (
		sqlBuilder.Ex{"city_id": city_id},
	).ToSQL();
	query = strings.Replace(query, geostring_from, geostring_to, 1);
	rows, err := db.Query(query);
	if err != nil {
		return nil, err;
	}
	defer rows.Close();
	var agents []*Agent;
	for rows.Next() {
		if err := rows.Err(); err != nil {			
			continue
		}
		agent, err := StructAgent(rows);
		if err == nil {			
			agents = append(agents, agent);
		}else{			
			continue
		}
	}	
	return &agents, nil
}

func FilterAgents(filterMap *pb.FilterAgentRequest, countOnly bool) (*[]*Agent, *int64, error) {
	var agents []*Agent;
	
	dbSQL := sqlBuilder.From(AGENT_TABLE).Select(sqlBuilder.COUNT("*"));

	if countOnly == false {		
		dbSQL = sqlBuilder.From(AGENT_TABLE).Select(GetAgentColumns()...);
	}

	if filterMap.CityId == 0 {
		return nil, nil, errors.New("City Id is required!");
	}else{
		dbSQL = dbSQL.Where(sqlBuilder.Ex {"city_id": filterMap.CityId})
	}

	if len(filterMap.Name) > 0 {
		dbSQL = dbSQL.Where(sqlBuilder.Ex {"name": "%"+filterMap.Name+"%"})
	}

	if len(filterMap.Phone) > 0 {
		dbSQL = dbSQL.Where(sqlBuilder.Ex {"phone": "%"+filterMap.Phone+"%"})
	}

	if filterMap.AgentId != 0 {
		dbSQL = dbSQL.Where(sqlBuilder.Ex {"id": filterMap.AgentId})
	}

	if filterMap.Status == AGENT_ONLINE_BUSY {
		//TODO: Not equalt to 
		dbSQL = dbSQL.Where(sqlBuilder.C("status").NotIn(AGENT_OFFLINE));
	}else if filterMap.Status != 9 {//TODO: 9 is temporary change later
		dbSQL = dbSQL.Where(sqlBuilder.Ex {"status": filterMap.Status})
	}

	if filterMap.Limit == 0 {
		dbSQL = dbSQL.Limit(10);
	}else{
		dbSQL = dbSQL.Limit(uint(filterMap.Limit));
	}

	if filterMap.Offset == 0 {
		dbSQL = dbSQL.Offset(0);
	}else{
		dbSQL = dbSQL.Offset(uint(filterMap.Offset));
	}

	if countOnly {
		query, _, _ := dbSQL.ToSQL();
		if len(filterMap.Name) > 0 {
			query = strings.Replace(query, "name\" =", "name\" ILIKE" , 1);
		}
		if len(filterMap.Phone) > 0 {
			query = strings.Replace(query, "phone\" =", "name\" ILIKE" , 1);
		}
		var count int64;
		err := db.QueryRow(query).Scan(&count);
		if err != nil {
			return nil, nil, err;
		}
		return nil, &count, nil;
	}
	
	query, _, _ := dbSQL.Order(sqlBuilder.I("id").Desc()).ToSQL();
	if len(filterMap.Name) > 0 {
		query = strings.Replace(query, "name\" =", "name\" ILIKE" , 1);
	}
	if len(filterMap.Phone) > 0 {
		query = strings.Replace(query, "phone\" =", "phone\" ILIKE" , 1);
	}
	query = strings.Replace(query, geostring_from, geostring_to, 1);
	rows, err := db.Query(query);
	if err != nil {
		return nil, nil, err;
	}
	
	defer rows.Close();
	for rows.Next() {		
		err := rows.Err();
		if err != nil {
			continue
		}
		agent, err := StructAgent(rows);
		agents = append(agents, agent);
	}
	return &agents, nil, nil;
}

func CountAgentBusyTask(agent_id int64) (int) {
	var count int;
	dbSQL := sqlBuilder.From(TASK_TABLE).Select(sqlBuilder.COUNT("*")).Where(
		sqlBuilder.Ex {"agent_id": agent_id},
		sqlBuilder.And(
	        sqlBuilder.C("status").Gte(TASK_ACCEPTED),
	        sqlBuilder.C("status").Lt(TASK_COMPLETED),
	    ),
	);
	query, _, _ := dbSQL.ToSQL();	
	err := db.QueryRow(query).Scan(&count);
	if err != nil {
		return 0;
	}
	return count;
}

func CountAgentAssignedTask(agent_id int64) (int) {
	var count int;
	dbSQL := sqlBuilder.From(TASK_TABLE).Select(sqlBuilder.COUNT("*")).Where(
		sqlBuilder.Ex {"agent_id": agent_id},
		sqlBuilder.Ex {"status": TASK_ASSIGNED},
	);
	query, _, _ := dbSQL.ToSQL();	
	err := db.QueryRow(query).Scan(&count);
	if err != nil {
		return 0;
	}
	return count;
}

func CountAgentCompletedTask(agent_id int64) (int) {
	var count int;
	dbSQL := sqlBuilder.From(TASK_TABLE).Select(sqlBuilder.COUNT("*")).Where(
		sqlBuilder.Ex {"agent_id": agent_id},
		sqlBuilder.C("status").Gte(TASK_COMPLETED),
	);
	query, _, _ := dbSQL.ToSQL();	
	err := db.QueryRow(query).Scan(&count);
	if err != nil {
		return 0;
	}
	return count;
}

func ValidateAgentStatus (agent_id int64) {
	agent, err := GetAgentById(agent_id);
	if err != nil {
		return;
	}
	updateStatus := -1;
	busyTaskCount := CountAgentBusyTask(agent_id);
	if busyTaskCount > 0 {
		UpdateAgentStatus(&Agent {
			Id: agent_id,
			Status: AGENT_BUSY,
		});
		updateStatus = AGENT_BUSY
	}else if agent.Status == AGENT_BUSY {
		UpdateAgentStatus(&Agent {
			Id: agent_id,
			Status: AGENT_ONLINE,
		});	
		updateStatus = AGENT_ONLINE
	}
	if updateStatus == -1 {
		return		
	}
	body, err := json.Marshal(&TaskUpdate_Realtime {
        Event 	 	: 	EVENT_AGENT_STATUS_UPDATE,      
        Agent_id    :   agent_id,
        City_id     :   agent.CityId,
        Status  	:   int32(updateStatus),
    });
    if err == nil {
        PublishMessage(body);
    }
}

func StructAgent (row *sql.Rows) (*Agent, error) {
	var agent Agent;
	err := row.Scan(	
		&agent.Id,	
		&agent.CityId,	
		&agent.Name,
		&agent.Phone,
		&agent.Passcode,
		&agent.Avatar,
		&agent.JSONLocation,
		&agent.Status,		
	);
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(agent.JSONLocation, &agent.Location);	
	if err != nil {
		return nil, err
	}
	return &agent, nil;
}