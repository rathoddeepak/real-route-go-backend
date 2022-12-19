/*
--------------------------------
 Author: Deepak Rathod
--------------------------------
 Date - 22 August 2022
 Phone: +917020814070
 Email: rathoddeepak537@gmail.com 
--------------------------------
 ---> Account Microservice <---
--------------------------------
--------------------------------
 --->   Logistics System   <---
--------------------------------
*/

package handler;

import (
	"fmt"
	"errors"
	"context"

	"encoding/json"

	pb "logisticsService/proto"
    db "logisticsService/database"
)

const AGENT_ALREADY_CREATED = "Agent Already Created!";
const OP_NOT_PERMITED = "Operation not permitted!";

const PointString string = "|ST_GeometryFromText('POINT(%v %v)')|";

func (lg *LogisticsService) CreateAgent (ctx context.Context, in *pb.CreateAgentRequest, out *pb.CreateAgentResponse) (error) {
	_, err := db.GetAgentByPhone(in.CityId, in.Phone);
	if err == nil {
		return errors.New(AGENT_ALREADY_CREATED);
	}
	agent := &db.Agent {
    	CityId   : in.CityId,
    	Name     : in.Name,
    	Phone    : in.Phone,
    	Passcode : in.Passcode,
    }
    agentId, err := db.InsertAgent(agent);
    if err != nil {
    	return err;
    }
    out.AgentId = *agentId;
    return nil;
}

func (lg *LogisticsService) UpdateAgent (ctx context.Context, in *pb.UpdateAgentRequest, out *pb.UpdateAgentResponse) (error) {
	agent := &db.Agent {
    	Id  	 : in.AgentId,
    	Name     : in.Name,
    	Status   : in.Status,
    	Passcode : in.Passcode,
    }
    err := db.UpdateAgent(agent);
    if err != nil {
    	return err;
    }
    out.Status = 200;
    out.Message = UPDATED_MSG;
    return nil;
}

func (lg *LogisticsService) UpdateAgentAvatar (ctx context.Context, in *pb.UpdateAgentRequest, out *pb.UpdateAgentResponse) (error) {
	//TODO delete previous image
	agent := &db.Agent {
    	Id  	 : in.AgentId,
    	Avatar   : in.Avatar,
    }
    err := db.UpdateAgentAvatar(agent);
    if err != nil {
    	return err;
    }
    out.Status = 200;
    out.Message = UPDATED_MSG;
    return nil;
}

func (lg *LogisticsService) UpdateAgentStatus (ctx context.Context, in *pb.UpdateAgentRequest, out *pb.UpdateAgentResponse) (error) {
	if in.Status != db.AGENT_OFFLINE && in.Status != db.AGENT_ONLINE {
		fmt.Println("h1")
		return errors.New(OP_NOT_PERMITED)
	}
	agent, err := db.GetAgentById(in.AgentId);
	if err != nil {
		return err
	}
	if agent.Status == in.Status {
		out.Status = in.Status;
		out.Message = UPDATED_MSG;
		return nil
	}
	busyTaskCount := db.CountAgentBusyTask(in.AgentId);
	if in.Status == db.AGENT_OFFLINE {		
		if busyTaskCount > 0 {
			return errors.New(OP_NOT_PERMITED)
		}
	}else if in.Status == db.AGENT_ONLINE {
		if busyTaskCount > 0 {
			in.Status = db.AGENT_BUSY
		}
	}
	if agent.Status == in.Status {
		out.Status = in.Status;
		out.Message = UPDATED_MSG;
		return nil
	}
    err = db.UpdateAgentStatus(&db.Agent {
    	Id  	 : in.AgentId,
    	Status   : in.Status,
    });
    if err != nil {
    	return err;
    }

    body, err := json.Marshal(&db.TaskUpdate_Realtime {
        Event 	 	: 	db.EVENT_AGENT_STATUS_UPDATE,
        Agent_id    :   in.AgentId,
        Status      :   in.Status,
        City_id     :   agent.CityId,
    });
    if err == nil {
        db.PublishMessage(body);
    }

    out.Status = in.Status;
    out.Message = UPDATED_MSG;
    return nil;
}

