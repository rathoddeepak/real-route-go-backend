/*
--------------------------------
 Author: Deepak Rathod
--------------------------------
 Date - 26 August 2022
 Phone: +917020814070
 Email: rathoddeepak537@gmail.com 
--------------------------------
 --->   Logistics System   <---
--------------------------------
*/

package handler;

import (
    "fmt"

    "time"
	"errors"
	"context"

    "encoding/json"

    "go-micro.dev/v4/logger"

    sqlBuilder "github.com/doug-martin/goqu/v9"

	pb "logisticsService/proto"
    db "logisticsService/database"

    "go-micro.dev/v4/broker"
)

func (lg *LogisticsService) MakeTaskFromRoute (ctx context.Context, in *pb.MakeTaskFromRouteRequest, out *pb.MakeTaskFromRouteResponse) (error) {    
    var dependent_id *int64;
 	route, err := db.GetRouteById(in.RouteId);
 	if err != nil {
        logger.Info("Route Fetch Error!", in.RouteId);
 		return err
 	}
    mH, err := lg.HubService.GetHubById(ctx, &pb.GetHubRequest{
        HubId:route.HubId,
    });
    if err != nil {
        return err;
    }
    nodes, err := db.GetNodesOfRoute(in.RouteId);
    if len(*nodes) == 0 {
        return errors.New("Insufficient nodes!");
    }
    created := time.Now().Unix();
    task_id, err := db.InsertRouteTask(mH.Hub.CityId, route);
    if err != nil {
        logger.Info("Insert Error!");
        return err
    }
    dependent_index := 0;
    for _, node := range *nodes {        
        if node.NodeType == db.HUB_NODE && node.NodeId == mH.Hub.Id {
            dependent_id, err = db.InsertRouteTaskPoint(&sqlBuilder.Record {
                "task_id"           : task_id,
                "hub_id"            : mH.Hub.Id,
                "name"              : mH.Hub.Name,
                "contact"           : mH.Hub.Phone,
                "address"           : mH.Hub.Address,
                "subscription_id"   : 0,
                "lat"               : mH.Hub.Lat,
                "lng"               : mH.Hub.Lng,
                "status"            : db.TASK_CREATED,
                "created"           : created,
                "task_type"         : node.ActionType,
                "agent_id"          : route.AgentId,
                "dependent_id"      : 0,
                "user_id"           : 0,
            });
            if err != nil {
                return err;
            }
            break;
        }
        dependent_index = dependent_index + 1;
    }

    if *dependent_id != 0 {
        newNodes := RemoveItem(*nodes, dependent_index);
        nodes = &newNodes;
    }

    var rows []interface{};
    for _, node := range *nodes {        
        var lng float64;
        var lat float64;
        var subscription_id int64;
        var hub_id int64;
        address, name, contact := "", "", "";
        if node.NodeType == db.HUB_NODE {
            h, err := lg.HubService.GetHubById(ctx, &pb.GetHubRequest{
                HubId:node.NodeId,
            });
            if err != nil {
                logger.Info(err);
                return err
            }
            hub_id = h.Hub.Id;
            name = h.Hub.Name;
            contact = h.Hub.Phone;
            address = h.Hub.Address;
            lat = h.Hub.Lat;
            lng = h.Hub.Lng;
        }else if node.NodeType == db.SUBSCRIPTION_NODE {
            scp, err := db.GetSubscriptionById(node.NodeId);
            if err != nil {
                logger.Info("Get Sub Error!");
                return err
            }
            a, err := lg.AccountService.GetAddressById(ctx, &pb.GetAddressRequest{
                AddressId:scp.AddressId,
            });
            if err != nil {
                logger.Info("Get Address Error!");
                return err
            }
            u, err := lg.AccountService.GetUserById(ctx, &pb.GetUserRequest{
                UserId:scp.UserId,
            });
            if err != nil {
                logger.Info("Get User Error!");
                return err
            }
            subscription_id = scp.Id;
            name = u.User.Name;
            contact = u.User.Phone;
            address = a.Address.Address;
            lat = a.Address.Lat
            lng = a.Address.Lng
        }
        rows = append(rows, sqlBuilder.Record {
            "task_id"         : task_id,
            "hub_id"          : hub_id,
            "name"            : name,
            "contact"         : contact,
            "address"         : address,
            "subscription_id" : subscription_id,
            "lat"             : lat,
            "lng"             : lng,
            "status"          : db.TASK_CREATED,
            "created"         : created,
            "task_type"       : node.ActionType,
            "agent_id"        : route.AgentId,
            "dependent_id"    : *dependent_id,
            "user_id"         : 0,
        });
    }   
    err = db.InsertRouteTaskPoints(&rows);
    if err != nil {
        go db.RollBackTaskAdd(*task_id);
        logger.Info("Insert Route Error!");
        return err
    }
    broadCast := func () {
        body, err := json.Marshal(&db.TaskUpdate_Realtime {            
            Task_id:*task_id,
            Event:db.EVENT_TASK_CREATED,
            City_id:mH.Hub.CityId,
        });        
        if err == nil {
            db.PublishMessage(body);
        }     
    }
    go broadCast();
    out.TaskId = *task_id;
    return nil;
}

