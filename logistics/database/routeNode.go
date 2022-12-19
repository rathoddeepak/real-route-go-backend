/*
--------------------------------
 Author: Deepak Rathod
--------------------------------
 Date - 23 August 2022
 Phone: +917020814070
 Email: rathoddeepak537@gmail.com 
--------------------------------
 --->   Logistics System   <---
--------------------------------
*/
package database;

import (
	"errors"
	_ "github.com/lib/pq"	
	"database/sql"
	pb "logisticsService/proto"
	sqlBuilder "github.com/doug-martin/goqu/v9"
)

type RouteNode struct {
	Id 	   		   int64 	 `json:"id"`
	RouteId 	   int64 	 `json:"route_id"`
	ActionType 	   int32     `json:"action_type"`
	NodeType       int32 	 `json:"node_type"`
	NodeId         int64	 `json:"node_id"`
	NodeName       string 	 `json:"node_name"`
}

func InsertRouteNodes(routeNode *[]*pb.RouteNode) error {
	var records []interface{};
	for _, node := range *routeNode {
		records = append(
			records,
			sqlBuilder.Record {
				"route_id":node.RouteId,
				"action_type":node.ActionType,
				"node_type":node.NodeType,
				"node_id":node.NodeId,			
				"seq":node.Seq,
			},
		);
	}
	sql, _, _ := sqlBuilder.Insert(ROUTE_NODE_TABLE).Rows(records...).ToSQL();
	rows, err := db.Query(sql);	
	if err != nil {
		return err;
	}
	rows.Close();
	return nil
}

func ClearRouteNodes (route_id int64) error {
	sql, _, _ := sqlBuilder.Delete(ROUTE_NODE_TABLE).Where(sqlBuilder.Ex {"route_id": route_id}).ToSQL();
	rows, err := db.Query(sql);
	defer rows.Close();
	if err != nil {
		return err;
	}
	return nil;
}

func GetNodesOfRoute(route_id int64) (*[]*RouteNode, error) {
	query, _, _ := sqlBuilder.Select(GetRouteNodeColumns()...).From(ROUTE_NODE_TABLE).Where (
		sqlBuilder.Ex{"route_id": route_id},
	).Order(sqlBuilder.I("seq").Asc()).ToSQL();
	rows, err := db.Query(query);
	if err != nil {
		return nil, err;
	}
	defer rows.Close();
	var routeNodes []*RouteNode;
	for rows.Next() {
		if err := rows.Err(); err != nil {			
			continue
		}
		routeNode, err := StructRouteNode(rows);
		if err == nil {			
			routeNodes = append(routeNodes, routeNode);
		}else{			
			continue
		}
	}	
	return &routeNodes, nil
}

func SubscriptionRouteNode(scp_id int64) (*RouteNode, error){	
	query, _, _ := sqlBuilder.Select(GetRouteNodeColumns()...).From(ROUTE_NODE_TABLE).Where (
		sqlBuilder.Ex{"node_type": SUBSCRIPTION_NODE},
		sqlBuilder.Ex{"node_id": scp_id},
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
	routeNode, err := StructRouteNode(row);
	if err != nil {
		return nil, err;
	}
	return routeNode, nil;
}
func StructRouteNode (row *sql.Rows) (*RouteNode, error) {
	var node RouteNode;
	err := row.Scan(	
		&node.Id,	
		&node.RouteId,
		&node.ActionType,
		&node.NodeType,
		&node.NodeId,
	);
	if err != nil {
		return nil, err
	}
	return &node, nil;
}