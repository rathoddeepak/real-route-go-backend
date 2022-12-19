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
	"errors"	
	_ "github.com/lib/pq"	
	"database/sql"
	sqlBuilder "github.com/doug-martin/goqu/v9"
)

type Slot struct {
	Id 	   			int64 		`json:"id"`
	CityId 			int64 		`json:"city_id"`
	Title 	   		string 		`json:"title"`
	StartHr         int32 		`json:"start_hr"`
	StartMin        int32 		`json:"start_min"`
	EndHr         	int32 		`json:"end_hr"`
	EndMin        	int32 		`json:"end_min"`
	Status   		int32 		`json:"status"`
}

func InsertSlot(slot *Slot) (*int64, error){
	sql, _, _ := sqlBuilder.Insert(SLOT_TABLE).Rows(
		sqlBuilder.Record {
			"title":slot.Title,
			"city_id":slot.CityId,
			"start_hr":slot.StartHr,
			"start_min":slot.EndMin,
			"end_hr":slot.EndHr,
			"end_min":slot.EndMin,
			"status":0,
		},
	).Returning("id").ToSQL();
	var slot_id int64;
	err := db.QueryRow(sql).Scan(&slot_id);
	if err != nil {
		return nil, err;
	}
	return &slot_id, nil
}

func UpdateSlot(slot *Slot) (error){	
	sql, _, _ := sqlBuilder.Update(SLOT_TABLE).Set(
		sqlBuilder.Record {			
			"title":slot.Title,
			"start_hr":slot.StartHr,
			"start_min":slot.EndMin,
			"end_hr":slot.EndHr,
			"end_min":slot.EndMin,
			"status":slot.Status,
		},
	).Where(sqlBuilder.Ex {"id": slot.Id}).ToSQL();
	rows, err := db.Query(sql);
	defer rows.Close();
	if err != nil {
		return err;
	}
	return nil;
}

func GetSlotById(slot_id int64) (*Slot, error){	
	query, _, _ := sqlBuilder.Select("id","title","city_id","start_hr","start_min","end_hr","end_min","status").From(SLOT_TABLE).Where (
		sqlBuilder.Ex{"id": slot_id},
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
	slot, err := StructSlot(row);
	if err != nil {
		return nil, err;
	}
	return slot, nil;
}

func GetSlotsOfCity (city_id int64) (*[]*Slot, error){
	query, _, _ := sqlBuilder.Select("id","title","city_id","start_hr","start_min","end_hr","end_min","status").From(SLOT_TABLE).Where (
		sqlBuilder.Ex{"city_id": city_id},
	).ToSQL();
	rows, err := db.Query(query);
	if err != nil {
		return nil, err;
	}
	defer rows.Close();
	var slots []*Slot;
	for rows.Next() {
		if err := rows.Err(); err != nil {			
			continue
		}
		slot, err := StructSlot(rows);
		if err == nil {			
			slots = append(slots, slot);
		}else{			
			continue
		}
	}	
	return &slots, nil
}

func DeleteSlotById(slot_id int64) (error){	
	sql, _, _ := sqlBuilder.Delete(SLOT_TABLE).Where(sqlBuilder.Ex {"id": slot_id}).ToSQL();
	rows, err := db.Query(sql);
	defer rows.Close();
	if err != nil {
		return err;
	}
	return nil;
}

func StructSlot (row *sql.Rows) (*Slot, error) {
	var slot Slot;
	err := row.Scan(	
		&slot.Id,	
		&slot.Title,	
		&slot.CityId,
		&slot.StartHr,	
		&slot.StartMin,	
		&slot.EndHr,	
		&slot.EndMin,
		&slot.Status,
	);
	if err != nil {
		return nil, err
	}
	return &slot, nil;
}