func (lg *LogisticsService) GetRouteTask (ctx context.Context, in *pb.GetRouteTaskRequest, out *pb.GetRouteTasksResponse) (error) {
    tasks, err := db.FilterRouteTasks(&db.TaskFilterMap {
        Status  : in.Status,        
        CityId  : in.CityId,
        StartStamp : in.StartStamp,
        EndStamp : in.EndStamp,
        AgentId : in.AgentId,
        UserId : in.UserId,
        LocationId : in.LocationId,
        RouteId : in.RouteId,
    });
    if err != nil {
        return err;
    }
    out.Tasks = *tasks;
    return nil;
}

func (lg *LogisticsService) AppHomeData (ctx context.Context, in *pb.GetRouteTaskRequest, out *pb.GetAppHomeResponse) (error) {
    //Task Will Be Marked As Completed When Auto Cancel
    _ = db.MarkCompleteForAutoCancelTask(in.CityId, time.Now().Unix());

    tasks, err := db.FilterRouteTasks(&db.TaskFilterMap {
        Status  : db.TASK_STATUS_PENDING,        
        CityId  : in.CityId,
        StartStamp : in.StartStamp,
        EndStamp : in.EndStamp,
        Limit : int64(10),
    });
    if err != nil {
        return err;
    }
    taskCount := db.CountCityCompletedTasks(in.CityId, in.StartStamp, in.EndStamp);
    pendingCount := db.CountCityPendingTasks(in.CityId, in.StartStamp, in.EndStamp);
    agentCount := db.CountCityAgents(in.CityId);
    out.TaskCount = int64(taskCount);
    out.PendingCount = int64(pendingCount);
    out.AgentCount = int64(agentCount);
    out.Tasks = *tasks;
    return nil;
}

func (lg *LogisticsService) GetRouteTaskById (ctx context.Context, in *pb.GetRouteTaskRequest, out *pb.GetRouteTaskResponse) (error) {
    task, err := db.GetRouteTaskById(in.TaskId);
    if err != nil {
        return err
    }
    out.Task = task;
    return nil;
}

