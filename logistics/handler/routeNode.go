/*
--------------------------------
 Author: Deepak Rathod
--------------------------------
 Date - 23 August 2022
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

    "context"

    "errors"

    pb "logisticsService/proto"
    db "logisticsService/database"
)

func (lg *LogisticsService) SaveRouteNode (ctx context.Context, in *pb.SaveRouteNodeRequest, out *pb.UpdateRouteResponse) (error) {
    db.ClearRouteNodes(in.RouteId);
    if(len(in.Nodes) < 2){
        return errors.New("Atleast Two Nodes are required!");
    }
    // TODO :: FOR HUB SECTION
    // if(in.Nodes[0].NodeType != db.HUB_NODE){
    //     return errors.New("Hub Should be on first node!");
    // }
    err := db.InsertRouteNodes(&in.Nodes);
    if err != nil {
        return err;
    }
    out.Status = 200;
    out.Message = "Nodes Created!";
    return nil;
}

func (lg *LogisticsService) GetRouteNodeOfRoute (ctx context.Context, in *pb.GetRouteNodeRequest, out *pb.GetRouteNodeResponse) (error) {
    route, err := db.GetRouteById(in.RouteId);
    if err != nil {
        fmt.Println("route err", err)
        return err;
    }

    mNodes, err := db.GetNodesOfRoute(in.RouteId);
    if err != nil {
        fmt.Println("nodes err", err)
        return err
    }    
    // hubNodeFound := false;
    var nodes []*pb.RouteNode;
    var scpIds []interface{};
    for _, mNode := range *mNodes {     
        var name string; 
        var contact string; 
        var address string; 
        var lat float64;
        var lng float64; 
        if mNode.NodeType == db.HUB_NODE {
            h, err := lg.HubService.GetHubById(ctx, &pb.GetHubRequest {HubId:mNode.NodeId});
            if err != nil {
                fmt.Println("hub ", err);
                continue
            }else{
                name = h.Hub.Name;
                lat = h.Hub.Lat;
                lng = h.Hub.Lng;
            }
            // hubNodeFound = true;
        }else if mNode.NodeType == db.SUBSCRIPTION_NODE {
            subscription, err := db.GetSubscriptionById(mNode.NodeId);
            if err != nil {
                fmt.Println("sub", err);
                continue;
            }
            u, err := lg.AccountService.GetUserById(ctx, &pb.GetUserRequest {UserId:subscription.UserId});            
            if err != nil {
                fmt.Println("ac", err);
                continue;
            }
            a, err := lg.AccountService.GetAddressById(ctx, &pb.GetAddressRequest {AddressId:subscription.AddressId});
            if err != nil {
                continue
            }
            name = u.User.Name;
            contact = u.User.Phone;
            address = a.Address.Address;
            lat = a.Address.Lat;
            lng = a.Address.Lng;
            scpIds = append(scpIds, fmt.Sprintf("%v", subscription.Id));
        }else if mNode.NodeType == db.LOCATION_NODE {
            ltn, err := db.GetLocationById(mNode.NodeId);
            if err != nil {
                fmt.Println("Location:  ", err);
                continue
            }else{
                name = ltn.Name;
                address = ltn.Address;
                contact = ltn.Contact;
                lat = ltn.Location.Coordinates[0];
                lng = ltn.Location.Coordinates[1];
            }
        }else if mNode.NodeType == db.ADDRESS_NODE {
            a, err := lg.AccountService.GetAddressById(ctx, &pb.GetAddressRequest {AddressId:mNode.NodeId});
            if err != nil {
                fmt.Println("ad", err);
                continue
            }
            u, err := lg.AccountService.GetUserById(ctx, &pb.GetUserRequest {UserId:a.Address.UserId});            
            if err != nil {
                fmt.Println("ac", err);
                continue;
            }
            name = u.User.Name;
            contact = u.User.Phone;
            address = a.Address.Address;
            lat = a.Address.Lat;
            lng = a.Address.Lng;
        } else{
            continue
        }
        node := makeProtoRouteNode(mNode, lat, lng, name, contact, address);
        nodes = append(nodes, node);
    }
    out.Connected = nodes;

    var disconnected []*pb.RouteNode;
    //TODO
    //For Subscription Purpose Adding Hub For Disconnected Route
    // if hubNodeFound == false {
    //     h, err := lg.HubService.GetHubById(ctx, &pb.GetHubRequest {HubId:route.HubId});
    //     if err != nil {
    //         fmt.Println(err);
    //         return err;
    //     }
    //     disconnected = append(disconnected, &pb.RouteNode {
    //         Id:0,
    //         RouteId:in.RouteId,
    //         ActionType:db.NODE_ACTION_PICKIUP,
    //         NodeType:db.HUB_NODE,
    //         NodeId:route.HubId,
    //         NodeName:h.Hub.Name,
    //         Lat:h.Hub.Lat,
    //         Lng:h.Hub.Lng,
    //     });
    // }
    //For Subscription Purpose Adding Subscription For Disconnected Route
    mSubscriptions, err := db.GetHubSubscriptionsExclude(route.HubId, scpIds);
    if err != nil {
        return err;
    }    
    for _, subscription := range *mSubscriptions {

        _, err := db.SubscriptionRouteNode(subscription.Id);
        if err == nil {
            continue
        }

        u, err := lg.AccountService.GetUserById(ctx, &pb.GetUserRequest {UserId:subscription.UserId});            
        if err != nil {
            fmt.Println(err);
            continue;
        }
        a, err := lg.AccountService.GetAddressById(ctx, &pb.GetAddressRequest {AddressId:subscription.AddressId});
        if err != nil {
            continue
        }
        disconnected = append(disconnected, &pb.RouteNode {
            Id:0,
            RouteId:in.RouteId,
            ActionType:db.NODE_ACTION_DELIVERY,
            NodeType:db.SUBSCRIPTION_NODE,
            NodeId:subscription.Id,
            NodeName:u.User.Name,
            Lat:a.Address.Lat,
            Lng:a.Address.Lng,
        });
    }
    out.Disconnected = disconnected;
    return nil;
}

func makeProtoRouteNode (node *db.RouteNode, lat float64, lng float64, name string, contact string, address string) (*pb.RouteNode){
    return &pb.RouteNode {
        Id         : node.Id,
        RouteId    : node.RouteId,
        ActionType : node.ActionType,
        NodeType   : node.NodeType,
        NodeId     : node.NodeId,
        NodeName   : name,
        NodeContact   : contact,
        NodeAddress   : address,
        Lat        : lat,
        Lng        : lng,
    }
}