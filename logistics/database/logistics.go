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
	"time"	
	"errors"	
	_ "github.com/lib/pq"	
	"database/sql"
	pb "logisticsService/proto"
	sqlBuilder "github.com/doug-martin/goqu/v9"
)

type Subscription struct {
	Id 	   			int64 		`json:"id"`
	UserId 			int64 		`json:"user_id"`
	ProductId 		int64		`json:"product_id"`
	StartFrom 	    int64 		`json:"start_from"`
	ExpEnd 	    	int64 		`json:"exp_end"`
	DlvIns 	    	int32 		`json:"dlv_ins"`
	DlvDays 	    int64 		`json:"dlv_days"`
	Sun 			float64 	`json:"sun"`
	Mon 			float64 	`json:"mon"`
	Tue 			float64 	`json:"tue"`
	Wed 			float64 	`json:"wed"`
	Thu 			float64 	`json:"thu"`
	Fri 			float64 	`json:"fri"`
	Sat 			float64 	`json:"sat"`	
	Pattern   		int64 		`json:"pattern"`
	Created   		int64 		`json:"created"`
	Status   		int32 		`json:"status"`
	AddressId   	int64 		`json:"address_id"`
	SlotId   		int64 		`json:"slot_id"`
	CityId   		int64 		`json:"city_id"`
	HubId   		int64 		`json:"hub_id"`
}

type GeoPoint struct {
	Type 		string 		`json:type`
	Coordinates []float64   `json:type`
}

func InsertSubscription(sub *Subscription) (*int64, error){
	sql, _, _ := sqlBuilder.Insert(SUBSCRIPTION_TABLE).Rows(
		sqlBuilder.Record {
			"user_id":sub.UserId,
			"product_id":sub.ProductId,
			"start_from":sub.StartFrom,
			"exp_end":sub.ExpEnd,
			"dlv_ins":sub.DlvIns,
			"dlv_days":sub.DlvDays,
			"sun":sub.Sun,
			"mon":sub.Mon,
			"tue":sub.Tue,
			"wed":sub.Wed,
			"thu":sub.Thu,
			"fri":sub.Fri,
			"sat":sub.Sat,
			"pattern":sub.Pattern,
			"address_id":sub.AddressId,
			"created":time.Now().Unix(),
			"slot_id":sub.SlotId,
			"city_id":sub.CityId,
			"hub_id":sub.HubId,
			"status":1,
		},
	).Returning("id").ToSQL();
	var subId int64;
	err := db.QueryRow(sql).Scan(&subId);
	if err != nil {
		return nil, err;
	}
	return &subId, nil
}

func UpdateSubscriptionStatus(sub *Subscription) (error){	
	sql, _, _ := sqlBuilder.Update(SUBSCRIPTION_TABLE).Set(
		sqlBuilder.Record {			
			"status" : sub.Status,
		},
	).Where(sqlBuilder.Ex {"id": sub.Id}).ToSQL();
	rows, err := db.Query(sql);
	defer rows.Close();
	if err != nil {
		return err;
	}
	return nil;
}

func UpdateSubscription(sub *Subscription) (error){	
	sql, _, _ := sqlBuilder.Update(SUBSCRIPTION_TABLE).Set(
		sqlBuilder.Record {			
			"start_from":sub.StartFrom,
			"dlv_ins":sub.DlvIns,
			"sun":sub.Sun,
			"mon":sub.Mon,
			"tue":sub.Tue,
			"wed":sub.Wed,
			"thu":sub.Thu,
			"fri":sub.Fri,
			"sat":sub.Sat,
			"pattern":sub.Pattern,
			"address_id":sub.AddressId,
			"slot_id":sub.SlotId,
			"exp_end":sub.ExpEnd,
		},
	).Where(sqlBuilder.Ex {"id": sub.Id}).ToSQL();
	rows, err := db.Query(sql);
	defer rows.Close();
	if err != nil {
		return err;
	}
	return nil;
}