//Plain Task Functions 
func (lg *LogisticsService) CreateTaskFromPoints (ctx context.Context, in *pb.CreateTaskRequest, out *pb.CreateTaskResponse) (error) {
    if len(in.Points) == 0 {
        return errors.New("Insufficient Points!");
    }

    taskPointStatus := db.TASK_CREATED;
    created := time.Now().Unix();
    status := int32(db.TASK_CREATED);
    if in.Task.AgentId != 0 {
        status = int32(db.TASK_ASSIGNED);
    }
    autoCancel := db.TASK_AUTO_CANCEL_INACTIVE;
    cityId := in.Task.CityId;
    if in.Template != nil {
        if in.Template.Name == "bus_khabri" {
            autoCancel = db.TASK_AUTO_CANCEL_ACTIVE;
            //TODO temp fix for city_id grabing...!
            cityResp, err := lg.CityService.GetCities(ctx, &pb.GetCityRequest {
                CompanyId: in.Template.CompanyId,
            });
            if err != nil {
                return errors.New("City Service Error!");
            }
            if len(cityResp.Cities) > 0 {
                cityId = cityResp.Cities[0].Id;
            }else{
                return errors.New("Atleast Create 1 City!");
            }

            if status == int32(db.TASK_ASSIGNED) {
                status = db.TASK_ACCEPTED;
                taskPointStatus = db.TASK_STARTED;
            }
        } else {
            return errors.New("Invalid Template!");
        }
    }
    task_id, err := db.InsertLogisticsTask(& db.Task {
        CityId : cityId,
        VehicleId : in.Task.VehicleId,
        StartAfter : in.Task.StartAfter,
        EndAfter   : in.Task.EndAfter,
        AgentId: in.Task.AgentId,
        Status : status,
        AutoCancel : int32(autoCancel),
    });
    if err != nil {
        logger.Info(err);
        return err
    }
    var rows []interface{};
    for _, taskPoint := range in.Points {        
        rows = append(rows, sqlBuilder.Record {
            "task_id"         : task_id,
            "hub_id"          : 0,
            "name"            : taskPoint.Name,
            "contact"         : taskPoint.Contact,
            "address"         : taskPoint.Address,
            "subscription_id" : 0,
            "lat"             : taskPoint.Lat,
            "lng"             : taskPoint.Lng,
            "status"          : taskPointStatus,
            "created"         : created,
            "task_type"       : taskPoint.TaskType,
            "agent_id"        : in.Task.AgentId,
            "dependent_id"    : 0,
            "user_id"         : 0,
        });
    }   
    err = db.InsertRouteTaskPoints(&rows);
    if err != nil {
        go db.RollBackTaskAdd(*task_id);
        logger.Info("Insert Route Error!");
        return err
    }
    broadCast := func () {
        body, err := json.Marshal(&db.TaskUpdate_Realtime {            
            Task_id:*task_id,
            Event:db.EVENT_TASK_CREATED,
            City_id:in.Task.CityId,
        });        
        if err == nil {
            db.PublishMessage(body);
        }     
    }
    go broadCast();
    out.TaskId = *task_id;
    return nil;
}

