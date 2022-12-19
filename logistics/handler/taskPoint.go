/*
--------------------------------
 Author: Deepak Rathod
--------------------------------
 Date - 30 August 2022
 Phone: +917020814070
 Email: rathoddeepak537@gmail.com 
--------------------------------
 --->   Logistics System   <---
--------------------------------
*/

package handler;

import (
	"context"

	pb "logisticsService/proto"
    db "logisticsService/database"
)

func (lg *LogisticsService) GetTaskPoints (ctx context.Context, in *pb.GetTaskPointRequest, out *pb.GetTaskPointsResponse) (error) {
    mPoints, err := db.GetTaskPoints(in.TaskId);
    if err != nil {
        return err
    }
    var points []*pb.TaskPoint;
    for _, mPoint := range *mPoints {
    	point := db.MakeProtoTaskPoint(mPoint);
    	points = append(points, point);
    }
    out.Points = points;
    return nil;
}

func (lg *LogisticsService) GetTaskPointById (ctx context.Context, in *pb.GetTaskPointRequest, out *pb.GetTaskPointResponse) (error) {
    task, err := db.GetTaskPointById(in.TaskPointId);
    if err != nil {
        return err
    }
    out.Point = db.MakeProtoTaskPoint(task);
    return nil;
}