func GetSubscriptionById(subscription_id int64) (*Subscription, error){	
	query, _, _ := sqlBuilder.Select(GetColumns()...).From(SUBSCRIPTION_TABLE).Where (
		sqlBuilder.Ex{"id": subscription_id},
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
	subscription, err := StructSubscription(row);
	if err != nil {
		return nil, err;
	}
	return subscription, nil;
}

func GetSubscriptions (user_id int64, status int64) (*[]Subscription, error){	
	query, _, _ := sqlBuilder.Select(GetColumns()...).From(SUBSCRIPTION_TABLE).Where (
		sqlBuilder.Ex{"user_id": user_id},
		sqlBuilder.Ex{"status": status},
	).ToSQL();
	rows, err := db.Query(query);
	if err != nil {
		return nil, err;
	}
	defer rows.Close();	
	var subscriptions []Subscription;
	for rows.Next() {
		if err := rows.Err(); err != nil {			
			continue
		}
		subscription, err := StructSubscription(rows);
		if err == nil {			
			subscriptions = append(subscriptions, *subscription);
		}else{			
			continue
		}
	}	
	return &subscriptions, nil
}

func GetHubSubscriptionsExclude (hub_id int64, exclude_ids []interface{}) (*[]Subscription, error){	
	dbSQL := sqlBuilder.Select(GetColumns()...).From(SUBSCRIPTION_TABLE);
	if len(exclude_ids) > 0 {
		dbSQL = dbSQL.Where(sqlBuilder.C("id").NotIn(exclude_ids...))
	}
	dbSQL = dbSQL.Where(sqlBuilder.Ex{"hub_id": hub_id});
	query, _, _ := dbSQL.ToSQL();	
	rows, err := db.Query(query);
	if err != nil {
		return nil, err;
	}
	defer rows.Close();	
	var subscriptions []Subscription;
	for rows.Next() {
		if err := rows.Err(); err != nil {			
			continue
		}
		subscription, err := StructSubscription(rows);
		if err == nil {			
			subscriptions = append(subscriptions, *subscription);
		}else{			
			continue
		}
	}	
	return &subscriptions, nil
}

func GetSubscriptionsBySlotId (slot_id int64) (*[]Subscription, error){
	query, _, _ := sqlBuilder.Select(GetColumns()...).From(SUBSCRIPTION_TABLE).Where (
		sqlBuilder.Ex{"slot_id": slot_id},
	).ToSQL();
	rows, err := db.Query(query);
	if err != nil {
		return nil, err;
	}
	defer rows.Close();
	var subscriptions []Subscription;
	for rows.Next() {
		if err := rows.Err(); err != nil {			
			continue
		}
		subscription, err := StructSubscription(rows);
		if err == nil {			
			subscriptions = append(subscriptions, *subscription);
		}else{			
			continue
		}
	}	
	return &subscriptions, nil
}

func CheckSubscriptionRunning(product_id int64, user_id int64) (bool){	
	query, _, _ := sqlBuilder.Select("id").From(SUBSCRIPTION_TABLE).Where (
		sqlBuilder.Ex{"product_id": product_id},
		sqlBuilder.Ex{"user_id": user_id},
	).ToSQL();
	row, err := db.Query(query);
	if err != nil {
		return false;
	}
	defer row.Close();
	present := row.Next();
	if present == false {
		return false;
	}
	if err := row.Err(); err != nil {
		return false;
	}	
	return true;
}

func FilterSubscriptions(filterMap *pb.FilterSubscriptionRequest, countOnly bool) (*[]*Subscription, *int64, error) {
	var subscriptions []*Subscription;
	
	dbSQL := sqlBuilder.From(SUBSCRIPTION_TABLE).Select(sqlBuilder.COUNT("*"));

	if countOnly == false {		
		dbSQL = sqlBuilder.From(SUBSCRIPTION_TABLE).Select(GetColumns()...);
	}

	if filterMap.CityId == 0 {
		return nil, nil, errors.New("City Id is required!");
	}else{
		dbSQL = dbSQL.Where(sqlBuilder.Ex {"city_id": filterMap.CityId})
	}

	if filterMap.SlotId != 0 {
		dbSQL = dbSQL.Where(sqlBuilder.Ex {"slot_id": filterMap.SlotId})
	}

	if filterMap.ProductId != 0 {
		dbSQL = dbSQL.Where(sqlBuilder.Ex {"product_id": filterMap.ProductId})
	}

	if filterMap.HubId != 0 {
		dbSQL = dbSQL.Where(sqlBuilder.Ex {"hub_id": filterMap.HubId})
	}

	if filterMap.FromDate != 0 {
		dbSQL = dbSQL.Where(sqlBuilder.C("exp_end").Gte(filterMap.FromDate))
		dbSQL = dbSQL.Where(sqlBuilder.C("start_from").Lte(filterMap.FromDate))
	}

	// if filterMap.ToDate != 0 {
		// dbSQL = dbSQL.Where(sqlBuilder.C("exp_end").Gte(filterMap.ToDate))
	// }

	if filterMap.SubscriptionId != 0 {
		dbSQL = dbSQL.Where(sqlBuilder.Ex {"id": filterMap.SubscriptionId})
	}

	if filterMap.Status != 9 {//TODO: 9 is temporary change later
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
		var count int64;
		err := db.QueryRow(query).Scan(&count);
		if err != nil {
			return nil, nil, err;
		}
		return nil, &count, nil;
	}
	
	query, _, _ := dbSQL.Order(sqlBuilder.I("id").Desc()).ToSQL();
	
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
		subscription, err := StructSubscription(rows);
		subscriptions = append(subscriptions, subscription);
	}
	return &subscriptions, nil, nil;
}

func StructSubscription (row *sql.Rows) (*Subscription, error) {
	var scp Subscription;
	err := row.Scan(	
		&scp.Id,	
		&scp.UserId,	
		&scp.ProductId,
		&scp.StartFrom,	
		&scp.DlvIns,	
		&scp.DlvDays,	
		&scp.Sun,	
		&scp.Mon,	
		&scp.Tue,	
		&scp.Wed,	
		&scp.Thu,	
		&scp.Fri,	
		&scp.Sat,	
		&scp.Pattern,	
		&scp.Created,	
		&scp.Status,
		&scp.AddressId,
		&scp.SlotId,
		&scp.CityId,
		&scp.HubId,
		&scp.ExpEnd,
	);
	if err != nil {
		return nil, err
	}
	return &scp, nil;
}