//Main Funciton Related To Updation of task!
func (lg *LogisticsService) UpdateRouteTask (ctx context.Context, in *pb.UpdateRouteTaskRequest, out *pb.UpdateRouteResponse) (error) {    
    var err error;
    var taskPoint *db.TaskPoint;
    var task_id int64;
    if in.UpdateStatus != db.TASK_ASSIGNED && in.UpdateStatus != db.TASK_ACCEPTED && in.TaskPointId == 0 {
        logger.Info("error here!")
        return errors.New(TASK_REQUIRED_ERR);
    }
    if in.TaskPointId != 0 {
        taskPoint, err = db.GetTaskPointById(in.TaskPointId);
        if err == nil {
            task_id = taskPoint.TaskId;
        }
    }
    
    cityObject, err := lg.CityService.GetCityById(ctx, &pb.GetCityRequest {
        CityId: in.CityId,
    });
    if err != nil {
        return errors.New("City not found!");
    }

    companyObject, err := lg.AccountService.GetCompany(ctx, &pb.GetCompanyRequest {
        CompanyId: cityObject.City.CompanyId,
    })

    if err != nil {
        return errors.New("Company not found!")
    }

    settings := companyObject.Company.Settings;

    var broadCast func();
    //Task Assigned
    if db.TASK_ASSIGNED == in.UpdateStatus {
        if in.TaskId == 0 {
            return errors.New(TASK_REQUIRED_ERR);
        }else if in.AgentId == 0 {
            return errors.New(AGENT_REQUIRED_ERR);
        }
        task, err := db.GetTaskById(in.TaskId);
        if err != nil {
            return errors.New(TASK_REQUIRED_ERR);
        }
        agent, err := db.GetAgentById(in.AgentId);
        if err != nil {
            return errors.New(AGENT_REQUIRED_ERR);
        }
        err = db.AssignTaskAgent(task, agent);
        if err != nil {
            return err
        }
        broadCast = func () {
            body, err := json.Marshal(&db.TaskUpdate_Realtime {
                Event       :   db.EVENT_TASK_ASSIGNED,
                Task_id     :   in.TaskId,
                City_id     :   in.CityId,
                Status      :   in.UpdateStatus,
                Agent_id    :   in.AgentId,
                Agent_name  :   agent.Name,
                Agent_phone :   agent.Phone,
                Agent_lat   :   agent.Location.Coordinates[0],
                Agent_lng   :   agent.Location.Coordinates[1],
            });
            if err == nil {
                db.PublishMessage(body);
            }

            //Sending Notification
            body, err = json.Marshal(&pb.SendNotificaitonRequest {
                Title    : fmt.Sprintf("Task #%v Assigned!", in.TaskId),
                Content  : "Press here to view details!",
                AgentIds : [] int64 {in.AgentId},
            });
            if err == nil {
                if err := broker.Publish(notificationTopic, &broker.Message{Body: body}); err != nil {
                    logger.Info(err);
                }
            } 
        }
    //Agent Task Accepted
    }else if db.TASK_ACCEPTED == in.UpdateStatus {
        err = db.UpdateAgentAcceptTask(in.TaskId);
        if err != nil {
            return err
        }
        broadCast = func () {
            //Sending Realtime Updates
            body, err := json.Marshal(&db.TaskUpdate_Realtime {            
                Event       :   db.EVENT_TASK_ACCEPTED,
                Task_id     :   in.TaskId,
                City_id     :   in.CityId,
                Status      :   in.UpdateStatus,             
            });
            if err == nil {
                db.PublishMessage(body);
            }
            //TODO:
            // if settings.Notification.OnStartTask == false {
            //     return
            // }
            // Sending Notification
            // body, err = json.Marshal(&pb.SendNotificaitonRequest {
            //     Title : fmt.Sprintf("Task #%v Accepted", in.TaskId),
            //     Content : "Press here to view details!",
            //     CityIds: [] int64 {in.CityId},
            // });
            // if err == nil {
            //     if err := broker.Publish(notificationTopic, &broker.Message{Body: body}); err != nil {
            //         logger.Info(err);
            //     }
            // }            
        }
    //Start Task
    }else if db.TASK_STARTED == in.UpdateStatus {
        err = db.UpdateAgentStartTask(in.TaskPointId);
        if err != nil {
            return err
        }
        broadCast = func () {
            body, err := json.Marshal(&db.TaskUpdate_Realtime {
                Event       :   db.EVENT_TASK_STARTED,
                Point_id    :   in.TaskPointId,
                Task_id     :   task_id,
                City_id     :   in.CityId,
                Status      :   in.UpdateStatus,             
            });
            if err == nil {
                db.PublishMessage(body);
            }


            if settings.Notification.OnStartTask == false {
                return
            }
            //Sending Notification
            taskPoint, err := db.GetTaskPointById(in.TaskPointId);
            if err != nil {
                return
            }
            body, err = json.Marshal(&pb.SendNotificaitonRequest {
                Title : fmt.Sprintf("Task %v Started!", taskPoint.Name),
                Content : "Press here to view details!",
                CityIds: [] int64 {in.CityId},
            });
            if err == nil {
                if err := broker.Publish(notificationTopic, &broker.Message{Body: body}); err != nil {
                    logger.Info(err);
                }
            }
        }        
    //Complete Task
    }else if db.TASK_COMPLETED == in.UpdateStatus {
        ac, err := db.UpdateAgentCompleteTask(in.TaskPointId);
        if err != nil {
            return err
        }
        event := db.EVENT_TASK_COMPLETED;
        if ac {
            event = db.EVENT_TASK_FULL_COMPLETED;
        }
        broadCast = func () {
            body, err := json.Marshal(&db.TaskUpdate_Realtime {
                Event       :   event,
                Point_id    :   in.TaskPointId,
                City_id     :   in.CityId,
                Task_id     :   task_id,
                Status      :   in.UpdateStatus,             
            });
            if err == nil {
                db.PublishMessage(body);
            }

            //Sending Notification
            taskPoint, err := db.GetTaskPointById(in.TaskPointId);
            if err != nil {
                return
            }
            body, err = json.Marshal(&pb.SendNotificaitonRequest {
                Title : fmt.Sprintf("Task %v Completed!", taskPoint.Name),
                Content : "Press here to view details!",
                CityIds: [] int64 {in.CityId},
            });
            if err == nil {
                if err := broker.Publish(notificationTopic, &broker.Message{Body: body}); err != nil {
                    logger.Info(err);
                }
            }            
        }
    }else if db.TASK_CUSTOMER_CANCEL == in.UpdateStatus {
        ac, err := db.UpdateCancelTask(db.TASK_CUSTOMER_CANCEL, in.TaskPointId);
        if err != nil {
            return err
        }
        event := db.EVENT_TASK_CUSTOMER_CANCEL;
        if ac {
            event = db.EVENT_TASK_FULL_COMPLETED;
        }
        broadCast = func () {
            body, err := json.Marshal(&db.TaskUpdate_Realtime {
                Event       :   event,    
                Point_id    :   in.TaskPointId,
                City_id     :   in.CityId,
                Task_id     :   task_id,
                Status      :   in.UpdateStatus,             
            });
            if err == nil {
                db.PublishMessage(body);
            }

            //Sending Notification
            taskPoint, err := db.GetTaskPointById(in.TaskPointId);
            if err != nil {
                return
            }
            body, err = json.Marshal(&pb.SendNotificaitonRequest {
                Title : fmt.Sprintf("Task %v Cancelled By Customer!", taskPoint.Name),
                Content : "Press here to view details!",
                CityIds: [] int64 {in.CityId},
            });
            if err == nil {
                if err := broker.Publish(notificationTopic, &broker.Message{Body: body}); err != nil {
                    logger.Info(err);
                }
            }   
        }
    }else if db.TASK_AGENT_CANCEL == in.UpdateStatus {
        ac, err := db.UpdateCancelTask(db.TASK_AGENT_CANCEL, in.TaskPointId);
        if err != nil {
            return err
        }
        event := db.EVENT_TASK_AGENT_CANCEL;
        if ac {
            event = db.EVENT_TASK_FULL_COMPLETED;
        }
        broadCast = func () {
            body, err := json.Marshal(&db.TaskUpdate_Realtime {
                Event       :   event,
                Point_id    :   in.TaskPointId,
                City_id     :   in.CityId,
                Task_id     :   task_id,
                Status      :   in.UpdateStatus,             
            });
            if err == nil {
                db.PublishMessage(body);
            }

            //Sending Notification
            taskPoint, err := db.GetTaskPointById(in.TaskPointId);
            if err != nil {
                return
            }
            body, err = json.Marshal(&pb.SendNotificaitonRequest {
                Title : fmt.Sprintf("Task %v Cancelled By Agent!", taskPoint.Name),
                Content : "Press here to view details!",
                CityIds: [] int64 {in.CityId},
            });
            if err == nil {
                if err := broker.Publish(notificationTopic, &broker.Message{Body: body}); err != nil {
                    logger.Info(err);
                }
            }           
        }
    }else if db.TASK_SYSTEM_CANCEL_PAID == in.UpdateStatus {
        ac, err := db.UpdateCancelTask(db.TASK_SYSTEM_CANCEL_PAID, in.TaskPointId);
        if err != nil {
            return err
        }
        event := db.EVENT_TASK_SYSTEM_CANCEL_PAID;
        if ac {
            event = db.EVENT_TASK_FULL_COMPLETED;
        }
        broadCast = func () {
            body, err := json.Marshal(&db.TaskUpdate_Realtime {
                Event       :   event,
                Point_id    :   in.TaskPointId,
                City_id     :   in.CityId,
                Task_id     :   task_id,
                Status      :   in.UpdateStatus,             
            });
            if err == nil {
                db.PublishMessage(body);
            }

            //Sending Notification
            taskPoint, err := db.GetTaskPointById(in.TaskPointId);
            if err != nil {
                return
            }
            body, err = json.Marshal(&pb.SendNotificaitonRequest {
                Title : fmt.Sprintf("Task %v Cancelled!", taskPoint.Name),
                Content : "Press here to view details!",
                CityIds: [] int64 {in.CityId},
            });
            if err == nil {
                if err := broker.Publish(notificationTopic, &broker.Message{Body: body}); err != nil {
                    logger.Info(err);
                }
            }   
        }
    }else if db.TASK_SYSTEM_CANCEL_NOT_PAID == in.UpdateStatus {
        ac, err := db.UpdateCancelTask(db.TASK_SYSTEM_CANCEL_NOT_PAID, in.TaskPointId);
        if err != nil {
            return err
        }
        event := db.EVENT_TASK_SYSTEM_CANCEL_NOT_PAID;
        if ac {
            event = db.EVENT_TASK_FULL_COMPLETED;
        }
        broadCast = func () {
            body, err := json.Marshal(&db.TaskUpdate_Realtime {
                Event       :   event,
                Point_id    :   in.TaskPointId,
                City_id     :   in.CityId,
                Task_id     :   task_id,
                Status      :   in.UpdateStatus,             
            });
            if err == nil {
                db.PublishMessage(body);
            }

            //Sending Notification
            taskPoint, err := db.GetTaskPointById(in.TaskPointId);
            if err != nil {
                return
            }
            body, err = json.Marshal(&pb.SendNotificaitonRequest {
                Title : fmt.Sprintf("Task %v Cancelled!", taskPoint.Name),
                Content : "Press here to view details!",
                CityIds: [] int64 {in.CityId},
            });
            if err == nil {
                if err := broker.Publish(notificationTopic, &broker.Message{Body: body}); err != nil {
                    logger.Info(err);
                }
            }
        }
    }else{
        return errors.New("Invalid status")
    }
    if err != nil {
        return err;
    }
    if broadCast != nil {
        go broadCast();
    }
    out.Status = 200;
    out.Message = UPDATED_MSG;
    return nil;
}

