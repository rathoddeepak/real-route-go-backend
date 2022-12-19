/*
--------------------------------
 Author: Deepak Rathod
--------------------------------
 Date - 27 August 2022
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
    "errors"
	_ "github.com/lib/pq"
    "database/sql"
    pb "logisticsService/proto"
	sqlBuilder "github.com/doug-martin/goqu/v9"
)

type Task struct {
    Id            int64  `json:"id"`
    CityId        int64  `json:"city_id"`
    RouteId       int64  `json:"route_id"`
    HubId         int64  `json:"hub_id"`
    VehicleId     int64  `json:"hub_id"`
    UserId        int64  `json:"user_id"`
    AgentId       int64  `json:"agent_id"`
    StartAfter    int64  `json:"start_after"`
    EndAfter      int64  `json:"end_after"`
    VisibleAfter  int64  `json:"visible_after"`
    Status        int32  `json:"status"`
    Created       int64  `json:"created"`
    AutoCancel    int32  `json:"auto_cancel"`
}

type TaskPoint struct {
    Id             int64    `json:"id"`
    TaskId         int64    `json:"task_id"`
    HubId          int64    `json:"hub_id"`
    UserId         int64    `json:"user_id"`
    AgentId        int64    `json:"agent_id"`
    SubscriptionId int64    `json:"subscription_id"`
    DependentId    int64    `json:"dependent_id"`
    TaskType       int32    `json:"task_type"`
    Name           string   `json:"name"`
    Contact        string   `json:"contact"`
    Address        string   `json:"address"`
    Lat            float64  `json:"lat"`
    Lng            float64  `json:"lng"`
    Status         int32    `json:"status"`
    Created        int64    `json:"created"`
}

type TaskFilterMap struct {
    CityId     int64
    AgentId    int64
    UserId     int64
    RouteId    int64
    LocationId int64
    StartStamp int64
    EndStamp   int64
    Limit      int64
    Offset     int64
    Status     int32    
}

func InsertRouteTaskPoints(nodes *[]interface {}) (error){
 	query, _, _ := sqlBuilder.Insert(TASK_POINT_TABLE).Rows(*nodes...).ToSQL();
    rows, err := db.Query(query);
    defer rows.Close();
    if err != nil {
        fmt.Println("error: ", err)
        return err
    }
    return nil
}

func InsertRouteTaskPoint(taskPoint *sqlBuilder.Record) (*int64,error){
    query, _, _ := sqlBuilder.Insert(TASK_POINT_TABLE).Rows(*taskPoint).Returning("id").ToSQL();
    var task_point_id int64;
    err := db.QueryRow(query).Scan(&task_point_id);
    if err != nil {
        return nil, err
    }
    return &task_point_id, nil
}

func InsertRouteTask (city_id int64, route *Route) (*int64, error) {
    created := time.Now().Unix();
    visibleAfter := 0;
    query, _, _ := sqlBuilder.Insert(TASK_TABLE).Rows(
        sqlBuilder.Record {
            "route_id"     : route.Id,
            "visible_after": visibleAfter,
            "status"       : TASK_ASSIGNED,
            "agent_id"     : route.AgentId,
            "city_id"      : city_id,
            "created"      : created,
            "user_id"      : 0,
            "hub_id"       : 0,
            "start_after"  : 0,
            "end_after"    : 0,
            "auto_cancel"  : 0,
        },
    ).Returning("id").ToSQL();
    var task_id int64;
    err := db.QueryRow(query).Scan(&task_id);
    if err != nil {
        return nil, err;
    }
    return &task_id, nil;
}

func InsertLogisticsTask (task *Task) (*int64, error) {
    created := time.Now().Unix();
    visibleAfter := 0;
    query, _, _ := sqlBuilder.Insert(TASK_TABLE).Rows(
        sqlBuilder.Record {
            "route_id"     : task.RouteId,
            "visible_after": visibleAfter,
            "status"       : task.Status,
            "agent_id"     : task.AgentId,
            "vehicle_id"   : task.VehicleId,
            "city_id"      : task.CityId,
            "created"      : created,
            "user_id"      : task.UserId,
            "hub_id"       : task.HubId,
            "start_after"  : task.StartAfter,
            "end_after"    : task.EndAfter,
            "auto_cancel"  : task.AutoCancel,
        },
    ).Returning("id").ToSQL();
    var task_id int64;
    err := db.QueryRow(query).Scan(&task_id);
    if err != nil {
        return nil, err;
    }
    return &task_id, nil;
}

func GetCityTasks (city_id int64) (*[]*Task, error) {
    query, _, _ := sqlBuilder.Select(GetTaskColumns()...).From(TASK_TABLE).Where (
        sqlBuilder.Ex{"city_id": city_id},
    ).Order(sqlBuilder.I("id").Desc()).ToSQL();
    rows, err := db.Query(query);
    if err != nil {
        return nil, err;
    }
    defer rows.Close();
    var tasks []*Task;
    for rows.Next() {
        if err := rows.Err(); err != nil {          
            continue
        }
        task, err := StructTask(rows);
        if err == nil {         
            tasks = append(tasks, task);
        }else{          
            continue
        }
    }   
    return &tasks, nil
}

func FilterTasks (filterMap *TaskFilterMap) (*[]*Task, error) {
    var tasks []*Task;

    dbSQL := sqlBuilder.From(TASK_TABLE).Select(GetTaskColumns()...);    

    if filterMap.CityId == 0 {
        return nil, errors.New("City Id is required!");
    }else{
        dbSQL = dbSQL.Where(sqlBuilder.Ex {"city_id": filterMap.CityId})
    }

    if filterMap.AgentId != 0 {
        dbSQL = dbSQL.Where(sqlBuilder.Ex {"agent_id": filterMap.AgentId})
    }

    if filterMap.RouteId != 0 {
        dbSQL = dbSQL.Where(sqlBuilder.Ex {"route_id": filterMap.RouteId})
    }

    // TODO: Filter Task By Location ID
    // if filterMap.LocationId != 0 {
    //     dbSQL = dbSQL.Where(sqlBuilder.Ex {"location_id": filterMap.LocationId})
    // }

    if filterMap.UserId != 0 {
        dbSQL = dbSQL.Where(sqlBuilder.Ex {"user_id": filterMap.UserId})
    }


    if filterMap.Status == TASK_STATUS_AGENT_PENDING {
        dbSQL = dbSQL.Where(sqlBuilder.And(
            sqlBuilder.C("status").Gte(TASK_ACCEPTED),
            sqlBuilder.C("status").Lt(TASK_COMPLETED),
        ))
    }else if filterMap.Status == TASK_STATUS_PENDING {
        dbSQL = dbSQL.Where(sqlBuilder.And(
            sqlBuilder.C("status").Gte(TASK_CREATED),
            sqlBuilder.C("status").Lt(TASK_COMPLETED),
        ))
    }else if filterMap.Status == TASK_STATUS_COMPLETED {
        dbSQL = dbSQL.Where(
            sqlBuilder.C("status").Gte(TASK_COMPLETED),
        )
    }else if filterMap.Status != 99 {//TODO: 9 is temporary change later
        dbSQL = dbSQL.Where(sqlBuilder.Ex {"status": filterMap.Status})
    }

    if filterMap.StartStamp != 0 && filterMap.EndStamp != 0 {
        dbSQL = dbSQL.Where(sqlBuilder.And(
            sqlBuilder.C("created").Gte(filterMap.StartStamp),
            sqlBuilder.C("created").Lte(filterMap.EndStamp),
        ))
    }

    if filterMap.Limit == 0 {
        dbSQL = dbSQL.Limit(100);//TODO: Task Limit Optimize
    }else{
        dbSQL = dbSQL.Limit(uint(filterMap.Limit));
    }

    if filterMap.Offset == 0 {
        dbSQL = dbSQL.Offset(0);
    }else{
        dbSQL = dbSQL.Offset(uint(filterMap.Offset));
    }
    
    query, _, _ := dbSQL.Order(sqlBuilder.I("id").Desc()).ToSQL();    
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
        task, err := StructTask(rows);
        tasks = append(tasks, task);
    }
    return &tasks, nil;

}

func GetTaskPoints (task_id int64) (*[]*TaskPoint, error) {
    query, _, _ := sqlBuilder.Select(GetTaskPointColumns()...).From(TASK_POINT_TABLE).Where (
        sqlBuilder.Ex{"task_id": task_id},
    ).Order(sqlBuilder.I("id").Asc()).ToSQL();
    rows, err := db.Query(query);
    if err != nil {
        return nil, err;
    }
    defer rows.Close();
    var tasks []*TaskPoint;
    for rows.Next() {
        if err := rows.Err(); err != nil {          
            continue
        }
        task, err := StructTaskPoint(rows);
        if err == nil {         
            tasks = append(tasks, task);
        }else{          
            continue
        }
    }   
    return &tasks, nil
}

func GetRouteTasks (city_id int64) (*[]*pb.RouteTask, error) {
    tasks, err := GetCityTasks(city_id);    
    if err != nil {
        return nil, err
    }
    routeTasks, err := FormalizeRouteTasks(tasks);
    return routeTasks, err;
}

func GetRouteTaskById (task_id int64) (*pb.RouteTask, error) {
    var tasks []*Task
    task, err := GetTaskById(task_id);
    if err != nil {
        return nil, err
    }
    tasks = append(tasks, task);
    routeTasks, err := FormalizeRouteTasks(&tasks);
    if len(*routeTasks) == 0 {
        return nil, errors.New(errText);
    }else if err != nil {
        return nil, err
    }
    return (*routeTasks)[0], nil;
}

func FilterRouteTasks (routeTaskFilterMap *TaskFilterMap) (*[]*pb.RouteTask, error) {
    tasks, err := FilterTasks(routeTaskFilterMap);
    if err != nil {
        return nil, err
    }
    routeTasks, err := FormalizeRouteTasks(tasks);
    return routeTasks, err;
}

func GetTaskById(task_id int64) (*Task, error){  
    query, _, _ := sqlBuilder.Select(GetTaskColumns()...).From(TASK_TABLE).Where (
        sqlBuilder.Ex{"id": task_id},
    ).ToSQL();
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
    task, err := StructTask(row);
    if err != nil {
        return nil, err;
    }
    return task, nil;
}

func GetTaskPointById(task_point_id int64) (*TaskPoint, error){  
    query, _, _ := sqlBuilder.Select(GetTaskPointColumns()...).From(TASK_POINT_TABLE).Where (
        sqlBuilder.Ex{"id": task_point_id},
    ).ToSQL();
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
    taskPoint, err := StructTaskPoint(row);
    if err != nil {
        return nil, err;
    }
    return taskPoint, nil;
}

func AssignTaskAgent (task *Task, agent *Agent) error {
    err := UpdateAgentStatus(&Agent {
        Id: agent.Id,
        Status: AGENT_BUSY,
    });
    if err != nil {
        return err
    }

    sql, _, _ := sqlBuilder.Update(TASK_TABLE).Set(
        sqlBuilder.Record {
            "agent_id" : agent.Id,
            "status"   : TASK_ASSIGNED,
        },
    ).Where(sqlBuilder.Ex {"id": task.Id}).ToSQL();
    rows, err := db.Query(sql);
    rows.Close();
    if err != nil {
        return err;
    }

    sql, _, _ = sqlBuilder.Update(TASK_POINT_TABLE).Set(
        sqlBuilder.Record {
            "agent_id" : agent.Id,
        },
    ).Where(sqlBuilder.Ex {"task_id": task.Id}).ToSQL();
    rows, err = db.Query(sql);
    rows.Close();
    if err != nil {
        return err;
    }

    return nil;
}

func UpdateTaskAgent (task_id int64, agent_id int64) error {
    sql, _, _ := sqlBuilder.Update(TASK_TABLE).Set(sqlBuilder.Record {
            "status"   : TASK_ASSIGNED,
            "agent_id" : agent_id,
    },).Where(sqlBuilder.Ex {"id": task_id}).ToSQL();

    rows, err := db.Query(sql);    
    if err != nil {
        return err;
    }else{
        rows.Close();
    }

    sql, _, _ = sqlBuilder.Update(TASK_POINT_TABLE).Set(sqlBuilder.Record {
            "status"   : TASK_CREATED,
            "agent_id" : agent_id,
    }).Where(sqlBuilder.Ex {"task_id": task_id}).ToSQL();
    rows, err = db.Query(sql);    
    if err != nil {
        return err;
    }else{
        rows.Close();
    }

    return nil;
}

func UpdateAgentAcceptTask (task_id int64) error {
    sql, _, _ := sqlBuilder.Update(TASK_TABLE).Set(sqlBuilder.Record {
            "status"   : TASK_ACCEPTED,
    },).Where(sqlBuilder.Ex {"id": task_id}).ToSQL();
    rows, err := db.Query(sql);
    rows.Close();
    if err != nil {
        return err;
    }
    return nil;
}

func MarkTaskPointTransferred (task_point_id int64) error {
    sql, _, _ := sqlBuilder.Update(TASK_POINT_TABLE).Set(sqlBuilder.Record {
            "status"   : TASK_TRANSFERRED,
    },).Where(sqlBuilder.Ex {"id": task_point_id}).ToSQL();
    rows, err := db.Query(sql);
    defer rows.Close();
    if err != nil {
        return err;
    }
    return nil;
}

func UpdateAgentStartTask (task_point_id int64) error {
    sql, _, _ := sqlBuilder.Update(TASK_POINT_TABLE).Set(sqlBuilder.Record {
            "status"   : TASK_STARTED,
    },).Where(sqlBuilder.Ex {"id": task_point_id}).ToSQL();
    rows, err := db.Query(sql);
    defer rows.Close();
    if err != nil {
        return err;
    }
    return nil;
}

func UpdateAgentCompleteTask (task_point_id int64) (bool,error) {
    taskPoint, err := GetTaskPointById(task_point_id);
    if err != nil {
        return false, err
    }
    if taskPoint.DependentId != 0 {
        dependent, _ := GetTaskPointById(taskPoint.DependentId);
        if dependent != nil && dependent.Status != TASK_COMPLETED {
            return false, errors.New(dependent.Name + " Not Completed!");
        }
    }
    
    allTaskCompleted, taskPoint, err := AllTaskCompleted(task_point_id);
    if err != nil {
        return false, err
    }
    sql, _, _ := sqlBuilder.Update(TASK_POINT_TABLE).Set(sqlBuilder.Record {
            "status"   : TASK_COMPLETED,
    }).Where(sqlBuilder.Ex {"id": task_point_id}).ToSQL();
    rows, err := db.Query(sql);
    rows.Close();
    if err != nil {
        return false, err
    }
    if allTaskCompleted {
        sql, _, _ := sqlBuilder.Update(TASK_TABLE).Set(sqlBuilder.Record {
            "status"   : TASK_COMPLETED,
        }).Where(sqlBuilder.Ex {"id": taskPoint.TaskId}).ToSQL();
        rows, err := db.Query(sql);
        defer rows.Close();
        if err != nil {
            return false, err
        }

        go ValidateAgentStatus(taskPoint.AgentId);
    }    
    return allTaskCompleted, nil;
}

func UpdateCancelTask (cancel_status int32, task_point_id int64) (bool, error) {
    allTaskCompleted, taskPoint, err := AllTaskCompleted(task_point_id);
    if err != nil {
        return false,  err
    }
    sql, _, _ := sqlBuilder.Update(TASK_POINT_TABLE).Set(sqlBuilder.Record {
            "status"   : cancel_status,
    }).Where(sqlBuilder.Ex {"id": task_point_id}).ToSQL();
    rows, err := db.Query(sql);
    rows.Close();
    if err != nil {
        return false, err;
    }
    if allTaskCompleted {
        sql, _, _ := sqlBuilder.Update(TASK_TABLE).Set(sqlBuilder.Record {
            "status"   : TASK_COMPLETED,
        }).Where(sqlBuilder.Ex {"id": taskPoint.TaskId}).ToSQL();
        rows, err := db.Query(sql);
        defer rows.Close();
        if err != nil {
            return false, err;
        }
        go ValidateAgentStatus(taskPoint.AgentId);
    }    
    return allTaskCompleted, nil;
}

func AllTaskCompleted(task_point_id int64) (bool, *TaskPoint, error) {
    taskPoint, err := GetTaskPointById(task_point_id);    
    if err != nil {
        return false, nil, err
    }
    points, err := GetTaskPoints(taskPoint.TaskId);
    if err != nil {
        return false, nil, err
    }
    allTaskCompleted := true;
    for _, point := range *points {
        if point.Id != task_point_id && point.Status < TASK_COMPLETED {
            allTaskCompleted = false;
            break;
        }
    }
    return allTaskCompleted, taskPoint, nil
}

func FormalizeRouteTasks(tasks *[]*Task) (*[]*pb.RouteTask, error) {
    var routeTasks []*pb.RouteTask;
    for _, task := range *tasks {
        var agent *Agent;
        if task.AgentId != 0 {
            agent, _ = GetAgentById(task.AgentId);
        }
        var routeTask pb.RouteTask;
        var taskPoints []*pb.TaskPoint;
        var metaData pb.RouteTaskMeta;
        var completed int64;
        var remaining int64;
        points, err := GetTaskPoints(task.Id);        
        if err != nil {
            return nil, err
        }
        for _, point := range *points {
            protoTaskPoint := MakeProtoTaskPoint(point);
            taskPoints = append(taskPoints, protoTaskPoint);
            if point.Status > TASK_STARTED {
                completed = completed + 1
            }else{
                remaining = remaining + 1
            }
        }
        if agent != nil {
            metaData.Agent = &pb.Agent{
                Id: agent.Id,
                Name : agent.Name,
                Phone : agent.Phone,
                Avatar: agent.Avatar,
                Lat : agent.Location.Coordinates[0],
                Lng : agent.Location.Coordinates[1],
            }
        }
        if task.VehicleId != 0 {
            vehicle, _ := GetVehicleById(task.VehicleId);
            if vehicle != nil {
                metaData.Vehicle = vehicle;
            }            
        }
        metaData.Total = int64(len(*points));
        metaData.Completed = completed;
        metaData.Remaining = remaining;
        metaData.Status = task.Status;
        
        routeTask.TaskId = task.Id;
        routeTask.Points = taskPoints;
        routeTask.MetaData = &metaData;
        routeTasks = append(routeTasks, &routeTask);
    }
    return &routeTasks, nil;
}

func MarkCompleteForAutoCancelTask(city_id int64, end_after_limit int64) error {
    query, _, _ := sqlBuilder.Select(GetTaskColumns()...).From(TASK_TABLE).Where (
        sqlBuilder.Ex{"city_id": city_id},
        sqlBuilder.Ex{"auto_cancel": TASK_AUTO_CANCEL_ACTIVE},
        sqlBuilder.C("end_after").Lt(end_after_limit),
    ).Order(sqlBuilder.I("id").Desc()).ToSQL();
    rows, err := db.Query(query);
    if err != nil {
        return err;
    }
    defer rows.Close();
    for rows.Next() {
        if err := rows.Err(); err != nil {          
            continue
        }
        task, err := StructTask(rows);
        if err == nil {

            //Update Tasks
            query, _, _ = sqlBuilder.Update(TASK_POINT_TABLE).Set(sqlBuilder.Record {
                    "status"   : TASK_COMPLETED,
            }).Where(sqlBuilder.Ex {"task_id": task.Id}).ToSQL();
            taskRows, err := db.Query(query);
            taskRows.Close();
            if err != nil {
                continue
            }

            //Update Task Points
            query, _, _ = sqlBuilder.Update(TASK_TABLE).Set(sqlBuilder.Record {
                "status"   : TASK_COMPLETED,
            }).Where(sqlBuilder.Ex {"id": task.Id}).ToSQL();
            taskPointRows, err := db.Query(query);
            taskPointRows.Close();
            if err != nil {
                continue
            }
        }
    }   
    return nil;
}
//Counting functions
func CountCityCompletedTasks(city_id int64, start_time int64, end_time int64) (int) {
    var count int;
    dbSQL := sqlBuilder.From(TASK_TABLE).Select(sqlBuilder.COUNT("*")).Where(
        sqlBuilder.Ex {"city_id": city_id},
        sqlBuilder.C("status").Gte(TASK_COMPLETED),
        sqlBuilder.And(
            sqlBuilder.C("created").Gte(start_time),
            sqlBuilder.C("created").Lte(end_time),
        ),
    );
    query, _, _ := dbSQL.ToSQL();   
    err := db.QueryRow(query).Scan(&count);
    if err != nil {
        return 0;
    }
    return count;
}

func CountCityPendingTasks(city_id int64, start_time int64, end_time int64) (int) {
    var count int;
    dbSQL := sqlBuilder.From(TASK_TABLE).Select(sqlBuilder.COUNT("*")).Where(
        sqlBuilder.Ex {"city_id": city_id},
        sqlBuilder.And(
            sqlBuilder.C("status").Gte(TASK_CREATED),
            sqlBuilder.C("status").Lt(TASK_COMPLETED),
        ),
        sqlBuilder.And(
            sqlBuilder.C("created").Gte(start_time),
            sqlBuilder.C("created").Lte(end_time),
        ),
    );
    query, _, _ := dbSQL.ToSQL();   
    err := db.QueryRow(query).Scan(&count);
    if err != nil {
        return 0;
    }
    return count;
}

func CountCityAgents(city_id int64) (int) {
    var count int;
    dbSQL := sqlBuilder.From(AGENT_TABLE).Select(sqlBuilder.COUNT("*")).Where(
        sqlBuilder.Ex {"city_id": city_id},
    );
    query, _, _ := dbSQL.ToSQL();   
    err := db.QueryRow(query).Scan(&count);
    if err != nil {
        return 0;
    }
    return count;
}

func RollBackTaskAdd (task_id int64) {
	query, _, _ := sqlBuilder.Delete(TASK_TABLE).Where(sqlBuilder.Ex {"id": task_id}).ToSQL();
	rows, _ := db.Query(query);
	defer rows.Close();
}

func StructTask (row *sql.Rows) (*Task, error) {
    var task Task;
    err := row.Scan(    
        &task.Id,  
        &task.CityId,
        &task.RouteId,
        &task.HubId,
        &task.UserId,
        &task.AgentId,
        &task.VehicleId,
        &task.StartAfter,
        &task.EndAfter,    
        &task.VisibleAfter,
        &task.Status,
        &task.Created,
        &task.AutoCancel,
    );
    if err != nil {
        return nil, err
    }
    return &task, nil;
}

func StructTaskPoint (row *sql.Rows) (*TaskPoint, error) {
    var task TaskPoint;
    err := row.Scan(    
        &task.Id,  
        &task.TaskId,
        &task.HubId,
        &task.UserId,
        &task.AgentId,
        &task.SubscriptionId,
        &task.DependentId,
        &task.TaskType,
        &task.Name,
        &task.Contact,
        &task.Address,
        &task.Lat,
        &task.Lng,
        &task.Status,
        &task.Created,        
    );
    if err != nil {
        return nil, err
    }
    return &task, nil;
}

func MakeProtoTaskPoint (point *TaskPoint) *pb.TaskPoint {
    return &pb.TaskPoint {
        Id:point.Id,
        TaskId:point.TaskId,
        HubId:point.HubId,
        UserId:point.UserId,
        SubscriptionId:point.SubscriptionId,
        Lat:point.Lat,
        Lng:point.Lng,
        TaskType:point.TaskType,
        Status:point.Status,
        Created:point.Created,
        Name:point.Name,
        Contact:point.Contact,
        Address:point.Address,
        AgentId:point.AgentId,
    }
}