func (lg *LogisticsService) UpdateAgentLocation (ctx context.Context, in *pb.UpdateAgentRequest, out *pb.UpdateAgentResponse) (error) {
	agent, err := db.GetAgentById(in.AgentId);
	if err != nil {
		return err
	}
	point := fmt.Sprintf(PointString, in.Lat, in.Lng);
    err = db.UpdateAgentLocation(point, in.AgentId);
    if err != nil {
    	return err;
    }

    body, err := json.Marshal(&db.TaskUpdate_Realtime {
        Event 	 	: 	db.EVENT_AGENT_LOCATION_UPDATE,      
        Agent_id    :   in.AgentId,
        Agent_lat   :   in.Lat,
        Agent_lng   :   in.Lng,
        City_id     :   agent.CityId,
    });
    if err == nil {
        db.PublishMessage(body);
    }
    
    out.Status = 200;
    out.Message = UPDATED_MSG;
    return nil;
}

func (lg *LogisticsService) GetAgentsOfCity (ctx context.Context, in *pb.GetAgentRequest, out *pb.GetAgentsResponse) (error) {
	mAgents, err := db.GetAgentsOfCity(in.CityId);
	if err != nil {
		return err
	}
	var agents []*pb.Agent;
	for _, mAgent := range *mAgents {
	   agent := makeProtoAgent(mAgent);
       agents = append(agents, agent);
    }
    out.Agents = agents;
	return nil;	 
}

func (lg *LogisticsService) FilterAgent (ctx context.Context, in *pb.FilterAgentRequest, out *pb.GetAgentsResponse) (error) {
	mAgents, _, err := db.FilterAgents(in, false);
	if err != nil {
		return err
	}
	var agents []*pb.Agent;
	for _, mAgent := range *mAgents {
	   agent := makeProtoAgent(mAgent);
       agents = append(agents, agent);
    }
    out.Agents = agents;
	return nil;	 
}

func (lg *LogisticsService) GetAgentById (ctx context.Context, in *pb.GetAgentRequest, out *pb.GetAgentResponse) (error) {
	mAgent, err := db.GetAgentById(in.AgentId);
	if err != nil {
		return err
	}
	agent := makeProtoAgent(mAgent);
    out.Agent = agent;
	return nil;
}

func (lg *LogisticsService) VerifyAgent (ctx context.Context, in *pb.GetAgentRequest, out *pb.GetAgentResponse) (error) {
	mAgent, err := db.VerifyAgent(in.Phone, in.Passcode);
	if err != nil {
		return err
	}
	agent := makeProtoAgent(mAgent);
    out.Agent = agent;
	return nil;
}

func (lg *LogisticsService) InitAgentHome (ctx context.Context, in *pb.GetAgentRequest, out *pb.AgentHomeData) (error) {
	agent, err := db.GetAgentById(in.AgentId);
	if err != nil {
		return err
	}
	pending, err := db.FilterRouteTasks(&db.TaskFilterMap {
		AgentId : in.AgentId,
		Status  : db.TASK_STATUS_PENDING,
		CityId  : agent.CityId,
	});
	if err != nil {
		return err
	}
	var assigned []*pb.AssignedTask;
	tasks, err := db.FilterTasks(&db.TaskFilterMap {
		AgentId : in.AgentId,
		Status  : db.TASK_ASSIGNED,
		CityId  : agent.CityId,
	});
	if err != nil {
		return err
	}
	for _, task := range *tasks {
		assigned = append(assigned, &pb.AssignedTask {
			Id: task.Id,
			Time: task.Created,
		});
	}

	out.Pending = *pending;
	out.Assigned = assigned;
	out.Status = agent.Status;
	out.NotchText = "/ODR: 10";
	return nil;
}

func (lg *LogisticsService) AgentProfileData (ctx context.Context, in *pb.GetAgentRequest, out *pb.AgentProfileResponse) (error) {
	agent, err := db.GetAgentById(in.AgentId);
	if err != nil {
		return err
	}
	assigned := db.CountAgentAssignedTask(in.AgentId);
	pending := db.CountAgentBusyTask(in.AgentId);
	completed := db.CountAgentCompletedTask(in.AgentId);

	out.Completed = int64(completed);
	out.Assigned = int64(assigned);
	out.Pending = int64(pending);
	out.Status = agent.Status;
	return nil;
}

func makeProtoAgent (agent *db.Agent) (*pb.Agent){
	return &pb.Agent {
		Id       : agent.Id,
		CityId   : agent.CityId,
		Name     : agent.Name,
		Phone    : agent.Phone,
		Passcode : agent.Passcode,
		Avatar   : agent.Avatar,
		Status   : agent.Status,
		Lat   	 : agent.Location.Coordinates[0],
		Lng      : agent.Location.Coordinates[1],
	}
}