func (lg *LogisticsService) AssignTaskAgent (ctx context.Context, in *pb.AssignTaskAgentRequest, out *pb.UpdateRouteResponse) (error) {    
    task, err := db.GetTaskById(in.TaskId);
    if err != nil {
        return errors.New(TASK_REQUIRED_ERR);
    }

    agent, err := db.GetAgentById(in.AgentId);
    if err != nil {
        return errors.New(AGENT_REQUIRED_ERR);
    }    

    if task.AgentId == in.AgentId {
        return errors.New("Already Same Agent!");
    }

    if task.CityId != agent.CityId {
        return errors.New("Internal Error!");   
    }
    err = db.UpdateTaskAgent(in.TaskId, in.AgentId);
    if err != nil {
        return err
    }
    broadCast := func () {
        body, err := json.Marshal(&db.TaskUpdate_Realtime {            
            Agent_id:in.AgentId,
            Agent_phone: agent.Phone,
            Agent_name: agent.Name,
            Task_id: in.TaskId,
            City_id: agent.CityId,
            Event:db.EVENT_TASK_AGENT_CHANGE,
        });
        if err == nil {
            db.PublishMessage(body);
        }     
    }
    go broadCast();
    out.Status = 200;
    out.Message = UPDATED_MSG;
    return nil;
}

