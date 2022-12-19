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

package database

import (
	_ "github.com/lib/pq"
	"go-micro.dev/v4/logger"
	"database/sql"
	"time"
	"errors"

	sqlBuilder "github.com/doug-martin/goqu/v9"
)

type Txn struct {
	Id 			 int64 		`json:"id"`
	UserId 	     int64		`json:"user_id"`
	TransType 	 bool 		`json:"trans_type"`
	Amount 	     float64 	`json:"amount"`
	Created      int64 		`json:created`
}

func CreateTXN (txn *Txn) (int64, error){
	timeUnix := time.Now().Unix();
	query, _, err := sqlBuilder.Insert(TXN_TABLE).Rows(
		sqlBuilder.Record {
			"user_id":txn.UserId,
			"trans_type":txn.TransType,
			"amount":txn.Amount,
			"created":timeUnix,
		},
	).ToSQL();
	if err != nil {
		logger.Fatal(err);
	}
	var lastInsertId int64;
	err = db.QueryRow(query + " RETURNING id").Scan(&lastInsertId);
	if err != nil {
		return 0, err
	}
	return lastInsertId, nil;
}

func GetTxnById (txn_id int64) (*Txn, error){
	query, _, _ := sqlBuilder.From(TXN_TABLE).Select("id", "user_id", "trans_type", "amount", "created").Where (
		sqlBuilder.Ex{"id": txn_id},
	).Limit(1).ToSQL();
	rows, err := db.Query(query);
	if err != nil {
		return nil, err;
	}
	defer rows.Close();
	present := rows.Next();
	if present == false {
		return nil, errors.New(errText);
	}
	if err := rows.Err(); err != nil {
		return nil, err;
	}

	txn, err := StructTxn(rows);
	if err != nil {
		return nil, err
	}
	return txn, nil
}

func GetUserTxn (user_id int64, limit int64, offset int64) (*[]*Txn, error){
	query, _, _ := sqlBuilder.From(TXN_TABLE).Select("id", "user_id", "trans_type", "amount", "created").Where (
		sqlBuilder.Ex{"user_id": user_id},
	).Limit(uint(limit)).Offset(uint(offset)).Order(sqlBuilder.I("id").Desc()).ToSQL();		
	rows, err := db.Query(query);
	if err != nil {
		return nil, err;
	}
	defer rows.Close();
	var txns []*Txn;
	for rows.Next() {		
		err := rows.Err();
		if err != nil {
			continue
		}
		txn, err := StructTxn(rows);
		txns = append(txns, txn);
	}
	return &txns, nil;
}

func StructTxn (row *sql.Rows) (*Txn, error) {
	var u Txn;
	err := row.Scan(
		&u.Id,&u.UserId,&u.TransType,&u.Amount,&u.Created,
	);	
	if err != nil {
		return nil, err
	}
	return &u, nil;
}