func (lg *LogisticsService) TransferTaskPoint (ctx context.Context, in *pb.TansferTaskPointRequest, out *pb.UpdateRouteResponse) (error) {
    taskPoint, err := db.GetTaskPointById(in.TaskPointId);
    if err != nil {
        return err;
    }

    task, err := db.GetTaskById(taskPoint.TaskId);
    if err != nil {
        return errors.New(TASK_REQUIRED_ERR);
    }

    if taskPoint.AgentId == in.AgentId {
        return errors.New("Already Same Agent!");
    }

    //TODO thinking
    // if taskPoint.Status == TASK_STATUS_COMPLETED {
    //     return errors.New("Cannot Transfer Completed Point!");
    // }

    err = db.MarkTaskPointTransferred(in.TaskPointId);
    if err != nil {
        return err
    }
    broadCast := func () {
        body, err := json.Marshal(&db.TaskUpdate_Realtime {            
            Point_id:in.TaskPointId,
            Task_id: task.Id,
            City_id: task.CityId,
            Status: db.TASK_TRANSFERRED,
            Event:db.EVENT_TASK_TRANSFERRED,
        });
        if err == nil {
            db.PublishMessage(body);
        }     
    }
    go broadCast();

    task_id, err := db.InsertLogisticsTask(&db.Task {
        CityId: task.CityId,
        AgentId: in.AgentId,
        Status: db.TASK_ASSIGNED,
    });
    if err != nil {
        return err;
    }
    created := time.Now().Unix();
    _, err = db.InsertRouteTaskPoint(&sqlBuilder.Record {
        "task_id"           : task_id,
        "hub_id"            : taskPoint.HubId,
        "name"              : taskPoint.Name,
        "contact"           : taskPoint.Contact,
        "address"           : taskPoint.Address,
        "subscription_id"   : taskPoint.SubscriptionId,
        "lat"               : taskPoint.Lat,
        "lng"               : taskPoint.Lng,
        "status"            : db.TASK_CREATED,
        "created"           : created,
        "task_type"         : taskPoint.TaskType,
        "agent_id"          : in.AgentId,
        "dependent_id"      : 0,
        "user_id"           : 0,
    });
    if err != nil {
        go db.RollBackTaskAdd(*task_id);
        return err
    }
    broadCast2 := func () {
        body, err := json.Marshal(&db.TaskUpdate_Realtime {            
            Task_id:*task_id,
            City_id: task.CityId,
            Event:db.EVENT_TASK_CREATED,
        });
        if err == nil {
            db.PublishMessage(body);
        }     
    }
    go broadCast2();
    out.Status = 200;
    out.Message = UPDATED_MSG;
    return nil;
}